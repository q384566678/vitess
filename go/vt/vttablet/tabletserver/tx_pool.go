/*
Copyright 2019 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tabletserver

import (
	"fmt"
	"sync"
	"time"

	"vitess.io/vitess/go/pools"

	"vitess.io/vitess/go/vt/servenv"

	"vitess.io/vitess/go/vt/vttablet/tabletserver/tx"

	"golang.org/x/net/context"

	"vitess.io/vitess/go/sync2"
	"vitess.io/vitess/go/timer"
	"vitess.io/vitess/go/trace"
	"vitess.io/vitess/go/vt/callerid"
	"vitess.io/vitess/go/vt/dbconfigs"
	"vitess.io/vitess/go/vt/log"
	"vitess.io/vitess/go/vt/vterrors"
	"vitess.io/vitess/go/vt/vttablet/tabletserver/tabletenv"
	"vitess.io/vitess/go/vt/vttablet/tabletserver/txlimiter"

	querypb "vitess.io/vitess/go/vt/proto/query"
	vtrpcpb "vitess.io/vitess/go/vt/proto/vtrpc"
)

// These consts identify how a transaction was resolved.
//const (
//	TxClose      = "close"
//	TxCommit     = "commit"
//	TxRollback   = "rollback"
//	TxKill       = "kill"
//	ConnInitFail = "initFail"
//)

const txLogInterval = 1 * time.Minute

type queries struct {
	setIsolationLevel string
	openTransaction   string
}

var (
	txIsolations = map[querypb.ExecuteOptions_TransactionIsolation]queries{
		querypb.ExecuteOptions_DEFAULT:                       {setIsolationLevel: "", openTransaction: "begin"},
		querypb.ExecuteOptions_REPEATABLE_READ:               {setIsolationLevel: "REPEATABLE READ", openTransaction: "begin"},
		querypb.ExecuteOptions_READ_COMMITTED:                {setIsolationLevel: "READ COMMITTED", openTransaction: "begin"},
		querypb.ExecuteOptions_READ_UNCOMMITTED:              {setIsolationLevel: "READ UNCOMMITTED", openTransaction: "begin"},
		querypb.ExecuteOptions_SERIALIZABLE:                  {setIsolationLevel: "SERIALIZABLE", openTransaction: "begin"},
		querypb.ExecuteOptions_CONSISTENT_SNAPSHOT_READ_ONLY: {setIsolationLevel: "REPEATABLE READ", openTransaction: "start transaction with consistent snapshot, read only"},
	}
)

// TxPool is the transaction pool for the query service.
type TxPool struct {
	env tabletenv.Env

	activePool         *StatefulConnectionPool
	transactionTimeout sync2.AtomicDuration
	ticks              *timer.Timer
	limiter            txlimiter.TxLimiter

	logMu   sync.Mutex
	lastLog time.Time
	txStats *servenv.TimingsWrapper
}

// NewTxPool creates a new TxPool. It's not operational until it's Open'd.
func NewTxPool(env tabletenv.Env, limiter txlimiter.TxLimiter) *TxPool {
	config := env.Config()
	transactionTimeout := time.Duration(config.Oltp.TxTimeoutSeconds * 1e9)
	axp := &TxPool{
		env:                env,
		activePool:         NewStatefulConnPool(env),
		transactionTimeout: sync2.NewAtomicDuration(transactionTimeout),
		ticks:              timer.NewTimer(transactionTimeout / 10),
		limiter:            limiter,
		txStats:            env.Exporter().NewTimings("Transactions", "Transaction stats", "operation"),
	}
	// Careful: conns also exports name+"xxx" vars,
	// but we know it doesn't export Timeout.
	env.Exporter().NewGaugeDurationFunc("TransactionTimeout", "Transaction timeout", axp.transactionTimeout.Get)
	return axp
}

// Open makes the TxPool operational. This also starts the transaction killer
// that will kill long-running transactions.
func (tp *TxPool) Open(appParams, dbaParams, appDebugParams dbconfigs.Connector) {
	tp.activePool.Open(appParams, dbaParams, appDebugParams)
	tp.ticks.Start(func() { tp.transactionKiller() })
}

// Close closes the TxPool. A closed pool can be reopened.
func (tp *TxPool) Close() {
	tp.ticks.Stop()
	tp.activePool.Close()
}

// AdjustLastID adjusts the last transaction id to be at least
// as large as the input value. This will ensure that there are
// no dtid collisions with future transactions.
func (tp *TxPool) AdjustLastID(id int64) {
	tp.activePool.AdjustLastID(id)
}

// RollbackNonBusy rolls back all transactions that are not in use.
// Transactions can be in use for situations like executing statements
// or in prepared state.
func (tp *TxPool) RollbackNonBusy(ctx context.Context) {
	for _, v := range tp.activePool.GetOutdated(time.Duration(0), "for transition") {
		tp.LocalConclude(ctx, v)
	}
}

func (tp *TxPool) transactionKiller() {
	defer tp.env.LogError()
	for _, conn := range tp.activePool.GetOutdated(tp.Timeout(), "for tx killer rollback") {
		log.Warningf("killing transaction (exceeded timeout: %v): %s", tp.Timeout(), conn.String())
		tp.env.Stats().KillCounters.Add("Transactions", 1)
		conn.Close()
		conn.conclude(fmt.Sprintf("exceeded timeout: %v", tp.Timeout()))
	}
}

// WaitForEmpty waits until all active transactions are completed.
func (tp *TxPool) WaitForEmpty() {
	tp.activePool.WaitForEmpty()
}

//NewTxProps creates a new TxProperties struct
func (tp *TxPool) NewTxProps(immediateCaller *querypb.VTGateCallerID, effectiveCaller *vtrpcpb.CallerID, autocommit bool) *TxProperties {
	return &TxProperties{
		StartTime:       time.Now(),
		EffectiveCaller: effectiveCaller,
		ImmediateCaller: immediateCaller,
		Autocommit:      autocommit,
		txStats:         tp.txStats,
	}
}

// GetAndLock fetches the connection associated to the transactionID and blocks it from concurrent use
// You must call Unlock on TxConnection once done.
func (tp *TxPool) GetAndLock(connID tx.ConnID, reason string) (*StatefulConnection, error) {
	conn, err := tp.activePool.GetAndLock(connID, reason)
	if err != nil {
		return nil, vterrors.Errorf(vtrpcpb.Code_ABORTED, "transaction %d: %v", connID, err)
	}
	return conn, nil
}

// Commit commits the transaction on the specified connection.
func (tp *TxPool) Commit(ctx context.Context, connID tx.ConnID) (string, error) {
	span, ctx := trace.NewSpan(ctx, "TxPool.Commit")
	defer span.Finish()
	conn, err := tp.GetAndLock(connID, "for commit")
	if err != nil {
		return "", err
	}
	return tp.LocalCommit(ctx, conn)
}

// LocalCommit commits the transaction on the connection. The connection will be either Release:ed or Unlock:ed,
// depending on if the connection is tainted or not.
func (tp *TxPool) LocalCommit(ctx context.Context, txConn *StatefulConnection) (string, error) {
	span, ctx := trace.NewSpan(ctx, "TxPool.LocalCommit")
	defer span.Finish()
	defer tp.txComplete(txConn, tx.TxCommit)
	if txConn.TxProps.Autocommit {
		return "", nil
	}

	if _, err := txConn.Exec(ctx, "commit", 1, false); err != nil {
		txConn.Close()
		return "", err
	}
	return "commit", nil
}

// Rollback rolls back the transaction on the specified connection.
func (tp *TxPool) Rollback(ctx context.Context, connID tx.ConnID) error {
	span, ctx := trace.NewSpan(ctx, "TxPool.Rollback")
	defer span.Finish()

	conn, err := tp.GetAndLock(connID, "for rollback")
	if err != nil {
		return err
	}
	return tp.localRollback(ctx, conn)
}

// Begin begins a transaction, and returns the associated connection and
// the statements (if any) executed to initiate the transaction. In autocommit
// mode the statement will be "".
// The connection returned is locked for the callee and its responsibility is to unlock the connection.
func (tp *TxPool) Begin(ctx context.Context, options *querypb.ExecuteOptions) (*StatefulConnection, string, error) {
	span, ctx := trace.NewSpan(ctx, "TxPool.Begin")
	defer span.Finish()
	beginQueries := ""

	immediateCaller := callerid.ImmediateCallerIDFromContext(ctx)
	effectiveCaller := callerid.EffectiveCallerIDFromContext(ctx)

	if !tp.limiter.Get(immediateCaller, effectiveCaller) {
		return nil, "", vterrors.Errorf(vtrpcpb.Code_RESOURCE_EXHAUSTED, "per-user transaction pool connection limit exceeded")
	}

	conn, err := tp.activePool.NewConn(ctx, options)
	if err != nil {
		switch err {
		case pools.ErrCtxTimeout:
			tp.LogActive()
			err = vterrors.Errorf(vtrpcpb.Code_RESOURCE_EXHAUSTED, "transaction pool aborting request due to already expired context")
		case pools.ErrTimeout:
			tp.LogActive()
			err = vterrors.Errorf(vtrpcpb.Code_RESOURCE_EXHAUSTED, "transaction pool connection limit exceeded")
		}
		return nil, "", err
	}
	err = func() error {
		autocommitTransaction := false
		if queries, ok := txIsolations[options.GetTransactionIsolation()]; ok {
			if queries.setIsolationLevel != "" {
				txQuery := "set transaction isolation level " + queries.setIsolationLevel
				if err := conn.execWithRetry(ctx, txQuery, 1, false); err != nil {
					return vterrors.Wrap(err, txQuery)
				}
				beginQueries = queries.setIsolationLevel + "; "
			}
			if err := conn.execWithRetry(ctx, queries.openTransaction, 1, false); err != nil {
				return vterrors.Wrap(err, queries.openTransaction)
			}
			beginQueries = beginQueries + queries.openTransaction
		} else if options.GetTransactionIsolation() == querypb.ExecuteOptions_AUTOCOMMIT {
			autocommitTransaction = true
		} else {
			return vterrors.Errorf(vtrpcpb.Code_INTERNAL, "don't know how to open a transaction of this type: %v", options.GetTransactionIsolation())
		}
		conn.TxProps = tp.NewTxProps(immediateCaller, effectiveCaller, autocommitTransaction)
		return nil
	}()
	if err != nil {
		conn.Close()
		conn.Release(tx.ConnInitFail)
		return nil, "", err
	}

	return conn, beginQueries, nil
}

// LocalConclude concludes a transaction started by Begin.
// If the transaction was not previously concluded, it's rolled back.
func (tp *TxPool) LocalConclude(ctx context.Context, conn *StatefulConnection) {
	if conn.dbConn == nil {
		return
	}
	span, ctx := trace.NewSpan(ctx, "TxPool.LocalConclude")
	defer span.Finish()
	_ = tp.localRollback(ctx, conn)
}

func (tp *TxPool) localRollback(ctx context.Context, txConn *StatefulConnection) error {
	if txConn.TxProps.Autocommit {
		tp.txComplete(txConn, tx.TxCommit)
		return nil
	}
	defer tp.txComplete(txConn, tx.TxRollback)
	if _, err := txConn.Exec(ctx, "rollback", 1, false); err != nil {
		txConn.Close()
		return err
	}
	return nil
}

// LogActive causes all existing transactions to be logged when they complete.
// The logging is throttled to no more than once every txLogInterval.
func (tp *TxPool) LogActive() {
	tp.logMu.Lock()
	defer tp.logMu.Unlock()
	if time.Since(tp.lastLog) < txLogInterval {
		return
	}
	tp.lastLog = time.Now()
	tp.activePool.ForAllTxProperties(func(props *TxProperties) {
		props.LogToFile = true
	})
}

// Timeout returns the transaction timeout.
func (tp *TxPool) Timeout() time.Duration {
	return tp.transactionTimeout.Get()
}

// SetTimeout sets the transaction timeout.
func (tp *TxPool) SetTimeout(timeout time.Duration) {
	tp.transactionTimeout.Set(timeout)
	tp.ticks.SetInterval(timeout / 10)
}

func (tp *TxPool) txComplete(conn *StatefulConnection, reason tx.ReleaseReason) {
	tp.log(conn, reason)
	if conn.tainted {
		conn.renewConnection()
	} else {
		conn.Release(reason)
	}
	tp.limiter.Release(conn.TxProps.ImmediateCaller, conn.TxProps.EffectiveCaller)
	conn.txClean()
}

func (tp *TxPool) log(txc *StatefulConnection, reason tx.ReleaseReason) {
	if txc.TxProps == nil {
		return //Nothing to log as no transaction exists on this connection.
	}
	txc.TxProps.Conclusion = reason.Name()
	txc.TxProps.EndTime = time.Now()

	username := callerid.GetPrincipal(txc.TxProps.EffectiveCaller)
	if username == "" {
		username = callerid.GetUsername(txc.TxProps.ImmediateCaller)
	}
	duration := txc.TxProps.EndTime.Sub(txc.TxProps.StartTime)
	txc.env.Stats().UserTransactionCount.Add([]string{username, reason.Name()}, 1)
	txc.env.Stats().UserTransactionTimesNs.Add([]string{username, reason.Name()}, int64(duration))
	txc.TxProps.txStats.Add(reason.Name(), duration)
	if txc.TxProps.LogToFile {
		log.Infof("Logged transaction: %s", txc.String())
	}
	tabletenv.TxLogger.Send(txc)
}

//TxProperties contains all information that is related to the currently running
//transaction on the connection
type TxProperties struct {
	EffectiveCaller *vtrpcpb.CallerID
	ImmediateCaller *querypb.VTGateCallerID

	StartTime time.Time
	EndTime   time.Time

	Queries []string

	Autocommit bool
	Conclusion string

	LogToFile bool

	txStats *servenv.TimingsWrapper
}

// RecordQuery records the query against this transaction.
func (tp *TxProperties) RecordQuery(query string) {
	tp.Queries = append(tp.Queries, query)
}

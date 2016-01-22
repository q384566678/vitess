<?php
// DO NOT EDIT! Generated by Protobuf-PHP protoc plugin 1.0
// Source: topodata.proto
//   Date: 2016-01-22 01:34:42

namespace Vitess\Proto\Topodata\Shard {

  class SourceShard extends \DrSlump\Protobuf\Message {

    /**  @var int */
    public $uid = null;
    
    /**  @var string */
    public $keyspace = null;
    
    /**  @var string */
    public $shard = null;
    
    /**  @var \Vitess\Proto\Topodata\KeyRange */
    public $key_range = null;
    
    /**  @var string[]  */
    public $tables = array();
    

    /** @var \Closure[] */
    protected static $__extensions = array();

    public static function descriptor()
    {
      $descriptor = new \DrSlump\Protobuf\Descriptor(__CLASS__, 'topodata.Shard.SourceShard');

      // OPTIONAL UINT32 uid = 1
      $f = new \DrSlump\Protobuf\Field();
      $f->number    = 1;
      $f->name      = "uid";
      $f->type      = \DrSlump\Protobuf::TYPE_UINT32;
      $f->rule      = \DrSlump\Protobuf::RULE_OPTIONAL;
      $descriptor->addField($f);

      // OPTIONAL STRING keyspace = 2
      $f = new \DrSlump\Protobuf\Field();
      $f->number    = 2;
      $f->name      = "keyspace";
      $f->type      = \DrSlump\Protobuf::TYPE_STRING;
      $f->rule      = \DrSlump\Protobuf::RULE_OPTIONAL;
      $descriptor->addField($f);

      // OPTIONAL STRING shard = 3
      $f = new \DrSlump\Protobuf\Field();
      $f->number    = 3;
      $f->name      = "shard";
      $f->type      = \DrSlump\Protobuf::TYPE_STRING;
      $f->rule      = \DrSlump\Protobuf::RULE_OPTIONAL;
      $descriptor->addField($f);

      // OPTIONAL MESSAGE key_range = 4
      $f = new \DrSlump\Protobuf\Field();
      $f->number    = 4;
      $f->name      = "key_range";
      $f->type      = \DrSlump\Protobuf::TYPE_MESSAGE;
      $f->rule      = \DrSlump\Protobuf::RULE_OPTIONAL;
      $f->reference = '\Vitess\Proto\Topodata\KeyRange';
      $descriptor->addField($f);

      // REPEATED STRING tables = 5
      $f = new \DrSlump\Protobuf\Field();
      $f->number    = 5;
      $f->name      = "tables";
      $f->type      = \DrSlump\Protobuf::TYPE_STRING;
      $f->rule      = \DrSlump\Protobuf::RULE_REPEATED;
      $descriptor->addField($f);

      foreach (self::$__extensions as $cb) {
        $descriptor->addField($cb(), true);
      }

      return $descriptor;
    }

    /**
     * Check if <uid> has a value
     *
     * @return boolean
     */
    public function hasUid(){
      return $this->_has(1);
    }
    
    /**
     * Clear <uid> value
     *
     * @return \Vitess\Proto\Topodata\Shard\SourceShard
     */
    public function clearUid(){
      return $this->_clear(1);
    }
    
    /**
     * Get <uid> value
     *
     * @return int
     */
    public function getUid(){
      return $this->_get(1);
    }
    
    /**
     * Set <uid> value
     *
     * @param int $value
     * @return \Vitess\Proto\Topodata\Shard\SourceShard
     */
    public function setUid( $value){
      return $this->_set(1, $value);
    }
    
    /**
     * Check if <keyspace> has a value
     *
     * @return boolean
     */
    public function hasKeyspace(){
      return $this->_has(2);
    }
    
    /**
     * Clear <keyspace> value
     *
     * @return \Vitess\Proto\Topodata\Shard\SourceShard
     */
    public function clearKeyspace(){
      return $this->_clear(2);
    }
    
    /**
     * Get <keyspace> value
     *
     * @return string
     */
    public function getKeyspace(){
      return $this->_get(2);
    }
    
    /**
     * Set <keyspace> value
     *
     * @param string $value
     * @return \Vitess\Proto\Topodata\Shard\SourceShard
     */
    public function setKeyspace( $value){
      return $this->_set(2, $value);
    }
    
    /**
     * Check if <shard> has a value
     *
     * @return boolean
     */
    public function hasShard(){
      return $this->_has(3);
    }
    
    /**
     * Clear <shard> value
     *
     * @return \Vitess\Proto\Topodata\Shard\SourceShard
     */
    public function clearShard(){
      return $this->_clear(3);
    }
    
    /**
     * Get <shard> value
     *
     * @return string
     */
    public function getShard(){
      return $this->_get(3);
    }
    
    /**
     * Set <shard> value
     *
     * @param string $value
     * @return \Vitess\Proto\Topodata\Shard\SourceShard
     */
    public function setShard( $value){
      return $this->_set(3, $value);
    }
    
    /**
     * Check if <key_range> has a value
     *
     * @return boolean
     */
    public function hasKeyRange(){
      return $this->_has(4);
    }
    
    /**
     * Clear <key_range> value
     *
     * @return \Vitess\Proto\Topodata\Shard\SourceShard
     */
    public function clearKeyRange(){
      return $this->_clear(4);
    }
    
    /**
     * Get <key_range> value
     *
     * @return \Vitess\Proto\Topodata\KeyRange
     */
    public function getKeyRange(){
      return $this->_get(4);
    }
    
    /**
     * Set <key_range> value
     *
     * @param \Vitess\Proto\Topodata\KeyRange $value
     * @return \Vitess\Proto\Topodata\Shard\SourceShard
     */
    public function setKeyRange(\Vitess\Proto\Topodata\KeyRange $value){
      return $this->_set(4, $value);
    }
    
    /**
     * Check if <tables> has a value
     *
     * @return boolean
     */
    public function hasTables(){
      return $this->_has(5);
    }
    
    /**
     * Clear <tables> value
     *
     * @return \Vitess\Proto\Topodata\Shard\SourceShard
     */
    public function clearTables(){
      return $this->_clear(5);
    }
    
    /**
     * Get <tables> value
     *
     * @param int $idx
     * @return string
     */
    public function getTables($idx = NULL){
      return $this->_get(5, $idx);
    }
    
    /**
     * Set <tables> value
     *
     * @param string $value
     * @return \Vitess\Proto\Topodata\Shard\SourceShard
     */
    public function setTables( $value, $idx = NULL){
      return $this->_set(5, $value, $idx);
    }
    
    /**
     * Get all elements of <tables>
     *
     * @return string[]
     */
    public function getTablesList(){
     return $this->_get(5);
    }
    
    /**
     * Add a new element to <tables>
     *
     * @param string $value
     * @return \Vitess\Proto\Topodata\Shard\SourceShard
     */
    public function addTables( $value){
     return $this->_add(5, $value);
    }
  }
}


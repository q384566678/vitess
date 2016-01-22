<?php
// DO NOT EDIT! Generated by Protobuf-PHP protoc plugin 1.0
// Source: vtgate.proto
//   Date: 2016-01-22 01:34:42

namespace Vitess\Proto\Vtgate\SplitQueryResponse {

  class Part extends \DrSlump\Protobuf\Message {

    /**  @var \Vitess\Proto\Query\BoundQuery */
    public $query = null;
    
    /**  @var \Vitess\Proto\Vtgate\SplitQueryResponse\KeyRangePart */
    public $key_range_part = null;
    
    /**  @var \Vitess\Proto\Vtgate\SplitQueryResponse\ShardPart */
    public $shard_part = null;
    
    /**  @var int */
    public $size = null;
    

    /** @var \Closure[] */
    protected static $__extensions = array();

    public static function descriptor()
    {
      $descriptor = new \DrSlump\Protobuf\Descriptor(__CLASS__, 'vtgate.SplitQueryResponse.Part');

      // OPTIONAL MESSAGE query = 1
      $f = new \DrSlump\Protobuf\Field();
      $f->number    = 1;
      $f->name      = "query";
      $f->type      = \DrSlump\Protobuf::TYPE_MESSAGE;
      $f->rule      = \DrSlump\Protobuf::RULE_OPTIONAL;
      $f->reference = '\Vitess\Proto\Query\BoundQuery';
      $descriptor->addField($f);

      // OPTIONAL MESSAGE key_range_part = 2
      $f = new \DrSlump\Protobuf\Field();
      $f->number    = 2;
      $f->name      = "key_range_part";
      $f->type      = \DrSlump\Protobuf::TYPE_MESSAGE;
      $f->rule      = \DrSlump\Protobuf::RULE_OPTIONAL;
      $f->reference = '\Vitess\Proto\Vtgate\SplitQueryResponse\KeyRangePart';
      $descriptor->addField($f);

      // OPTIONAL MESSAGE shard_part = 3
      $f = new \DrSlump\Protobuf\Field();
      $f->number    = 3;
      $f->name      = "shard_part";
      $f->type      = \DrSlump\Protobuf::TYPE_MESSAGE;
      $f->rule      = \DrSlump\Protobuf::RULE_OPTIONAL;
      $f->reference = '\Vitess\Proto\Vtgate\SplitQueryResponse\ShardPart';
      $descriptor->addField($f);

      // OPTIONAL INT64 size = 4
      $f = new \DrSlump\Protobuf\Field();
      $f->number    = 4;
      $f->name      = "size";
      $f->type      = \DrSlump\Protobuf::TYPE_INT64;
      $f->rule      = \DrSlump\Protobuf::RULE_OPTIONAL;
      $descriptor->addField($f);

      foreach (self::$__extensions as $cb) {
        $descriptor->addField($cb(), true);
      }

      return $descriptor;
    }

    /**
     * Check if <query> has a value
     *
     * @return boolean
     */
    public function hasQuery(){
      return $this->_has(1);
    }
    
    /**
     * Clear <query> value
     *
     * @return \Vitess\Proto\Vtgate\SplitQueryResponse\Part
     */
    public function clearQuery(){
      return $this->_clear(1);
    }
    
    /**
     * Get <query> value
     *
     * @return \Vitess\Proto\Query\BoundQuery
     */
    public function getQuery(){
      return $this->_get(1);
    }
    
    /**
     * Set <query> value
     *
     * @param \Vitess\Proto\Query\BoundQuery $value
     * @return \Vitess\Proto\Vtgate\SplitQueryResponse\Part
     */
    public function setQuery(\Vitess\Proto\Query\BoundQuery $value){
      return $this->_set(1, $value);
    }
    
    /**
     * Check if <key_range_part> has a value
     *
     * @return boolean
     */
    public function hasKeyRangePart(){
      return $this->_has(2);
    }
    
    /**
     * Clear <key_range_part> value
     *
     * @return \Vitess\Proto\Vtgate\SplitQueryResponse\Part
     */
    public function clearKeyRangePart(){
      return $this->_clear(2);
    }
    
    /**
     * Get <key_range_part> value
     *
     * @return \Vitess\Proto\Vtgate\SplitQueryResponse\KeyRangePart
     */
    public function getKeyRangePart(){
      return $this->_get(2);
    }
    
    /**
     * Set <key_range_part> value
     *
     * @param \Vitess\Proto\Vtgate\SplitQueryResponse\KeyRangePart $value
     * @return \Vitess\Proto\Vtgate\SplitQueryResponse\Part
     */
    public function setKeyRangePart(\Vitess\Proto\Vtgate\SplitQueryResponse\KeyRangePart $value){
      return $this->_set(2, $value);
    }
    
    /**
     * Check if <shard_part> has a value
     *
     * @return boolean
     */
    public function hasShardPart(){
      return $this->_has(3);
    }
    
    /**
     * Clear <shard_part> value
     *
     * @return \Vitess\Proto\Vtgate\SplitQueryResponse\Part
     */
    public function clearShardPart(){
      return $this->_clear(3);
    }
    
    /**
     * Get <shard_part> value
     *
     * @return \Vitess\Proto\Vtgate\SplitQueryResponse\ShardPart
     */
    public function getShardPart(){
      return $this->_get(3);
    }
    
    /**
     * Set <shard_part> value
     *
     * @param \Vitess\Proto\Vtgate\SplitQueryResponse\ShardPart $value
     * @return \Vitess\Proto\Vtgate\SplitQueryResponse\Part
     */
    public function setShardPart(\Vitess\Proto\Vtgate\SplitQueryResponse\ShardPart $value){
      return $this->_set(3, $value);
    }
    
    /**
     * Check if <size> has a value
     *
     * @return boolean
     */
    public function hasSize(){
      return $this->_has(4);
    }
    
    /**
     * Clear <size> value
     *
     * @return \Vitess\Proto\Vtgate\SplitQueryResponse\Part
     */
    public function clearSize(){
      return $this->_clear(4);
    }
    
    /**
     * Get <size> value
     *
     * @return int
     */
    public function getSize(){
      return $this->_get(4);
    }
    
    /**
     * Set <size> value
     *
     * @param int $value
     * @return \Vitess\Proto\Vtgate\SplitQueryResponse\Part
     */
    public function setSize( $value){
      return $this->_set(4, $value);
    }
  }
}


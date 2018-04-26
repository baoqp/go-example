package bplus_tree

import (
	"util"
)


var HeaderSize = uint64(32)
var KVHeaderSize = uint64(24)

//  bp_key_s  bp_key_t
type Key struct {
	length     uint64
	value      []byte
	prevOffset uint64
	prevLength uint64
}

func NewKey(key []byte) *Key {
	return &Key {
		value : key,
		length: uint64(len(key)),
	}
}

// bp_value_t
type Value Key

func NewValue(value []byte) *Value{
	return &Value {
		value : value,
		length: uint64(len(value)),
	}
}

// bp__kv_s bp__kv_t
type KV struct {
	length    uint64  // 数据长度
	value     []byte
	offset    uint64  // 在文件中的offset
	config    uint64
	allocated bool
}

func kvSize(kv *KV) uint64 {
	return KVHeaderSize + kv.length
}

func kvCopy(source *KV, target *KV, alloc bool) {
	if alloc {
		target.value = make([]byte, len(source.value))
		copy(target.value, source.value)
		target.allocated = true
	} else {
		target.value = source.value // 复制引用而不是值
	}
	target.length = source.length
	target.offset = source.offset
	target.config = source.config

}

func valueLoad(tree *Tree, offset uint64, length uint64, value *Value) error {

	// read data from disk first
	bufLen := length
	buff, err := tree.writerRead(DefaultComp, offset, &bufLen)

	if err != nil {
		return err
	}
	value.value = make([]byte, bufLen-16)
	// first 16 bytes are representing previous value
	value.prevOffset = util.DecodeUint64(buff, 0)
	value.prevLength = util.DecodeUint64(buff, 8)
	copy(value.value, buff[16:])
	value.length = bufLen - 16

	return nil
}

func valueSave(tree *Tree, value *Value, previous *KV, offset *uint64,
	length *uint64) error {

	buff := make([]byte, value.length+16)
	if previous != nil {
		util.EncodeUint64(buff, 0, previous.offset)
		util.EncodeUint64(buff, 8, previous.length)
	} else {
		util.EncodeUint64(buff, 0, uint64(0))
		util.EncodeUint64(buff, 8, uint64(0))
	}

	copy(buff[16: 16+value.length], value.value)

	*length = value.length + 16
	return tree.writerWrite(DefaultComp, buff, offset, length)
}

package bplus_tree

import (
	"unsafe"
	"util"
)

const (
	HeaderSize = 24
)

//  bp_key_s  bp_key_t
type Key struct {
	length     uint64
	value      []byte
	prevOffset uint64
	prevLength uint64
}

// bp_value_t
type Value Key

// bp__kv_s bp__kv_t
type KV struct {
	length    uint64
	value     []byte
	offset    uint64
	config    uint64
	allocated bool
}

func kvSize(kv *KV) uint64 {
	return HeaderSize + kv.length
}

func kvCopy(source *KV, target *KV, alloc bool) error {
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

	return nil
}

func valueLoad(db *DB, offset uint64, length uint64, value *Value) error {
	//cast db to writer
	p := unsafe.Pointer(db)
	w := (*Writer)(p)

	// read data from disk first
	bufLen := length
	buff, err := writerRead(w, DefaultComp, offset, &bufLen)
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

func valueSave(db *DB, value *Value, previous *KV, offset *uint64,
	length *uint64) error {

	//cast db to writer
	p := unsafe.Pointer(db)
	w := (*Writer)(p)

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
	return writerWrite(w, DefaultComp, buff, offset, length)
}

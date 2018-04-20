package bplus_tree

const (
	HeaderSize = 24
)

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


func kvClone(source *KV, alloc bool ) *KV{
	clone := new(KV)

	if alloc {
		clone.value = make([]byte, len(source.value))
		copy(clone.value, source.value)
		clone.allocated = true
	} else {
		clone.value = source.value // TODO 复制引用而不是值
	}
	clone.length = source.length
	clone.offset = source.offset
	clone.config = source.config
	return clone
}
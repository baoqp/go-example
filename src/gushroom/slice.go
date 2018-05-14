package gushroom

type KeySlice struct {
	pageNo pageId
	data   []byte
}

func (ks *KeySlice) AssignPageNo(pageNo pageId) {
	ks.pageNo = pageNo
}

func (ks *KeySlice) Assign(pageNo pageId, data []byte) {
	ks.pageNo = pageNo
	ks.data = data
}

type DataSlice struct {
	len  uint16
	data []byte
}

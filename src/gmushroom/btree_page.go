package gmushroom

const (
	ROOT   = 0
	BRANCH = 1
	LEAF   = 2
)

type BTreePage struct {
	pageNo   PageId
	first    PageId
	totalKay uint16
	keyLen   uint8
	level    uint8
	typ      byte
	dirty    byte
	occupy   byte
	lock     byte
	reders   byte
	data     []byte
}

func NewBTreePage() *BTreePage {
	return &BTreePage{
		typ:    LEAF,
		dirty:  1,
		occupy: 1,
		lock:   1,
		reders: 3,
	}
}

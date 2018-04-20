package bplus_tree

type PageType int

const (
	kPage PageType = 0
	kLeaf PageType = 1
)

type SearchType int

const (
	kNotLoad SearchType = 0
	kLoad    SearchType = 1
)

// bp__page_s bp__page_t
type Page struct {
	typ      PageType
	length   uint64
	byteSize uint64
	offset   uint64
	config   uint64
	buff     []byte
	isHead   bool
	keys     []*KV
}




func pageCreate(db *DB, typ PageType, offset uint64, config uint64) *Page {
	page := new(Page)
	page.typ = typ
	if typ == kLeaf {
		page.length = 0
		page.byteSize = 0
	} else {
		/* TODO non-leaf pages always have left element */
		page.length = 1
		kv := &KV{
			value:nil,
			offset:0,
			length:0,
			config:0,
			allocated:0,
		}
		page.keys = append(page.keys, kv)
		page.byteSize = kvSize(page.keys[0])
	}

	page.offset = offset
	page.config = config
	page.buff = nil
	page.isHead = false

	return page
}

// TODO 可以考虑采用对象池来缓冲page对象
func pageDestroy(db *DB, page *Page) {

}

func pageClone(db *DB, page *Page) *Page {
	clone := pageCreate(db, page.typ, page.offset, page.config)
	clone.isHead = page.isHead
	clone.length = 0

	for i:=0; i<len(page.keys); i++ {
		clone.keys = append(clone.keys, kvClone(page.keys[i], true))
		clone.length ++
	}
	clone.byteSize = page.byteSize
	return clone
}



type PageSearchRes struct {
	chid *Page
	index uint64
	cmp int
}

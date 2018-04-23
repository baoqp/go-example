package bplus_tree

import (
	"unsafe"
	"util"
)

type PageType int

const (
	kPage PageType = 0
	kLeaf PageType = 1
)

type SearchType int

const (
	kNotLoad         SearchType = 0
	kLoad            SearchType = 1
	DefaultSeachType SearchType = kLoad
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

type PageSearchRes struct {
	child *Page
	index uint64
	cmp   int
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
			value:     nil,
			offset:    0,
			length:    0,
			config:    0,
			allocated: false,
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

// TODO 可以考虑采用对象池来缓冲page对象，避免频繁创建对象
func pageDestroy(db *DB, page *Page) {
	for i:=uint64(0); i<page.length; i++ {
		if page.keys[i].allocated {
			page.keys[i].value = nil
		}
	}

	if page.buff != nil {
		page.buff = nil
	}

}

func pageClone(db *DB, page *Page) *Page {
	clone := pageCreate(db, page.typ, page.offset, page.config)
	clone.isHead = page.isHead
	clone.length = 0

	for i := 0; i < len(page.keys); i++ {
		kv := new(KV)
		kvCopy(page.keys[i], kv, true)
		clone.keys = append(clone.keys, kv)
		clone.length ++
	}
	clone.byteSize = page.byteSize
	return clone
}

// 从文件中读取到内存中
func pageRead(db *DB, page *Page) error {
	//cast db to writer
	p := unsafe.Pointer(db)
	w := (*Writer)(p)

	// Read page size and leaf flag
	size := page.config >> 1
	if page.config&1 > 0 {
		page.typ = kLeaf
	} else {
		page.typ = kPage
	}

	buff, err := writerRead(w, DefaultComp, page.offset, &size)
	if err != nil {
		return err
	}

	// parse data
	i := 0
	o := uint64(0)
	for ; o < size; {
		page.keys[i].length = util.DecodeUint64(buff, 0)
		page.keys[i].offset = util.DecodeUint64(buff, 8)
		page.keys[i].config = util.DecodeUint64(buff, 16)
		page.keys[i].value = buff[24:]
		page.keys[i].allocated = false
		i++
	}

	page.length = uint64(i)
	page.byteSize = size
	page.buff = buff

	return nil
}

func pageLoad(db *DB, offset uint64, config uint64) (*Page, error) {
	newPage := pageCreate(db, kPage, offset, config)
	if err := pageRead(db, newPage); err != nil {
		pageDestroy(db, newPage)
		return nil, err
	}
	return newPage, nil
}

func pageSave(db *DB, page *Page) error {
	//cast db to writer
	p := unsafe.Pointer(db)
	w := (*Writer)(p)

	util.Assert(page.typ == kLeaf || page.length != 0,
		"wrong page type or page.length is 0")

	// Allocate space for serialization (header + keys)
	buff := make([]byte, page.byteSize)

	o := uint64(0)
	for i := uint64(0); i < page.length; i++ {
		util.Assert(o+kvSize(page.keys[i]) <= page.byteSize,
			"no enough buff for page.keys")
		util.EncodeUint64(buff, 0, page.keys[i].length)
		util.EncodeUint64(buff, 8, page.keys[i].offset)
		util.EncodeUint64(buff, 16, page.keys[i].config)
		copy(buff[24:24+page.keys[i].length], page.keys[i].value)
		o += kvSize(page.keys[i])
	}

	util.Assert(o == page.byteSize,
		"sum of all kv size not equals to page.byteSize")
	err := writerWrite(w, DefaultComp, buff, &page.offset, &page.config)
	if err != nil {
		return err
	}
	if page.typ == kLeaf {
		page.config = page.config<<1 | uint64(1)
	} else {
		page.config = page.config<<1 | uint64(0)
	}

	return nil
}

func pageLoadValue(db *DB, page *Page, index uint64, value *Value) error {
	return valueLoad(db, page.keys[index].offset, page.keys[index].config, value)
}

type UpdateCallback func(arg []byte, previous *Value, value *Value) error

func pageSaveValue(db *DB, page *Page, index uint64, cmp int, key *Key,
	value *Value, callback UpdateCallback, arg []byte) error {

	var previous *KV = nil
	var tmp = new(KV)
	// replace item with same key from page
	if cmp == 0 {
		previous = new(KV)
		if callback != nil {
			var prevValue Value
			err := pageLoadValue(db, page, index, &prevValue)
			if err != nil {
				return err
			}

			err = callback(arg, &prevValue, value)
			if err != nil {
				return EUPDATECONFLICT
			}
		}
		previous.offset = page.keys[index].offset
		previous.length = page.keys[index].length
		pageRemoveIdx(db, page, index);
	}

	tmp.value = key.value
	tmp.length = key.length

	err := valueSave(db, value, previous, &tmp.offset, &tmp.config)

	if err != nil {
		return err
	}

	pageShiftr(db, page, index)

	err = kvCopy(tmp, page.keys[index], true)
	if err != nil {
		// shift keys back
		pageShiftl(db, page, index)
		return err
	}
	page.byteSize += kvSize(tmp)
	page.length ++

	return nil
}

func pageSearch(db *DB, page *Page, key *Key, searchType SearchType,
	result *PageSearchRes) error {

	util.Assert(page.typ == kLeaf || page.length != 0,
		"wrong page type or page.length is 0")

	var i uint64 = 0
	if page.typ == kLeaf {
		i = 1
	}

	cmp := -1

	for i < page.length {
		p := unsafe.Pointer(page.keys[i])
		k := (*Key)(p)
		cmp = db.comparaCb(k, key)

		if cmp >= 0 {
			break
		}

		i ++
	}

	result.cmp = cmp

	if page.typ == kLeaf {
		result.index = i
		result.child = nil
		return nil
	} else {
		util.Assert(i > 0, "find idx is not  0")

		if searchType == kLoad {
			child, err := pageLoad(db, page.keys[i].offset, page.keys[i].config)
			if err != nil {
				return err
			}
			result.child = child
		} else {
			result.child = nil
		}

		return nil
	}
}

func pageGet(db *DB, page *Page, key *Key, value *Value) error {
	var res PageSearchRes
	err := pageSearch(db, page, key, DefaultSeachType, &res)
	if err != nil {
		return err
	}

	if res.child == nil {
		if res.cmp != 0 {
			return ENOTFOUND
		}

		return pageLoadValue(db, page, res.index, value)
	} else {
		err := pageGet(db, res.child, key, value)
		pageDestroy(db, res.child)
		return err
	}

}

type FilterCallback func(arg []byte, key *Key) bool
type RangeCallback func(arg []byte, key *Key, value *Value)

func pageGetRange(db *DB, page *Page, start *Key, end *Key,
	filter FilterCallback, rangeCb RangeCallback, arg []byte) error {

	var startRes, endRes PageSearchRes
	err := pageSearch(db, page, start, kNotLoad, &startRes)
	if err != nil {
		return err
	}

	err = pageSearch(db, page, end, kNotLoad, &endRes)
	if err != nil {
		return err
	}

	if page.typ == kLeaf {
		// on leaf pages end-key should always be greater or equal than first key
		if endRes.cmp > 0 && endRes.index == 0 {
			return nil
		}
		if endRes.cmp < 0 {
			endRes.index --
		}
	}

	//go through each page item
	for i := startRes.index; i <= endRes.index; i++ {
		p := unsafe.Pointer(page.keys[i])
		key := (*Key)(p)
		if !filter(arg, key) {
			continue
		}

		if page.typ == kLeaf {

			child, err := pageLoad(db, page.keys[i].offset, page.keys[i].config)
			if err != nil {
				return err
			}

			err = pageGetRange(db, child, start, end, filter, rangeCb, arg)
			pageDestroy(db, child)
			if err != nil {
				return err
			}
		} else {
			var value Value
			err := pageLoadValue(db, page, i, &value)
			if err != nil {
				return err
			}

			p := unsafe.Pointer(page.keys[i])
			key := (*Key)(p)
			rangeCb(arg, key, &value)
		}
	}

	return nil
}

func pageInsert(db *DB, page *Page, key *Key, value *Value,
	updataCb UpdateCallback, arg []byte) error {

	var err error
	var res PageSearchRes

	err = pageSearch(db, page, key, kLoad, &res)
	if err != nil {
		return err
	}

	if res.child == nil {
		// TODO store value in db file to get offset and config ???
		err = pageSaveValue(db, page, res.index, res.cmp, key, value, updataCb, arg)
		if err != nil {
			return err
		}
	} else {
		// Insert kv in child page
		err = pageInsert(db, res.child, key, value, updataCb, arg)
		if err == ESPLITPAGE {
			err = pageSplit(db, page, res.index, res.child)
		} else if err == nil {
			page.keys[res.index].offset = res.child.offset
			page.keys[res.index].config = res.child.config
		}

		pageDestroy(db, res.child)
		res.child = nil

		if err != nil {
			return err
		}
	}

	if page.length == db.header.pageSize {
		if page.isHead {
			newHead, err := pageSplitHead(db, page)
			if err != nil {
				return err
			} else {
				page = newHead
			}
		} else {
			return ESPLITPAGE
		}

	}

	util.Assert(page.length < db.header.pageSize,
		"page.length is not smaller than pageSize")

	return pageSave(db, page)
}

func pageBulkInsert(db *DB, page *Page, limit *Key, count *uint64, keys []*Key,
	values []*Value, updateCb UpdateCallback, arg []byte) error {

	var res PageSearchRes
	var err error
	for *count > 0 && (limit == nil || db.comparaCb(limit, keys[0]) > 0 ) {
		err = pageSearch(db, page, keys[0], kLoad, &res)
		if err != nil {
			return err
		}

		if res.child == nil {
			// store value in db file to get offset and config
			err = pageSaveValue(db, page, res.index, res.cmp, keys[0], values[0],
				updateCb, arg)
			// gnore update conflicts, to handle situations where
			// only one kv failed in a bulk
			if err != nil && err != EUPDATECONFLICT {
				return err
			}
			keys = keys[1:]
			values = values[1:]
		} else { // we're in regular page
			var newLimit *Key
			if res.index+1 < page.length {
				p := unsafe.Pointer(page.keys[res.index+1])
				newLimit = (*Key)(p)
			}
			err = pageBulkInsert(db, res.child, newLimit, count, keys, values,
				updateCb, arg)

			if err == ESPLITPAGE {
				err = pageSplit(db, page, res.index, res.child)
			} else if err == nil {
				page.keys[res.index].offset = res.child.offset
				page.keys[res.index].config = res.child.config
			}

			pageDestroy(db, res.child)
			res.child = nil

			if err != nil {
				return err
			}
		}

		if page.length == db.header.pageSize {
			if page.isHead {
				newHead, err := pageSplitHead(db, page)
				if err != nil {
					return err
				} else {
					page = newHead
				}
			} else {
				return ESPLITPAGE
			}
		}

		util.Assert(page.length < db.header.pageSize,
			"page.length is not smaller than pageSize")

	}

	return pageSave(db, page)
}

type RemoveCallback func(arg []byte, value *Value) error

func pageRemove(db *DB, page *Page, key *Key, removeCb RemoveCallback,
	arg []byte) error {

	var err error
	var res PageSearchRes
	err = pageSearch(db, page, key, kLoad, &res)
	if err != nil {
		return err
	}

	if res.child == nil {
		if res.cmp != 0 {
			return ENOTFOUND
		}

		if removeCb != nil {
			var prevVal Value
			err = pageLoadValue(db, page, res.index, &prevVal)
			if err != nil {
				return err
			}

			err = removeCb(arg, &prevVal)
			if err != nil {
				return EREMOVECONFLICT
			}
		}

		pageRemoveIdx(db, page, res.index)

		if page.length == 0 && !page.isHead {
			return EEMPTYPAGE
		}
	} else {
		// Insert kv in child page
		err = pageRemove(db, res.child, key, removeCb, arg)
		if err != nil && err != EEMPTYPAGE {
			return err
		}

		// kv was inserted but page is full now
		if err == EEMPTYPAGE {
			pageRemoveIdx(db, page, res.index)
			pageDestroy(db, res.child)
			res.child = nil

			// only one item left - lift kv from last child to current page
			if page.length == 1 {
				page.offset = page.keys[0].offset
				page.config = page.keys[0].config

				// remove child to free memory
				pageRemoveIdx(db, page, 0)

				//and load child as current page
				err = pageRead(db, page)
				if err != nil {
					return err
				}
			}
		} else {
			// Update offsets in page
			page.keys[res.index].offset = res.child.offset
			page.keys[res.index].config = res.child.config

			pageDestroy(db, res.child)
			res.child = nil
		}

	}

	return pageSave(db, page)
}


func pageCopy(source *DB, target *DB, page *Page) error {

	for i:=uint64(0); i<page.length; i++ {
		if page.typ == kPage {

			child, err := pageLoad(source, page.keys[i].offset, page.keys[i].config)
			if err != nil {
				return err
			}

			if err := pageCopy(source, target, child); err != nil {
				return err
			}

			page.keys[i].offset = child.offset
			page.keys[i].config = child.config
			pageDestroy(source, child)
		} else {
			var value Value
			err := pageLoadValue(source, page, i, &value)
			if err != nil {
				return err
			}

			page.keys[i].config = value.length
			err = valueSave(target, &value, nil, &page.keys[i].offset,
				&page.keys[i].config)
			if err != nil {
				return err
			}
		}
	}

	return nil
}


func pageSplit(db *DB, parent *Page, index uint64, child *Page) error {
	var err error
	left := pageCreate(db, child.typ, 0, 0)
	right := pageCreate(db, child.typ, 0, 0)
	middle := db.header.pageSize >> 1

	var middleKey KV
	err = kvCopy(child.keys[middle], &middleKey, true)
	if err != nil {
		return err
	}

	// non-leaf nodes has byte_size > 0 nullify it first
	var i = uint64(0)
	left.byteSize = 0
	left.length = 0
	left.keys = make([]*KV, middle)
	for ; i < middle; i++ {
		err = kvCopy(child.keys[i], left.keys[left.length], true)
		if err != nil {
			goto fatal
		}
		left.length ++
		left.byteSize += kvSize(child.keys[i])
	}

	right.byteSize = 0
	right.length = 0
	right.keys = make([]*KV, db.header.pageSize-middle)
	for ; i < db.header.pageSize; i++ {
		err = kvCopy(child.keys[i], right.keys[right.length], true)
		if err != nil {
			goto fatal
		}
		right.length ++
		right.byteSize += kvSize(child.keys[i])
	}

	// save left and right parts to get offsets
	err = pageSave(db, left)
	if err != nil {
		return err
	}
	err = pageSave(db, right)
	if err != nil {
		return err
	}

	// store offsets with middle key
	middleKey.offset = right.offset
	middleKey.config = right.config

	// insert middle key into parent page
	pageShiftr(db, parent, index+1)
	kvCopy(&middleKey, parent.keys[index+1], false)

	parent.byteSize += kvSize(&middleKey)
	parent.length ++

	// change left element
	parent.keys[index].offset = left.offset
	parent.keys[index].config = left.config

	return nil
fatal:
	pageDestroy(db, left)
	pageDestroy(db, right)
	return err
}


func pageSplitHead(db *DB, page *Page) (*Page, error) {
	newHead := pageCreate(db, 0, 0, 0)
	newHead.isHead = true
	err := pageSplit(db, newHead, 0, page)
	if err != nil {
		pageDestroy(db, newHead)
		return nil, err
	}

	db.header.page = newHead
	pageDestroy(db, page)

	return newHead, nil
}

func pageRemoveIdx(db *DB, page *Page, index uint64) error {
	util.Assert(index < page.length,
		"idx is not small than page.length")

	page.byteSize -= kvSize(page.keys[index])
	if page.keys[index].allocated {
		page.keys[index].value = nil
	}

	// Shift all keys left
	pageShiftl(db, page, index)
	page.length --

	return nil
}

func pageShiftr(db *DB, page *Page, index uint64) {
	if page.length > 0 {
		for i := page.length - 1; i >= index; i-- {
			kvCopy(page.keys[i], page.keys[i+1], false)
			if i == 0 {
				break
			}
		}
	}
}

func pageShiftl(db *DB, page *Page, index uint64) {
	for i := index + 1; i < page.length; i++ {
		kvCopy(page.keys[i], page.keys[i-1], false)
	}
}

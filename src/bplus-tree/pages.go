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
	keys     []KV // 指向子节点或数据的指针
}

type PageSearchRes struct {
	child *Page
	index uint64
	cmp   int
}

func (page *Page) destroy() {
	if len(page.keys) > 0 {
		for _, key := range page.keys {
			if key.allocated {
				key.value = key.value[:0]
			}
		}
		page.keys = page.keys[:0]
	}

	if len(page.buff) > 0 {
		page.buff = page.buff[:0]
	}
}

func (page *Page) clone(tree *Tree) *Page {
	clone := tree.pageCreate(page.typ, page.offset, page.config)
	clone.isHead = page.isHead
	clone.length = page.length
	clone.keys = make([]KV, clone.length)
	for i := 0; i < int(page.length); i++ {
		kvCopy(&page.keys[i], &clone.keys[i], true)
	}
	clone.byteSize = page.byteSize
	return clone
}

// 从文件反序列化到内存中
func (page *Page) read(tree *Tree) error {

	// Read page size and leaf flag
	size := page.config >> 1
	if page.config&1 > 0 {
		page.typ = kLeaf
	} else {
		page.typ = kPage
	}

	buff, err := tree.writerRead(DefaultComp, page.offset, &size)
	if err != nil {
		return err
	}

	// parse data
	i := 0
	o := uint64(0)
	for ; o < size; {

		if i >= len(page.keys) {
			page.keys = append(page.keys, KV{})
		}

		page.keys[i].length = util.DecodeUint64(buff, int(o+0))
		page.keys[i].offset = util.DecodeUint64(buff, int(o+8))
		page.keys[i].config = util.DecodeUint64(buff, int(o+16))
		page.keys[i].value = buff[int(o+24): int(o+24+page.keys[i].length)]
		page.keys[i].allocated = false

		o += 24 + page.keys[i].length
		i++
	}

	page.length = uint64(i)
	page.byteSize = size
	page.buff = buff

	return nil
}

func (page *Page) save(tree *Tree) error {
	util.Assert(page.typ == kLeaf || page.length != 0,
		"wrong page type or page.length is 0")

	// Allocate space for serialization (header + keys)
	buff := make([]byte, page.byteSize)
	keys := page.keys
	o := uint64(0)
	for i := uint64(0); i < page.length; i++ {
		util.Assert(o+kvSize(&keys[i]) <= page.byteSize,
			"no enough buff for page.keys")
		util.EncodeUint64(buff, int(o+0), keys[i].length)
		util.EncodeUint64(buff, int(o+8), keys[i].offset)
		util.EncodeUint64(buff, int(o+16), keys[i].config)
		copy(buff[int(o+24):int(o+24+page.keys[i].length)], keys[i].value)
		o += 24 + keys[i].length
	}

	util.Assert(o == page.byteSize,
		"sum of all kv size not equals to page.byteSize")
	page.config = page.byteSize
	err := tree.writerWrite(DefaultComp, buff, &page.offset, &page.config)

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

type UpdateCallback func(arg []byte, previous *Value, value *Value) error

func (page *Page) saveValue(tree *Tree, index uint64, cmp int, key *Key,
	value *Value, callback UpdateCallback, arg []byte) error {

	var previous *KV = nil
	var tmp = new(KV)
	// replace item with same key from page
	if cmp == 0 { // cmp >= 0
		previous = new(KV)
		if callback != nil {
			var prevValue Value
			err := pageLoadValue(tree, page, index, &prevValue)
			if err != nil {
				return err
			}

			err = callback(arg, &prevValue, value)
			if err != nil {
				return EUPDATECONFLICT
			}
		}
		previous.offset = page.keys[index].offset // 记录相同key老数据在文件中的位置和长度
		previous.length = page.keys[index].length
		page.removeIdx(index)
	}

	tmp.value = key.value
	tmp.length = key.length

	// 插入或替换的数据都是写到文件的末尾，执行compact操作来删除被替换的老数据以减小数据文件的大小
	err := valueSave(tree, value, previous, &tmp.offset, &tmp.config)

	if err != nil {
		return err
	}

	page.shiftr(index)

	kvCopy(tmp, &page.keys[index], true) // keys 为 null

	page.byteSize += kvSize(tmp)
	page.length ++

	return nil
}

func (page *Page) search(tree *Tree, key *Key, searchType SearchType, result *PageSearchRes) error {

	util.Assert(page.typ == kLeaf || page.length > 0,
		"wrong page type or page.length is 0")

	var i uint64 = 0
	if page.typ == kPage {
		i = 1
	}

	cmp := -1

	for i < page.length {
		p := unsafe.Pointer(&page.keys[i])
		k := (*Key)(p)
		cmp = tree.comparaCb(k, key)

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

		if cmp != 0 {
			i --
		}

		if searchType == kLoad {
			// 加载子页面并在其中搜索
			child, err := pageLoad(tree, page.keys[i].offset, page.keys[i].config)
			if err != nil {
				return err
			}
			result.child = child
		} else {
			result.child = nil
		}
		result.index = i
		return nil
	}
}

func (page *Page) get(tree *Tree, key *Key, value *Value) error {
	var res PageSearchRes
	err := page.search(tree, key, DefaultSeachType, &res)
	if err != nil {
		return err
	}

	if res.child == nil {
		if res.cmp != 0 {
			return ENOTFOUND
		}
		return pageLoadValue(tree, page, res.index, value)
	} else {
		err := res.child.get(tree, key, value)
		res.child.destroy()
		return err
	}

}

type FilterCallback func(arg []byte, key *Key) bool
type RangeCallback func(arg []byte, key *Key, value *Value)

func (page *Page) getRange(tree *Tree, start *Key, end *Key, filter FilterCallback, rangeCb RangeCallback,
	arg []byte) error {

	var startRes, endRes PageSearchRes
	err := page.search(tree, start, kNotLoad, &startRes)
	if err != nil {
		return err
	}

	err = page.search(tree, end, kNotLoad, &endRes)
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
		p := unsafe.Pointer(&page.keys[i])
		key := (*Key)(p)
		if !filter(arg, key) {
			continue
		}

		if page.typ == kLeaf {

			child, err := pageLoad(tree, page.keys[i].offset, page.keys[i].config)
			if err != nil {
				return err
			}

			err = child.getRange(tree, start, end, filter, rangeCb, arg)
			child.destroy()
			if err != nil {
				return err
			}
		} else {
			var value Value
			err := pageLoadValue(tree, page, i, &value)
			if err != nil {
				return err
			}

			p := unsafe.Pointer(&page.keys[i])
			key := (*Key)(p)
			rangeCb(arg, key, &value)
		}
	}

	return nil
}

func (page *Page) insert(tree *Tree, key *Key, value *Value, updataCb UpdateCallback,
	arg []byte) error {

	var err error
	var res PageSearchRes

	err = page.search(tree, key, kLoad, &res)
	if err != nil {
		return err
	}

	if res.child == nil { // 叶子节点
		// cmp != 0 说明有元素插入，需要扩大slice
		if res.cmp != 0 {
			page.keys = append(page.keys, KV{})
		}

		// store value in Tree file to get offset and config
		err = page.saveValue(tree, res.index, res.cmp, key, value, updataCb, arg)
		if err != nil {
			return err
		}
	} else {
		// Insert kv in child page
		err = res.child.insert(tree, key, value, updataCb, arg)
		if err == ESPLITPAGE {
			err = page.split(tree, res.index, res.child)
		} else if err == nil {
			page.keys[res.index].offset = res.child.offset
			page.keys[res.index].config = res.child.config
		}

		res.child.destroy()
		res.child = nil

		if err != nil {
			return err
		}
	}

	err = page.save(tree)
	if err != nil {
		return err
	}

	// 子节点分裂时会在父节点插入一个key, 此时父节点也可能变满，需要再次判断
	if page.length == tree.header.pageSize {
		if page.isHead {
			_, err := page.splitHead(tree)
			if err != nil {
				return err
			}
		} else {
			return ESPLITPAGE
		}

	}

	util.Assert(page.length < tree.header.pageSize,
		"page.length is not smaller than pageSize")
	// 每次插入数据后，由于page中keys改变了，所以page也要序列化到文件中
	return nil

}

func (page *Page) bulkInsert(tree *Tree, limit *Key, count *uint64, keys []*Key,
	values []*Value, updateCb UpdateCallback, arg []byte) error {

	var res PageSearchRes
	var err error
	for *count > 0 && (limit == nil || tree.comparaCb(limit, keys[0]) > 0 ) {
		err = page.search(tree, keys[0], kLoad, &res)
		if err != nil {
			return err
		}

		if res.child == nil {
			// store value in Tree file to get offset and config
			err = page.saveValue(tree, res.index, res.cmp, keys[0], values[0],
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
				p := unsafe.Pointer(&page.keys[res.index+1])
				newLimit = (*Key)(p)
			}
			err = res.child.bulkInsert(tree, newLimit, count, keys, values, updateCb, arg)

			if err == ESPLITPAGE {
				err = page.split(tree, res.index, res.child)
			} else if err == nil {
				page.keys[res.index].offset = res.child.offset
				page.keys[res.index].config = res.child.config
			}

			res.child.destroy()
			res.child = nil

			if err != nil {
				return err
			}
		}

		if page.length == tree.header.pageSize {
			if page.isHead {
				newHead, err := page.splitHead(tree)
				if err != nil {
					return err
				} else {
					page = newHead
				}
			} else {
				return ESPLITPAGE
			}
		}

		util.Assert(page.length < tree.header.pageSize,
			"page.length is not smaller than pageSize")

	}

	return page.save(tree)
}

type RemoveCallback func(arg []byte, value *Value) error

func (page *Page) remove(tree *Tree, key *Key, removeCb RemoveCallback, arg []byte) error {

	var err error
	var res PageSearchRes
	err = page.search(tree, key, kLoad, &res)
	if err != nil {
		return err
	}

	if res.child == nil {
		if res.cmp != 0 {
			return ENOTFOUND
		}

		if removeCb != nil {
			var prevVal Value
			err = pageLoadValue(tree, page, res.index, &prevVal)
			if err != nil {
				return err
			}

			err = removeCb(arg, &prevVal)
			if err != nil {
				return EREMOVECONFLICT
			}
		}

		page.removeIdx(res.index)

		if page.length == 0 && !page.isHead {
			return EEMPTYPAGE
		}
	} else {
		// Insert kv in child page
		err = res.child.remove(tree, key, removeCb, arg)
		if err != nil && err != EEMPTYPAGE {
			return err
		}

		// kv was inserted but page is full now
		if err == EEMPTYPAGE {
			page.removeIdx(res.index)
			res.child.destroy()
			res.child = nil

			// only one item left - lift kv from last child to current page
			if page.length == 1 {
				page.offset = page.keys[0].offset
				page.config = page.keys[0].config
				// remove child to free memory
				page.removeIdx(0)

				//and load child as current page
				err = page.read(tree)
				if err != nil {
					return err
				}
			}
		} else {
			// Update offsets in page
			page.keys[res.index].offset = res.child.offset
			page.keys[res.index].config = res.child.config

			res.child.destroy()
			res.child = nil
		}
	}

	return page.save(tree)
}

func (page *Page) split(tree *Tree, index uint64, child *Page) error {
	var err error
	left := tree.pageCreate(child.typ, 0, 0)
	right := tree.pageCreate(child.typ, 0, 0)
	middle := tree.header.pageSize >> 1

	var middleKey KV
	kvCopy(&child.keys[middle], &middleKey, true)

	// non-leaf nodes has byte_size > 0 nullify it first
	var i = uint64(0)
	left.byteSize = 0
	left.length = 0
	left.keys = make([]KV, middle)
	for ; i < middle; i++ {
		kvCopy(&child.keys[i], &left.keys[left.length], true)
		left.length ++
		left.byteSize += kvSize(&child.keys[i])
	}

	right.byteSize = 0
	right.length = 0
	right.keys = make([]KV, tree.header.pageSize-middle)
	for i = middle; i < tree.header.pageSize; i++ {
		kvCopy(&child.keys[i], &right.keys[right.length], true)
		right.length ++
		right.byteSize += kvSize(&child.keys[i])
	}

	// save left and right parts to get offsets
	err = left.save(tree)
	if err != nil {
		return err
	}
	err = right.save(tree)
	if err != nil {
		return err
	}

	// store offsets with middle key
	middleKey.offset = right.offset
	middleKey.config = right.config

	// insert middle key into parent page
	page.keys = append(page.keys, KV{})
	page.shiftr(index + 1)
	kvCopy(&middleKey, &page.keys[index+1], false)

	page.byteSize += kvSize(&middleKey)
	page.length ++

	// change left element
	page.keys[index].offset = left.offset
	page.keys[index].config = left.config

	return nil
}

func (page *Page) splitHead(tree *Tree) (*Page, error) {
	newHead := tree.pageCreate(0, 0, 0)
	newHead.isHead = true
	err := newHead.split(tree, 0, page)
	if err != nil {
		newHead.destroy()
		return nil, err
	}

	page.destroy()
	*tree.header.page = *newHead
	*page = *newHead

	return newHead, nil
}

// 删除index位置的key
func (page *Page) removeIdx(index uint64) error {
	util.Assert(index < page.length,
		"idx is not small than page.length")

	page.byteSize -= kvSize(&page.keys[index])
	if page.keys[index].allocated {
		page.keys[index].value = nil
	}

	// Shift all keys left
	page.shiftl(index)
	page.length --

	return nil
}

func (page *Page) shiftr(index uint64) {
	if page.length > 0 {
		for i := page.length - 1; i >= index; i-- {
			source := &page.keys[i]
			target := &page.keys[i+1]
			kvCopy(source, target, false)
			if i == 0 {
				break
			}
		}
	}
}

func (page *Page) shiftl(index uint64) {
	for i := index + 1; i < page.length; i++ {
		kvCopy(&page.keys[i], &page.keys[i-1], false)
	}
}

func pageLoad(tree *Tree, offset uint64, config uint64) (*Page, error) {
	newPage := tree.pageCreate(kPage, offset, config)
	if err := newPage.read(tree); err != nil {
		newPage.destroy()
		return nil, err
	}
	return newPage, nil
}

func pageLoadValue(Tree *Tree, page *Page, index uint64, value *Value) error {
	return valueLoad(Tree, page.keys[index].offset, page.keys[index].config, value)
}

func pageCopy(source *Tree, target *Tree, page *Page) error {

	for i := uint64(0); i < page.length; i++ {
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
			child.destroy()
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

	return page.save(target)
}

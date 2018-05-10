package bplus_tree_aof

import (
	"sync"
	"github.com/pkg/errors"
	"util"
	"os"
	"fmt"
)

const PADDING = 64

const PageSize = 64

type ComparaCallback func(a *Key, b *Key) int

type TreeHeader struct {
	offset   uint64
	config   uint64
	pageSize uint64
	hash     uint64
	page     *Page
}


type Tree struct {
	*Writer
	rmLock    sync.RWMutex
	header    *TreeHeader
	comparaCb ComparaCallback
}


type CompType int

const (
	kNotCompressed CompType = 0
	kCompressed    CompType = 1
	DefaultComp    CompType = kNotCompressed
)

func open(filename string, isCompact bool) (*Tree, error) {
	var err error

	tree := new(Tree)
	tree.header = new(TreeHeader)

	tree.rmLock.Lock()
	tree.Writer, err = tree.creatreWriter(filename, isCompact)
	if err != nil {
		return nil, err
	}

	tree.header.page = nil
	err = tree.init()
	tree.rmLock.Unlock()
	return tree, err
}

func (t *Tree) init() error {
	//  Load head.
	err := t.writerFind(kNotCompressed, HeaderSize,
		nil, treeReadHead, treeWriteHead)

	if err == nil {
		t.comparaCb = defaultCompareCb
	}
	return err
}

func (t *Tree) get(key *Key, value *Value) error {
	t.rmLock.RLock()
	err := t.header.page.get(t, key, value)
	t.rmLock.RUnlock()
	return err
}

func (t *Tree) getPrevious(value *Value, previous *Value) error {
	t.rmLock.RLock()
	if value.prevOffset == 0 && value.prevLength == 0 {
		return ENOTFOUND
	}
	err := valueLoad(t, value.prevOffset, value.prevOffset, previous)
	t.rmLock.RUnlock()
	return err
}

func (t *Tree) update(key *Key, value *Value, updateCb UpdateCallback, arg []byte) error {
	var err error
	t.rmLock.Lock()
	err = t.header.page.insert(t, key, value, updateCb, arg)
	if err == nil {
		err = treeWriteHead(t, nil) // TODO 待优化，不是每次都要重写Header
	}
	t.rmLock.Unlock()
	return err
}

func (t *Tree) bulkUpdate(count uint64, keys []*Key, values []*Value, updateCb UpdateCallback, arg []byte) error {
	var err error
	t.rmLock.Lock()
	left := count
	err = t.header.page.bulkInsert(t, nil, &left, keys, values, updateCb, arg)
	if err == nil {
		err = treeWriteHead(t, nil)
	}
	t.rmLock.Unlock()
	return err
}

func (t *Tree) remove(key *Key, removeCb RemoveCallback, arg []byte) error {
	var err error
	t.rmLock.Lock()
	err = t.header.page.remove(t, key, removeCb, arg)
	if err == nil {
		err = treeWriteHead(t, nil)
	}
	t.rmLock.Unlock()
	return err

}

// 把原来的树的各个节点一次读出来写入新的文件中
func (t *Tree) compact() error {
	var err error

	compactName := t.compactName
	compactExists, err := util.Exists(compactName)
	if err != nil{
		return err
	}
	if compactExists {
		err = os.Remove(compactName)
		if err != nil {
			return err
		}
	}

	compacted, err := open(t.originalName, true)
	if err != nil {
		return err
	}
	util.Assert(compacted.fileName == compactName,
		fmt.Sprintf("compact file name mismatch, %s = %s", compacted.fileName, compactName))

	if compacted.header.page != nil {
		compacted.header.page.destroy()
	}

	t.rmLock.RLock()
	compacted.header.page = t.header.page.clone(compacted)

	err = pageCopy(t, compacted, compacted.header.page)
	if err != nil {
		return err
	}

	err = treeWriteHead(compacted, nil)
	if err != nil {
		return err
	}

	err = t.destroyWriter()
	if err != nil {
		return err
	}

	err = t.deleteTreeFile()
	if err != nil {
		return err
	}

	t.header.page.destroy()
	t.Writer = compacted.Writer
	t.header = compacted.header

	t.rmLock.RUnlock()

	return nil
}

func (t *Tree) getFilteredRange(start *Key, end *Key, callback FilterCallback,
	rangeCallback RangeCallback, arg []byte) error {

	var err error
	t.rmLock.Lock()
	err = t.header.page.getRange(t, start, end, callback, rangeCallback, arg)
	t.rmLock.Unlock()
	return err
}

func defaultFilterCallback(arg []byte, key *Key) bool {
	return true
}

func (t *Tree) getRange(start *Key, end *Key,
	rangeCallback RangeCallback, arg []byte) error {
	var err error
	t.rmLock.Lock()
	err = t.header.page.getRange(t, start, end, defaultFilterCallback,
		rangeCallback, arg)
	t.rmLock.Unlock()
	return err
}

func (t *Tree) pageCreate(typ PageType, offset uint64, config uint64) *Page {
	page := new(Page)
	page.typ = typ
	if typ == kLeaf {
		page.length = 0
		page.byteSize = 0
	} else {
		// non-leaf pages always have one left-most element
		page.length = 1
		kv := KV{
			value:     nil,
			offset:    0,
			length:    0,
			config:    0,
			allocated: false,
		}
		page.keys = append(page.keys, kv)
		page.byteSize = kvSize(&page.keys[0])
	}

	page.offset = offset
	page.config = config
	page.buff = nil
	page.isHead = false
 	return page
}

func treeReadHead(t *Tree, data []byte) error {
	var err error
	t.header.offset = util.DecodeUint64(data, 0)
	t.header.config = util.DecodeUint64(data, 8)
	t.header.pageSize = util.DecodeUint64(data, 16)
	t.header.hash = util.DecodeUint64(data, 24)
	data = data[:0]

	if computeHashl(t.header.offset) != t.header.hash {
		return errors.New("hash inconsistent ")
	}

	// 载入b+树的根
	t.header.page, err = pageLoad(t, t.header.offset, t.header.config)
	if err != nil {
		return err
	}
	t.header.page.isHead = true
	return nil
}

func treeWriteHead(t *Tree, data []byte) error {

	if t.header.page == nil {
		t.header.pageSize = PageSize
		t.header.page = t.pageCreate(kLeaf, 0, 1)
		t.header.page.isHead = true
	}

	t.header.offset = t.header.page.offset
	t.header.config = t.header.page.config
	t.header.hash = computeHashl(t.header.offset)

	buff := make([]byte, HeaderSize)
	util.EncodeUint64(buff, 0, t.header.offset)
	util.EncodeUint64(buff, 8, t.header.config)
	util.EncodeUint64(buff, 16, t.header.pageSize)
	util.EncodeUint64(buff, 24, t.header.hash)

	size := uint64(HeaderSize)
	var offset uint64
	return t.writerWrite(kNotCompressed, buff, &offset, &size)
}

func defaultCompareCb(a *Key, b *Key) int {
	var len uint64

	if a.length < b.length {
		len = a.length
	} else {
		len = b.length
	}

	for i := uint64(0); i < len; i++ {
		if a.value[i] > b.value[i] {
			return 1
		}

		if a.value[i] < b.value[i] {
			return -1
		}

	}

	if a.length > b.length {
		return 1
	}

	if a.length < b.length {
		return -1
	}

	return 0

}

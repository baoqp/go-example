package bplus_tree

import (
	"os"
	"sync"
	"unsafe"
	"github.com/pkg/errors"
)

const PADDING = 64

type ComparaCallback func(a *Key, b *Key) int

type TreeHeader struct {
	offset   uint64
	config   uint64
	pageSize uint64
	hash     uint64

	page *Page
}

// bp_db_s
type DB struct {
	// 实际是tree
	file      *os.File
	fileName  string
	fileSize  uint64
	padding   [PADDING]byte // TODO
	rmLock    sync.RWMutex
	header    TreeHeader
	comparaCb ComparaCallback
}

// https://github.com/embedded2016/bplus-tree

type CompType int

const (
	kNotCompressed CompType = 0
	kCompressed    CompType = 1
	DefaultComp    CompType = kNotCompressed
)

func open(tree *DB, filename string) error {
	var err error
	tree.rmLock.Lock()
	err = writerCreate((*Writer)(unsafe.Pointer(tree)), filename)
	if err != nil {
		return err
	}
	tree.header.page = nil
	err = initTree(tree)
	tree.rmLock.Unlock()
	return err
}

func treeReadHead(w *Writer, data []byte) error {
	var err error
	t := (*DB)(unsafe.Pointer(w))
	head := *(*TreeHeader)(unsafe.Pointer(&data))

	t.header.offset = head.offset // ntohll ???
	t.header.config = head.config
	t.header.pageSize = head.pageSize
	t.header.hash = head.hash

	data = data[:0]

	if computeHashl(t.header.offset) != t.header.hash {
		return errors.New("hash inconsistent ")
	}

	t.header.page, err = pageLoad(t, t.header.offset, t.header.config)
	if err != nil {
		return err
	}
	t.header.page.isHead = true
	return nil
}

func treeWriteHead(w *Writer, data []byte) error {
	t := (*DB)(unsafe.Pointer(w))

	if t.header.page == nil {
		t.header.pageSize = 64
		t.header.page = pageCreate(t, kLeaf, 0, 1)
		t.header.page.isHead = true
	}

	t.header.offset = t.header.page.offset
	t.header.config = t.header.page.config
	t.header.hash = computeHashl(t.header.offset)

	nhead := new(TreeHeader)
	nhead.offset = t.header.offset
	nhead.config = t.header.config
	nhead.pageSize = t.header.pageSize
	nhead.hash = t.header.hash

	size := uint64(HeaderSize)
	var offset uint64
	return writerWrite(w, kNotCompressed, *(*[]byte)(unsafe.Pointer(nhead)),
		&offset, &size)
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

func initTree(tree *DB) error {

	err := writerFind((*Writer)(unsafe.Pointer(tree)), kNotCompressed, HeaderSize,
		*(*[]byte)(unsafe.Pointer(&tree.header)), treeReadHead, treeWriteHead)

	if err == nil {
		tree.comparaCb = defaultCompareCb
	}

	return err
}

func close(tree *DB) error {
	tree.rmLock.Lock()
	tree.file.Close()
	if tree.header.page != nil {
		pageDestroy(tree, tree.header.page)
		tree.header.page = nil
	}

	return nil
}

func get(tree *DB, key *Key, value *Value) error {
	tree.rmLock.RLock()
	err := pageGet(tree, tree.header.page, key, value)
	tree.rmLock.RUnlock()
	return err
}

func getPrevious(tree *DB, value *Value, previous *Value) error {
	tree.rmLock.RLock()
	if value.prevOffset == 0 && value.prevLength == 0 {
		return ENOTFOUND
	}
	err := valueLoad(tree, value.prevOffset, value.prevOffset, previous)
	tree.rmLock.RUnlock()
	return err
}

func update(tree *DB, key *Key, value *Value, updateCb UpdateCallback, arg []byte) error {
	var err error
	tree.rmLock.Lock()
	err = pageInsert(tree, tree.header.page, key, value, updateCb, arg)
	if err == nil { // TODO
		err = treeWriteHead((*Writer)(unsafe.Pointer(tree)), nil)
	}
	tree.rmLock.Unlock()
	return err
}

func bulkUpdate(tree *DB, count uint64, key []*Key, value []*Value, updateCb UpdateCallback, arg []byte) error {
	var err error
	tree.rmLock.Lock()
	left := count
	err = pageBulkInsert(tree, tree.header.page, nil, &left, key, value, updateCb, arg)
	if err == nil {
		err = treeWriteHead((*Writer)(unsafe.Pointer(tree)), nil)
	}
	tree.rmLock.Unlock()
	return err
}

func set(tree *DB, key *Key, value *Value) error {
	return update(tree, key, value, nil, nil)
}

func bulkSet(tree *DB, count uint64, key []*Key, value []*Value) error {
	return bulkUpdate(tree, count, key, value, nil, nil)
}

func removev(tree *DB, key *Key, removeCb RemoveCallback, arg []byte) error {
	var err error
	tree.rmLock.Lock()

	err = pageRemove(tree, tree.header.page, key, removeCb, arg)
	if err == nil {
		err = treeWriteHead((*Writer)(unsafe.Pointer(tree)), nil)
	}
	tree.rmLock.Unlock()
	return err

}

func remove(tree *DB, key *Key) error {
	return removev(tree, key, nil, nil)
}

// TODO
func compact(tree *DB) error {
	return nil
}

func getFilteredRange(tree *DB, start *Key, end *Key, callback FilterCallback,
	rangeCallback RangeCallback, arg []byte) error {

	var err error
	tree.rmLock.Lock()
	err = pageGetRange(tree, tree.header.page, start, end, callback, rangeCallback, arg)
	tree.rmLock.Unlock()
	return err
}

func defaultFilterCallback(arg []byte, key *Key) bool {
	return true
}

func getRange(tree *DB, start *Key, end *Key,
	rangeCallback RangeCallback, arg []byte) error {
	var err error
	tree.rmLock.Lock()
	err = pageGetRange(tree, tree.header.page, start, end, defaultFilterCallback,
		rangeCallback, arg)
	tree.rmLock.Unlock()
	return err
}

package bplus_tree

import (
	"os"
	"sync"
)

const PADDING = 64

type comparaCb func(a *Key, b *Key) int

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
	comparaCb *comparaCb
}

//  bp_key_s  bp_key_t
type Key struct {
	length uint64
	value  []byte

	prevOffset uint64
	prevLength uint64
}

// https://github.com/embedded2016/bplus-tree
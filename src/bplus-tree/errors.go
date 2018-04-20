package bplus_tree

import "github.com/pkg/errors"

var (
	EFILE           = errors.New("0x101")
	EFILEREAD_OOB   = errors.New("0x102")
	EFILEREAD       = errors.New("0x103")
	EFILEWRITE      = errors.New("0x104")
	EFILEFLUSH      = errors.New("0x105")
	EFILERENAME     = errors.New("0x106")
	ECOMPACT_EXISTS = errors.New("0x107")

	ECOMP   = errors.New("0x201")
	EDECOMP = errors.New("0x202")

	EALLOC  = errors.New("0x301")
	EMUTEX  = errors.New("0x302")
	ERWLOCK = errors.New("0x303")

	ENOTFOUND       = errors.New("0x401")
	ESPLITPAGE      = errors.New("0x402")
	EEMPTYPAGE      = errors.New("0x403")
	EUPDATECONFLICT = errors.New("0x404")
	EREMOVECONFLICT = errors.New("0x405")
)

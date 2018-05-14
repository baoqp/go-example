package gushroom

import "github.com/pkg/errors"

const (
	LengthByte = 2
	PageSize   = 4096
)

var DataMaxLen uint16 = PageSize - (4 + 2 + 2 + 1)
var DataNotEnoughErr = errors.New("dataPage's data is not enough")


type DataPage struct {
	pageNo pageId
	total  uint16
	curr   uint16
	dirty  uint8
	data   []byte // []byte在内存中不是和struct分配在一起的
}



func NewPage(pageNo pageId) *DataPage {
	dp := &DataPage{}
	dp.data = make([]byte, 0, DataMaxLen)
	dp.Reset(pageNo)
	return dp
}


func (dp *DataPage) Reset(pageNo pageId) {
	dp.pageNo = pageNo
	dp.total = 0
	dp.curr = 0
	dp.dirty = 0
}


func (dp *DataPage) PutData(slice DataSlice) error {
	len := uint16(len(slice) + LengthByte)
	if dp.curr + len > DataMaxLen {

	}

	var res  uint32 = uint32(dp.pageNo)
	res <<= 12
	res |= uint32(dp.curr) & 0xFFF
	// TODO
}
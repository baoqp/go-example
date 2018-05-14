package gushroom

import (
	"github.com/pkg/errors"
	"util"
	"os"
)

const (
	LengthByte = 2
	PageSize   = 4096
)

var DataMaxLen uint16 = PageSize - (4 + 2 + 2 + 1)
var DataNotEnoughErr = errors.New("dataPage's data is not enough")

type DataPage struct {
	pageNo pageId
	total  uint16 // 写入次数
	curr   uint16 // data当前索引
	dirty  bool
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
	dp.dirty = false
}

func (dp *DataPage) PutData(slice DataSlice) (uint32, error) {

	len := uint16(slice.len + LengthByte)
	if dp.curr+len > DataMaxLen {
		return 0, DataNotEnoughErr
	}

	// 把pageNo和offset编码在res中
	res := uint32(dp.pageNo)
	res <<= 12
	res |= uint32(dp.curr) & 0xFFF

	util.EncodeUint16(dp.data, int(dp.curr), slice.len)
	copy(dp.data[dp.curr+LengthByte:], slice.data)
	dp.curr += len

	dp.total ++
	dp.dirty = true

	return res, nil
}

func (dp *DataPage) GetData(pageNo pageId) *DataSlice {
	util.Assert(dp.pageNo == (pageNo >> 12), "GetData wrong pageNo")
	pos := uint32(pageNo) & 0xFFF

	var len uint16
	util.DecodeUint16(dp.data, int(pos), &len)

	return &DataSlice {
		len:len,
		data:dp.data[pos + LengthByte : uint16(pos) + LengthByte + len], // TODO 不需要拷贝???
	}
}


func (dp *DataPage) Write(fd *os.File) {
	fd.WriteAt()
}

























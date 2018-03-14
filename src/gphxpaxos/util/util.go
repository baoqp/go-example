package util

import (
	"time"
	"math/rand"
	"os"
	"encoding/binary"
	"hash/crc32"
)

func Rand(up int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(up)
}

//--------------------------------文件操作--------------------------------//
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// ---------------------------byte[]和类型的转换-----------------------------//
func DecodeUint64(buffer []byte, offset int, ret *uint64) {
	*ret = binary.LittleEndian.Uint64(buffer[offset:])
}

func EncodeUint64(buffer []byte, offset int, ret uint64) {
	binary.LittleEndian.PutUint64(buffer[offset:], ret)
}

func DecodeInt32(buffer []byte, offset int, ret *int32) {
	tmp := binary.LittleEndian.Uint32(buffer[offset:])
	*ret = int32(tmp)
}

func EncodeInt32(buffer []byte, offset int, ret int32) {
	binary.LittleEndian.PutUint32(buffer[offset:], uint32(ret))
}

func DecodeUint32(buffer []byte, offset int, ret *uint32) {
	*ret = binary.LittleEndian.Uint32(buffer[offset:])
}

func EncodeUint32(buffer []byte, offset int, ret uint32) {
	binary.LittleEndian.PutUint32(buffer[offset:], ret)
}

func DecodeUint16(buffer []byte, offset int, ret *uint16) {
	*ret = binary.LittleEndian.Uint16(buffer[offset:])
}

func EncodeUint16(buffer []byte, offset int, ret uint16) {
	binary.LittleEndian.PutUint16(buffer[offset:], ret)
}

//----------------------------------加密和编码------------------------------------//
func Crc32(crc uint32, value []byte, skiplen int) uint32 { // crc32编码
	vlen := len(value)
	data := value[:vlen-skiplen]
	return crc32.Update(crc, crc32.IEEETable, []byte(data))
}



// ---------------------------------时间相关操作---------------------------------//
func NowTimeMs() uint64 {
	return uint64(time.Now().UnixNano() / 1000000)
}

func SleepMs(ms int32) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
package util

import (
	"time"
	"os"
	"bytes"
	"encoding/binary"
	"strings"
	"strconv"
	"math/rand"
	"hash/crc32"
	"runtime"
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

func DeleteDir(path string) error {
	return os.RemoveAll(path)
}

//---------------------------------[]byte操作-------------------------------------//

func AppendBytes(inputs ...[]byte) [] byte {
	return bytes.Join(inputs, []byte(""))
}

func CopyBytes(src []byte) [] byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}

// ---------------------------[]byte和类型的转换-----------------------------//
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

func Inet_addr(ipaddr string) uint32 {
	var (
		ip                 = strings.Split(ipaddr, ".")
		ip1, ip2, ip3, ip4 uint64
		ret                uint32
	)
	ip1, _ = strconv.ParseUint(ip[0], 10, 8)
	ip2, _ = strconv.ParseUint(ip[1], 10, 8)
	ip3, _ = strconv.ParseUint(ip[2], 10, 8)
	ip4, _ = strconv.ParseUint(ip[3], 10, 8)
	ret = uint32(ip4)<<24 + uint32(ip3)<<16 + uint32(ip2)<<8 + uint32(ip1)
	return ret
}



//打印调用栈
func Pstack() string {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	return string(buf[0:n])
}
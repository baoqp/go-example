package util

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"gphxpaxos/comm"
	"encoding/binary"
	"fmt"
)

var GROUPIDXLEN = binary.Size(int32(0))
var HEADLEN_LEN = binary.Size(uint16(0))
var CHECKSUM_LEN = binary.Size(uint32(0))

// format: groupId(int) + header_len(uint16) + header + body + crc32 checksum(uint32)
func packBaseMsg(body []byte, cmd int32) (buffer []byte, header *comm.Header, err error) {
	groupIdx := int32(5)
	gid := uint64(10)

	header = &comm.Header{
		Cmdid:   proto.Int32(cmd),
		Gid:     proto.Uint64(gid), // buffer len + checksum len
		Rid:     proto.Uint64(0),
		Version: proto.Int32(comm.Version),
	}


	headerBuf, err := proto.Marshal(header)
	if err != nil {
		log.Errorf("header Marshal fail:%v", err)
		return
	}

	groupIdxBuf := make([]byte, GROUPIDXLEN)
	EncodeInt32(groupIdxBuf, 0, groupIdx)

	headerLenBuf := make([] byte, HEADLEN_LEN)
	EncodeUint16(headerLenBuf, 0, uint16(len(headerBuf)))

	buffer = AppendBytes(groupIdxBuf, headerLenBuf, headerBuf, body)

	ckSum := Crc32(0, buffer, comm.NET_CRC32SKIP)
	ckSumBuf := make([]byte, CHECKSUM_LEN)
	EncodeUint32(ckSumBuf, 0, ckSum)

	buffer = AppendBytes(buffer, ckSumBuf)

	return
}

// TODO to be checked
func unpackBaseMsg(buffer []byte, header *comm.Header) (body []byte, err error) {

	headStartPos := GROUPIDXLEN + HEADLEN_LEN

	var bufferLen = len(buffer)

	if bufferLen < headStartPos {
		log.Error("no head")
		err = comm.ErrInvalidMsg
		return
	}

	var headLen uint16
	DecodeUint16(buffer, GROUPIDXLEN, &headLen)

	if bufferLen < headStartPos+int(headLen) {
		log.Error("msg head lost ")
		err = comm.ErrInvalidMsg
		return
	}

	bodyStartPos := headStartPos + int(headLen)

	proto.Unmarshal(buffer[headStartPos:bodyStartPos], header)

	if bodyStartPos+CHECKSUM_LEN > bufferLen {
		log.Errorf("no checksum, body start pos %d, buffer size %d \r\n", bodyStartPos, bufferLen)
		err = comm.ErrInvalidMsg
		return
	}

	var ckSum uint32
	DecodeUint32(buffer, bufferLen-CHECKSUM_LEN, &ckSum)

	calCkSum := Crc32(0, buffer[:bufferLen-CHECKSUM_LEN], comm.NET_CRC32SKIP)
	if calCkSum != ckSum {
		log.Errorf("data bring ckSum %d not equal to cal ckSum %d \r\n", ckSum, calCkSum)
		err = comm.ErrInvalidMsg
		return
	}

	body = buffer[bodyStartPos: bufferLen-CHECKSUM_LEN]
	err = nil
	return
}

func TestUtil(t *testing.T) {
	body := []byte("hello world")
	cmd := int32(8)

	buffer, header, err := packBaseMsg(body, cmd)
	fmt.Println(header)

	if err != nil {
		fmt.Println("packmsg error ")
	}

	header_ := &comm.Header{}
	body_ ,err := unpackBaseMsg(buffer, header_)
	if err != nil {
		fmt.Println("unpackmsg error ")
	}

	fmt.Println(string(body_))
	fmt.Println(header_)

	fmt.Println(uint64(-1))

}

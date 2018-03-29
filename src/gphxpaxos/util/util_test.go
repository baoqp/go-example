package util

import (
	"testing"

	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"gphxpaxos/comm"
	"gphxpaxos/util"
	"google.golang.org/grpc/balancer/base"
)

// format: groupId(int) + header_len(uint16) + header + body + crc32 checksum(uint32)
func packBaseMsg(body []byte, cmd int32) (buffer []byte, header *comm.Header, err error) {
	groupIdx := 5
	gid := uint64(10)

	h := &comm.Header{
		Cmdid:   proto.Int32(cmd),
		Gid:     proto.Uint64(gid), // buffer len + checksum len
		Rid:     proto.Uint64(0),
		Version: proto.Int32(comm.Version),
	}
	*header = *h // TODO

	headerBuf, err := proto.Marshal(header)
	if err != nil {
		log.Errorf("header Marshal fail:%v", err)
		return
	}

	groupIdxBuf := make([]byte, GROUPIDXLEN)
	util.EncodeInt32(groupIdxBuf, 0, int32(groupIdx))

	headerLenBuf := make([] byte, HEADLEN_LEN)
	util.EncodeUint16(headerLenBuf, 0, uint16(len(headerBuf)))

	buffer = util.AppendBytes(groupIdxBuf, headerLenBuf, headerBuf, body)

	ckSum := util.Crc32(0, buffer, comm.NET_CRC32SKIP)
	ckSumBuf := make([]byte, CHECKSUM_LEN)
	util.EncodeUint32(ckSumBuf, 0, ckSum)

	buffer = util.AppendBytes(buffer, ckSumBuf)

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
	util.DecodeUint16(buffer, GROUPIDXLEN, &headLen)

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
	util.DecodeUint32(buffer, bufferLen-CHECKSUM_LEN, &ckSum)

	calCkSum := util.Crc32(0, buffer[:bufferLen-CHECKSUM_LEN], comm.NET_CRC32SKIP)
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
	paths, _ := IterDir("D:\\tmp\\seaweedfs")
	fmt.Println(paths)

}

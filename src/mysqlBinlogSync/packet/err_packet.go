package packet

import (
	"bytes"
	"mysqlBinlogSync/util"
	"mysqlBinlogSync/comm"
)

// https://dev.mysql.com/doc/internals/en/packet-ERR_Packet.html
type ErrPacket struct {
	CapabilityFlags uint32
	Header          byte
	ErrorCode       uint16
	SqlStateMarker  string
	SqlState        string
	ErrorMessage    string
}

func (errPacket *ErrPacket) Read(data []byte) error {
	errPacket.Header = data[0]
	data = data[1:]
	util.DecodeUint16(data, 0, &errPacket.ErrorCode)
	data = data[2:]
	if errPacket.CapabilityFlags & comm.CLIENT_PROTOCOL_41 > 0 {
		errPacket.SqlStateMarker = string(data[:1])
		errPacket.SqlState = string(data[1:6])
		data = data[6:]
	}

	msgEnd := bytes.IndexByte(data, 0x00) //
	if msgEnd != -1 { // serverVersion 以00结尾
		errPacket.ErrorMessage = string(data[:msgEnd])
	} else {
		errPacket.ErrorMessage = string(data[:])
	}


	return nil
}

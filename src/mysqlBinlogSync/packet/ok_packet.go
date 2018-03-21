package packet

import (
	"bytes"
	"mysqlBinlogSync/util"
	"mysqlBinlogSync/comm"
)

// https://dev.mysql.com/doc/internals/en/packet-OK_Packet.html
type OKPacket struct {
	CapabilityFlags uint32
	Header          byte
	AffectedRows    uint64
	LastInsertId    uint64
	StatusFlag      uint16
	Warnings uint16
	Info string
}

func (okPacket *OKPacket) Read(data []byte) error {
	var n int
	okPacket.Header = data[0]
	data = data[1:]
	okPacket.AffectedRows, _, n = util.LengthEncodedInt(data)
	data = data[n:]
	okPacket.LastInsertId, _, n = util.LengthEncodedInt(data)
	data = data[n:]

	if okPacket.CapabilityFlags & comm.CLIENT_PROTOCOL_41 > 0 {
		util.DecodeUint16(data, 0, &okPacket.StatusFlag)
		util.DecodeUint16(data, 2, &okPacket.Warnings)
		data = data[4:]
	} else if okPacket.CapabilityFlags & comm.CLIENT_TRANSACTIONS > 0 {
		util.DecodeUint16(data, 0, &okPacket.StatusFlag)
		data = data[2:]
	}

	// CLIENT_SESSION_TRACK is not supported
	if  okPacket.CapabilityFlags & comm.CLIENT_SESSION_TRACK > 0 {

	} else {
		infoEnd := bytes.IndexByte(data, 0x00) // serverVersion 以00结尾
		okPacket.Info = string(data[:infoEnd])
	}

	return nil
}

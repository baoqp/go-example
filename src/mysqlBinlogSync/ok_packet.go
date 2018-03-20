package mysqlBinlogSync

import "bytes"

// https://dev.mysql.com/doc/internals/en/packet-OK_Packet.html
type OKPacket struct {
	capabilityFlags uint32
	header          byte
	affectedRows    uint64
	lastInsertId    uint64
	statusFlag      uint16
	warnings uint16
	info string
}

func (okPacket *OKPacket) read(data []byte) error {
	var n int
	okPacket.header = data[0]
	data = data[1:]
	okPacket.affectedRows, _, n = LengthEncodedInt(data)
	data = data[n:]
	okPacket.lastInsertId, _, n = LengthEncodedInt(data)
	data = data[n:]

	if okPacket.capabilityFlags&CLIENT_PROTOCOL_41 > 0 {
		DecodeUint16(data, 0, &okPacket.statusFlag)
		DecodeUint16(data, 2, &okPacket.warnings)
		data = data[4:]
	} else if okPacket.capabilityFlags&CLIENT_TRANSACTIONS > 0 {
		DecodeUint16(data, 0, &okPacket.statusFlag)
		data = data[2:]
	}

	// CLIENT_SESSION_TRACK is not supported
	if  okPacket.capabilityFlags& CLIENT_SESSION_TRACK > 0 {

	} else {
		infoEnd := bytes.IndexByte(data, 0x00) // serverVersion 以00结尾
		okPacket.info = string(data[:infoEnd])
	}

	return nil
}

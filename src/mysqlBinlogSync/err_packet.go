package mysqlBinlogSync

import "bytes"

// https://dev.mysql.com/doc/internals/en/packet-ERR_Packet.html
type ErrPacket struct {
	capabilityFlags uint32
	header          byte
	errorCode       uint16
	sqlStateMarker  string
	sqlState        string
	errorMessage    string
}

func (errPacket *ErrPacket) read(data []byte) error {

	errPacket.header = data[0]
	data = data[1:]
	DecodeUint16(data, 1, &errPacket.errorCode)
	data = data[3:]
	if errPacket.capabilityFlags&CLIENT_PROTOCOL_41 > 0 {
		errPacket.sqlStateMarker = string(data[:1])
		errPacket.sqlState = string(data[1:6])
		data = data[6:]
	}

	msgEnd := bytes.IndexByte(data, 0x00) // serverVersion 以00结尾
	errPacket.errorMessage = string(data[:msgEnd])

	return nil
}

package mysqlBinlogSync

// https://dev.mysql.com/doc/internals/en/packet-ERR_Packet.html
type ErrPacket struct {
	header         byte
	errorCode      uint16
	sqlStateMarker string
	sqlState       string
	errorMessage   string
}

func (errPacket *ErrPacket) read(data []byte) error {

	return nil
}

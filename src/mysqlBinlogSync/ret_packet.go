package mysqlBinlogSync

type RetPacket struct {
	*OKPacket
	*ErrPacket
	isOk bool
}

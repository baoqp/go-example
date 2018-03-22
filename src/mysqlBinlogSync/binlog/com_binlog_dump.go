package binlog

import (
	"mysqlBinlogSync/comm"
	"mysqlBinlogSync/util"
)

// https://dev.mysql.com/doc/internals/en/com-binlog-dump.html
type ComBinlogDump struct {
	Position
	Flags uint16
	ServerId uint32
}

func (cbd *ComBinlogDump) Write() ([]byte, error) {
	//4 (header)+  1 + 4 + 2 + 4 + len(name) + 1
	var n int

	length := 16 + len(cbd.Name)

	data := make([]byte, length)
	pos := 4
	data[pos] = comm.COM_BINLOG_DUMP
	pos++

	util.EncodeUint32(data, pos, cbd.Pos)
	pos += 4

	util.EncodeUint16(data, pos, cbd.Flags)
	pos += 2

	//replication rank, not used
	util.EncodeUint32(data, pos, cbd.ServerId)
	pos += 4

	n = copy(data[pos:], cbd.Name)
	pos += n

	data[pos] = 0

	return data, nil
}



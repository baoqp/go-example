package binlog

import (
	"mysqlBinlogSync/comm"
	"mysqlBinlogSync/util"

)

// https://dev.mysql.com/doc/internals/en/com-register-slave.html
type ComRegisterSlave struct {
	ServerId uint32
	HostName string
	User     string
	Password string
	Port     uint16
	MasterId uint32
}


func (crs *ComRegisterSlave) Write() ([]byte, error) {
	//4 (header)+  1 + 4 + 1 + len(crs.HostName) + 1 + len(crs.User) + 1 + len(crs.Password) + 2 + 4 + 4
	var n int

	length := 22 + len(crs.HostName) + len(crs.User) + len(crs.Password)

	data := make([]byte, length)
	pos := 4
	data[pos] = comm.COM_REGISTER_SLAVE
	pos++

	util.EncodeUint32(data, pos, crs.ServerId)
	pos += 4

	data[pos] = uint8(len(crs.HostName))
	pos++
	n = copy(data[pos:], crs.HostName)
	pos += n

	data[pos] = uint8(len(crs.User))
	pos++
	n = copy(data[pos:], crs.User)
	pos += n

	data[pos] = uint8(len(crs.Password))
	pos++
	n = copy(data[pos:], crs.Password)
	pos += n

	util.EncodeUint16(data, pos, crs.Port)
	pos += 2

	//replication rank, not used
	util.EncodeUint32(data, pos, 0)
	pos += 4

	util.EncodeUint32(data, pos, crs.MasterId)

	return data, nil
}

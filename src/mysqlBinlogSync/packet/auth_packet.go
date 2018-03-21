package packet

import (
	"mysqlBinlogSync/util"
	"mysqlBinlogSync/comm"
)

type AuthPacket struct {
	CapabilityFlags uint32
	Salt            []byte
	User            string
	Passwd          string
	Auth            []byte
	DBName          string
}

func (authPacket *AuthPacket) Write() ([]byte, error) {

	authPacket.Auth = util.EncodePassword(authPacket.Salt, []byte(authPacket.Passwd))

	length := authPacket.calcPacketSize()
	data := make([]byte, length+4) // 4是每个packet都有的头部

	if len(authPacket.DBName) > 0 {
		authPacket.CapabilityFlags |= comm.CLIENT_CONNECT_WITH_DB
	}

	data[4] = byte(authPacket.CapabilityFlags)
	data[5] = byte(authPacket.CapabilityFlags >> 8)
	data[6] = byte(authPacket.CapabilityFlags >> 16)
	data[7] = byte(authPacket.CapabilityFlags >> 24)
	data[12] = byte(comm.DEFAULT_COLLATION_ID)

	// 23个filler
	pos := 13 + 23
	if len(authPacket.User) > 0 {
		pos += copy(data[pos:], authPacket.User)
	}
	data[pos] = 0x00
	pos++

	// auth [length encoded integer]
	data[pos] = byte(len(authPacket.Auth))
	pos += 1 + copy(data[pos+1:], authPacket.Auth)

	// db [null terminated string]
	if len(authPacket.DBName) > 0 {
		pos += copy(data[pos:], authPacket.DBName)
		data[pos] = 0x00
		pos++
	}

	// Assume native client during response
	pos += copy(data[pos:], "mysql_native_password")
	data[pos] = 0x00

	return data, nil
}

func (authPacket *AuthPacket) calcPacketSize() int {
	size := 32 // 4+4+1+23;
	size += len(authPacket.User) + 1
	size += len(authPacket.Auth) + 1
	size += len(authPacket.DBName) + 1
	size += 21 + 1 // mysql_native_password + null-terminated
	return size
}

package mysqlBinlogSync

import (
	"crypto/tls"
)

type AuthPacket struct {
	capabilityFlags uint32
	salt []byte
	user          string
	passwd        string
	auth []byte
	dbName        string
}

func (authPacket *AuthPacket) write() ([]byte, error) {

	authPacket.auth = EncodePassword(authPacket.salt, []byte(authPacket.passwd))

	length := authPacket.calcPacketSize()
	data := make([]byte, length+4) // 4是每个packet都有的头部

	if len(authPacket.dbName) > 0 {
		authPacket.capabilityFlags |= CLIENT_CONNECT_WITH_DB
	}

	data[4] = byte(authPacket.capabilityFlags)
	data[5] = byte(authPacket.capabilityFlags >> 8)
	data[6] = byte(authPacket.capabilityFlags >> 16)
	data[7] = byte(authPacket.capabilityFlags >> 24)
	data[12] = byte(DEFAULT_COLLATION_ID)

	// 23个filler
	pos := 13 + 23
	if len(authPacket.user) > 0 {
		pos += copy(data[pos:], authPacket.user)
	}
	data[pos] = 0x00
	pos++

	// auth [length encoded integer]
	data[pos] = byte(len(authPacket.auth))
	pos += 1 + copy(data[pos+1:], authPacket.auth)

	// db [null terminated string]
	if len(authPacket.dbName) > 0 {
		pos += copy(data[pos:], authPacket.dbName)
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
	size += len(authPacket.user) + 1
	size += len(authPacket.auth) + 1
	size += len(authPacket.dbName) + 1
	size += 21 + 1  // mysql_native_password + null-terminated
	return size
}

package mysqlBinlogSync

import (
	"fmt"
	"bytes"
	"errors"
)

// https://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::Handshake

type HandshakePacket struct {
	protocolVersion byte
	serverVersion   string
	connectionId    uint32
	capabilityFlags uint32
	characterSet    uint32
	statusFlag      uint16
	salt            []byte
}

func NewHandshakePacket() *HandshakePacket {
	return &HandshakePacket{}
}


func (handshakePacket *HandshakePacket) read(data []byte) error {

	handshakePacket.protocolVersion = data[0]

	if handshakePacket.protocolVersion < MinProtocolVersion {
		return errors.New("unsupported protocol version")
	}

	data = data[1:]
	serverVersionEnd := bytes.IndexByte(data, 0x00) // serverVersion 以00结尾
	handshakePacket.serverVersion = string(data[:serverVersionEnd])
	data = data[serverVersionEnd+1:] // 跳过 00
	DecodeUint32(data, 0, &handshakePacket.connectionId)
	data = data[4:]
	handshakePacket.salt = data[0:8]
	data = data[9:]
	var capabilityFlags uint16
	DecodeUint16(data, 0, &capabilityFlags)
	handshakePacket.capabilityFlags = uint32(capabilityFlags)
	data = data[2:]

	if len(data) > 0 {
		handshakePacket.characterSet = uint32(data[0])
		data = data[1:]
		DecodeUint16(data, 0, &handshakePacket.statusFlag)
		data = data[2:]
		DecodeUint16(data, 0, &capabilityFlags)
		handshakePacket.capabilityFlags = uint32(capabilityFlags)<<16 | handshakePacket.capabilityFlags
		data = data[13:]
		handshakePacket.salt = append(handshakePacket.salt, data[:12]...)
	}

	return nil
}

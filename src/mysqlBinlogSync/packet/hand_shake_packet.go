package packet

import (
	"bytes"
	"errors"
	"mysqlBinlogSync/comm"
	"mysqlBinlogSync/util"
)

// https://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::Handshake

type HandshakePacket struct {
	ProtocolVersion byte
	ServerVersion   string
	ConnectionId    uint32
	CapabilityFlags uint32
	CharacterSet    uint32
	StatusFlag      uint16
	Salt            []byte
}

func NewHandshakePacket() *HandshakePacket {
	return &HandshakePacket{}
}

func (handshakePacket *HandshakePacket) Read(data []byte) error {

	handshakePacket.ProtocolVersion = data[0]

	if handshakePacket.ProtocolVersion < comm.MinProtocolVersion {
		return errors.New("unsupported protocol version")
	}

	data = data[1:]
	serverVersionEnd := bytes.IndexByte(data, 0x00) // serverVersion 以00结尾
	handshakePacket.ServerVersion = string(data[:serverVersionEnd])
	data = data[serverVersionEnd+1:] // 跳过 00
	util.DecodeUint32(data, 0, &handshakePacket.ConnectionId)
	data = data[4:]
	handshakePacket.Salt = data[0:8]
	data = data[9:]
	var capabilityFlags uint16
	util.DecodeUint16(data, 0, &capabilityFlags)
	handshakePacket.CapabilityFlags = uint32(capabilityFlags)
	data = data[2:]

	if len(data) > 0 {
		handshakePacket.CharacterSet = uint32(data[0])
		data = data[1:]
		util.DecodeUint16(data, 0, &handshakePacket.StatusFlag)
		data = data[2:]
		util.DecodeUint16(data, 0, &capabilityFlags)
		handshakePacket.CapabilityFlags = uint32(capabilityFlags)<<16 | handshakePacket.CapabilityFlags
		data = data[13:]
		handshakePacket.Salt = append(handshakePacket.Salt, data[:12]...)
	}

	return nil
}

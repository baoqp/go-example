package binlog

import (
	"fmt"
	"gphxpaxos/util"
)

const (
	EventHeaderSize = 19
)

type EventHeader struct {
	Timestamp uint32
	EventType EventType
	ServerID  uint32
	EventSize uint32
	LogPos    uint32
	Flags     uint16
}

func (header *EventHeader) Read(data []byte) error {
	if len(data) < EventHeaderSize {
		return fmt.Errorf("header size too short %d, must 19", len(data))
	}

	pos := 0

	util.DecodeUint32(data, pos, &header.Timestamp)
	pos += 4

	header.EventType = EventType(data[pos])
	pos++

	util.DecodeUint32(data, pos, &header.ServerID)
	pos += 4

	util.DecodeUint32(data, pos, &header.EventSize)
	pos += 4

	util.DecodeUint32(data, pos, &header.LogPos)
	pos += 4

	util.DecodeUint16(data, pos, &header.Flags)
	pos += 2

	if header.EventSize < uint32(EventHeaderSize) {
		return fmt.Errorf("invalid event size %d, must >= 19", header.EventSize)
	}

	return nil
}

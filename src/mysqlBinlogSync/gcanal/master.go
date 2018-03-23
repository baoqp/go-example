package gcanal

import (
	"sync"
	"mysqlBinlogSync/binlog"
)

type masterInfo struct {
	sync.RWMutex
	pos binlog.Position
}

func (m *masterInfo) Update(pos binlog.Position) {
	m.Lock()
	m.pos = pos
	m.Unlock()
}

func (m *masterInfo) Position() binlog.Position {
	m.RLock()
	defer m.RUnlock()

	return m.pos
}

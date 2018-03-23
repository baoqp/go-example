package gcanal

import (
	"mysqlBinlogSync/binlog"
)


type RowAction string
const(
	UpdateAction RowAction = "update"
	InsertAction RowAction = "insert"
	DeleteAction RowAction = "delete"
)


type EventHandler interface {
	OnRotate(rotateEvent *binlog.RotateEvent) error
	OnDDL(nextPos binlog.Position, queryEvent *binlog.QueryEvent) error
	OnRow(action RowAction, e *binlog.RowsEvent) error
	OnXID(nextPos binlog.Position) error

	// OnPosSynced Use your own way to sync position. When force is true, sync position immediately.
	OnPosSynced(pos binlog.Position, force bool) error
	String() string
}

type DummyEventHandler struct {
}

func (h *DummyEventHandler) OnRotate(*binlog.RotateEvent) error { return nil }
func (h *DummyEventHandler) OnDDL(binlog.Position, *binlog.QueryEvent) error {
	return nil
}
func (h *DummyEventHandler) OnRow(action RowAction, e *binlog.RowsEvent) error { return nil }
func (h *DummyEventHandler) OnXID(binlog.Position) error   { return nil }

func (h *DummyEventHandler) OnPosSynced(binlog.Position, bool) error { return nil }
func (h *DummyEventHandler) String() string                          { return "DummyEventHandler" }

// `SetEventHandler` registers the sync handler, you must register your
// own handler before starting Canal.
func (gCanal *GCanal) SetEventHandler(h EventHandler) {
	gCanal.eventHandler = h
}

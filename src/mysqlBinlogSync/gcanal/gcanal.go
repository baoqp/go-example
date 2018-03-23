package gcanal

import (
	"sync"
	"mysqlBinlogSync/binlog"
	log "github.com/sirupsen/logrus"
	"fmt"
	"context"
)

type GCanal struct {
	m            sync.Mutex
	master       *masterInfo
	syncer       *binlog.BinlogSyncer
	syncCfg      *binlog.SyncConfig
	eventHandler EventHandler

	ctx    context.Context
	cancel context.CancelFunc
}

func NewGCanal(config *binlog.SyncConfig, startPos binlog.Position, handler EventHandler) (*GCanal, error) {
	gCanal := &GCanal{}

	gCanal.ctx, gCanal.cancel = context.WithCancel(context.Background())
	gCanal.syncCfg = config
	gCanal.syncer = binlog.NewBinlogSyncer(config)

	if handler == nil {
		gCanal.eventHandler = &DummyEventHandler{}
	} else {
		gCanal.eventHandler = handler
	}

	gCanal.master = &masterInfo{}
	gCanal.master.Update(startPos)

	return gCanal, nil
}

func (gCanal *GCanal) startSyncer() (*binlog.BinlogStreamer, error) {
	pos := gCanal.master.Position()
	s, err := gCanal.syncer.StartSync(pos)
	if err != nil {
		return nil, fmt.Errorf("start sync replication at binlog %v error %v", pos, err)
	}
	log.Infof("start sync binlog at binlog file %v", pos)
	return s, nil

}

func (gCanal *GCanal) Run() error {
	return gCanal.runSyncBinlog()
}

func (gCanal *GCanal) runSyncBinlog() error {
	s, err := gCanal.startSyncer()
	if err != nil {
		return err
	}

	savePos := false
	force := false
	for {
		ev, err := s.GetEvent(gCanal.ctx)

		if err != nil {
			return err
		}
		savePos = false
		force = false
		pos := gCanal.master.Position()

		pos.Pos = ev.Header.LogPos

		// We only save position with RotateEvent and XIDEvent.
		// For RowsEvent, we can't save the position until meeting XIDEvent
		// which tells the whole transaction is over.
		// TODO : 但是这样XIDEvent之前的RowsEvent还是会消费的，如果存在回滚等会不会存在数据一致性问题。
		// TODO: If we meet any DDL query, we must save too.
		switch e := ev.Event.(type) {
		case *binlog.RotateEvent:
			pos.Name = string(e.NextLogName)
			pos.Pos = uint32(e.Position)
			log.Infof("rotate binlog to %s", pos)
			savePos = true
			force = true
			if err = gCanal.eventHandler.OnRotate(e); err != nil { // TODO 回调处理错误应该由调用方处理
				return err
			}
		case *binlog.RowsEvent:
			// we only focus row based event
			var action RowAction
			switch ev.Header.EventType {
			case binlog.WRITE_ROWS_EVENTv1, binlog.WRITE_ROWS_EVENTv2:
				action = InsertAction
			case binlog.DELETE_ROWS_EVENTv1, binlog.DELETE_ROWS_EVENTv2:
				action = DeleteAction
			case binlog.UPDATE_ROWS_EVENTv1, binlog.UPDATE_ROWS_EVENTv2:
				action = UpdateAction
			default:
				return fmt.Errorf("%s not supported now", ev.Header.EventType)
			}

			err = gCanal.eventHandler.OnRow(action, e)
			if err != nil {
				// TODO
			}
			continue

		case *binlog.XIDEvent:
			savePos = true
			if err := gCanal.eventHandler.OnXID(pos); err != nil {
				return err
			}
		case *binlog.QueryEvent:
			// TODO ddl 操作

		default:
			continue
		}

		if savePos {
			gCanal.master.Update(pos)
			gCanal.eventHandler.OnPosSynced(pos, force)
		}
	}

	return nil
}

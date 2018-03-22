package binlog

import (
	"mysqlBinlogSync/conn"
	"time"
	"os"
	"fmt"
	"mysqlBinlogSync/packet"
	"sync"
	log "github.com/sirupsen/logrus"
	"mysqlBinlogSync/util"
	"context"
	"mysqlBinlogSync/comm"
	"errors"
	"bytes"
)

type SyncStatus int

const (
	Starting  SyncStatus = iota
	Preparing
	Running
	Stop
)

type SyncConfig struct {
	Host        string
	Port        uint16
	User        string
	Password    string
	DBName      string
	ReadTimeout time.Duration
	ServerId    uint32
	MasterId    uint32
	Localhost   string
	UseDecimal  bool
	ParseTime   bool
}

type BinlogSyncer struct {
	cfg     *SyncConfig
	conn    *conn.Connection
	status  SyncStatus
	lock    sync.RWMutex
	wg      sync.WaitGroup
	nextPos Position
	ctx     context.Context
	cancel  context.CancelFunc
	parser  *EventParser
}

func NewBinlogSyncer(cfg *SyncConfig) *BinlogSyncer {
	bs := &BinlogSyncer{cfg: cfg}
	bs.parser = NewEventParser()
	bs.parser.useDecimal = cfg.UseDecimal
	bs.parser.parseTime = cfg.ParseTime
	bs.ctx, bs.cancel = context.WithCancel(context.Background())
	return bs
}

func (bs *BinlogSyncer) StartSync(pos Position) (*BinlogStreamer, error) {
	log.Infof("begin to sync binlog from position %s", pos)

	bs.lock.Lock()
	defer bs.lock.Unlock()

	bs.status = Starting

	if err := bs.prepareSyncPos(pos); err != nil {
		return nil, err
	}

	return bs.startInfluxStream(), nil
}

func (bs *BinlogSyncer) startInfluxStream() *BinlogStreamer {
	bs.status = Running
	s := newBinlogStreamer()
	bs.wg.Add(1)
	bs.onStream(s)
	return s
}

// 接收binlog数据，相当于一个acceptor
func (bs *BinlogSyncer) onStream(s *BinlogStreamer) {
	defer func() {
		if e := recover(); e != nil {
			s.closeWithError(fmt.Errorf("Err: %v\n Stack: %s", e, util.Pstack()))
		}
		bs.wg.Done()
	}()

	for {
		data, err := bs.conn.ReadPacket()

		if err != nil {
			log.Error(err)

			//reconnect,  last nextPos or nextGTID we got.
			if len(bs.nextPos.Name) == 0 {
				// we can't get the correct position, close.
				s.closeWithError(err)
				return
			}

			// TODO: add a max retry count.
			for {
				select {
				case <-bs.ctx.Done():
					s.close()
					return
				case <-time.After(time.Second):
					if err = bs.retrySync(); err != nil {
						log.Errorf("retry sync err: %v, wait 1s and retry again", err)
						continue
					}
				}

				break
			}

			// we connect the server and begin to re-sync again.
			continue
		}

		//set read timeout
		switch data[0] {
		case comm.OK_HEADER:
			// TODO 这里最好开个goroutine
			if err = bs.parseEvent(s, data); err != nil {
				s.closeWithError(err)
				return
			}
		case comm.ERR_HEADER:
			errPacket := &packet.ErrPacket{}
			_ = errPacket.Read(data)
			err = fmt.Errorf("ErrPacket with ErrCode:%d, ErrMsg:%s", errPacket.ErrorCode, errPacket.ErrorMessage)
			log.Error(err)
			s.closeWithError(err)
			return
		case comm.EOF_HEADER:
			// Refer http://dev.mysql.com/doc/internals/en/packet-EOF_Packet.html
			// In the MySQL client/server protocol, EOF and OK packets serve the same purpose.
			// Some users told me that they received EOF packet here, but I don't know why.
			// So we only log a message and retry ReadPacket.
			log.Info("receive EOF packet, retry ReadPacket")
			continue
		default:
			log.Errorf("invalid stream header %c", data[0])
			continue
		}
	}
}

func (bs *BinlogSyncer) parseEvent(s *BinlogStreamer, data []byte) error {
	//skip OK byte, 0x00
	data = data[1:]

	e, err := bs.parser.Parse(data)

	if _, cast := e.Event.(*RowsEvent); cast {
		buff := new(bytes.Buffer)
		e.Event.Write(buff)
		log.Info(string(buff.Bytes()))
	}


	if err != nil {
		return err
	}

	if e.Header.LogPos > 0 {
		// Some events like FormatDescriptionEvent return 0, ignore.
		bs.nextPos.Pos = e.Header.LogPos
	}
	switch event := e.Event.(type) {
	case *RotateEvent:
		bs.nextPos.Name = string(event.NextLogName)
		bs.nextPos.Pos = uint32(event.Position)
		log.Infof("rotate to %s", bs.nextPos)
	default:

	}

	needStop := false
	select {
	case s.ch <- e:
	case <-bs.ctx.Done():
		needStop = true
	}

	if needStop {
		return errors.New("sync is been closing...")
	}

	return nil
}

func (bs *BinlogSyncer) retrySync() error {
	bs.lock.Lock()
	defer bs.lock.Unlock()

	bs.parser.Reset()

	log.Infof("begin to re-sync from %s", bs.nextPos)
	if err := bs.prepareSyncPos(bs.nextPos); err != nil {
		return err
	}

	return nil
}

func (bs *BinlogSyncer) prepareSyncPos(pos Position) error {

	bs.status = Preparing

	// 先注册
	if err := bs.RegisterAsSlave(); err != nil {
		return err
	}

	// always start from position 4 TODO  why ???
	if pos.Pos < 4 {
		pos.Pos = 4
	}

	if err := bs.WriteBinlogDumpCommand(pos); err != nil {
		return err
	}

	return nil
}

func (bs *BinlogSyncer) WriteBinlogDumpCommand(p Position) error {
	bs.conn.ResetSequence()
	cbd := &ComBinlogDump{
		Position: p,
		Flags:    BINLOG_DUMP_NEVER_STOP,
		ServerId: bs.cfg.ServerId,
	}

	data, _ := cbd.Write()
	return bs.conn.WritePacket(data)
}

func (bs *BinlogSyncer) RegisterAsSlave() error {
	var err error
	if bs.conn != nil {
		bs.conn.Close()
	}

	addr := fmt.Sprintf("%s:%d", bs.cfg.Host, bs.cfg.Port)

	bs.conn, err = conn.Connect(addr, bs.cfg.User, bs.cfg.Password, bs.cfg.DBName)

	if err != nil {
		return err
	}

	// TODO connection setting

	// TODO check old connection and kill it

	// send RegisterSlaveCommand
	comRegisterSlave := &ComRegisterSlave{
		ServerId: bs.cfg.ServerId,
		HostName: bs.localHostname(),
		User:     bs.cfg.User,
		Password: bs.cfg.Password,
		Port:     bs.cfg.Port,
		MasterId: bs.cfg.MasterId,
	}

	data, _ := comRegisterSlave.Write()

	bs.conn.ResetSequence() // 每次写command都需要reset sequence
	if err = bs.conn.WritePacket(data); err != nil {
		bs.conn.Close()
		return err
	}

	// 读取server端回复
	var retPacket *packet.RetPacket
	if retPacket, err = bs.conn.ReadRet(); err != nil {
		bs.conn.Close()
		return err
	}

	if !retPacket.IsOk {
		bs.conn.Close()
		return fmt.Errorf("auth failed, errorCode:%d, errorMsg:%s",
			retPacket.ErrPacket.ErrorCode, retPacket.ErrPacket.ErrorMessage)
	}

	return nil
}

func (bs *BinlogSyncer) localHostname() string {
	if len(bs.cfg.Localhost) == 0 {
		h, _ := os.Hostname()
		return h
	}
	return bs.cfg.Localhost
}

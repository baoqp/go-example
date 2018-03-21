package binlog

import (
	"mysqlBinlogSync/conn"
	"time"
	"mysqlBinlogSync/command"
	"os"
	"fmt"
	"mysqlBinlogSync/packet"
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
}

type BinlogSyncer struct {
	cfg  *SyncConfig
	conn *conn.Connection
}

func NewBinlogSyncer(cfg *SyncConfig) *BinlogSyncer {
	return &BinlogSyncer{cfg: cfg}
}


func(bs *BinlogSyncer) prepareSyncPos(pos Position) error {
	// 先注册
	if err := bs.RegisterAsSlave(); err != nil {
		return err
	}

	// always start from position 4 TODO  why ???
	if pos.Pos < 4 {
		pos.Pos = 4
	}




	return nil
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
	comRegisterSlave := &command.ComRegisterSlave{
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

package mysqlBinlogSync

import (
	"testing"
	"fmt"
	"mysqlBinlogSync/binlog"
)

func Test(t *testing.T) {
	/*_, err := conn.Connect("45.77.120.142:3306","root", "bqp0205", "dev" )
		if err != nil {
			fmt.Println(err)
		}*/
	// TODO 目前必须设置master的binlog选项 global binlog_checksum='NONE';
	cfg := &binlog.SyncConfig{
		Host:       "45.77.120.142",
		Port:       3306,
		User:       "root",
		Password:   "bqp0205",
		DBName:     "dev",
		ServerId:   2,
		MasterId:   1,
		UseDecimal: true,
		ParseTime:  true,
	}

	bs := binlog.NewBinlogSyncer(cfg)

	_, err := bs.StartSync(binlog.Position{"mysql-bin.000004", 150})

	if err != nil {
		fmt.Println(err)
	}

}

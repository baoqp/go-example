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


	cfg := &binlog.SyncConfig{
		Host:     "45.77.120.142",
		Port:     3306,
		User:     "root",
		Password: "bqp0205",
		DBName:   "dev",
		ServerId: 2,
		MasterId: 1,
	}

	bs := binlog.NewBinlogSyncer(cfg)

	err := bs.RegisterAsSlave()

	if err != nil {
		fmt.Println(err)
	}

}

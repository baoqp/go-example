package mysqlBinlogSync

import (
	"testing"
	"fmt"
	"mysqlBinlogSync/binlog"
	"time"
	"mysqlBinlogSync/gcanal"
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
		UseDecimal: false,
		ParseTime:  true,
	}

	gcanal, _ := gcanal.NewGCanal(cfg, binlog.Position{"mysql-bin.000004", 150}, &TestHandler{})

	err := gcanal.Run()

	if err != nil {
		fmt.Println(err)
	}
}

/**
测试表
CREATE TABLE `t_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `address` text,
  `sex` tinyint(1) DEFAULT '1',
  `createAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8
 */
type User struct {
	id      int32
	name    string
	address string
	sex     int8
	create  time.Time
}

type TestHandler struct {
}

func (handler *TestHandler) OnRotate(rotateEvent *binlog.RotateEvent) error {

	fmt.Printf("--TestHandler, OnTRotate, Pos:%d, NextLogName %s \r\n", rotateEvent.Position, string(rotateEvent.NextLogName))

	return nil
}

func (handler *TestHandler) OnDDL(nextPos binlog.Position, queryEvent *binlog.QueryEvent) error {
	return nil
}

func (handler *TestHandler) OnRow(action gcanal.RowAction, e *binlog.RowsEvent) error {
	schema := string(e.Table.Schema)
	table := string(e.Table.Table)
	fmt.Printf("--Receive Row Event from schema:%s, table %s, with action %s \r\n ", schema, table, action)

	var users = make([]*User, len(e.Rows))
	for i, row := range e.Rows {
		user := &User{}
		user.id = row[0].(int32)
		user.name = row[1].(string)
		user.address = string(row[2].([]byte))
		user.sex = row[3].(int8)
		user.create = row[4].(time.Time)
		users[i] = user
		fmt.Println(user)
	}
	return nil
}

func (handler *TestHandler) OnXID(nextPos binlog.Position) error {

	return nil
}

func (handler *TestHandler) OnPosSynced(pos binlog.Position, force bool) error {
	return nil
}

func (handler *TestHandler) String() string {
	return "Test"
}

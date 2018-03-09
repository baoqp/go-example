package server

import (
	"github.com/chrislusf/raft"
	"fmt")

type OPCommand struct {
	Type   []byte
	Bucket []byte
	Key    []byte
	Value  []byte
}

func NewOPCommand(Type []byte, Bucket []byte, Key []byte, Value []byte) *OPCommand {
	return &OPCommand{Type:Type, Bucket:Bucket, Key:Key, Value:Value}
}

func (opCommand *OPCommand) CommandName() string {
	return "OP"
}

func (opCommand *OPCommand) Apply( server raft.Server) (interface{}, error) {

	fmt.Println("Apply Command " + string(opCommand.Type), " ", string(opCommand.Bucket), " ",
		string(opCommand.Key), " ",string(opCommand.Value))

	kvServer := server.Context().(*KVServer)

	var err error = nil
	var value []byte = nil

	switch string(opCommand.Type) {
	case "GET":
		value, err = kvServer.GET(opCommand.Bucket, opCommand.Key)
	case "PUT", "POST":
		err = kvServer.Put(opCommand.Bucket, opCommand.Key, opCommand.Value)
	case "DELETE":
		err = kvServer.Delete(opCommand.Bucket, opCommand.Key)
	}

	return value, err
}



package server

import "github.com/chrislusf/raft"

type OPCommand struct {
	Bucket []byte
	Key    []byte
	Value  []byte
}

func NewOPCommand(Bucket []byte, Key []byte, Value []byte) *OPCommand {
	return &OPCommand{Bucket:Bucket, Key:Key, Value:Value}
}

func (opCommand *OPCommand) CommandName() string {
	return "OP"
}


func (opCommand *OPCommand) Apply( server raft.Server) (interface{}, error) {


	return nil, nil
}
package graft

import "github.com/golang/protobuf/proto/proto3_proto"



type ProtoBufFile struct {
	path string
}

func NewProtoBufFile(path string) *ProtoBufFile {
	return &ProtoBufFile{path:path}
}

// TODO
func (protoBufFile *ProtoBufFile) save(message proto3_proto.Message) {

}

// TODO
func (protoBufFile *ProtoBufFile) load(message *proto3_proto.Message)  {

}
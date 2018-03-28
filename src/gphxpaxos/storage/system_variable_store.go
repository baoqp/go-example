package storage

import (
	"gphxpaxos/comm"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

type SystemVariablesStore struct {
	logstorage LogStorage
}

func NewSystemVariablesStore(logstorage LogStorage) *SystemVariablesStore {
	return &SystemVariablesStore{
		logstorage: logstorage,
	}
}

func (s *SystemVariablesStore) Write(writeOptions *WriteOptions, groupId int,
	variables *comm.SystemVariables) error {

	buffer, err := proto.Marshal(variables)

	if err != nil {
		log.Errorf("Variables.Serialize fail")
		return nil
	}

	err = s.logstorage.SetSystemVariables(writeOptions, groupId, buffer)

	if err != nil {
		log.Errorf("DB.Put fail, groupidx %d bufferlen %zu ret %v",
			groupId, len(buffer), err)
		return err
	}

	return nil
}


func (s *SystemVariablesStore) Read(groupId int, variables *comm.SystemVariables) error {
	buffer, err := s.logstorage.GetSystemVariables(groupId)
	if err != nil {
		return err
	}
	// TODO not found error

	 return proto.Unmarshal(buffer, variables)
}
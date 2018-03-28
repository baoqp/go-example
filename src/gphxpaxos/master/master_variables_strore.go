package master

import (
	log "github.com/sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"gphxpaxos/storage"
	"gphxpaxos/comm"
)

type MasterVariablesStore struct {
	logstorage storage.LogStorage
}

func NewMasterVariablesStore(logstorage storage.LogStorage) *MasterVariablesStore {
	return &MasterVariablesStore{
		logstorage: logstorage,
	}
}

func (s *MasterVariablesStore) Write(writeOptions *storage.WriteOptions, groupId int,
	variables *comm.MasterVariables) error {

	buffer, err := proto.Marshal(variables)

	if err != nil {
		log.Errorf("Variables.Serialize fail")
		return nil
	}

	err = s.logstorage.SetMasterVariables(writeOptions, groupId, buffer)

	if err != nil {
		log.Errorf("DB.Put fail, groupidx %d bufferlen %zu ret %v",
			groupId, len(buffer), err)
		return err
	}

	return nil
}

func (s *MasterVariablesStore) Read(groupId int, variables *comm.MasterVariables) error {
	buffer, err := s.logstorage.GetMasterVariables(groupId)
	if err != nil {
		return err
	}
	// TODO not found error

	return proto.Unmarshal(buffer, variables)
}

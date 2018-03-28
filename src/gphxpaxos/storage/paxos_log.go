package storage

import (
	"gphxpaxos/comm"
	"github.com/golang/protobuf/proto"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type PaxosLog struct {
	logStorage LogStorage
}

func NewPaxosLog(logStorage LogStorage) *PaxosLog {
	return &PaxosLog{logStorage: logStorage}
}

func (paxosLog *PaxosLog) WriteLog(options *WriteOptions, groupIdx int,
	instanceId uint64, value []byte) error {

	state := &comm.AcceptorStateData{
		InstanceID:     &instanceId, //也可以使用 proto.Uint64 包装下
		AcceptedValue:  value,
		PromiseID:      &comm.UINT64_0,
		PromiseNodeID:  &comm.UINT64_0,
		AcceptedID:     &comm.UINT64_0,
		AcceptedNodeID: &comm.UINT64_0,
	}
	return paxosLog.WriteState(options, groupIdx, instanceId, state)
}

func (paxosLog *PaxosLog) ReadLog(groupIdx int, instanceId uint64) ([]byte, error) {
	var state = &comm.AcceptorStateData{}
	err := paxosLog.ReadState(groupIdx, instanceId, state)

	if err != nil {
		log.Errorf("Read log error ")
		return nil, err
	}

	value := state.AcceptedValue

	return value, nil
}

func (paxosLog *PaxosLog) GetMaxInstanceIdFromLog(groupIdx int) (uint64, error) {
	instanceId, err := paxosLog.logStorage.GetMaxInstanceId(groupIdx)
	if err != nil {
		log.Errorf("db.getmax fail, error:%v", err)
		return comm.INVALID_INSTANCEID, err
	}

	return instanceId, nil
}

func (paxosLog *PaxosLog) WriteState(options *WriteOptions, groupIdx int, instanceId uint64,
	state *comm.AcceptorStateData) error {

	value, err := proto.Marshal(state)
	if err != nil {
		return fmt.Errorf("marshal state error")
	}

	err = paxosLog.logStorage.Put(options, groupIdx, instanceId, value)

	if err != nil {
		log.Errorf("write state error")
		return err
	}
	return nil
}

func (paxosLog *PaxosLog) ReadState(groupIdx int, instanceId uint64, state *comm.AcceptorStateData) (error) {
	value, err := paxosLog.logStorage.Get(groupIdx, instanceId)

	if err != nil {
	}
	err = proto.Unmarshal(value, state)

	if err != nil {
		log.Errorf("Read State error caused by unmarshal error ")
		return err
	}

	return nil
}

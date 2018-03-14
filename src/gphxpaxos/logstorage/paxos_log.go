package logstorage

import (
	"gphxpaxos/comm"
	"github.com/gogo/protobuf/proto"
	"fmt"
	log "github.com/sirupsen/logrus"
	"log"
)

type PaxosLog struct {
	logStorage LogStorage
}

func (paxosLog *PaxosLog) WriteLog(options *WriteOptions, groupIdx int,
	instanceId uint64, value []byte) error {

	state := &comm.AcceptorStateData{
		InstanceID:     &instanceId, // TODO  proto 生成的代码里面是*uint64类型，为什么不是uint64
		AcceptedValue:  value,
		PromiseID:      &comm.UINT64_0,
		PromiseNodeID:  &comm.UINT64_0,
		AcceptedID:     &comm.UINT64_0,
		AcceptedNodeID: &comm.UINT64_0,
	}
	return paxosLog.WriteState(options, groupIdx, instanceId, state)
}

func (paxosLog *PaxosLog) ReadLog(groupIdx int, instanceId uint64) ([]byte, error) {
	var state comm.AcceptorStateData
	err := paxosLog.ReadState(groupIdx, instanceId, &state)

	if err != nil {
		log.Error("Read log error ")
		return nil, err
	}

	value := state.AcceptedValue

	return value, nil
}

func (paxosLog *PaxosLog) GetMaxInstanceIDFromLog(groupIdx int) (uint64, error) {

}

func (paxosLog *PaxosLog) WriteState(options *WriteOptions, groupIdx int, instanceId uint64, state *comm.AcceptorStateData) error {
	value, err := proto.Marshal(state)
	if err != nil {
		return fmt.Errorf("marshal state error")
	}

	err = paxosLog.logStorage.Put(options, groupIdx, instanceId, value)

	if err != nil {
		log.Error("write state error")
		return err
	}
	return nil
}

func (paxosLog *PaxosLog) ReadState(groupIdx int, instanceId uint64, state *comm.AcceptorStateData) error {
	value, err := paxosLog.logStorage.Get(groupIdx, instanceId)

	if err != nil {
		return err
	}

	err = proto.Unmarshal(value, state)

	if err != nil {
		log.Error("Read State error caused by unmarshal error ")
	}

	return err
}

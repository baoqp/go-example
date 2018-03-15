package algorithm

import (
	"gphxpaxos/config"
	"sync"
)

// 统计reply数据
type MsgCounter struct {
	config *config.Config

	receiveMsgNodeIDMaps         map[uint64]bool
	rejectMsgNodeIDMaps          map[uint64]bool
	promiseOrAcceptMsgNodeIDMaps map[uint64]bool
	mutex sync.Mutex
}

func NewMsgCounter(config *config.Config) *MsgCounter {
	counter := &MsgCounter{
		config: config,
	}

	return counter
}

func (msgCounter *MsgCounter) StartNewRound() {
	msgCounter.receiveMsgNodeIDMaps = make(map[uint64]bool, 0)
	msgCounter.rejectMsgNodeIDMaps = make(map[uint64]bool, 0)
	msgCounter.promiseOrAcceptMsgNodeIDMaps = make(map[uint64]bool, 0)
}

func (msgCounter *MsgCounter) AddReceive(nodeId uint64) {
	msgCounter.mutex.Lock()
	msgCounter.receiveMsgNodeIDMaps[nodeId] = true
	msgCounter.mutex.Unlock()
}

func (msgCounter *MsgCounter) AddReject(nodeId uint64) {
	msgCounter.mutex.Lock()
	msgCounter.rejectMsgNodeIDMaps[nodeId] = true
	msgCounter.mutex.Unlock()
}

func (msgCounter *MsgCounter) AddPromiseOrAccept(nodeId uint64) {
	msgCounter.mutex.Lock()
	msgCounter.promiseOrAcceptMsgNodeIDMaps[nodeId] = true
	msgCounter.mutex.Unlock()
}

func (msgCounter *MsgCounter) IsPassedOnThisRound() bool {
	msgCounter.mutex.Lock()
	defer msgCounter.mutex.Unlock()
	return len(msgCounter.promiseOrAcceptMsgNodeIDMaps) >= msgCounter.config.GetMajorityCount()
}

func (msgCounter *MsgCounter) GetPassedCount() int {
	msgCounter.mutex.Lock()
	defer msgCounter.mutex.Unlock()
	return len(msgCounter.promiseOrAcceptMsgNodeIDMaps)
}

func (msgCounter *MsgCounter) IsRejectedOnThisRound() bool {
	msgCounter.mutex.Lock()
	defer msgCounter.mutex.Unlock()
	return len(msgCounter.rejectMsgNodeIDMaps) >= msgCounter.config.GetMajorityCount()
}

func (msgCounter *MsgCounter) IsAllReceiveOnThisRound() bool {
	msgCounter.mutex.Lock()
	defer msgCounter.mutex.Unlock()
	return len(msgCounter.receiveMsgNodeIDMaps) == msgCounter.config.GetNodeCount()
}

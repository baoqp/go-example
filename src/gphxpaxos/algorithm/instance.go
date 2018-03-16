package algorithm

import (
	log "github.com/sirupsen/logrus"
	"gphxpaxos/network"
	"sync"
	"gphxpaxos/config"
	"fmt"
	"time"
	"gphxpaxos/logstorage"
	"gphxpaxos/util"
	"gphxpaxos/node"
	"gphxpaxos/comm"
	"container/list"
	"gphxpaxos/checkpoint"
	"github.com/golang/protobuf/proto"
)

const (
	RETRY_QUEUE_MAX_LEN = 300
)

type CommitMsg struct {
}

type Instance struct {
	config     *config.Config
	logStorage logstorage.LogStorage
	paxosLog   *logstorage.PaxosLog
	committer  *Committer
	commitctx  *CommitContext
	proposer   *Proposer
	learner    *Learner
	acceptor   *Acceptor
	name       string
	factory    *node.SMFac

	transport network.MsgTransport

	timerThread *util.TimerThread

	endChan chan bool
	end     bool

	commitChan   chan CommitMsg
	paxosMsgChan chan *comm.PaxosMsg

	retryMsgList *list.List

	ckMnger      *checkpoint.CheckpointManager
	lastChecksum uint32
	mutex        sync.Mutex
}

func NewInstance(config *config.Config, logStorage logstorage.LogStorage, transport network.MsgTransport, useCkReplayer bool) (*Instance, error) {
	instance := &Instance{
		config:       config,
		logStorage:   logStorage,
		transport:    transport,
		paxosLog:     logstorage.NewPaxosLog(logStorage),
		factory:      node.NewSMFac(config.GetMyGroupId()),
		timerThread:  util.NewTimerThread(),
		endChan:      make(chan bool),
		commitChan:   make(chan CommitMsg),
		paxosMsgChan: make(chan *comm.PaxosMsg, 100),
		retryMsgList: list.New(),
	}

	instance.acceptor = NewAcceptor(instance)

	instance.ckMnger = checkpoint.NewCheckpointManager(config, instance.factory, logStorage, useCkReplayer)
	instance.ckMnger.Init()
	cpInstanceId := instance.ckMnger.GetCheckpointInstanceID() + 1

	log.Info("acceptor OK, log.instanceid %d checkpoint.instanceid %d", instance.acceptor.GetInstanceId(), cpInstanceId)
	nowInstanceId := cpInstanceId
	if nowInstanceId < instance.acceptor.GetInstanceId() {
		err := instance.PlayLog(nowInstanceId, instance.acceptor.GetInstanceId())
		if err != nil {
			return nil, err
		}
		nowInstanceId = instance.acceptor.GetInstanceId()
	} else {
		if nowInstanceId > instance.acceptor.GetInstanceId() {
			instance.acceptor.InitForNewPaxosInstance(false)
		}
		instance.acceptor.setInstanceId(nowInstanceId)
	}

	log.Info("now instance id: %d", nowInstanceId)

	instance.commitctx = newCommitContext(instance)
	instance.committer = newCommitter(instance)
	// learner must create before proposer
	instance.learner = NewLearner(instance)

	instance.proposer = NewProposer(instance)
	instance.proposer.setStartProposalID(instance.acceptor.GetAcceptorState().GetPromiseNum().proposalId + 1)

	instance.name = fmt.Sprintf("%s-%d", config.GetOptions().MyNode.String(), config.GetMyNodeId())

	maxInstanceId, err := loglogstorage.GetMaxInstanceID()
	log.Debug("max instance id:%d:%v， propose id:%d", maxInstanceId, err, instance.proposer.GetInstanceId())

	instance.ckMnger.SetMinChosenInstanceID(nowInstanceId)
	err = instance.InitLastCheckSum()
	if err != nil {
		return nil, err
	}
	instance.learner.Reset_AskforLearn_Noop(comm.GetAskforLearnInterval())

	instance.ckMnger.Start()

	util.StartRoutine(instance.main)

	return instance, nil
}


// instance main loop
func (instance *Instance) main() {
	end := false
	for !end {
		timer := time.NewTimer(100 * time.Millisecond)
		select {
		case <-instance.endChan:
			end = true
			break
		case <-instance.commitChan:
			instance.onCommit()
			break
		case msg := <-instance.paxosMsgChan:
			instance.OnReceivePaxosMsg(msg, false)
			break
		case <-timer.C:
			break
		}

		timer.Stop()
		instance.dealRetryMsg()
	}
}

func (instance *Instance) Stop() {
	instance.end = true
	instance.endChan <- true

	// instance.transport.Close()
	close(instance.paxosMsgChan)
	close(instance.commitChan)
	close(instance.endChan)
	instance.timerThread.Stop()
}

func (instance *Instance) Status(instanceId uint64) (Status, []byte) {
	if instanceId < instance.acceptor.GetInstanceId() {
		value, _, _ := instance.GetInstanceValue(instanceId)
		return Decided, value
	}

	return Pending, nil
}

func (instance *Instance) InitLastCheckSum() error {
	acceptor := instance.acceptor
	ckMnger := instance.ckMnger

	if acceptor.GetInstanceId() == 0 {
		instance.lastChecksum = 0
		return nil
	}

	if acceptor.GetInstanceId() <= ckMnger.GetMinChosenInstanceID() {
		instance.lastChecksum = 0
		return nil
	}

	state, err := instance.paxosLog.ReadState(acceptor.GetInstanceId() - 1)
	if err != nil && err != comm.ErrKeyNotFound {
		return err
	}

	if err == comm.ErrKeyNotFound {
		log.Error("last checksum not exist, now instance id %d", instance.acceptor.GetInstanceId())
		instance.lastChecksum = 0
		return nil
	}

	instance.lastChecksum = state.GetChecksum()
	log.Info("OK, last checksum %d", instance.lastChecksum)

	return nil
}

func (instance *Instance) PlayLog(beginInstanceId uint64, endInstanceId uint64) error {
	if beginInstanceId < instance.ckMnger.GetMinChosenInstanceID() {
		log.Error("now instanceid %d small than chosen instanceid %d", beginInstanceId, instance.ckMnger.GetMinChosenInstanceID())
		return comm.ErrInvalidInstanceId
	}

	for instanceId := beginInstanceId; instanceId < endInstanceId; instanceId++ {
		state, err := instance.paxosLog.ReadState(instanceId)
		if err != nil {
			log.Error("read instance %d log fail %v", instanceId, err)
			return err
		}

		err = instance.factory.Execute(instanceId, state.GetAcceptedValue(), nil)
		if err != nil {
			log.Error("execute instanceid %d fail:%v", instanceId, err)
			return err
		}
	}

	return nil
}

func (instance *Instance) NowInstanceId() uint64 {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()

	return instance.acceptor.GetInstanceId() - 1
}

// try to propose a value, return instanceid end error
func (instance *Instance) Propose(value []byte) (uint64, error) {
	log.Debug("[%s]try to propose value %s", instance.name, string(value))
	return instance.committer.NewValue(value)
}

func (instance *Instance) dealRetryMsg() {
	len := instance.retryMsgList.Len()
	hasRetry := false
	for i := 0; i < len; i++ {
		obj := instance.retryMsgList.Front()
		msg := obj.Value.(*comm.PaxosMsg)
		msgInstanceId := msg.GetInstanceID()
		nowInstanceId := instance.GetNowInstanceId()

		if msgInstanceId > nowInstanceId {
			break
		} else if msgInstanceId == nowInstanceId+1 {
			if hasRetry {
				instance.OnReceivePaxosMsg(msg, true)
				log.Debug("[%s]retry msg i+1 instanceid %d", msgInstanceId)
			} else {
				break
			}
		} else if msgInstanceId == nowInstanceId {
			instance.OnReceivePaxosMsg(msg, false)
			log.Debug("[%s]retry msg instanceid %d", msgInstanceId)
			hasRetry = true
		}

		instance.retryMsgList.Remove(obj)
	}
}

func (instance *Instance) addRetryMsg(msg *comm.PaxosMsg) {
	if instance.retryMsgList.Len() > RETRY_QUEUE_MAX_LEN {
		obj := instance.retryMsgList.Front()
		instance.retryMsgList.Remove(obj)
	}
	instance.retryMsgList.PushBack(msg)
}

func (instance *Instance) clearRetryMsg() {
	instance.retryMsgList = list.New()
}

func (instance *Instance) GetNowInstanceId() uint64 {
	return instance.acceptor.GetInstanceId()
}

func (instance *Instance) sendCommitMsg() {
	instance.commitChan <- CommitMsg{}
}

// handle commit message
func (instance *Instance) onCommit() {
	if !instance.commitctx.isNewCommit() {
		return
	}

	if !instance.learner.IsImLatest() {
		return
	}

	if instance.config.IsIMFollower() {
		log.Error("[%s]I'm follower, skip commit new value", instance.name)
		instance.commitctx.setResultOnlyRet(gpaxos.PaxosTryCommitRet_Follower_Cannot_Commit)
		return
	}

	commitValue := instance.commitctx.getCommitValue()
	if len(commitValue) > comm.GetMaxValueSize() {
		log.Error("[%s]value size %d to large, skip commit new value", instance.name, len(commitValue))
		instance.commitctx.setResultOnlyRet(gpaxos.PaxosTryCommitRet_Value_Size_TooLarge)
	}

	timeOutMs := instance.commitctx.StartCommit(instance.proposer.GetInstanceId())

	log.Debug("[%s]start commit instance %d, timeout:%d", instance.String(), instance.proposer.GetInstanceId(), timeOutMs)
	instance.proposer.NewValue(instance.commitctx.getCommitValue(), timeOutMs)
}

func (instance *Instance) String() string {
	return instance.name
}

func (instance *Instance) GetLastChecksum() uint32 {
	return 0
}

func (instance *Instance) GetInstanceValue(instanceId uint64) ([]byte, int32, error) {
	if instanceId >= instance.acceptor.GetInstanceId() {
		return nil, -1, gpaxos.Paxos_GetInstanceValue_Value_Not_Chosen_Yet
	}

	state, err := instance.paxosLog.ReadState(instanceId)
	if err != nil {
		return nil, -1, err
	}

	value, smid := instance.factory.UnpackPaxosValue(state.GetAcceptedValue())
	return value, smid, nil
}

func (instance *Instance) isCheckSumValid(msg *comm.PaxosMsg) bool {
	return true
}

func (instance *Instance) NewInstance(isMyCommit bool) {
	instance.acceptor.NewInstance(isMyCommit)
	instance.proposer.NewInstance(isMyCommit)
	instance.learner.NewInstance(isMyCommit)
}

func (instance *Instance) receiveMsgForLearner(msg *comm.PaxosMsg) error {
	log.Info("[%s]recv msg %d for learner", instance.name, msg.GetMsgType())
	learner := instance.learner
	msgType := msg.GetMsgType()

	switch msgType {
	case comm.MsgType_PaxosLearner_AskforLearn:
		learner.OnAskforLearn(msg)
		break
	case comm.MsgType_PaxosLearner_SendLearnValue:
		learner.OnSendLearnValue(msg)
		break
	case comm.MsgType_PaxosLearner_ProposerSendSuccess:
		learner.OnProposerSendSuccess(msg)
		break
	case comm.MsgType_PaxosLearner_SendNowInstanceID:
		learner.OnSendNowInstanceId(msg)
		break
	case comm.MsgType_PaxosLearner_ConfirmAskforLearn:
		learner.OnConfirmAskForLearn(msg)
		break
	case comm.MsgType_PaxosLearner_SendLearnValue_Ack:
		learner.OnSendLearnValue_Ack(msg)
		break
	case comm.MsgType_PaxosLearner_AskforCheckpoint:
		learner.OnAskforCheckpoint(msg)
		break
	}
	if learner.IsLearned() {
		commitCtx := instance.commitctx
		isMyCommit, _ := commitCtx.IsMyCommit(msg.GetNodeID(), learner.GetInstanceId(), learner.GetLearnValue())
		if isMyCommit {
			log.Debug("[%s]instance %d is my commit", instance.name, learner.GetInstanceId())
		} else {
			log.Debug("[%s]instance %d is not my commit", instance.name, learner.GetInstanceId())
		}

		commitCtx.setResult(gpaxos.PaxosTryCommitRet_OK, learner.GetInstanceId(), learner.GetLearnValue())

		instance.NewInstance(isMyCommit)

		log.Info("[%s]new paxos instance has started, Now instance id:proposer %d, acceptor %d, learner %d",
			instance.name, instance.proposer.GetInstanceId(), instance.acceptor.GetInstanceId(), instance.learner.GetInstanceId())
	}
	return nil
}

func (instance *Instance) receiveMsgForProposer(msg *comm.PaxosMsg) error {
	if instance.config.IsIMFollower() {
		log.Error("[%s]follower skip %d msg", instance.name, msg.GetMsgType())
		return nil
	}

	msgInstanceId := msg.GetInstanceID()
	proposerInstanceId := instance.proposer.GetInstanceId()

	if msgInstanceId != proposerInstanceId {
		log.Error("[%s]msg instance id %d not same to proposer instance id %d",
			instance.name, msgInstanceId, proposerInstanceId)
		return nil
	}

	msgType := msg.GetMsgType()
	if msgType == comm.MsgType_PaxosPrepareReply {
		return instance.proposer.OnPrepareReply(msg)
	} else if msgType == comm.MsgType_PaxosAcceptReply {
		return instance.proposer.OnAcceptReply(msg)
	}

	return comm.ErrInvalidMsg
}

// handle msg type which for acceptor
func (instance *Instance) receiveMsgForAcceptor(msg *comm.PaxosMsg, isRetry bool) error {
	if instance.config.IsIMFollower() {
		log.Error("[%s]follower skip %d msg", instance.name, msg.GetMsgType())
		return nil
	}

	msgInstanceId := msg.GetInstanceID()
	acceptorInstanceId := instance.acceptor.GetInstanceId()

	log.Info("[%s]msg instance %d, acceptor instance %d", instance.name, msgInstanceId, acceptorInstanceId)
	// msgInstanceId == acceptorInstanceId + 1  means acceptor instance has been approved
	// so just learn it
	if msgInstanceId == acceptorInstanceId+1 {
		newMsg := &comm.PaxosMsg{}
		util.CopyStruct(newMsg, *msg)
		newMsg.InstanceID = proto.Uint64(acceptorInstanceId)
		newMsg.MsgType = proto.Int(comm.MsgType_PaxosLearner_ProposerSendSuccess)
		log.Debug("learn it, node id: %d:%d", newMsg.GetNodeID(), msg.GetNodeID())
		instance.receiveMsgForLearner(newMsg)
	}

	msgType := msg.GetMsgType()

	// msg instance == acceptorInstanceId means this msg is what acceptor processing
	// so call the acceptor function to handle it
	if msgInstanceId == acceptorInstanceId {
		if msgType == comm.MsgType_PaxosPrepare {
			return instance.acceptor.onPrepare(msg)
		} else if msgType == comm.MsgType_PaxosAccept {
			return instance.acceptor.onAccept(msg)
		}

		// never reach here
		log.Error("wrong msg type %d", msgType)
		return comm.ErrInvalidMsg
	}

	// ignore retry msg
	if isRetry {
		log.Debug("ignore retry msg")
		return nil
	}

	// ignore expired msg
	if msgInstanceId <= acceptorInstanceId {
		log.Debug("[%s]ignore expired %d msg from %d, now %d", instance.name, msgInstanceId, msg.GetNodeID(), acceptorInstanceId)
		return nil
	}

	if msgInstanceId < instance.learner.getSeenLatestInstanceId() {
		log.Debug("ignore has learned msg")
		return nil
	}

	if msgInstanceId < acceptorInstanceId+RETRY_QUEUE_MAX_LEN {
		//need retry msg precondition
		//  1. prepare or accept msg
		//  2. msg.instanceid > nowinstanceid.
		//    (if < nowinstanceid, this msg is expire)
		//  3. msg.instanceid >= seen latestinstanceid.
		//    (if < seen latestinstanceid, proposer don't need reply with this instanceid anymore.)
		//  4. msg.instanceid close to nowinstanceid.
		instance.addRetryMsg(msg)
	} else {
		instance.clearRetryMsg()
	}
	return nil
}

func (instance *Instance) OnReceivePaxosMsg(msg *comm.PaxosMsg, isRetry bool) error {
	proposer := instance.proposer
	learner := instance.learner
	msgType := msg.GetMsgType()

	log.Info("[%s]instance id %d, msg instance id:%d, msgtype: %d, from: %d, my node id:%d, latest instanceid %d",
		instance.name, proposer.GetInstanceId(), msg.GetInstanceID(), msgType, msg.GetNodeID(),
		instance.config.GetMyNodeId(), learner.getSeenLatestInstanceId())

	// handle msg for acceptor
	if msgType == comm.MsgType_PaxosPrepare || msgType == comm.MsgType_PaxosAccept {
		if !instance.config.IsValidNodeID(msg.GetNodeID()) {
			instance.config.AddTmpNodeOnlyForLearn(msg.GetNodeID())
			log.Error("[%s]is not valid node id", instance.name)
			return nil
		}

		if !instance.isCheckSumValid(msg) {
			log.Error("[%s]checksum invalid", instance.name)
			return comm.ErrInvalidMsg
		}

		return instance.receiveMsgForAcceptor(msg, isRetry)
	}

	// handle paxos prepare and accept reply msg
	if (msgType == comm.MsgType_PaxosPrepareReply || msgType == comm.MsgType_PaxosAcceptReply) {
		return instance.receiveMsgForProposer(msg)
	}

	// handler msg for learner
	if (msgType == comm.MsgType_PaxosLearner_AskforLearn ||
		msgType == comm.MsgType_PaxosLearner_SendLearnValue ||
		msgType == comm.MsgType_PaxosLearner_ProposerSendSuccess ||
		msgType == comm.MsgType_PaxosLearner_ConfirmAskforLearn ||
		msgType == comm.MsgType_PaxosLearner_SendNowInstanceID ||
		msgType == comm.MsgType_PaxosLearner_SendLearnValue_Ack ||
		msgType == comm.MsgType_PaxosLearner_AskforCheckpoint) {
		if !instance.isCheckSumValid(msg) {
			return comm.ErrInvalidMsg
		}

		return instance.receiveMsgForLearner(msg)
	}

	log.Error("invalid msg %d", msgType)
	return comm.ErrInvalidMsg
}

func (instance *Instance) OnTimeout(timer *util.Timer) {
	if timer.TimerType == PrepareTimer {
		instance.proposer.onPrepareTimeout()
		return
	}

	if timer.TimerType == AcceptTimer {
		instance.proposer.onAcceptTimeout()
		return
	}

	if timer.TimerType == LearnerTimer {
		instance.learner.AskforLearn_Noop()
		return
	}
}

func (instance *Instance) OnReceiveMsg(buffer []byte, cmd int32) error {
	if instance.end {
		return nil
	}
	if cmd == comm.MsgCmd_PaxosMsg {
		var msg comm.PaxosMsg
		err := proto.Unmarshal(buffer, &msg)
		if err != nil {
			log.Error("[%s]unmarshal msg error %v", instance.name, err)
			return err
		}
		instance.paxosMsgChan <- &msg
	}

	return nil
}

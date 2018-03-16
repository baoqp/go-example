package algorithm

import (
	log "github.com/sirupsen/logrus"
	"gphxpaxos/config"
	"gphxpaxos/comm"
	"gphxpaxos/logstorage"
	"github.com/golang/protobuf/proto"
	"gphxpaxos/util"
	"gphxpaxos/checkpoint"
)

//-------------------------------------------------LearnerState-----------------------------------------//

type LearnerState struct {
	config       *config.Config
	learnedValue []byte
	isLearned    bool
	newChecksum  uint32
	paxosLog     *logstorage.PaxosLog
}

func NewLearnerState(instance *Instance) *LearnerState {
	state := &LearnerState{
		config:   instance.config,
		paxosLog: instance.paxosLog,
	}

	state.Init()

	return state
}

func (learnerState *LearnerState) GetNewChecksum() uint32 {
	return learnerState.newChecksum
}

func (learnerState *LearnerState) LearnValueWithoutWrite(instanceId uint64, value []byte, checksum uint32) {
	learnerState.learnedValue = value
	learnerState.isLearned = true
	learnerState.newChecksum = checksum
}

func (learnerState *LearnerState) LearnValue(instanceId uint64, learnedBallot BallotNumber, value []byte, lastChecksum uint32) error {
	if instanceId == 0 && lastChecksum == 0 {
		learnerState.newChecksum = 0
	} else if len(value) > 0 {
		learnerState.newChecksum = util.Crc32(lastChecksum, value, comm.CRC32_SKIP)
	}

	state := comm.AcceptorStateData{
		InstanceID:     proto.Uint64(instanceId),
		AcceptedValue:  proto.NewBuffer(value).Bytes(),
		PromiseID:      proto.Uint64(learnedBallot.proposalId),
		PromiseNodeID:  proto.Uint64(learnedBallot.nodeId),
		AcceptedID:     proto.Uint64(learnedBallot.proposalId),
		AcceptedNodeID: proto.Uint64(learnedBallot.nodeId),
		Checksum:       proto.Uint32(learnerState.newChecksum),
	}

	options := logstorage.WriteOptions{
		Sync: false,
	}

	err := learnerState.paxosLog.WriteState(&options, learnerState.config.GetMyGroupId(), instanceId, &state)
	if err != nil {
		log.Error("storage writestate fail, instanceid %d valuelen %d err %v",
			instanceId, len(value), err)
		return err
	}

	learnerState.LearnValueWithoutWrite(instanceId, value, learnerState.newChecksum)

	return nil
}

func (learnerState *LearnerState) GetLearnValue() []byte {
	return learnerState.learnedValue
}

func (learnerState *LearnerState) IsLearned() bool {
	return learnerState.isLearned
}

func (learnerState *LearnerState) Init() {
	learnerState.learnedValue = nil
	learnerState.isLearned = false
	learnerState.newChecksum = 0
}

//-------------------------------------------------Learner-----------------------------------------//
type Learner struct {
	Base
	instance                         *Instance
	paxosLog                         *logstorage.PaxosLog
	acceptor                         *Acceptor
	state                            *LearnerState
	isImLearning                     bool
	highestSeenInstanceID            uint64
	highestSeenInstanceID_fromNodeID uint64
	lastAckInstanceId                uint64
	askforlearnNoopTimerID           uint32
	sender                           *LearnerSender
	timerThread                      *util.TimerThread
	ckReceiver                       *CheckpointReceiver
	ckSender                         *CheckpointSender
	ckMnger                          *checkpoint.CheckpointManager
	factory                          *statemachine.StatemachineFactory
}

func NewLearner(instance *Instance) *Learner {
	learner := &Learner{
		Base:                             newBase(instance),
		paxosLog:                         instance.paxosLog,
		acceptor:                         instance.acceptor,
		isImLearning:                     false,
		highestSeenInstanceID:            0,
		highestSeenInstanceID_fromNodeID: comm.NULL_NODEID,
		lastAckInstanceId:                0,
		state:                            NewLearnerState(instance),
		instance:                         instance,
		timerThread:                      instance.timerThread,
		ckReceiver:                       NewCheckpointReceiver(instance.config, instance.logStorage),
		ckMnger:                          instance.ckMnger,
		factory:                          instance.factory,
	}
	learner.sender = NewLearnerSender(instance, learner)

	learner.InitForNewPaxosInstance(false)

	return learner
}

func (self *Learner) InitForNewPaxosInstance(isMyCommit bool) {
	self.state.Init()
}

func (self *Learner) NewInstance(isMyComit bool) {
	self.Base.newInstance()
	self.InitForNewPaxosInstance(isMyComit)
	log.Debug("[%s]now learner instance id %d", self.instance.String(), self.Base.GetInstanceId())
}

func (self *Learner) Init() {
	self.sender.Start()
}

func (self *Learner) IsLearned() bool {
	return self.state.IsLearned()
}

func (self *Learner) GetLearnValue() []byte {
	return self.state.GetLearnValue()
}

func (self *Learner) Stop() {
	self.sender.Stop()
}

func (self *Learner) IsImLatest() bool {
	return self.GetInstanceId()+1 >= self.highestSeenInstanceID
}

func (self *Learner) GetSeenLatestInstanceID() uint64 {
	return self.highestSeenInstanceID
}

func (self *Learner) SetSeenInstanceID(instanceId uint64, fromNodeId uint64) {
	if instanceId > self.highestSeenInstanceID {
		self.highestSeenInstanceID = instanceId
		self.highestSeenInstanceID_fromNodeID = fromNodeId
	}
}

func (self *Learner) GetNewChecksum() uint32 {
	return self.state.GetNewChecksum()
}

func (self *Learner) Reset_AskforLearn_Noop(timeout uint32) {
	if self.askforlearnNoopTimerID > 0 {
		self.timerThread.DelTimer(self.askforlearnNoopTimerID)
	}

	self.askforlearnNoopTimerID = self.timerThread.AddTimer(timeout, LearnerTimer, self.instance)
}

func (self *Learner) AskforLearn_Noop() {
	self.Reset_AskforLearn_Noop(comm.GetAskforLearnInterval())
	self.isImLearning = false
	self.askforLearn()
}

func (self *Learner) askforLearn() {
	log.Info("start learn")

	base := self.Base

	msg := &comm.PaxosMsg{
		InstanceID: proto.Uint64(self.GetInstanceId()),
		NodeID:     proto.Uint64(self.config.GetMyNodeId()),
		MsgType:    proto.Int32(comm.MsgType_PaxosLearner_AskforLearn),
	}

	if self.config.IsIMFollower() {
		msg.ProposalNodeID = proto.Uint64(self.config.GetFollowToNodeID())
	}

	log.Info("end instanceid %d, mynodeid %d", msg.GetInstanceID(), msg.GetNodeID())

	base.broadcastMessage(msg, BroadcastMessage_Type_RunSelf_None)
	//self.BroadcastMessageToTempNode(msg, comm.Message_SendType_UDP)
}

func (self *Learner) OnAskforLearn(msg *comm.PaxosMsg) {
	log.Info("start msg.instanceid %d now.instanceid %d msg.fromnodeid %d",
		msg.GetInstanceID(), self.GetInstanceId(), msg.GetNodeID())

	self.SetSeenInstanceID(msg.GetInstanceID(), msg.GetNodeID())
	if msg.GetProposalNodeID() == self.config.GetMyNodeId() {
		log.Info("found a node %d follow me", msg.GetNodeID())
		self.config.AddFollowerNode(msg.GetNodeID())
	}

	if msg.GetInstanceID() >= self.GetInstanceId() {
		return
	}

	/*
	if msg.GetInstanceID() >= self.CpMng.GetMinChosenInstanceID() {
	  if !self.LearnerSender.Prepare(msg.GetInstanceID(), msg.GetNodeID()) {
		log.Error("learner sender working for others")

		if msg.GetInstanceID() == self.GetInstanceId()-1 {
		  log.Info("instanceid only difference one, just send this value to other")
		  var state comm.AcceptorStateData
		  err := self.PaxosLog.ReadState(self.config.GetMyGroupIdx(), msg.GetInstanceID(), &state)
		  if err == nil {
			ballot := newBallotNumber(state.GetAcceptedID(), state.GetAcceptedNodeID())
			self.SendLearnValue(msg.GetNodeID(), msg.GetInstanceID(), ballot, state.GetAcceptedValue(), 0, false)
		  }
		}

		return
	  }
	}
	*/

	self.sendNowInstanceID(msg.GetInstanceID(), msg.GetNodeID())
}

func (self *Learner) sendNowInstanceID(instanceId uint64, sendNodeId uint64) {
	msg := &comm.PaxosMsg{
		InstanceID:    proto.Uint64(instanceId),
		NodeID:        proto.Uint64(self.config.GetMyNodeId()),
		MsgType:       proto.Int32(comm.MsgType_PaxosLearner_SendNowInstanceID),
		NowInstanceID: proto.Uint64(self.GetInstanceId()),
		//MinChosenInstanceID:
	}

	if self.GetInstanceId()-instanceId > 50 {
		/*
		systemVarBuffer, err := self.config.GetSystemVSM().GetCheckpointBuffer()
		if err == nil {
		  msg.SystemVariables = util.CopyBytes(systemVarBuffer)
		}

		masterVarBuffer, err := self.config.GetMasterSM().GetCheckpointBuffer()
		if err == nil {
		  msg.MasterVariables = util.CopyBytes(masterVarBuffer)
		}
		*/
	}

	self.sendPaxosMessage(sendNodeId, msg)
}

func (self *Learner) OnSendNowInstanceId(msg *comm.PaxosMsg) {
	instance := self.instance
	instanceId := self.instanceId

	log.Info("[%s]start msg.instanceid %d now.instanceid %d msg.from_nodeid %d msg.maxinstanceid %d",
		instance.String(), msg.GetInstanceID(), self.instanceId, msg.GetNodeID(), msg.GetNowInstanceID())

	self.SetSeenInstanceID(msg.GetNowInstanceID(), msg.GetNodeID())

	if msg.GetInstanceID() != instanceId {
		log.Error("[%s]lag msg instanceid %d", instance.String(), msg.GetInstanceID())
		return
	}

	if msg.GetNowInstanceID() <= instanceId {
		log.Error("[%s]lag msg instanceid %d", instance.String(), msg.GetNowInstanceID())
		return
	}

	if msg.GetMinChosenInstanceID() > instanceId {

	} else if (!self.isImLearning) {
		self.confirmAskForLearn(msg.GetNodeID())
	}
}

func (self *Learner) confirmAskForLearn(sendNodeId uint64) {
	msg := &comm.PaxosMsg{
		InstanceID: proto.Uint64(self.instanceId),
		NodeID:     proto.Uint64(self.config.GetMyNodeId()),
		MsgType:    proto.Int32(comm.MsgType_PaxosLearner_ConfirmAskforLearn),
	}
	self.sendPaxosMessage(sendNodeId, msg)
	self.isImLearning = true
}

func (self *Learner) OnConfirmAskForLearn(msg *comm.PaxosMsg) {
	log.Info("start msg.instanceid %d msg.from nodeid %d", msg.GetInstanceID(), msg.GetNodeID())

	if !self.sender.Confirm(msg.GetInstanceID(), msg.GetNodeID()) {
		log.Error("learner sender confirm fail,maybe is lag msg")
		return
	}

	log.Info("ok, success confirm")
}

func (self *Learner) OnAskforCheckpoint(msg *comm.PaxosMsg) {

}

func (self *Learner) SendLearnValue(sendNodeId uint64, learnInstanceId uint64,
	ballot BallotNumber, value []byte, cksum uint32, needAck bool) error {
	var paxosMsg = &comm.PaxosMsg{
		MsgType:        proto.Int32(comm.MsgType_PaxosLearner_SendLearnValue),
		InstanceID:     proto.Uint64(learnInstanceId),
		NodeID:         proto.Uint64(self.config.GetMyNodeId()),
		ProposalNodeID: proto.Uint64(ballot.nodeId),
		ProposalID:     proto.Uint64(ballot.proposalId),
		Value:          value,
		LastChecksum:   proto.Uint32(cksum),
	}

	if needAck {
		paxosMsg.Flag = proto.Uint32(comm.PaxosMsgFlagType_SendLearnValue_NeedAck)
	}

	return self.sendPaxosMessage(sendNodeId, paxosMsg)
}

func (self *Learner) OnSendLearnValue(msg *comm.PaxosMsg) {
	log.Info("START Msg.InstanceID %d Now.InstanceID %d Msg.ballot_proposalid %d Msg.ballot_nodeid %d Msg.ValueSize %d",
		msg.GetInstanceID(), self.GetInstanceId(), msg.GetProposalID(),
		msg.GetNodeID(), len(msg.Value))

	if msg.GetInstanceID() > self.GetInstanceId() {
		log.Debug("[Latest Msg] i can't learn")
		return
	}

	if msg.GetInstanceID() < self.GetInstanceId() {
		log.Debug("[Lag Msg] no need to learn")
	} else {
		ballot := NewBallotNumber(msg.GetProposalID(), msg.GetProposalNodeID())
		err := self.state.LearnValue(msg.GetInstanceID(), *ballot, msg.GetValue(), self.GetLastChecksum())
		if err != nil {
			log.Error("LearnState.LearnValue fail:%v", err)
			return
		}
		log.Info("END LearnValue OK, proposalid %d proposalid_nodeid %d valueLen %d",
			msg.GetProposalID(), msg.GetNodeID(), len(msg.Value))
	}

	if msg.GetFlag() == comm.PaxosMsgFlagType_SendLearnValue_NeedAck {
		self.Reset_AskforLearn_Noop(comm.GetAskforLearnInterval())
		self.SendLearnValue_Ack(msg.GetNodeID())
	}
}

func (self *Learner) SendLearnValue_Ack(sendNodeId uint64) {
	log.Info("START LastAck.Instanceid %d Now.Instanceid %d", self.lastAckInstanceId, self.GetInstanceId())

	if self.GetInstanceId() < self.lastAckInstanceId+comm.GetLeanerReceiver_Ack_Lead() {
		log.Info("no need ack")
		return
	}

	self.lastAckInstanceId = self.GetInstanceId()

	msg := &comm.PaxosMsg{
		InstanceID: proto.Uint64(self.GetInstanceId()),
		MsgType:    proto.Int32(comm.MsgType_PaxosLearner_SendLearnValue_Ack),
		NodeID:     proto.Uint64(self.config.GetMyNodeId()),
	}

	self.sendPaxosMessage(sendNodeId, msg)

	log.Info("END.OK")
}

func (self *Learner) OnSendLearnValue_Ack(msg *comm.PaxosMsg) {
	log.Info("Msg.Ack.Instanceid %d Msg.from_nodeid %d", msg.GetInstanceID(), msg.GetNodeID())
	self.sender.Ack(msg.GetInstanceID(), msg.GetNodeID())
}

func (self *Learner) getSeenLatestInstanceId() uint64 {
	return self.highestSeenInstanceID
}

func (self *Learner) ProposerSendSuccess(instanceId uint64, proposalId uint64) {
	msg := &comm.PaxosMsg{
		MsgType:      proto.Int32(comm.MsgType_PaxosLearner_ProposerSendSuccess),
		InstanceID:   proto.Uint64(instanceId),
		NodeID:       proto.Uint64(self.config.GetMyNodeId()),
		ProposalID:   proto.Uint64(proposalId),
		LastChecksum: proto.Uint32(self.GetLastChecksum()),
	}

	self.broadcastMessage(msg, BroadcastMessage_Type_RunSelf_First)
}

func (self *Learner) OnProposerSendSuccess(msg *comm.PaxosMsg) {
	log.Info("[%s]OnProposerSendSuccess Msg.InstanceID %d Now.InstanceID %d Msg.ProposalID %d "+
		"State.AcceptedID %d State.AcceptedNodeID %d, Msg.from_nodeid %d",
		self.instance.String(), msg.GetInstanceID(), self.GetInstanceId(), msg.GetProposalID(),
		self.acceptor.GetAcceptorState().acceptedNum.proposalId,
		self.acceptor.GetAcceptorState().acceptedNum.nodeId,
		msg.GetNodeID())

	if msg.GetInstanceID() != self.GetInstanceId() {
		log.Debug("instance id %d not same as msg instance id %d", self.GetInstanceId(), msg.GetInstanceID())
		return
	}

	if self.acceptor.GetAcceptorState().acceptedNum.IsNull() {
		log.Debug("not accepted any proposal")
		return
	}

	ballot := NewBallotNumber(msg.GetProposalID(), msg.GetNodeID())
	if !self.acceptor.GetAcceptorState().acceptedNum.EQ(ballot) {
		log.Debug("[%s]proposal ballot %s not same to accepted ballot %s", self.instance.String(),
			self.acceptor.GetAcceptorState().acceptedNum.String(), ballot.String())
		return
	}

	self.state.LearnValueWithoutWrite(msg.GetInstanceID(),
		self.acceptor.GetAcceptorState().GetAcceptedValue(),
		self.acceptor.GetAcceptorState().GetChecksum())

	log.Info("learn value instanceid %d ok", msg.GetInstanceID())
	//self.TransmitToFollower()
}

func (self *Learner) SendCheckpointBegin(sendNodeId uint64, uuid uint64,
	sequence uint64, ckInstanceId uint64) error {
	ckMsg := &comm.CheckpointMsg{
		MsgType:              proto.Int32(comm.CheckpointMsgType_SendFile),
		NodeID:               proto.Uint64(self.config.GetMyNodeId()),
		Flag:                 proto.Int32(comm.CheckpointSendFileFlag_BEGIN),
		UUID:                 proto.Uint64(uuid),
		Sequence:             proto.Uint64(sequence),
		CheckpointInstanceID: proto.Uint64(ckInstanceId),
	}

	return self.sendCheckpointMessage(sendNodeId, ckMsg)
}

func (self *Learner) SendCheckpointEnd(sendNodeId uint64, uuid uint64,
	sequence uint64, ckInstanceId uint64) error {
	ckMsg := &comm.CheckpointMsg{
		MsgType:              proto.Int32(comm.CheckpointMsgType_SendFile),
		NodeID:               proto.Uint64(self.config.GetMyNodeId()),
		Flag:                 proto.Int32(comm.CheckpointSendFileFlag_END),
		UUID:                 proto.Uint64(uuid),
		Sequence:             proto.Uint64(sequence),
		CheckpointInstanceID: proto.Uint64(ckInstanceId),
	}

	return self.sendCheckpointMessage(sendNodeId, ckMsg)
}

func (self *Learner) SendCheckpoint(sendNodeId uint64, uuid uint64,
	sequence uint64, ckInstanceId uint64, ckssum uint32,
	filePath string, smid int32, offset uint64, buffer []byte) error {
	ckMsg := &comm.CheckpointMsg{
		MsgType:              proto.Int32(comm.CheckpointMsgType_SendFile),
		NodeID:               proto.Uint64(self.config.GetMyNodeId()),
		Flag:                 proto.Int32(comm.CheckpointSendFileFlag_ING),
		UUID:                 proto.Uint64(uuid),
		Sequence:             proto.Uint64(sequence),
		CheckpointInstanceID: proto.Uint64(ckInstanceId),
		Checksum:             proto.Uint32(ckssum),
		FilePath:             proto.String(filePath),
		SMID:                 proto.Int(int(smid)),
		Offset:               proto.Uint64(offset),
		Buffer:               buffer,
	}

	return self.sendCheckpointMessage(sendNodeId, ckMsg)
}

func (self *Learner) OnSendCheckpointBegin(ckMsg *comm.CheckpointMsg) error {
	err := self.ckReceiver.NewReceiver(ckMsg.GetNodeID(), ckMsg.GetUUID())
	if err != nil {
		return err
	}

	err = self.ckMnger.SetMinChosenInstanceID(ckMsg.GetCheckpointInstanceID())
	if err != nil {
		return err
	}

	return nil
}

func (self *Learner) OnSendCheckpointIng(ckMsg *comm.CheckpointMsg) error {
	return self.ckReceiver.ReceiveCheckpoint(ckMsg)
}

func (self *Learner) OnSendCheckpointEnd(ckMsg *comm.CheckpointMsg) error {
	if !self.ckReceiver.IsReceiverFinish(ckMsg.GetNodeID(), ckMsg.GetUUID(), ckMsg.GetSequence()) {
		log.Error("receive end msg but receiver not finish")
		return comm.ErrInvalidMsg
	}

	smList := self.factory.GetSMList()
	for _, sm := range smList {
		smid := sm.SMID()
		if smid == gpaxos.SYSTEM_V_SMID || smid == gpaxos.MASTER_V_SMID {
			continue
		}

		//tmpDirPath := self.ckReceiver.GetTmpDirPath(smid)

	}

	return nil
}

func (self *Learner) OnSendCheckpoint(ckMsg *comm.CheckpointMsg) {
	log.Info("start uuid %d flag %d sequence %d cpi %d checksum %d smid %d offset %d filepath %s",
		ckMsg.GetUUID(), ckMsg.GetFlag(), ckMsg.GetSequence(), ckMsg.GetCheckpointInstanceID(),
		ckMsg.GetChecksum(), ckMsg.GetSMID(), ckMsg.GetOffset(), ckMsg.GetFilePath())

	var err error
	switch ckMsg.GetFlag() {
	case comm.CheckpointSendFileFlag_BEGIN:
		err = self.OnSendCheckpointBegin(ckMsg)
		break
	case comm.CheckpointSendFileFlag_ING:
		err = self.OnSendCheckpointIng(ckMsg)
		break
	case comm.CheckpointSendFileFlag_END:
		err = self.OnSendCheckpointEnd(ckMsg)
		break
	}

	if err != nil {
		log.Error("[FATAL]rest checkpoint receiver and reset askforlearn")
		self.ckReceiver.Reset()
		self.Reset_AskforLearn_Noop(5000)
		self.SendCheckpointAck(ckMsg.GetNodeID(), ckMsg.GetUUID(), ckMsg.GetSequence(), comm.CheckpointSendFileAckFlag_Fail)
	} else {
		self.SendCheckpointAck(ckMsg.GetNodeID(), ckMsg.GetUUID(), ckMsg.GetSequence(), comm.CheckpointSendFileAckFlag_OK)
		self.Reset_AskforLearn_Noop(10000)
	}
}

func (self *Learner) SendCheckpointAck(sendNodeId uint64, uuid uint64, sequence uint64, flag int) error {
	ckMsg := &comm.CheckpointMsg{
		MsgType:  proto.Int32(comm.CheckpointMsgType_SendFile_Ack),
		NodeID:   proto.Uint64(self.config.GetMyNodeId()),
		UUID:     proto.Uint64(uuid),
		Sequence: proto.Uint64(sequence),
		Flag:     proto.Int32(comm.CheckpointSendFileFlag_ING),
	}

	return self.sendCheckpointMessage(sendNodeId, ckMsg)
}

func (self *Learner) OnSendCheckpointAck(ckMsg *comm.CheckpointMsg) {
	log.Info("START flag %d", ckMsg.GetFlag())

	if self.ckSender != nil && !self.ckSender.IsEnd() {
		if ckMsg.GetFlag() == comm.CheckpointSendFileAckFlag_OK {
			self.ckSender.Ack(ckMsg.GetNodeID(), ckMsg.GetUUID(), ckMsg.GetSequence())
		} else {
			self.ckSender.End()
		}
	}
}

/*
func (self *Learner) TransmitToFollower() {
  if self.config.GetMyFollowerCount() == 0 {
    return
  }

  acceptor := self.acceptor
  msg := &comm.PaxosMsg{
    MsgType:        proto.Int32(comm.MsgType_PaxosLearner_SendLearnValue),
    InstanceID:     proto.Uint64(self.GetInstanceId()),
    NodeID:         proto.Uint64(self.config.GetMyNodeId()),
    ProposalNodeID: proto.Uint64(acceptor.GetAcceptorState().acceptedNum.nodeId),
    ProposalID:     proto.Uint64(acceptor.GetAcceptorState().acceptedNum.proposalId),
    Value:          acceptor.GetAcceptorState().GetAcceptedValue(),
    LastChecksum:   proto.Uint32(self.GetLastChecksum()),
  }

  self.broadcastMessageToFollower(msg)

  log.Info("OK")
}
*/

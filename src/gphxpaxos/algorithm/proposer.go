package algorithm

import (
	"gphxpaxos/config"
	log "github.com/sirupsen/logrus"
	"gphxpaxos/util"
	"gphxpaxos/comm"
	"github.com/golang/protobuf/proto"
	"math/rand"
)

const (
	PAUSE   = iota
	PREPARE
	ACCEPT
)

//-----------------------------------------------ProposerState-------------------------------------------//
type ProposerState struct {
	config *config.Config

	value      []byte
	proposalId uint64

	// save the highest other propose id,
	// next propose id = max(proposalId, highestOtherProposalId) + 1
	highestOtherProposalId uint64

	// save pre-accept ballot number
	highestOtherPreAcceptBallot BallotNumber

	state int
}

func newProposalState(config *config.Config) *ProposerState {
	proposalState := new(ProposerState)
	proposalState.config = config
	proposalState.proposalId = 1

	return proposalState.init()
}

func (proposerState *ProposerState) init() *ProposerState {
	proposerState.highestOtherProposalId = 0
	proposerState.value = nil
	proposerState.state = PAUSE

	return proposerState
}

func (proposerState *ProposerState) getState() int {
	return proposerState.state
}

func (proposerState *ProposerState) setState(state int) {
	proposerState.state = state
}

func (proposerState *ProposerState) setStartProposalId(proposalId uint64) {
	proposerState.proposalId = proposalId
}

// 更新proposalId
func (proposerState *ProposerState) newPrepare() {
	log.Info("start proposalid %d highestother %d mynodeid %d",
		proposerState.proposalId, proposerState.highestOtherProposalId, proposerState.config.GetMyNodeId())

	// next propose id = max(proposalId, highestOtherProposalId) + 1
	maxProposalId := proposerState.highestOtherProposalId
	if proposerState.proposalId > proposerState.highestOtherProposalId {
		maxProposalId = proposerState.proposalId
	}

	proposerState.proposalId = maxProposalId + 1

	log.Info("end proposalid %d", proposerState.proposalId)
}

func (proposerState *ProposerState) AddPreAcceptValue(otherPreAcceptBallot BallotNumber, otherPreAcceptValue []byte) {

	if otherPreAcceptBallot.IsNull() {
		return
	}

	// update value only when the ballot >  highestOtherPreAcceptBallot
	if otherPreAcceptBallot.GT(&proposerState.highestOtherPreAcceptBallot) {
		proposerState.highestOtherPreAcceptBallot = otherPreAcceptBallot
		proposerState.value = util.CopyBytes(otherPreAcceptValue)
	}
}

func (proposerState *ProposerState) GetProposalId() uint64 {
	return proposerState.proposalId
}

func (proposerState *ProposerState) GetValue() []byte {
	return proposerState.value
}

func (proposerState *ProposerState) SetValue(value []byte) {
	proposerState.value = util.CopyBytes(value)
}

func (proposerState *ProposerState) SetOtherProposalId(otherProposalId uint64) {
	if otherProposalId > proposerState.highestOtherProposalId {
		proposerState.highestOtherProposalId = otherProposalId
	}
}

func (proposerState *ProposerState) ResetHighestOtherPreAcceptBallot() {
	proposerState.highestOtherPreAcceptBallot.Reset()
}


//-------------------------------------------Proposer---------------------------------------------//

type Proposer struct {
	Base

	config               *config.Config
	state                *ProposerState
	msgCounter           *MsgCounter
	learner              *Learner
	preparing            bool
	prepareTimerId       uint32
	acceptTimerId        uint32
	lastPrepareTimeoutMs uint32
	lastAcceptTimeoutMs  uint32
	canSkipPrepare       bool
	wasRejectBySomeone   bool
	timerThread          *util.TimerThread
	timeOutMs 			uint32
	lastStartTimeMs   		 uint64
}

func NewProposer(instance *Instance) *Proposer {
	proposer := &Proposer{
		Base:        newBase(instance),
		config:      instance.config,
		state:       newProposalState(instance.config),
		msgCounter:  NewMsgCounter(instance.config),
		learner:     instance.learner,
		timerThread: instance.timerThread,
	}

	proposer.InitForNewPaxosInstance(false)

	return proposer
}

func (proposer *Proposer) InitForNewPaxosInstance(isMyCommit bool) {
	if !isMyCommit {
		return
	}
	proposer.msgCounter.StartNewRound()
	proposer.state.init()

	proposer.exitPrepare()
	proposer.exitAccept()
}

func (proposer *Proposer) NewInstance(isMyComit bool) {
	proposer.Base.newInstance()
	proposer.InitForNewPaxosInstance(isMyComit)
}

func (proposer *Proposer) setStartProposalID(proposalId uint64) {
	proposer.state.setStartProposalId(proposalId)
}

func (proposer *Proposer) isWorking() bool {
	return proposer.prepareTimerId > 0 || proposer.acceptTimerId > 0
}

func (proposer *Proposer) NewValue(value []byte, timeOutMs uint32) {
	if len(proposer.state.GetValue()) == 0 {
		proposer.state.SetValue(value)
	}

	proposer.lastPrepareTimeoutMs = comm.GetInsideOptions().GetStartPrepareTimeoutMs()
	proposer.lastAcceptTimeoutMs = comm.GetInsideOptions().GetStartAcceptTimeoutMs()
	proposer.timeOutMs = timeOutMs
	proposer.lastStartTimeMs = util.NowTimeMs()

	if proposer.canSkipPrepare && !proposer.wasRejectBySomeone {
		log.Info("skip prepare,directly start accept")
		proposer.accept()
	} else {
		proposer.prepare(proposer.wasRejectBySomeone)
	}
}

func (proposer *Proposer) isTimeout() bool {
	now := util.NowTimeMs()
	diff := now - proposer.lastStartTimeMs
	log.Debug("[%s]diff %d, timeout %d",proposer.instance.String(), diff, proposer.timeOutMs)
	if uint32(diff) >= proposer.timeOutMs {
		proposer.timeOutMs = 0
	}
	if proposer.timeOutMs <= 0 {
		log.Debug("[%s]instance %d timeout", proposer.instance.String(), proposer.instanceId)
		proposer.instance.commitctx.setResult(gpaxos.PaxosTryCommitRet_Timeout, proposer.instanceId, []byte(""))
		return true
	}
	proposer.timeOutMs -= uint32(diff)

	proposer.lastStartTimeMs = now
	return false
}

func (proposer *Proposer) prepare(needNewBallot bool) {
	if proposer.isTimeout() {
		return
	}
	base := proposer.Base
	state := proposer.state

	// first reset all state
	proposer.exitAccept()
	proposer.state.setState(PREPARE)
	proposer.canSkipPrepare = false
	proposer.wasRejectBySomeone = false
	proposer.state.ResetHighestOtherPreAcceptBallot()

	if needNewBallot {
		proposer.state.newPrepare()
	}

	log.Info("[%s]start prepare now.instanceid %d mynodeid %d state.proposal id %d state.valuelen %d new %v",
		proposer.instance.String(),proposer.GetInstanceId(), proposer.config.GetMyNodeId(), state.GetProposalId(), len(state.GetValue()), needNewBallot)

	// pack paxos prepare msg and broadcast
	msg := &comm.PaxosMsg{
		MsgType:    proto.Int32(comm.MsgType_PaxosPrepare),
		InstanceID: proto.Uint64(base.GetInstanceId()),
		NodeID:     proto.Uint64(proposer.config.GetMyNodeId()),
		ProposalID: proto.Uint64(state.GetProposalId()),
	}

	proposer.msgCounter.StartNewRound()
	proposer.addPrepareTimer(proposer.lastPrepareTimeoutMs)

	base.broadcastMessage(msg, BroadcastMessage_Type_RunSelf_First)
}

func (proposer *Proposer) exitAccept() {
	if proposer.acceptTimerId != 0 {
		proposer.timerThread.DelTimer(proposer.acceptTimerId)
		proposer.acceptTimerId = 0
	}
}

func (proposer *Proposer) exitPrepare() {
	if proposer.prepareTimerId != 0 {
		proposer.timerThread.DelTimer(proposer.prepareTimerId)
		proposer.prepareTimerId = 0
	}
}

func (proposer *Proposer) addPrepareTimer(timeOutMs uint32) {
	if proposer.prepareTimerId != 0 {
		proposer.timerThread.DelTimer(proposer.prepareTimerId)
		proposer.prepareTimerId = 0
	}

	if timeOutMs > proposer.timeOutMs {
		timeOutMs = uint32(proposer.timeOutMs)
	}

	proposer.prepareTimerId = proposer.timerThread.AddTimer(timeOutMs, comm.Timer_Proposer_Prepare_Timeout, proposer.instance)
	proposer.lastPrepareTimeoutMs *= 2
	if proposer.lastPrepareTimeoutMs > comm.GetInsideOptions().GetMaxPrepareTimeoutMs() {
		proposer.lastPrepareTimeoutMs = comm.GetInsideOptions().GetMaxPrepareTimeoutMs()
	}
}

func (proposer *Proposer) addAcceptTimer(timeOutMs uint32) {
	if proposer.acceptTimerId != 0 {
		proposer.timerThread.DelTimer(proposer.acceptTimerId)
		proposer.acceptTimerId = 0
	}

	if timeOutMs > proposer.timeOutMs {
		timeOutMs = uint32(proposer.timeOutMs)
	}
	proposer.acceptTimerId = proposer.timerThread.AddTimer(timeOutMs, comm.Timer_Proposer_Accept_Timeout, proposer.instance)
	proposer.lastAcceptTimeoutMs *= 2
	if proposer.lastAcceptTimeoutMs > comm.GetInsideOptions().GetMaxAcceptTimeoutMs() {
		proposer.lastAcceptTimeoutMs = comm.GetInsideOptions().GetMaxAcceptTimeoutMs()
	}
}

func (proposer *Proposer) OnPrepareReply(msg *comm.PaxosMsg) error {
	log.Info("[%s]OnPrepareReply from %d", proposer.instance.String(), msg.GetNodeID())

	if proposer.state.state != PREPARE {
		log.Error("[%s]proposer state not PREPARE", proposer.instance.String())
		return nil
	}

	if msg.GetProposalID() != proposer.state.GetProposalId() {
		log.Error("[%s]msg proposal id %d not same to proposer proposal id",
			proposer.instance.String(), msg.GetProposalID(), proposer.state.GetProposalId())
		return nil
	}

	proposer.msgCounter.AddReceive(msg.GetNodeID())

	if msg.GetRejectByPromiseID() == 0 {
		ballot := NewBallotNumber(msg.GetPreAcceptID(), msg.GetPreAcceptNodeID())
		proposer.msgCounter.AddPromiseOrAccept(msg.GetNodeID())
		proposer.state.AddPreAcceptValue(*ballot, msg.GetValue())
		log.Debug("[%s]prepare accepted", proposer.instance.String())
	} else {
		proposer.msgCounter.AddReject(msg.GetNodeID())
		proposer.wasRejectBySomeone = true
		proposer.state.SetOtherProposalId(msg.GetRejectByPromiseID())
		log.Debug("[%s]prepare rejected", proposer.instance.String())
	}

	log.Debug("[%s]%d prepare pass count:%d, major count:%d", proposer.instance.String(), proposer.GetInstanceId(),
		proposer.msgCounter.GetPassedCount(), proposer.config.GetMajorityCount())

	if proposer.msgCounter.IsPassedOnThisRound() {
		proposer.canSkipPrepare = true
		proposer.exitPrepare()
		proposer.accept()
	} else if proposer.msgCounter.IsRejectedOnThisRound() || proposer.msgCounter.IsAllReceiveOnThisRound(){
		proposer.addPrepareTimer(uint32(rand.Intn(30) + 10)) // 重试
	}

	return nil
}

func (proposer *Proposer) accept() {
	if proposer.isTimeout() {
		return
	}
	base := proposer.Base
	state := proposer.state

	log.Info("[%s]start accept %s", proposer.instance.String(), string(state.GetValue()))

	proposer.exitAccept()
	proposer.state.setState(ACCEPT)

	msg := &comm.PaxosMsg{
		MsgType:    proto.Int32(comm.MsgType_PaxosAccept),
		InstanceID: proto.Uint64(base.GetInstanceId()),
		NodeID:     proto.Uint64(proposer.config.GetMyNodeId()),
		ProposalID: proto.Uint64(state.GetProposalId()),
		Value:			state.GetValue(),
		LastChecksum:proto.Uint32(base.GetLastChecksum()),
	}

	proposer.msgCounter.StartNewRound()

	proposer.addAcceptTimer(proposer.lastAcceptTimeoutMs)

	base.broadcastMessage(msg, BroadcastMessage_Type_RunSelf_Final)
}

func (proposer *Proposer) OnAcceptReply(msg *comm.PaxosMsg) error {
	state := proposer.state
	log.Info("[%s]START msg.proposalId %d, state.proposalId %d, msg.from %d, rejectby %d",
		proposer.instance.String(), msg.GetProposalID(), state.GetProposalId(), msg.GetNodeID(), msg.GetRejectByPromiseID())

	base := proposer.Base

	if state.state != ACCEPT {
		log.Error("[%s]proposer state not ACCEPT", proposer.instance.String())
		return nil
	}

	if msg.GetProposalID() != state.GetProposalId() {
		log.Error("[%s]msg proposal id %d not same to proposer proposal id",
			proposer.instance.String(), msg.GetProposalID(), proposer.state.GetProposalId())
		return nil
	}

	msgCounter := proposer.msgCounter
	if msg.GetRejectByPromiseID() == 0 {
		log.Debug("[%s]accept accepted", proposer.instance.String())
		msgCounter.AddPromiseOrAccept(msg.GetNodeID())
	} else {
		log.Debug("[%s]accept rejected", proposer.instance.String())
		msgCounter.AddReject(msg.GetNodeID())
		proposer.wasRejectBySomeone = true
		state.SetOtherProposalId(msg.GetRejectByPromiseID())
	}

	if msgCounter.IsPassedOnThisRound() {
		proposer.exitAccept()
		proposer.learner.ProposerSendSuccess(base.GetInstanceId(), state.GetProposalId())
		log.Info("[%s]instance %d passed", proposer.instance.String(), msg.GetInstanceID())
	} else {
		proposer.addAcceptTimer(uint32(rand.Intn(30) + 10))
	}

	log.Info("OnAcceptReply END")
	return nil
}

func (proposer *Proposer) onPrepareTimeout() {
	proposer.prepare(proposer.wasRejectBySomeone)
}

func (proposer *Proposer) onAcceptTimeout() {
	proposer.prepare(proposer.wasRejectBySomeone)
}



package algorithm

import (
	log "github.com/sirupsen/logrus"
	"bytes"
	"sync"
	"gphxpaxos/comm"
	"gphxpaxos/util"
	"gphxpaxos/config"
	"gphxpaxos/smbase"
)

//---------------------------------------CommitContext--------------------------------------------//
type CommitContext struct {
	instanceId          uint64
	commitEnd           bool
	value               []byte
	stateMachineContext *smbase.SMCtx
	mutex               sync.Mutex
	commitRet           error

	// the start and end time of commit(in ms)
	start uint64
	end   uint64

	instance *Instance

	timeoutMs uint32

	// wait result channel
	wait chan bool
}

func newCommitContext(instance *Instance) *CommitContext {
	context := &CommitContext{
		value:    nil,
		instance: instance,
		wait:     make(chan bool),
	}
	//context.newCommit(nil, nil)

	return context
}

func (commitContext *CommitContext) newCommit(value []byte, timeoutMs uint32, context *smbase.SMCtx) {
	commitContext.mutex.Lock()
	defer commitContext.mutex.Unlock()

	commitContext.instanceId = comm.INVALID_INSTANCEID
	commitContext.commitEnd = false
	commitContext.value = value
	commitContext.stateMachineContext = context
	commitContext.end = 0
	commitContext.start = util.NowTimeMs()
	commitContext.timeoutMs = timeoutMs
}

func (commitContext *CommitContext) isNewCommit() bool {
	return commitContext.instanceId == comm.INVALID_INSTANCEID && commitContext.value != nil
}

func (commitContext *CommitContext) StartCommit(instanceId uint64) uint32 {
	commitContext.mutex.Lock()
	commitContext.instanceId = instanceId
	commitContext.mutex.Unlock()

	return commitContext.timeoutMs
}

func (commitContext *CommitContext) getCommitValue() [] byte {
	return commitContext.value
}

// 是否是本节点自己提交的消息
func (commitContext *CommitContext) IsMyCommit(nodeId uint64, instanceId uint64, learnValue []byte) (bool, *smbase.SMCtx) {
	commitContext.mutex.Lock()
	defer commitContext.mutex.Unlock()

	isMyCommit := false
	var ctx *smbase.SMCtx

	if nodeId != commitContext.instance.config.GetMyNodeId() {
		log.Debug("[%s]%d not my instance id", commitContext.instance.String(), nodeId)
		return false, nil
	}

	isMyCommit = true
	// TODO
	/*
	if !commitContext.commitEnd && commitContext.instanceId == instanceId {
		if bytes.Compare(commitContext.value, learnValue) == 0 {
			isMyCommit = true
		} else {
			log.Debug("[%s]%d not my value", commitContext.instance.String(), instanceId)
			isMyCommit = false
		}
	}
	*/

	if isMyCommit {
		ctx = commitContext.stateMachineContext
	} else {
		log.Debug("[%s]%d not my commit %v:%d", commitContext.instance.String(), instanceId,
			commitContext.commitEnd, commitContext.instanceId)
	}

	return isMyCommit, ctx
}

func (commitContext *CommitContext) setResultOnlyRet(commitret error) {
	commitContext.setResult(commitret, comm.INVALID_INSTANCEID, []byte(""))
}

func (commitContext *CommitContext) setResult(commitret error, instanceId uint64, learnValue []byte) {
	commitContext.mutex.Lock()
	defer commitContext.mutex.Unlock()

	if commitContext.commitEnd || commitContext.instanceId != instanceId {
		log.Errorf("[%s]set result error, commitContext instance id %d,msg instance id %d",
			commitContext.instance.String(), commitContext.instanceId, instanceId)
		return
	}

	commitContext.commitRet = commitret
	if commitContext.commitRet == comm.PaxosTryCommitRet_OK {
		if bytes.Compare(commitContext.value, learnValue) != 0 {
			commitContext.commitRet = comm.PaxosTryCommitRet_Conflict
		}
	}

	commitContext.commitEnd = true
	commitContext.value = nil

	log.Debug("[%s]set commit result instance %d ret %v", commitContext.instance.String(),
		instanceId, commitContext.commitRet)

	commitContext.wait <- true
}

func (commitContext *CommitContext) getResult() (uint64, error) {
	select {
	case <-commitContext.wait:
		break
	}

	if commitContext.commitRet == comm.PaxosTryCommitRet_OK {
		return commitContext.instanceId, commitContext.commitRet
	}

	return 0, commitContext.commitRet
}

//---------------------------------------Committer--------------------------------------------//

const (
	MaxTryCount = 3
)

type Committer struct {
	config    *config.Config
	commitCtx *CommitContext
	factory   *smbase.SMFac
	instance  *Instance

	timeoutMs   uint32
	lastLogTime uint64

	waitLock util.Waitlock
}

func newCommitter(instance *Instance) *Committer {
	return &Committer{
		config:    instance.config,
		commitCtx: instance.commitctx,
		factory:   instance.factory,
		instance:  instance,
	}
}
func (committer *Committer) SetMaxHoldThreads(maxHoldThreads int32) {
	committer.waitLock.MaxWaitCount = maxHoldThreads
}

func (committer *Committer) SetTimeoutMs(timeout uint32) {
	committer.timeoutMs = timeout
}

func (committer *Committer) NewValue(value []byte) (uint64, error) {
	committer.timeoutMs = config.GetMaxCommitTimeoutMs()
	return committer.NewValueGetID(value, nil)
}

func (committer *Committer) NewValueGetID(value []byte, context *smbase.SMCtx) (uint64, error) {
	err := comm.PaxosTryCommitRet_OK
	var instanceid uint64
	for i := 0; i < MaxTryCount && committer.timeoutMs > 0; i++ {
		instanceid, err = committer.newValueGetIDNoRetry(value, context)

		if err != comm.PaxosTryCommitRet_Conflict && err != comm.PaxosTryCommitRet_WaitTimeout {
			break
		}

		if context != nil && context.SMID == comm.MASTER_V_SMID {
			break
		}
	}

	return instanceid, err
}

func (committer *Committer) newValueGetIDNoRetry(value []byte, context *smbase.SMCtx) (uint64, error) {

	lockUseTime, err := committer.waitLock.Lock(int(committer.timeoutMs))

	if err == util.Waitlock_Timeout {
		return 0, comm.PaxosTryCommitRet_WaitTimeout
	}

	if committer.timeoutMs <= uint32(200 + lockUseTime) {
		committer.waitLock.Unlock()
		committer.timeoutMs = 0
		return 0, comm.PaxosTryCommitRet_Timeout
	}

	leftTimeoutMs := committer.timeoutMs - uint32(lockUseTime)

	var smid = int32(0)
	if context != nil {
		smid = context.SMID
	}

	packValue := committer.factory.PackPaxosValue(value,  smid)
	committer.commitCtx.newCommit(packValue, leftTimeoutMs, context)
	committer.instance.sendCommitMsg()

	instanceId, err := committer.commitCtx.getResult()

	committer.waitLock.Unlock()
	return instanceId, err
}

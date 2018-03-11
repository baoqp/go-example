package graft

import (
	"sync"
)

type Node interface {
	nodeId() NodeId

	leaderId() PeerId

	isLeader() bool

	//启动
	init(options *NodeOptions) int

	// 停止本地副本，done回调用于清理资源或者回应客户端等
	shutdown(done *Closure)

	// 阻塞线程知道join成功
	join()

	// [Thread-safe and wait-free]
	// apply task to the replicated-state-machine
	//
	// About the ownership:
	// |task.data|: for the performance consideration, we will take away the
	//              content. If you want keep the content, copy it before call
	//              this function
	// |task.done|: If the data is successfully committed to the raft group. We
	//              will pass the ownership to StateMachine::on_apply.
	//              Otherwise we will specify the error and call it.
	//
	apply(task *Task)

	listPeers() []PeerId

	addPeer(peerId *PeerId, done *Closure)

	removePeer(peerId *PeerId, done *Closure)

	changePeers(newPeers *Configuration, done *Closure)

	// Reset the configuration of this node individually, without any repliation
	// to other peers before this node beomes the leader. This function is
	// supposed to be inovoked when the majority of the replication group are
	// dead and you'd like to revive the service in the consideration of
	// availability.
	// Notice that neither consistency nor consensus are guaranteed in this
	// case, BE CAREFULE when dealing with this method.
	resetPeers(newPeers *Configuration)

	resetElectionTimeOut(timeOutMS int)

	// Try transferring leadership to |peer|.
	// If peer is ANY_PEER, a proper follower will be chosen as the leader for
	// the next term.
	// Returns 0 on success, -1 otherwise.
	transerLeadershipTo(peer *PeerId)

	// Read the first committed user log from the given index.
	// Return OK on success and user_log is assigned with the very data. Be awared
	// that the user_log may be not the exact log at the given index, but the
	// first available user log from the given index to last_committed_index.
	// Otherwise, appropriate errors are returned:
	//     - return ELOGDELETED when the log has been deleted;
	//     - return ENOMOREUSERLOG when we can't get a user log even reaching last_committed_index.
	// [NOTE] in consideration of safety, we use last_applied_index instead of last_committed_index
	// in code implementation.
	readCommittedUserLog(index uint64, userLog *UserLog) Status
}

type Stage int

// 节点状态
const (
	STAGE_NONE        Stage = iota
	STAGE_CATCHING_UP
	STAGE_JOINT
	STAGE_STABLE
)

type ConfigurationCtx struct {
	node        NodeImpl
	stage       Stage
	nchanges    int
	version     uint64
	newPeers    []PeerId
	oldPeers    []PeerId
	addingPeers []PeerId
	done        *Closure
}

type LogEntryAndClosure struct {
	typ      EntryType
	id       LogId
	peers    []PeerId
	oldPeers []PeerId
	data     []byte
}

type StopTransferArg struct {
}

// 实现Node接口
type NodeImpl struct {
	state               State
	currentTerm         uint64
	lastLeaderTimeStamp uint64
	leaderId            PeerId
	votedId             PeerId
	voteCtx             Ballot
	preVoteCtx          Ballot
	conf                ConfigurationEntry

	groupId  GroupId
	serverId PeerId
	options  NodeOptions

	mutex                 sync.Mutex
	confCtx               ConfigurationCtx
	LogStorage            *LogStorage
	metaStorage           *RaftMetaStorage
	closureQueue          *ClosureQueue
	configManager         *ConfigurationManager
	logManager            *LogManager
	fsmCaller             *FSMCaller
	ballotBox             *BallotBox
	replicatorGroup       *ReplicatorGroup
	shutdownContinuations []Closure
	// electionTimer
	// voteTimer
	// stepdownTimer

	wakingCandidate ReplicaId

	//bthread::ExecutionQueueId<LogEntryAndClosure> _apply_queue_id;
	//bthread::ExecutionQueue<LogEntryAndClosure>::scoped_ptr_t _apply_queue;
}

func (this *NodeImpl) bootstrap(options *BootstrapOptions) {
	if options.groupConf == nil {
		panic("bootstraping an empty node makes no sense")
	}

	bootstrapLogTerm := 0 // TODO

	this.bootstrap(options)

	bootstrapId := &LogId{index: options.lastLogIndex, term: bootstrapLogTerm}

	configManager := &ConfigurationManager{}

	fsmCaller := &FSMCaller{}

	this.initLogStorage()

}
func (this *NodeImpl) initLogStorage() {
	logStorage := CreateLogStorage(this.options.logUri)
	this.logManager = &LogManager{}



}

package graft

const (
	IP_ANY = "0.0.0.0"
)


// 节点状态
const (
	STAGE_NONE        Stage = iota
	STAGE_CATCHING_UP
	STAGE_JOINT
	STAGE_STABLE
)


const(
	raftElectionDelayMS = 100
	raftElectionHeartbeatFactor = 10
)

const (
	raftSync = true
	raftCreateParentDirectories = true
	raftMetaSync = true
)

const (
	sRaftMeta = "raft_meta"
)
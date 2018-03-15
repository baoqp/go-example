package node

// 状态机上下文
type SMCtx struct {
	SMID int
	pCtx interface{}
}

// 状态机接口
type StateMachine interface {

	Execute(groupIdx int, instanceId uint64, paxosValue[]byte, context *SMCtx) error

	SMID() int32

	ExecuteForCheckpoint(groupIdx int, instanceId uint64, paxosValue []byte)error

	GetCheckpointInstanceID(groupIdx int) uint64

	LockCheckpointState() error

	GetCheckpointState(groupIdx int, dirPath *string, fileList []string) error

	UnLockCheckpointState()

	LoadCheckpointState(groupIdx int, checkpointTmpFileDirPath string,
		fileList []string, checkpointInstanceID uint64) error

	BeforePropose(groupIdx int, value []byte)

	NeedCallBeforePropose() bool
}


type InsideSM interface {
	StateMachine

	GetCheckpointBuffer(cpBuffer *string) error
	UpdateByCheckpoint(cpBuffer *string, change bool) error
}


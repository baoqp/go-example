package graft


type ErrorType int

type Closure interface {
	Run() int
}

const (
	ERROR_TYPE_NONE          ErrorType = iota
	ERROR_TYPE_LOG
	ERROR_TYPE_STABLE
	ERROR_TYPE_SNAPSHOT
	ERROR_TYPE_STATE_MACHINE
)

type Error struct {
	typ    ErrorType
	status string
}

// Basic message structure of raft
type Task struct {
	data         []byte
	done         *Closure
	expectedTerm uint64
}



// This class encapsulates the parameter of on_start_following and on_stop_following interfaces.
type LeaderChangeContext struct {

}

// 状态机
type StateMachine interface {
	onApply(iter *Iterator)

	onShutdown()

	onLeaderStart(term uint64)

	onLeaderStop(status string)

	onError(e Error)

	onConfigurationCommitted(conf Configuration)

	onStopFollowing(ctx *LeaderChangeContext)

	onStartFollowing(ctx LeaderChangeContext);
}

type State int

// 状态
const (
	STATE_LEADER        State = iota
	STATE_TRANSFERRING
	STATE_CANDIDATE
	STATE_FOLLOWER
	STATE_ERROR
	STATE_UNINITIALIZED
	STATE_SHUTTING
	STATE_SHUTDOWN
	STATE_END
)


type Status string

func IsActiveState(s State) bool {
	return s < STATE_ERROR
}


type UserLog struct {
	data []byte
	index uint64
}
































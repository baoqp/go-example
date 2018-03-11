package graft

type NodeOptions struct {

	*BootstrapOptions


	electionTimeoutMS int

	catchupMargin int

	initialConf Configuration

	fsm *StateMachine

	// If |node_owns_fsm| is true. |fms| would be destroyed when the backing
	// Node is no longer referenced. Default: false
	nodeOwnsFsm bool

	// Describe a specific LogStorage in format ${type}://${parameters}
	logUri string

	// Describe a specific RaftMetaStorage in format ${type}://${parameters}
	raftMetaUri string
}

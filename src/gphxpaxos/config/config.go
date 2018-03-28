package config

import (
	"math"
	"gphxpaxos/comm"
	"gphxpaxos/master"
	"gphxpaxos/util"
	"gphxpaxos/smbase"
)

// TODO
type Config struct {
	myNodeId  uint64
	nodeCount int
	myGroupId int

	isFollower          bool
	followToNodeId      uint64
	systemVSM           *smbase.SystemVSM
	masterStateMachine  *master.MasterStateMachine
	myFollowerMap       map[uint64]uint64
	tmpNodeOnlyForLearn map[uint64]uint64

	options  *comm.Options
	majorCnt int
}

func NewConfig(options *comm.Options, groupId int ) *Config {
	return &Config{
		options:  options,
		myGroupId: groupId,
		majorCnt: int(math.Floor(float64(len(options.NodeInfoList))/2)) + 1,
	}
}

func (config *Config) GetOptions() *comm.Options {
	return config.options
}

func (config *Config) LogSync() bool {
	return true
}

func (config *Config) SyncInterval() int32 {
	return 5
}

func (config *Config) GetMyGroupId() int {
	return config.myGroupId
}

func (config *Config) GetGid() uint64 {
	return 0
}

func (config *Config) GetMyNodeId() uint64 {
	return config.options.MyNodeInfo.NodeId
}

func (config *Config) GetMajorityCount() int {
	return config.majorCnt
}

func (config *Config) GetNodeCount() int {
	return 0
}

func (config *Config) IsIMFollower() bool {
	return config.isFollower
}

func (config *Config) GetFollowToNodeID() uint64 {
	return config.followToNodeId
}

func (config *Config) GetMyFollowerCount() int32 {
	return int32(len(config.myFollowerMap))
}

func (config *Config) AddFollowerNode(followerNodeId uint64) {
	config.myFollowerMap[followerNodeId] = util.NowTimeMs() + uint64(comm.GetInsideOptions().GetAskforLearnInterval()*3)
}

func (config *Config) AddTmpNodeOnlyForLearn(nodeId uint64) {

}

func (config *Config) GetSystemVSM() *smbase.SystemVSM {
	return config.systemVSM
}

func (config *Config) SetMasterSM(masterSM *master.MasterStateMachine) *master.MasterStateMachine {
	config.masterStateMachine = masterSM
}

func (config *Config) GetMasterSM() *master.MasterStateMachine {
	return config.masterStateMachine
}

func (config *Config) CheckConfig() bool {
	return true
}

func (config *Config) GetIsUseMembership() bool {
	return false
}

func (config *Config) IsValidNodeID(nodeId uint64) bool {
	return true
}

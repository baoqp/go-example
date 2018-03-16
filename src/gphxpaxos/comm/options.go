package comm

import (
	"fmt"
	"gphxpaxos/util"
	"gphxpaxos/logstorage"
	"gphxpaxos/network"
	"gphxpaxos/node"
)

type NodeInfoList []*NodeInfo

type NodeInfo struct {
	Ip     string
	Port   int
	NodeId uint64
}

func (nodeInfo *NodeInfo) String() string {
	return fmt.Sprintf("%s:%d", nodeInfo.Ip, nodeInfo.Port)
}

func makeNodeId(nodeInfo *NodeInfo) {
	ip := util.Inet_addr(nodeInfo.Ip)
	nodeInfo.NodeId = uint64(ip)<<32 | uint64(nodeInfo.Port)
}

func NewNodeInfo(ip string, port int) *NodeInfo {
	nodeInfo := &NodeInfo{
		Ip:   ip,
		Port: port,
	}
	makeNodeId(nodeInfo)
	return nodeInfo
}

type FollowerNodeInfoList []*FollowerNodeInfo

type FollowerNodeInfo struct {
	myNode     NodeInfo
	followNode NodeInfo
}

// 两个回调函数
type MembershipChangeCallback func(int, list *NodeInfoList)
type MasterChangeCallback func(int, *NodeInfo, uint64)

// group的状态机数据
type GroupSMInfo struct {
	groupIdx    int
	SMList      []node.StateMachine
	isUseMaster bool // 是否使用内置的状态机来进行master选举
}

type GroupSMInfoList []*GroupSMInfo

type Options struct {
	LogStorage                   logstorage.LogStorage
	LogStoragePath               string
	Sync                         bool
	SyncInternal                 int
	NetWork                      network.NetWork
	GroupCount                   int
	UseMemebership               bool
	MyNodeInfo                   *NodeInfo
	NodeInfoList                 NodeInfoList
	MembershipChangeCallback     *MembershipChangeCallback
	MasterChangeCallback         *MasterChangeCallback
	GroupSMINfoList              GroupSMInfoList
	FollowerNodeInfoList         FollowerNodeInfoList
	UseCheckpointReplayer        bool
	UseBatchPropose              bool
	OpenChangeValueBeforePropose bool
}

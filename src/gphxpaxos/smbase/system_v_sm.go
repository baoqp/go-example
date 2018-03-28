package smbase

import (
	"gphxpaxos/comm"
	"gphxpaxos/storage"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"errors"
)

// 实现InsideSM接口
type SystemVSM struct {
	myGroupId                int
	systemVariables           *comm.SystemVariables
	systemStore              *storage.SystemVariablesStore
	nodeIdSet                map[uint64]struct{} // 需要一个set, 使用map表示
	myNodeId                 uint64
	membershipChangeCallback comm.MembershipChangeCallback
}

func NewSystemVSM(groupId int, myNodeId uint64, logstorage storage.LogStorage,
	membershipChangeCallback comm.MembershipChangeCallback) *SystemVSM {

	return &SystemVSM{
		myGroupId:                groupId,
		myNodeId:                 myNodeId,
		systemStore:              storage.NewSystemVariablesStore(logstorage),
		membershipChangeCallback: membershipChangeCallback,
	}
}

func (systemVSM *SystemVSM) Init() error {
	var variable = &comm.SystemVariables{}
	err := systemVSM.systemStore.Read(systemVSM.myGroupId, variable)
	if err != nil && err != comm.ErrKeyNotFound {
		return err
	} else if err != nil && err == comm.ErrKeyNotFound {
		systemVSM.systemVariables.Gid = proto.Uint64(0)
		systemVSM.systemVariables.Version = proto.Uint64(-1) // TODO set uint64 = -1 ???
		log.Info("variables not exist")
	} else {
		systemVSM.RefleshNodeID()
	}
	return nil
}

func (systemVSM *SystemVSM) UpdateSystemVariables(variables *comm.SystemVariables) error {
	writeOpt := &storage.WriteOptions{Sync:true}
	err := systemVSM.systemStore.Write(writeOpt, systemVSM.myGroupId, variables)
	if err != nil {
		return err
	}
	systemVSM.systemVariables = variables
	systemVSM.RefleshNodeID()
	return nil
}


func (systemVSM *SystemVSM) Execute(groupId int, instanceId uint64, value []byte, ctx *SMCtx) error {
	var variables = &comm.SystemVariables{}
	err := proto.Unmarshal(value, variables)
	if err != nil {
		log.Errorf("Variables.ParseFromArray fail:%v", err)
		return err
	}

	var smret error
	if ctx != nil && ctx.pCtx != nil {
		smret = (ctx.pCtx).(error)
	}

	if variables.GetGid() != 0 && variables.GetGid() != systemVSM.systemVariables.GetGid() {
		log.Errorf("modify.gid %d not equal to now.gid %d", variables.GetGid(), systemVSM.systemVariables.GetGid())
		return errors.New("bad gid")
	}

	if variables.GetVersion() != systemVSM.systemVariables.GetVersion() {
		log.Errorf("modify.version %d not equal to now.version %d", variables.GetVersion(), systemVSM.systemVariables.GetVersion())
		if smret != nil {
			smret = comm.Paxos_MembershipOp_GidNotSame
		}
		return nil
	}

	variables.Version = proto.Uint64(instanceId)
	err = systemVSM.UpdateSystemVariables(variables)
	if err != nil {
		return err
	}

	log.Info("OK, new version %d gid %d", systemVSM.systemVariables.GetVersion(), systemVSM.systemVariables.GetGid())
	smret = nil
	return nil

}


func (systemVSM *SystemVSM) GetGid() uint64 {
	return systemVSM.systemVariables.GetGid()
}

func (systemVSM *SystemVSM) GetMembership(nodes comm.NodeInfoList, version *uint64) {
	*version = systemVSM.systemVariables.GetVersion()

	for i := 0; i < len(systemVSM.systemVariables.MemberShip); i++ {
		node := systemVSM.systemVariables.MemberShip[i]
		tmp :=  &comm.NodeInfo{
			NodeId: node.GetNodeid(),
		}
		nodes = append(nodes, tmp)
	}
}

func (systemVSM *SystemVSM) Membership_OPValue(nodes comm.NodeInfoList, version uint64, value []byte) error {

	variables := &comm.SystemVariables{
		Version: proto.Uint64(version),
		Gid:     proto.Uint64(systemVSM.systemVariables.GetGid()),
	}

	for _, node := range nodes {
		tmp := &comm.PaxosNodeInfo{
			Rid:    proto.Uint64(0),
			Nodeid: proto.Uint64(node.NodeId),
		}

		systemVSM.systemVariables.MemberShip = append(systemVSM.systemVariables.MemberShip, tmp)
	}

	value, err := proto.Marshal(variables)
	if err != nil {
		log.Errorf("Variables.Serialize fail: %v", err)
		return err
	}

	return nil
}


func (systemVSM *SystemVSM) CreateGid_OPValue(gid uint64) (  []byte, error) {
	variables := proto.Clone(systemVSM.systemVariables).(*comm.SystemVariables)
	variables.Gid = proto.Uint64(gid)
	value, err := proto.Marshal(variables)
	if err != nil {
		log.Errorf("Variables.Serialize fail: %v", err)
		return nil, err
	}
	return value, nil
}


func (systemVSM *SystemVSM) AddNodeIDList(nodes comm.NodeInfoList) {
	if systemVSM.systemVariables.GetGid() != 0 {
		return
	}

	systemVSM.nodeIdSet = make(map[uint64]struct{})
	systemVSM.systemVariables.MemberShip = make([]*comm.PaxosNodeInfo, 0)

	for _, node := range nodes {
		tmp := &comm.PaxosNodeInfo{
			Rid:    proto.Uint64(0),
			Nodeid: proto.Uint64(node.NodeId),
		}

		systemVSM.systemVariables.MemberShip = append(systemVSM.systemVariables.MemberShip, tmp)
	}

	systemVSM.RefleshNodeID()
}


func (systemVSM *SystemVSM) RefleshNodeID() {
	systemVSM.nodeIdSet = make(map[uint64]struct{})
	var infolist []*comm.NodeInfo
	membership := systemVSM.systemVariables.MemberShip
	for i := 0; i < len(membership); i++ {
		paxosNodeInfo := membership[i]
		tmpNode := &comm.NodeInfo{NodeId: *paxosNodeInfo.Nodeid} // TODO
		systemVSM.nodeIdSet[tmpNode.NodeId] = struct{}{}
		infolist = append(infolist, tmpNode)
	}

	if systemVSM.membershipChangeCallback != nil {
		systemVSM.membershipChangeCallback(systemVSM.myGroupId, comm.NodeInfoList(infolist))
	}

}

func (systemVSM *SystemVSM) GetNodeCount() int {
	return len(systemVSM.nodeIdSet)
}

func (systemVSM *SystemVSM) GetMajorityCount() int {
	return int(systemVSM.GetNodeCount() / 2.0 + 1)
}


func (systemVSM *SystemVSM) IsValidNodeID(nodeId uint64) bool {
	if systemVSM.systemVariables.GetGid() == 0 {
		return true
	}

	_, ok := systemVSM.nodeIdSet[nodeId]
	return ok
}

func (systemVSM *SystemVSM) IsIMInMembership() bool {
	_, ok := systemVSM.nodeIdSet[systemVSM.myNodeId]
	return ok
}


func (systemVSM *SystemVSM) GetCheckpointBuffer() ([]byte, error) {
	// TODO 使用的地方需要判断是否为空
	if systemVSM.systemVariables.GetVersion() == uint64(-1) ||
		systemVSM.systemVariables.GetGid() == 0{

		return nil, nil
	}

	value, err := proto.Marshal(systemVSM.systemVariables)
	if err != nil {
		return nil, err
	}

	return value, nil
}

var VersionGidErr = errors.New("variables.version not init or gid not same")

func (systemVSM *SystemVSM) UpdateByCheckpoint(value []byte) (bool, error) {

	if len(value) == 0 {
		return false, nil
	}

	change := false

	var varaible = &comm.SystemVariables{}
	err := proto.Unmarshal(value, varaible)
	if err != nil {
		return false, err
	}

	if *varaible.Version == uint64(-1) {
		return false, VersionGidErr
	}

	if *varaible.Gid != 0 && varaible.GetGid() != systemVSM.systemVariables.GetGid() {
		return false, VersionGidErr
	}

	if systemVSM.systemVariables.GetVersion() != uint64(-1) &&
		*varaible.Version <= *systemVSM.systemVariables.Version {
		return false, nil
	}

	change = true
	err = systemVSM.UpdateSystemVariables(varaible)
	if err != nil {
		return change, err
	}
	return change, nil
}

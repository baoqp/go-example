package master

import (
	"sync"
	"gphxpaxos/comm"
	"gphxpaxos/storage"
	log "github.com/sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"gphxpaxos/util"
	"github.com/chrislusf/glow/resource/service_discovery/master"
)

// 实现InsideSM接口
type MasterStateMachine struct {
	myGroupId            int
	myNodeId             uint64
	mvStore              *MasterVariablesStore
	masterNodeId         uint64
	masterVersion        uint64
	leaseTime            int
	absExpireTime        uint64
	mutex                sync.Mutex
	masterChangeCallback comm.MasterChangeCallback
}

func NewMasterStateMachine(groupId int, myNodeId uint64, logstorage storage.LogStorage,
	masterChangeCallback comm.MasterChangeCallback) *MasterStateMachine {

	return &MasterStateMachine{
		myGroupId:            groupId,
		myNodeId:             myNodeId,
		mvStore:              NewMasterVariablesStore(logstorage),
		masterNodeId:         0,
		masterVersion:        uint64(-1),
		leaseTime:            0,
		absExpireTime:        0,
		masterChangeCallback: masterChangeCallback,
	}
}

func (masterSM *MasterStateMachine) Init() error {
	var variables comm.MasterVariables
	err := masterSM.mvStore.Read(masterSM.myGroupId, &variables)
	if err != nil && err != comm.ErrKeyNotFound {
		log.Errorf("Master variables read from store fail %v", err)
		return err
	}

	if err == comm.ErrKeyNotFound {
		log.Info("no master variables exist")
	} else {
		masterSM.masterVersion = variables.GetVersion()
		if variables.GetVersion() == masterSM.myNodeId {
			masterSM.masterNodeId = comm.NULL_NODEID
			masterSM.absExpireTime = 0
		} else {
			masterSM.masterNodeId = variables.GetMasterNodeid()
			masterSM.absExpireTime = util.NowTimeMs() + uint64(variables.GetLeaseTime())
		}
	}

	log.Info("OK, master nodeid %d version %d expiretime %d",
		masterSM.masterNodeId, masterSM.masterVersion, masterSM.absExpireTime)
	return nil
}

func (masterSM *MasterStateMachine) UpdateMasterToStore(masterNodeId uint64, version uint64,
	leaseTime int32) error {

	variables := comm.MasterVariables{
		MasterNodeid: proto.Uint64(masterNodeId),
		Version:      proto.Uint64(version),
		LeaseTime:    proto.Uint32(uint32(leaseTime)),
	}

	options := &storage.WriteOptions{
		Sync: false,
	}

	return masterSM.mvStore.Write(options, masterSM.myGroupId, variables)
}

func (masterSM *MasterStateMachine) LearnMaster(instanceId uint64, operator *MasterOperator,
	absMasterTimeout uint64) error {

	masterSM.mutex.Lock()
	defer masterSM.mutex.Unlock()

	if operator.GetLastversion() != 0 &&
		instanceId > masterSM.masterVersion &&
		operator.GetLastversion() != masterSM.masterVersion {

		log.Errorf("other last version %d not same to my last version %d, instanceid %d",
			operator.GetLastversion(), masterSM.masterVersion, instanceId)

		log.Errorf("try to fix, set my master version %d as other last version %d, instanceid %d",
			masterSM.masterVersion, operator.GetLastversion(), instanceId)
		masterSM.masterVersion = operator.GetLastversion()

	}

	if operator.GetVersion() != masterSM.masterVersion {
		log.Errorf("version conflit, op version %d now master version %d",
			operator.GetVersion(), masterSM.masterVersion)

		return nil
	}

	err := masterSM.UpdateMasterToStore(operator.GetNodeid(), instanceId, operator.GetTimeout())
	if err != nil {
		log.Errorf("UpdateMasterToStore fail %v", err)
		return err
	}

	// TODO

	masterSM.masterNodeId = operator.GetNodeid()
	if masterSM.masterNodeId == masterSM.myNodeId {
		masterSM.absExpireTime = absMasterTimeout
		log.Info("Be master success, absexpiretime %d", masterSM.absExpireTime)
	} else {
		masterSM.absExpireTime = util.NowTimeMs() + uint64(operator.GetTimeout())
		log.Info("Other be master, absexpiretime %d", masterSM.absExpireTime)
	}

	masterSM.leaseTime = operator.GetTimeout()
	masterSM.masterVersion = instance

	log.Info("OK, masternodeid %d version %d abstimeout %d",
		masterSM.masterNodeId, masterSM.masterVersion, masterSM.absExpireTime)

	return nil
}

func (masterSM *MasterStateMachine) SafeGetMaster(masterNodeId *uint64, masterVersion *uint64) {
	masterSM.mutex.Lock()

	if util.NowTimeMs() >= masterSM.absExpireTime {
		*masterNodeId = comm.NULL_NODEID
	} else {
		*masterNodeId = masterSM.masterNodeId
	}
	*masterVersion = masterSM.masterVersion
	masterSM.mutex.Unlock()
}

func (masterSM *MasterStateMachine) GetMaster() uint64 {
	if util.NowTimeMs() >= masterSM.absExpireTime {
		return comm.NULL_NODEID
	}

	return masterSM.masterNodeId
}

func (masterSM *MasterStateMachine) GetMasterWithVersion(version *uint64) uint64 {
	masterNodeId := comm.NULL_NODEID
	masterSM.SafeGetMaster(&masterNodeId, version)
	return masterNodeId
}

func (masterSM *MasterStateMachine) IsIMMaster() bool {
	return masterSM.GetMaster() == masterSM.myNodeId
}

func (masterSM *MasterStateMachine) Execute(groupIdx int32, instanceId uint64, value []byte,
	ctx *gpaxos.StateMachineContext) error {
	var operator master.MasterOperator
	err := proto.Unmarshal(value, &operator)
	if err != nil {
		log.Errorf("oMasterOper data wrong %v", err)
		return err
	}

	if operator.GetOperator() == MasterOperatorType_Complete {
		var absMasterTimeout uint64 = 0
		if ctx != nil && ctx.Context != nil {
			absMasterTimeout = *(ctx.Context.(*uint64))
		}

		log.Info("absmaster timeout %v", absMasterTimeout)

		err = masterSM.LearnMaster(instanceId, operator, absMasterTimeout)
		if err != nil {
			return err
		}
	} else {
		log.Errorf("unknown op %d", operator.GetOperator())
		return nil
	}

	return nil
}

func (masterSM *MasterStateMachine) MakeOpValue(nodeId uint64, version uint64,
	timeout int32, op uint32) ([]byte, error) {
	operator := master.MasterOperator{
		Nodeid:   proto.Uint64(nodeId),
		Version:  proto.Uint64(version),
		Timeout:  proto.Int32(timeout),
		Operator: proto.Uint32(op),
		Sid:      proto.Uint32(util.Rand()),
	}

	return proto.Marshal(operator)
}

func (masterSM *MasterStateMachine) GetCheckpointBuffer() ([]byte, error) {
	if masterSM.masterVersion == -1 {
		return nil, nil
	}

	v := comm.MasterVariables{
		MasterNodeid: proto.Uint64(masterSM.masterNodeId),
		Version:      proto.Uint64(masterSM.masterVersion),
		LeaseTime:    proto.Uint32(uint32(masterSM.leaseTime)),
	}

	return proto.Marshal(v)
}

func (masterSM *MasterStateMachine) UpdateByCheckpoint(buffer []byte, change *bool) error {
	if len(buffer) == 0 {
		return nil
	}

	var variables comm.MasterVariables
	err := proto.Unmarshal(buffer, &variables)
	if err != nil {
		log.Errorf("Variables.ParseFromArray fail: %v", err)
		return err
	}

	if variables.GetVersion() <= masterSM.masterVersion && masterSM.masterVersion != -1 {
		log.Info("lag checkpoint, no need update, cp.version %d now.version %d",
			variables.GetVersion(), masterSM.masterVersion)
		return nil
	}

	err = masterSM.UpdateMasterToStore(variables.GetMasterNodeid(), variables.GetVersion(), int32(variables.GetLeaseTime()))
	if err != nil {
		return err
	}

	log.Info("ok, cp.version %d cp.masternodeid %d old.version %d old.masternodeid %d",
		variables.GetVersion(), variables.GetMasterNodeid(), masterSM.masterVersion, masterSM.masterNodeId)

	masterSM.masterVersion = variables.GetVersion()
	if variables.GetMasterNodeid() == masterSM.myNodeId {
		masterSM.masterNodeId = comm.NULL_NODEID
		masterSM.absExpireTime = 0
	} else {
		masterSM.masterNodeId = variables.GetMasterNodeid()
		masterSM.absExpireTime = util.NowTimeMs() + uint64(variables.GetLeaseTime())
	}

	return nil
}

// TODO
func (masterStateMachine *MasterStateMachine) GetCheckpointBuffer() (string, error) {
	return "", nil
}

func (masterStateMachine *MasterStateMachine) UpdateByCheckpoint(systemVariables []byte) (bool, error) {
	return true, nil
}

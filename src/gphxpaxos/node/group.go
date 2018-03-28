package node

import (
	"gphxpaxos/network"
	"gphxpaxos/config"
	"gphxpaxos/algorithm"
	"gphxpaxos/storage"
	"gphxpaxos/comm"
	"gphxpaxos/master"
	"gphxpaxos/checkpoint"
	"gphxpaxos/smbase"
)

type Group struct {
	communicate *network.Communicate
	config      *config.Config
	instance    *algorithm.Instance
}

func NewGroup(logstorage storage.LogStorage, network_ network.NetWork, masterSM *master.MasterStateMachine,
	groupId int, options *comm.Options) (*Group, error) {

	group := &Group{}
	group.config = config.NewConfig(options, groupId)
	group.config.SetMasterSM(masterSM)
	group.communicate = network.NewCommunicate(group.config, options.MyNodeInfo.NodeId, network_)
	//func NewInstance(config *config.Config, logstorage storage.LogStorage, transport network.MsgTransport,
	//useCkReplayer bool) (*Instance, error) {
	var err error
	group.instance, err = algorithm.NewInstance(group.config, logstorage, group.communicate, options.UseCheckpointReplayer)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (group *Group) Stop() {
	group.instance.Stop()
}

func (group *Group) GetConfig() *config.Config {
	return group.config
}

func (group *Group) GetInstance() *algorithm.Instance {
	return group.instance
}

func (group *Group) GetCheckpointCleaner() *checkpoint.Cleaner {
	return group.instance.GetCheckpointCleaner()
}


func (group *Group) GetCheckpointReplayer() *checkpoint.Replayer {
	return group.instance.GetCheckpointReplayer()
}

func (group *Group) AddStateMachine(sm smbase.StateMachine) {
	group.instance.AddStateMachine(sm)
}


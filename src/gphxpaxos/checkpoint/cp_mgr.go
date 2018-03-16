package checkpoint

import (
	"gphxpaxos/config"
	"gphxpaxos/logstorage"
	"gphxpaxos/node"
)

type CheckpointManager struct {
	config     *config.Config
	logStorage logstorage.LogStorage
	factory    *node.SMFac
	cleaner    *Cleaner
	replayer   *Replayer

	minChosenInstanceId    uint64
	maxChosenInstanceId    uint64
	inAskforCheckpointMode bool
	useCheckpointReplayer  bool

	needAskSet               map[uint64]bool
	lastAskforCheckpointTime uint64
}

func NewCheckpointManager(config *config.Config, factory *node.SMFac,
	logStorage logstorage.LogStorage, useReplayer bool) *CheckpointManager {

	mnger := &CheckpointManager{
		config:                config,
		logStorage:            logStorage,
		factory:               factory,
		useCheckpointReplayer: useReplayer,
	}

	mnger.cleaner = NewCleaner(config, factory, logStorage, mnger)
	if useReplayer {
		mnger.replayer = NewReplayer(config, factory, logStorage, mnger)
	}
	return mnger
}

func (checkpointManager *CheckpointManager) Init() error {
	instanceId, err := checkpointManager.logStorage.GetMinChosenInstanceId(checkpointManager.config.GetMyGroupId())
	if err != nil {
		return err
	}

	checkpointManager.minChosenInstanceId = instanceId
	err = checkpointManager.cleaner.FixMinChosenInstanceID(checkpointManager.minChosenInstanceId)
	if err != nil {
		return err
	}
	return nil
}

func (checkpointManager *CheckpointManager) Start() {
	checkpointManager.cleaner.Start()
	if checkpointManager.useCheckpointReplayer {
		checkpointManager.replayer.Start()
	}
}

func (checkpointManager *CheckpointManager) Stop() {
	if checkpointManager.useCheckpointReplayer {
		checkpointManager.replayer.Stop()
	}
	checkpointManager.cleaner.Stop()
}

func (checkpointManager *CheckpointManager) GetRelayer() *Replayer {
	return checkpointManager.replayer
}

func (checkpointManager *CheckpointManager) PrepareForAskforCheckpoint(sendNodeId uint64) error {
	checkpointManager.needAskSet[sendNodeId] = true
	if checkpointManager.lastAskforCheckpointTime == 0 {
		checkpointManager.lastAskforCheckpointTime = util.NowTimeMs()
	}

	now := util.NowTimeMs()
	if now >= checkpointManager.lastAskforCheckpointTime+60000 {

	} else {
		if len(checkpointManager.needAskSet) < checkpointManager.config.GetMajorityCount() {

		}
	}

	checkpointManager.lastAskforCheckpointTime = 0
	checkpointManager.inAskforCheckpointMode = true

	return nil
}

func (checkpointManager *CheckpointManager) GetMinChosenInstanceID() uint64 {
	return checkpointManager.minChosenInstanceId
}

func (checkpointManager *CheckpointManager) GetMaxChosenInstanceID() uint64 {
	return checkpointManager.maxChosenInstanceId
}

func (checkpointManager *CheckpointManager) SetMaxChosenInstanceID(instanceId uint64) {
	checkpointManager.maxChosenInstanceId = instanceId
}

func (checkpointManager *CheckpointManager) SetMinChosenInstanceID(instanceId uint64) error {
	/*
	options := storage.WriteOptions{
		Sync:true,
	}
	*/
	err := checkpointManager.logStorage.SetMinChosenInstanceID(instanceId)
	if err != nil {
		return err
	}

	checkpointManager.minChosenInstanceId = instanceId
	return nil
}

func (checkpointManager *CheckpointManager) GetCheckpointInstanceID() uint64 {
	return checkpointManager.factory.GetCheckpointInstanceID()
}

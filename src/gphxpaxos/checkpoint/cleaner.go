package checkpoint

import (
	"gphxpaxos/config"
	"time"
	log "github.com/sirupsen/logrus"
	"gphxpaxos/storage"
	"gphxpaxos/node"
	"gphxpaxos/util"
	"gphxpaxos/comm"
)

const (
	DELETE_SAVE_INTERVAL = 10
)

type Cleaner struct {
	config     *config.Config
	logStorage storage.LogStorage
	factory    *node.SMFac
	ckmnger    *CheckpointManager

	holdCount          uint64
	lastSaveInstanceId uint64

	isPaused bool
	isEnd    bool
	isStart  bool
	canRun   bool
}

func NewCleaner(config *config.Config, factory *node.SMFac,
	logStorage storage.LogStorage, mnger *CheckpointManager) *Cleaner {
	cleaner := &Cleaner{
		config:     config,
		logStorage: logStorage,
		factory:    factory,
		ckmnger:    mnger,
		holdCount:  100000,
	}

	return cleaner
}

func (cleaner *Cleaner) Start() {
	util.StartRoutine(cleaner.main)
}

func (cleaner *Cleaner) Stop() {
	cleaner.isEnd = true
}

func (cleaner *Cleaner) Pause() {
	cleaner.canRun = false
}

func (cleaner *Cleaner) main() {
	cleaner.isStart = true
	cleaner.Continue()

	deleteQps := comm.GetInsideOptions().GetCleanerDeleteQps()
	sleepMs := 1
	deleteInterval := 1
	if deleteQps < 1000 {
		sleepMs = int(1000 / deleteQps)
	} else {
		deleteInterval = int(deleteQps/1000) + 1
	}

	deleteCnt := 0
	for {
		if cleaner.isEnd {
			break
		}

		if !cleaner.canRun {
			cleaner.isPaused = true
			time.Sleep(1000 * time.Millisecond)
			continue
		}

		instanceId := cleaner.ckmnger.GetMinChosenInstanceID()
		maxInstanceId := cleaner.ckmnger.GetMinChosenInstanceID()
		cpInstanceId := cleaner.factory.GetCheckpointInstanceID(cleaner.config.GetMyGroupId()) + 1

		//
		for instanceId+cleaner.holdCount < cpInstanceId && instanceId+cleaner.holdCount < maxInstanceId {
			err := cleaner.DeleteOne(instanceId)
			if err != nil {
				log.Error("delete system fail, instanceid %d", instanceId)
				break
			}

			instanceId += 1
			deleteCnt += 1
			if deleteCnt >= deleteInterval { // TODO ???
				deleteCnt = 0
				time.Sleep(time.Duration(sleepMs) * time.Millisecond)
			}
		}

		if cpInstanceId == 0 {
			log.Info("sleep a while, max deleted instanceid %d checkpoint instanceid(no checkpoint) now instance id %d",
				instanceId, cleaner.ckmnger.GetMaxChosenInstanceID())
		} else {
			log.Info("sleep a while, max deleted instanceid %d checkpoint instanceid %d now instance id %d",
				instanceId, cpInstanceId, cleaner.ckmnger.GetMaxChosenInstanceID())
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func (cleaner *Cleaner) FixMinChosenInstanceID(oldMinChosenInstanceId uint64) error {
	cpInstanceId := cleaner.factory.GetCheckpointInstanceID(cleaner.config.GetMyGroupId()) + 1
	fixMinChosenInstanceId := oldMinChosenInstanceId

	for instanceId := oldMinChosenInstanceId; instanceId < oldMinChosenInstanceId+DELETE_SAVE_INTERVAL; instanceId++ {
		if instanceId >= cpInstanceId {
			break
		}

		_, err := cleaner.logStorage.Get(cleaner.config.GetMyGroupId(), instanceId)

		if err == nil {
			fixMinChosenInstanceId = instanceId + 1
		}
	}

	if fixMinChosenInstanceId > oldMinChosenInstanceId {
		err := cleaner.ckmnger.SetMinChosenInstanceID(fixMinChosenInstanceId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cleaner *Cleaner) DeleteOne(instanceId uint64) error {
	options := &storage.WriteOptions{
		Sync: false,
	}

	err := cleaner.logStorage.Del(options, cleaner.config.GetMyGroupId(), instanceId)
	if err != nil {
		return err
	}

	cleaner.ckmnger.SetMinChosenInstanceID(instanceId)
	if instanceId >= cleaner.lastSaveInstanceId+DELETE_SAVE_INTERVAL {
		err := cleaner.ckmnger.SetMinChosenInstanceID(instanceId + 1)
		if err != nil {
			log.Error("SetMinChosenInstanceID fail, now delete instanceid %d", instanceId)
			return err
		}
		cleaner.lastSaveInstanceId = instanceId

		log.Info("delete %d instance done, now minchosen instanceid %d", DELETE_SAVE_INTERVAL, instanceId+1)
	}

	return nil
}

func (cleaner *Cleaner) Continue() {
	cleaner.isPaused = false
	cleaner.canRun = true
}

package algorithm

import (
	"time"
	"gphxpaxos/config"
	"gphxpaxos/checkpoint"
	log "github.com/sirupsen/logrus"
	"os"
	"gphxpaxos/util"
	"math"
	"gphxpaxos/comm"
	"gphxpaxos/smbase"
)

const tmpBufferLen = 102400

type CheckpointSender struct {
	sendNodeId        uint64
	config            *config.Config
	learner           *Learner
	factory           *smbase.SMFac
	ckMnger           *checkpoint.CheckpointManager
	uuid              uint64
	sequence          uint64
	isEnd             bool
	isEnded           bool
	ackSequence       uint64
	absLastAckTime    uint64
	alreadySendedFile map[string]bool
	tmpBuffer         []byte
}

func NewCheckpointSender(sendNodeId uint64, config *config.Config, learner *Learner,
	factory *smbase.SMFac,
	ckmnger *checkpoint.CheckpointManager) *CheckpointSender {
	cksender := &CheckpointSender{
		sendNodeId: sendNodeId,
		config:     config,
		learner:    learner,
		factory:    factory,
		ckMnger:    ckmnger,
		uuid:       config.GetMyNodeId() ^ learner.GetInstanceId() + uint64(util.Rand(math.MaxInt32)),
		tmpBuffer:  make([]byte, tmpBufferLen),
	}

	util.StartRoutine(cksender.main)
	return cksender
}

func (checkpointSender *CheckpointSender) Stop() {
	if !checkpointSender.isEnded {
		checkpointSender.isEnd = true
	}
}

func (checkpointSender *CheckpointSender) IsEnd() bool {
	return checkpointSender.isEnded
}

func (checkpointSender *CheckpointSender) End() {
	checkpointSender.isEnd = true
}

func (checkpointSender *CheckpointSender) main() {
	checkpointSender.absLastAckTime = util.NowTimeMs()

	needContinue := false
	for !checkpointSender.ckMnger.GetRelayer().IsPaused() {
		if checkpointSender.isEnd {
			checkpointSender.isEnded = true
			return
		}

		needContinue = true
		checkpointSender.ckMnger.GetRelayer().Pause()
		log.Debug("wait replayer pause")
		util.SleepMs(100)
	}

	err := checkpointSender.LockCheckpoint()
	if err == nil {
		checkpointSender.SendCheckpoint()
		checkpointSender.UnlockCheckpoint()
	}

	if needContinue {
		checkpointSender.ckMnger.GetRelayer().Continue()
	}

	log.Info("Checkpoint.Sender [END]")
	checkpointSender.isEnded = true
}

func (checkpointSender *CheckpointSender) LockCheckpoint() error {
	smList := checkpointSender.factory.GetSMList()
	lockSmList := make([] smbase.StateMachine, 0)

	var err error
	for _, sm := range smList {
		err = sm.LockCheckpointState()
		if err != nil {
			break
		}
		lockSmList = append(lockSmList, sm)
	}

	if err != nil {
		for _, sm := range lockSmList {
			sm.UnLockCheckpointState()
		}
	}

	return err
}

func (checkpointSender *CheckpointSender) UnlockCheckpoint() {
	smList := checkpointSender.factory.GetSMList()

	for _, sm := range smList {
		sm.UnLockCheckpointState()
	}
}

func (checkpointSender *CheckpointSender) SendCheckpoint() error {
	learner := checkpointSender.learner
	err := learner.SendCheckpointBegin(checkpointSender.sendNodeId, checkpointSender.uuid, checkpointSender.sequence,
		checkpointSender.factory.GetCheckpointInstanceId(checkpointSender.config.GetMyGroupId()))
	if err != nil {
		log.Errorf("SendCheckpoint fail: %v \r\n", err)
		return err
	}

	checkpointSender.sequence += 1

	smList := checkpointSender.factory.GetSMList()
	for _, sm := range smList {
		err = checkpointSender.SendCheckpointForSM(sm)
		if err != nil {
			return err
		}
	}

	err = learner.SendCheckpointEnd(checkpointSender.sendNodeId, checkpointSender.uuid, checkpointSender.sequence,
		checkpointSender.factory.GetCheckpointInstanceId(checkpointSender.config.GetMyGroupId()))

	if err != nil {
		log.Errorf("SendCheckpointEnd fail: %v \r\n", err)
	}

	return err
}

func (checkpointSender *CheckpointSender) SendCheckpointForSM(statemachine smbase.StateMachine) error {
	var dirPath string
	var fileList = make([]string, 0)

	err := statemachine.GetCheckpointState(checkpointSender.config.GetMyGroupId(), &dirPath, fileList)
	if err != nil {
		return err
	}

	if len(dirPath) == 0 {
		return nil
	}

	if dirPath[len(dirPath)-1] != '/' {
		dirPath += "/"
	}

	for _, file := range fileList {
		err = checkpointSender.SendFile(statemachine, dirPath, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (checkpointSender *CheckpointSender) SendFile(statemachine smbase.StateMachine, dir string, file string) error {
	path := dir + file

	_, exist := checkpointSender.alreadySendedFile[path]
	if exist {
		return nil
	}

	fd, err := os.Open(path)
	if err != nil {
		return err
	}

	var offset uint64 = 0
	for {
		readLen, err := fd.Read(checkpointSender.tmpBuffer)

		if err != nil {
			fd.Close()
			return err
		}
		if readLen == 0 {
			break
		}

		err = checkpointSender.SendBuffer(statemachine.SMID(), statemachine.GetCheckpointInstanceId(checkpointSender.config.GetMyGroupId()),
			path, offset, checkpointSender.tmpBuffer, readLen)
		if err != nil {
			fd.Close()
			return err
		}

		if readLen < tmpBufferLen {
			return nil
		}

		offset += uint64(readLen)
	}

	checkpointSender.alreadySendedFile[path] = true
	fd.Close()
	return nil
}

func (checkpointSender *CheckpointSender) SendBuffer(smid int32, ckInstanceId uint64, file string,
	offser uint64, buffer []byte, bufLen int) error {
	ckSum := util.Crc32(0, buffer[:bufLen], comm.CRC32_SKIP)

	for {
		if checkpointSender.isEnd {
			return nil
		}

		err := checkpointSender.CheckAck(checkpointSender.sequence)
		if err != nil {
			return err
		}

		err = checkpointSender.learner.SendCheckpoint(checkpointSender.sendNodeId, checkpointSender.uuid, checkpointSender.sequence,
			ckInstanceId, ckSum, file, smid, offser, buffer)
		if err != nil {
			util.SleepMs(30000)
		} else {
			checkpointSender.sequence += 1
			break
		}
	}

	return nil
}

func (checkpointSender *CheckpointSender) Ack(sendNodeId uint64, uuid uint64, sequence uint64) {
	if sendNodeId != checkpointSender.sendNodeId {
		return
	}

	if checkpointSender.uuid != uuid {
		return
	}

	if checkpointSender.ackSequence != sequence {
		return
	}

	checkpointSender.ackSequence += 1
	checkpointSender.absLastAckTime = util.NowTimeMs()
}

func (checkpointSender *CheckpointSender) CheckAck(sendSequence uint64) error {
	for sendSequence > checkpointSender.ackSequence+100 {
		now := util.NowTimeMs()
		var passTime uint64
		if now > checkpointSender.absLastAckTime {
			passTime = now - checkpointSender.absLastAckTime
		}

		if checkpointSender.isEnd {

		}

		if passTime > 200 {

		}

		time.Sleep(20 * time.Millisecond)
	}

	return nil
}

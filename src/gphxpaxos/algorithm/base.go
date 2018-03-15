package algorithm

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"encoding/binary"
	"gphxpaxos/config"
	"gphxpaxos/network"
	"github.com/golang/protobuf/proto"
	"gphxpaxos/comm"
	"gphxpaxos/util"
)

//------------------------------------BallotNumber-----------------------------//

type BallotNumber struct {
	proposalId uint64
	nodeId     uint64
}

func NewBallotNumber(proposalId uint64, nodeId uint64) *BallotNumber {
	return &BallotNumber{
		proposalId: proposalId,
		nodeId:     nodeId,
	}
}

func (ballotNumber *BallotNumber) String() string {
	return fmt.Sprintf("%d:%d", ballotNumber.proposalId, ballotNumber.nodeId)
}

// >=
func (ballotNumber *BallotNumber) GE(other *BallotNumber) bool {
	if ballotNumber.proposalId == other.proposalId {
		return ballotNumber.nodeId >= other.nodeId
	}

	return ballotNumber.proposalId >= other.proposalId
}

// !=
func (ballotNumber *BallotNumber) NE(other *BallotNumber) bool {
	return ballotNumber.proposalId != other.proposalId ||
		ballotNumber.nodeId != other.nodeId
}

// ==
func (ballotNumber *BallotNumber) EQ(other *BallotNumber) bool {
	return !ballotNumber.NE(other)
}

// GT
func (ballotNumber *BallotNumber) GT(other *BallotNumber) bool {
	if ballotNumber.proposalId == other.proposalId {
		return ballotNumber.nodeId > other.nodeId
	}

	return ballotNumber.proposalId > other.proposalId
}

func (ballotNumber *BallotNumber) IsNull() bool {
	return ballotNumber.proposalId == 0
}

func (ballotNumber *BallotNumber) Clone(bn *BallotNumber) {
	ballotNumber.nodeId = bn.nodeId
	ballotNumber.proposalId = bn.proposalId
}

func (ballotNumber *BallotNumber) Reset() {
	ballotNumber.nodeId = 0
	ballotNumber.proposalId = 0
}

//-----------------------------------------------Base-------------------------------------------------//

const (
	BroadcastMessage_Type_RunSelf_First = 1
	BroadcastMessage_Type_RunSelf_Final = 2
	BroadcastMessage_Type_RunSelf_None  = 3
)

var GROUPIDXLEN = binary.Size(int(0))
var HEADLEN_LEN = binary.Size(uint16(0))
var CHECKSUM_LEN = binary.Size(uint32(0))

type Base struct {
	instanceId uint64
	config     *config.Config
	transport  network.MsgTransport
	instance   *Instance
	isTestMode bool
}

func init() {
	HEADLEN_LEN =  binary.Size(uint16(0))
	CHECKSUM_LEN =  binary.Size(uint32(0))
}

func newBase(instance *Instance) Base {
	var instanceId uint64 = 1
	maxInstanceId, err := instance.logStorage.GetMaxInstanceID()
	if err == nil {
		instanceId = maxInstanceId + 1
	}

	return Base{
		config:     instance.config,
		transport:  instance.transport,
		instance:   instance,
		instanceId: instanceId,
		isTestMode: false,
	}
}

func (base *Base) GetInstanceId() uint64 {
	return base.instanceId
}

func (base *Base) setInstanceId(instanceId uint64) {
	base.instanceId = instanceId
}

func (base *Base) newInstance() {
	base.instanceId++
}

func (base *Base) GetLastChecksum() uint32 {
	return base.instance.GetLastChecksum()
}

func (base *Base) packPaxosMsg(paxosMsg *comm.PaxosMsg) ([]byte, *comm.Header, error) {
	body, err := proto.Marshal(paxosMsg)
	if err != nil {
		log.Error("paxos msg Marshal fail:%v", err)
		return nil, nil, err
	}

	return base.packBaseMsg(body, comm.MsgCmd_PaxosMsg)
}

func (base *Base) packCheckpointMsg(msg *comm.CheckpointMsg) ([]byte, *comm.Header, error) {
	body, err := proto.Marshal(msg)
	if err != nil {
		log.Error("checkpoint msg Marshal fail:%v", err)
		return nil, nil, err
	}

	return base.packBaseMsg(body, comm.MsgCmd_CheckpointMsg)
}

// format: groupId(int) + headerlen(uint16) + header + body + crc32 checksum(uint32)
func (base *Base) packBaseMsg(body []byte, cmd int32) (buffer []byte, header *comm.Header, err error) {
	groupIdx := base.config.GetMyGroupId()

	h := &comm.Header{
		Cmdid: proto.Int32(cmd),
		// buffer len + checksum len
		Gid:     proto.Uint64(base.config.GetGid()),
		Rid:     proto.Uint64(0),
		Version: proto.Int32(comm.Version),
	}
	header = h

	headerBuf, err := proto.Marshal(header)
	if err != nil {
		log.Error("header Marshal fail:%v", err)
		return
	}

	groupIdxBuf := make([]byte, GROUPIDXLEN)
	util.EncodeInt32(groupIdxBuf, 0, int32(groupIdx))

	headerLenBuf := make([] byte, HEADLEN_LEN)
	util.EncodeUint16(headerLenBuf, 0, uint16(len(headerBuf)))

	buffer = util.AppendBytes(groupIdxBuf, headerLenBuf, headerBuf, body)

	ckSum := util.Crc32(0, buffer, comm.NET_CRC32SKIP)
	cksumBuf := make([]byte, CHECKSUM_LEN)
	util.EncodeUint32(cksumBuf, 0, ckSum)

	buffer = util.AppendBytes(buffer, cksumBuf)

	return
}

// TODO to be checked
func (base *Base) unpackBaseMsg(buffer []byte, header *comm.Header) (body []byte, err error) {

	headStartPos := GROUPIDXLEN + HEADLEN_LEN

	var bufferLen = len(buffer)

	if bufferLen < headStartPos {
		log.Error("no head")
		err = comm.ErrInvalidMsg
		return
	}

	var headLen uint16
	util.DecodeUint16(buffer, GROUPIDXLEN, &headLen)

	if bufferLen < headStartPos+int(headLen) {
		log.Error("msg head lost ")
		err = comm.ErrInvalidMsg
		return
	}

	bodyStartPos := headStartPos + int(headLen)

	proto.Unmarshal(buffer[headStartPos:bodyStartPos], header)

	if bodyStartPos + CHECKSUM_LEN > bufferLen {
		log.Error("no checksum, body start pos %d, buffersize %d", bodyStartPos, bufferLen)
		err = comm.ErrInvalidMsg
		return
	}

	var cksum uint32
	util.DecodeUint32(buffer, bufferLen-CHECKSUM_LEN, &cksum)

	calCksum := util.Crc32(0, buffer[:bufferLen-CHECKSUM_LEN], comm.NET_CRC32SKIP)
	if calCksum != cksum {
		log.Error("data bring cksum %d not equal to cal cksum %d", cksum, calCksum)
		err = comm.ErrInvalidMsg
		return
	}

	body = buffer[bodyStartPos : bufferLen-CHECKSUM_LEN]
	err = nil
	return
}

func (base *Base) sendCheckpointMessage(sendToNodeId uint64, msg *comm.CheckpointMsg, sendType int) error {
	if sendToNodeId == base.config.GetMyNodeId() {
		return nil
	}

	buffer, _, err := base.packCheckpointMsg(msg)
	if err != nil {
		return err
	}

	return base.transport.SendMessage(base.config.GetMyGroupId(), sendToNodeId, buffer, sendType)
}

func (base *Base) sendPaxosMessage(sendToNodeId uint64, msg *comm.PaxosMsg, sendType int) error {
	if sendToNodeId == base.config.GetMyNodeId() {
		base.instance.OnReceivePaxosMsg(msg, false)
		return nil
	}

	buffer, _, err := base.packPaxosMsg(msg)
	if err != nil {
		log.Error("pack paxos msg error %v", err)
		return err
	}

	return base.transport.SendMessage(base.config.GetMyGroupId(), sendToNodeId, buffer, sendType)
}

func (base *Base) broadcastMessage(msg *comm.PaxosMsg, runType int, sendType int) error {
	if base.isTestMode {
		return nil
	}

	if runType == BroadcastMessage_Type_RunSelf_First {
		err := base.instance.OnReceivePaxosMsg(msg, false)
		if err != nil {
			return err
		}
	}

	buffer, _, err := base.packPaxosMsg(msg)
	if err != nil {
		return err
	}

	err = base.transport.BroadcastMessage(base.config.GetMyGroupId(), buffer, sendType)

	if runType == BroadcastMessage_Type_RunSelf_Final {
		base.instance.OnReceivePaxosMsg(msg, false)
	}

	return err
}


func (base *Base) BroadcastMessageToFollower(msg *comm.PaxosMsg, sendType int) error {

  value, _, err := base.packPaxosMsg(msg)
  if err != nil {
    return err
  }

  return base.transport.BroadcastMessageFollower(base.config.GetMyGroupId(), value, sendType)
}

func (base *Base) BroadcastMessageToTempNode(msg *comm.PaxosMsg, sendType int) error {

  value, _, err := base.packPaxosMsg(msg)
	  if err != nil {
    return err
  }

  return base.transport.BroadcastMessageTempNode(base.config.GetMyGroupId(), value, sendType)
}


func (base *Base) setAsTestMode() {
	base.isTestMode = true
}

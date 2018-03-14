package comm

import (
	"encoding/binary"
	"errors"
	"math"
)

const (
	CRC32_SKIP    = 0
	NET_CRC32SKIP = 0
)

const (
	MAX_QUEUE_MEM_SIZE = 209715200

	// enum MsgCmd
	MsgCmd_PaxosMsg      = 1
	MsgCmd_CheckpointMsg = 2

	// enum PaxosMsgType
	MsgType_PaxosPrepare                     = 1
	MsgType_PaxosPrepareReply                = 2
	MsgType_PaxosAccept                      = 3
	MsgType_PaxosAcceptReply                 = 4
	MsgType_PaxosLearner_AskforLearn         = 5
	MsgType_PaxosLearner_SendLearnValue      = 6
	MsgType_PaxosLearner_ProposerSendSuccess = 7
	MsgType_PaxosProposal_SendNewValue       = 8
	MsgType_PaxosLearner_SendNowInstanceID   = 9
	MsgType_PaxosLearner_ComfirmAskforLearn  = 10
	MsgType_PaxosLearner_SendLearnValue_Ack  = 11
	MsgType_PaxosLearner_AskforCheckpoint    = 12
	MsgType_PaxosLearner_OnAskforCheckpoint  = 13

	// enum PaxosMsgFlagType
	PaxosMsgFlagType_SendLearnValue_NeedAck = 1

	// enum CheckpointMsgType
	CheckpointMsgType_SendFile     = 1
	CheckpointMsgType_SendFile_Ack = 2

	//enum CheckpointSendFileFlag
	CheckpointSendFileFlag_BEGIN = 1
	CheckpointSendFileFlag_ING   = 2
	CheckpointSendFileFlag_END   = 3

	//enum CheckpointSendFileAckFlag

	CheckpointSendFileAckFlag_OK   = 1
	CheckpointSendFileAckFlag_Fail = 2

	//enum TimerType

	Timer_Proposer_Prepare_Timeout = 1
	Timer_Proposer_Accept_Timeout  = 2
	Timer_Learner_Askforlearn_noop = 3
	Timer_Instance_Commit_Timeout  = 4
)

var (
	INT32SIZE  = binary.Size(int32(0))
	INT64SIZE  = binary.Size(int64(0))
	UINT32SIZE = binary.Size(uint32(0))
	UINT64SIZE = binary.Size(uint64(0))
	UINT16SIZE = binary.Size(uint16(0))
)

// 自定义错误
var (
	ErrKeyNotFound       = errors.New("key not found")
	ErrGetFail           = errors.New("get fail")
	ErrInvalidGroupIndex = errors.New("invalid group index")
	ErrInvalidInstanceId = errors.New("invalid instanceid")
	ErrInvalidMetaFileId = errors.New("invalid meta file id")
	ErrInvalidMsg        = errors.New("invalid msg")
	ErrFileNotExist      = errors.New("file not exist")
	ErrDbNotInit         = errors.New("db not init yet")
	ErrNodeNotFound      = errors.New("node not found")
	ErrWriteFileFail     = errors.New("write file fail")
)

var INVALID_INSTANCEID uint64 = math.MaxUint64
var NULL_NODEID uint64 = math.MaxUint64 - 1

var UINT64_0 = uint64(0)

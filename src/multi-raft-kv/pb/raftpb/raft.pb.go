// Code generated by protoc-gen-go. DO NOT EDIT.
// source: raft/raft.proto

/*
Package raftpb is a generated protocol buffer package.

It is generated from these files:
	raft/raft.proto

It has these top-level messages:
	Entry
	SnapshotMetadata
	Snapshot
	Message
	HardState
	ConfState
	ConfChange
	RaftMessage
	StoreIdent
	ACKMessage
	SnapshotMessageHeader
	SnapshotMessage
	SnapshotChunkMessage
	SnapshotAckMessage
	SnapshotAskMessage
	CellLocalState
	RaftLocalState
	RaftTruncatedState
	RaftApplyState
*/
package raftpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import metapb "metapb"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type EntryType int32

const (
	EntryType_EntryNormal     EntryType = 0
	EntryType_EntryConfChange EntryType = 1
)

var EntryType_name = map[int32]string{
	0: "EntryNormal",
	1: "EntryConfChange",
}
var EntryType_value = map[string]int32{
	"EntryNormal":     0,
	"EntryConfChange": 1,
}

func (x EntryType) Enum() *EntryType {
	p := new(EntryType)
	*p = x
	return p
}
func (x EntryType) String() string {
	return proto.EnumName(EntryType_name, int32(x))
}
func (x *EntryType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(EntryType_value, data, "EntryType")
	if err != nil {
		return err
	}
	*x = EntryType(value)
	return nil
}
func (EntryType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type MessageType int32

const (
	MessageType_MsgHup            MessageType = 0
	MessageType_MsgBeat           MessageType = 1
	MessageType_MsgProp           MessageType = 2
	MessageType_MsgApp            MessageType = 3
	MessageType_MsgAppResp        MessageType = 4
	MessageType_MsgVote           MessageType = 5
	MessageType_MsgVoteResp       MessageType = 6
	MessageType_MsgSnap           MessageType = 7
	MessageType_MsgHeartbeat      MessageType = 8
	MessageType_MsgHeartbeatResp  MessageType = 9
	MessageType_MsgUnreachable    MessageType = 10
	MessageType_MsgSnapStatus     MessageType = 11
	MessageType_MsgCheckQuorum    MessageType = 12
	MessageType_MsgTransferLeader MessageType = 13
	MessageType_MsgTimeoutNow     MessageType = 14
	MessageType_MsgReadIndex      MessageType = 15
	MessageType_MsgReadIndexResp  MessageType = 16
	MessageType_MsgPreVote        MessageType = 17
	MessageType_MsgPreVoteResp    MessageType = 18
)

var MessageType_name = map[int32]string{
	0:  "MsgHup",
	1:  "MsgBeat",
	2:  "MsgProp",
	3:  "MsgApp",
	4:  "MsgAppResp",
	5:  "MsgVote",
	6:  "MsgVoteResp",
	7:  "MsgSnap",
	8:  "MsgHeartbeat",
	9:  "MsgHeartbeatResp",
	10: "MsgUnreachable",
	11: "MsgSnapStatus",
	12: "MsgCheckQuorum",
	13: "MsgTransferLeader",
	14: "MsgTimeoutNow",
	15: "MsgReadIndex",
	16: "MsgReadIndexResp",
	17: "MsgPreVote",
	18: "MsgPreVoteResp",
}
var MessageType_value = map[string]int32{
	"MsgHup":            0,
	"MsgBeat":           1,
	"MsgProp":           2,
	"MsgApp":            3,
	"MsgAppResp":        4,
	"MsgVote":           5,
	"MsgVoteResp":       6,
	"MsgSnap":           7,
	"MsgHeartbeat":      8,
	"MsgHeartbeatResp":  9,
	"MsgUnreachable":    10,
	"MsgSnapStatus":     11,
	"MsgCheckQuorum":    12,
	"MsgTransferLeader": 13,
	"MsgTimeoutNow":     14,
	"MsgReadIndex":      15,
	"MsgReadIndexResp":  16,
	"MsgPreVote":        17,
	"MsgPreVoteResp":    18,
}

func (x MessageType) Enum() *MessageType {
	p := new(MessageType)
	*p = x
	return p
}
func (x MessageType) String() string {
	return proto.EnumName(MessageType_name, int32(x))
}
func (x *MessageType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(MessageType_value, data, "MessageType")
	if err != nil {
		return err
	}
	*x = MessageType(value)
	return nil
}
func (MessageType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type ConfChangeType int32

const (
	ConfChangeType_ConfChangeAddNode    ConfChangeType = 0
	ConfChangeType_ConfChangeRemoveNode ConfChangeType = 1
	ConfChangeType_ConfChangeUpdateNode ConfChangeType = 2
)

var ConfChangeType_name = map[int32]string{
	0: "ConfChangeAddNode",
	1: "ConfChangeRemoveNode",
	2: "ConfChangeUpdateNode",
}
var ConfChangeType_value = map[string]int32{
	"ConfChangeAddNode":    0,
	"ConfChangeRemoveNode": 1,
	"ConfChangeUpdateNode": 2,
}

func (x ConfChangeType) Enum() *ConfChangeType {
	p := new(ConfChangeType)
	*p = x
	return p
}
func (x ConfChangeType) String() string {
	return proto.EnumName(ConfChangeType_name, int32(x))
}
func (x *ConfChangeType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(ConfChangeType_value, data, "ConfChangeType")
	if err != nil {
		return err
	}
	*x = ConfChangeType(value)
	return nil
}
func (ConfChangeType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type SnapshotState int32

const (
	SnapshotState_Accept   SnapshotState = 0
	SnapshotState_Reject   SnapshotState = 1
	SnapshotState_Received SnapshotState = 2
)

var SnapshotState_name = map[int32]string{
	0: "Accept",
	1: "Reject",
	2: "Received",
}
var SnapshotState_value = map[string]int32{
	"Accept":   0,
	"Reject":   1,
	"Received": 2,
}

func (x SnapshotState) Enum() *SnapshotState {
	p := new(SnapshotState)
	*p = x
	return p
}
func (x SnapshotState) String() string {
	return proto.EnumName(SnapshotState_name, int32(x))
}
func (x *SnapshotState) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(SnapshotState_value, data, "SnapshotState")
	if err != nil {
		return err
	}
	*x = SnapshotState(value)
	return nil
}
func (SnapshotState) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type PeerState int32

const (
	PeerState_Normal    PeerState = 0
	PeerState_Applying  PeerState = 1
	PeerState_Tombstone PeerState = 2
)

var PeerState_name = map[int32]string{
	0: "Normal",
	1: "Applying",
	2: "Tombstone",
}
var PeerState_value = map[string]int32{
	"Normal":    0,
	"Applying":  1,
	"Tombstone": 2,
}

func (x PeerState) Enum() *PeerState {
	p := new(PeerState)
	*p = x
	return p
}
func (x PeerState) String() string {
	return proto.EnumName(PeerState_name, int32(x))
}
func (x *PeerState) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(PeerState_value, data, "PeerState")
	if err != nil {
		return err
	}
	*x = PeerState(value)
	return nil
}
func (PeerState) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type Entry struct {
	Term             *uint64    `protobuf:"varint,2,opt,name=Term" json:"Term,omitempty"`
	Index            *uint64    `protobuf:"varint,3,opt,name=Index" json:"Index,omitempty"`
	Type             *EntryType `protobuf:"varint,1,opt,name=Type,enum=raftpb.EntryType" json:"Type,omitempty"`
	Data             []byte     `protobuf:"bytes,4,opt,name=Data" json:"Data,omitempty"`
	XXX_unrecognized []byte     `json:"-"`
}

func (m *Entry) Reset()                    { *m = Entry{} }
func (m *Entry) String() string            { return proto.CompactTextString(m) }
func (*Entry) ProtoMessage()               {}
func (*Entry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Entry) GetTerm() uint64 {
	if m != nil && m.Term != nil {
		return *m.Term
	}
	return 0
}

func (m *Entry) GetIndex() uint64 {
	if m != nil && m.Index != nil {
		return *m.Index
	}
	return 0
}

func (m *Entry) GetType() EntryType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return EntryType_EntryNormal
}

func (m *Entry) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type SnapshotMetadata struct {
	ConfState        *ConfState `protobuf:"bytes,1,opt,name=conf_state,json=confState" json:"conf_state,omitempty"`
	Index            *uint64    `protobuf:"varint,2,opt,name=index" json:"index,omitempty"`
	Term             *uint64    `protobuf:"varint,3,opt,name=term" json:"term,omitempty"`
	XXX_unrecognized []byte     `json:"-"`
}

func (m *SnapshotMetadata) Reset()                    { *m = SnapshotMetadata{} }
func (m *SnapshotMetadata) String() string            { return proto.CompactTextString(m) }
func (*SnapshotMetadata) ProtoMessage()               {}
func (*SnapshotMetadata) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *SnapshotMetadata) GetConfState() *ConfState {
	if m != nil {
		return m.ConfState
	}
	return nil
}

func (m *SnapshotMetadata) GetIndex() uint64 {
	if m != nil && m.Index != nil {
		return *m.Index
	}
	return 0
}

func (m *SnapshotMetadata) GetTerm() uint64 {
	if m != nil && m.Term != nil {
		return *m.Term
	}
	return 0
}

type Snapshot struct {
	Data             []byte            `protobuf:"bytes,1,opt,name=data" json:"data,omitempty"`
	Metadata         *SnapshotMetadata `protobuf:"bytes,2,opt,name=metadata" json:"metadata,omitempty"`
	XXX_unrecognized []byte            `json:"-"`
}

func (m *Snapshot) Reset()                    { *m = Snapshot{} }
func (m *Snapshot) String() string            { return proto.CompactTextString(m) }
func (*Snapshot) ProtoMessage()               {}
func (*Snapshot) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Snapshot) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Snapshot) GetMetadata() *SnapshotMetadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type Message struct {
	Type             *MessageType `protobuf:"varint,1,opt,name=type,enum=raftpb.MessageType" json:"type,omitempty"`
	To               *uint64      `protobuf:"varint,2,opt,name=to" json:"to,omitempty"`
	From             *uint64      `protobuf:"varint,3,opt,name=from" json:"from,omitempty"`
	Term             *uint64      `protobuf:"varint,4,opt,name=term" json:"term,omitempty"`
	LogTerm          *uint64      `protobuf:"varint,5,opt,name=logTerm" json:"logTerm,omitempty"`
	Index            *uint64      `protobuf:"varint,6,opt,name=index" json:"index,omitempty"`
	Entries          []*Entry     `protobuf:"bytes,7,rep,name=entries" json:"entries,omitempty"`
	Commit           *uint64      `protobuf:"varint,8,opt,name=commit" json:"commit,omitempty"`
	Snapshot         *Snapshot    `protobuf:"bytes,9,opt,name=snapshot" json:"snapshot,omitempty"`
	Reject           *bool        `protobuf:"varint,10,opt,name=reject" json:"reject,omitempty"`
	RejectHint       *uint64      `protobuf:"varint,11,opt,name=rejectHint" json:"rejectHint,omitempty"`
	Context          []byte       `protobuf:"bytes,12,opt,name=context" json:"context,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Message) GetType() MessageType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return MessageType_MsgHup
}

func (m *Message) GetTo() uint64 {
	if m != nil && m.To != nil {
		return *m.To
	}
	return 0
}

func (m *Message) GetFrom() uint64 {
	if m != nil && m.From != nil {
		return *m.From
	}
	return 0
}

func (m *Message) GetTerm() uint64 {
	if m != nil && m.Term != nil {
		return *m.Term
	}
	return 0
}

func (m *Message) GetLogTerm() uint64 {
	if m != nil && m.LogTerm != nil {
		return *m.LogTerm
	}
	return 0
}

func (m *Message) GetIndex() uint64 {
	if m != nil && m.Index != nil {
		return *m.Index
	}
	return 0
}

func (m *Message) GetEntries() []*Entry {
	if m != nil {
		return m.Entries
	}
	return nil
}

func (m *Message) GetCommit() uint64 {
	if m != nil && m.Commit != nil {
		return *m.Commit
	}
	return 0
}

func (m *Message) GetSnapshot() *Snapshot {
	if m != nil {
		return m.Snapshot
	}
	return nil
}

func (m *Message) GetReject() bool {
	if m != nil && m.Reject != nil {
		return *m.Reject
	}
	return false
}

func (m *Message) GetRejectHint() uint64 {
	if m != nil && m.RejectHint != nil {
		return *m.RejectHint
	}
	return 0
}

func (m *Message) GetContext() []byte {
	if m != nil {
		return m.Context
	}
	return nil
}

type HardState struct {
	Term             *uint64 `protobuf:"varint,1,opt,name=term" json:"term,omitempty"`
	Vote             *uint64 `protobuf:"varint,2,opt,name=vote" json:"vote,omitempty"`
	Commit           *uint64 `protobuf:"varint,3,opt,name=commit" json:"commit,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *HardState) Reset()                    { *m = HardState{} }
func (m *HardState) String() string            { return proto.CompactTextString(m) }
func (*HardState) ProtoMessage()               {}
func (*HardState) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *HardState) GetTerm() uint64 {
	if m != nil && m.Term != nil {
		return *m.Term
	}
	return 0
}

func (m *HardState) GetVote() uint64 {
	if m != nil && m.Vote != nil {
		return *m.Vote
	}
	return 0
}

func (m *HardState) GetCommit() uint64 {
	if m != nil && m.Commit != nil {
		return *m.Commit
	}
	return 0
}

type ConfState struct {
	Nodes            []uint64 `protobuf:"varint,1,rep,name=nodes" json:"nodes,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *ConfState) Reset()                    { *m = ConfState{} }
func (m *ConfState) String() string            { return proto.CompactTextString(m) }
func (*ConfState) ProtoMessage()               {}
func (*ConfState) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ConfState) GetNodes() []uint64 {
	if m != nil {
		return m.Nodes
	}
	return nil
}

type ConfChange struct {
	ID               *uint64         `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
	Type             *ConfChangeType `protobuf:"varint,2,opt,name=Type,enum=raftpb.ConfChangeType" json:"Type,omitempty"`
	NodeID           *uint64         `protobuf:"varint,3,opt,name=NodeID" json:"NodeID,omitempty"`
	Context          []byte          `protobuf:"bytes,4,opt,name=Context" json:"Context,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *ConfChange) Reset()                    { *m = ConfChange{} }
func (m *ConfChange) String() string            { return proto.CompactTextString(m) }
func (*ConfChange) ProtoMessage()               {}
func (*ConfChange) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *ConfChange) GetID() uint64 {
	if m != nil && m.ID != nil {
		return *m.ID
	}
	return 0
}

func (m *ConfChange) GetType() ConfChangeType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return ConfChangeType_ConfChangeAddNode
}

func (m *ConfChange) GetNodeID() uint64 {
	if m != nil && m.NodeID != nil {
		return *m.NodeID
	}
	return 0
}

func (m *ConfChange) GetContext() []byte {
	if m != nil {
		return m.Context
	}
	return nil
}

type RaftMessage struct {
	CellID    *uint64           `protobuf:"varint,1,opt,name=cellID" json:"cellID,omitempty"`
	FromPeer  *metapb.Peer      `protobuf:"bytes,2,opt,name=fromPeer" json:"fromPeer,omitempty"`
	ToPeer    *metapb.Peer      `protobuf:"bytes,3,opt,name=toPeer" json:"toPeer,omitempty"`
	Message   *Message          `protobuf:"bytes,4,opt,name=message" json:"message,omitempty"`
	CellEpoch *metapb.CellEpoch `protobuf:"bytes,5,opt,name=cellEpoch" json:"cellEpoch,omitempty"`
	// true means to_peer is a tombstone peer and it should remove itself.
	IsTombstone *bool `protobuf:"varint,6,opt,name=isTombstone" json:"isTombstone,omitempty"`
	// Cell key range [start_key, end_key).
	Start            []byte `protobuf:"bytes,7,opt,name=start" json:"start,omitempty"`
	End              []byte `protobuf:"bytes,8,opt,name=end" json:"end,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *RaftMessage) Reset()                    { *m = RaftMessage{} }
func (m *RaftMessage) String() string            { return proto.CompactTextString(m) }
func (*RaftMessage) ProtoMessage()               {}
func (*RaftMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *RaftMessage) GetCellID() uint64 {
	if m != nil && m.CellID != nil {
		return *m.CellID
	}
	return 0
}

func (m *RaftMessage) GetFromPeer() *metapb.Peer {
	if m != nil {
		return m.FromPeer
	}
	return nil
}

func (m *RaftMessage) GetToPeer() *metapb.Peer {
	if m != nil {
		return m.ToPeer
	}
	return nil
}

func (m *RaftMessage) GetMessage() *Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (m *RaftMessage) GetCellEpoch() *metapb.CellEpoch {
	if m != nil {
		return m.CellEpoch
	}
	return nil
}

func (m *RaftMessage) GetIsTombstone() bool {
	if m != nil && m.IsTombstone != nil {
		return *m.IsTombstone
	}
	return false
}

func (m *RaftMessage) GetStart() []byte {
	if m != nil {
		return m.Start
	}
	return nil
}

func (m *RaftMessage) GetEnd() []byte {
	if m != nil {
		return m.End
	}
	return nil
}

type StoreIdent struct {
	ClusterID        *uint64 `protobuf:"varint,1,opt,name=clusterID" json:"clusterID,omitempty"`
	StoreID          *uint64 `protobuf:"varint,2,opt,name=storeID" json:"storeID,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *StoreIdent) Reset()                    { *m = StoreIdent{} }
func (m *StoreIdent) String() string            { return proto.CompactTextString(m) }
func (*StoreIdent) ProtoMessage()               {}
func (*StoreIdent) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *StoreIdent) GetClusterID() uint64 {
	if m != nil && m.ClusterID != nil {
		return *m.ClusterID
	}
	return 0
}

func (m *StoreIdent) GetStoreID() uint64 {
	if m != nil && m.StoreID != nil {
		return *m.StoreID
	}
	return 0
}

type ACKMessage struct {
	Seq              *uint64 `protobuf:"varint,1,opt,name=seq" json:"seq,omitempty"`
	To               *uint64 `protobuf:"varint,2,opt,name=to" json:"to,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *ACKMessage) Reset()                    { *m = ACKMessage{} }
func (m *ACKMessage) String() string            { return proto.CompactTextString(m) }
func (*ACKMessage) ProtoMessage()               {}
func (*ACKMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *ACKMessage) GetSeq() uint64 {
	if m != nil && m.Seq != nil {
		return *m.Seq
	}
	return 0
}

func (m *ACKMessage) GetTo() uint64 {
	if m != nil && m.To != nil {
		return *m.To
	}
	return 0
}

type SnapshotMessageHeader struct {
	Cell             *metapb.Cell `protobuf:"bytes,1,opt,name=cell" json:"cell,omitempty"`
	FromPeer         *metapb.Peer `protobuf:"bytes,2,opt,name=fromPeer" json:"fromPeer,omitempty"`
	ToPeer           *metapb.Peer `protobuf:"bytes,3,opt,name=toPeer" json:"toPeer,omitempty"`
	Term             *uint64      `protobuf:"varint,4,opt,name=term" json:"term,omitempty"`
	Index            *uint64      `protobuf:"varint,5,opt,name=index" json:"index,omitempty"`
	Seq              *uint64      `protobuf:"varint,6,opt,name=seq" json:"seq,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *SnapshotMessageHeader) Reset()                    { *m = SnapshotMessageHeader{} }
func (m *SnapshotMessageHeader) String() string            { return proto.CompactTextString(m) }
func (*SnapshotMessageHeader) ProtoMessage()               {}
func (*SnapshotMessageHeader) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *SnapshotMessageHeader) GetCell() *metapb.Cell {
	if m != nil {
		return m.Cell
	}
	return nil
}

func (m *SnapshotMessageHeader) GetFromPeer() *metapb.Peer {
	if m != nil {
		return m.FromPeer
	}
	return nil
}

func (m *SnapshotMessageHeader) GetToPeer() *metapb.Peer {
	if m != nil {
		return m.ToPeer
	}
	return nil
}

func (m *SnapshotMessageHeader) GetTerm() uint64 {
	if m != nil && m.Term != nil {
		return *m.Term
	}
	return 0
}

func (m *SnapshotMessageHeader) GetIndex() uint64 {
	if m != nil && m.Index != nil {
		return *m.Index
	}
	return 0
}

func (m *SnapshotMessageHeader) GetSeq() uint64 {
	if m != nil && m.Seq != nil {
		return *m.Seq
	}
	return 0
}

type SnapshotMessage struct {
	Header           *SnapshotMessageHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Chunk            *SnapshotChunkMessage  `protobuf:"bytes,2,opt,name=chunk" json:"chunk,omitempty"`
	Ack              *SnapshotAckMessage    `protobuf:"bytes,3,opt,name=ack" json:"ack,omitempty"`
	Ask              *SnapshotAskMessage    `protobuf:"bytes,4,opt,name=ask" json:"ask,omitempty"`
	XXX_unrecognized []byte                 `json:"-"`
}

func (m *SnapshotMessage) Reset()                    { *m = SnapshotMessage{} }
func (m *SnapshotMessage) String() string            { return proto.CompactTextString(m) }
func (*SnapshotMessage) ProtoMessage()               {}
func (*SnapshotMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *SnapshotMessage) GetHeader() *SnapshotMessageHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *SnapshotMessage) GetChunk() *SnapshotChunkMessage {
	if m != nil {
		return m.Chunk
	}
	return nil
}

func (m *SnapshotMessage) GetAck() *SnapshotAckMessage {
	if m != nil {
		return m.Ack
	}
	return nil
}

func (m *SnapshotMessage) GetAsk() *SnapshotAskMessage {
	if m != nil {
		return m.Ask
	}
	return nil
}

type SnapshotChunkMessage struct {
	Data             []byte  `protobuf:"bytes,1,opt,name=data" json:"data,omitempty"`
	First            *bool   `protobuf:"varint,2,opt,name=first" json:"first,omitempty"`
	Last             *bool   `protobuf:"varint,3,opt,name=last" json:"last,omitempty"`
	FileSize         *uint64 `protobuf:"varint,4,opt,name=fileSize" json:"fileSize,omitempty"`
	CheckSum         *uint64 `protobuf:"varint,5,opt,name=checkSum" json:"checkSum,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *SnapshotChunkMessage) Reset()                    { *m = SnapshotChunkMessage{} }
func (m *SnapshotChunkMessage) String() string            { return proto.CompactTextString(m) }
func (*SnapshotChunkMessage) ProtoMessage()               {}
func (*SnapshotChunkMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *SnapshotChunkMessage) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *SnapshotChunkMessage) GetFirst() bool {
	if m != nil && m.First != nil {
		return *m.First
	}
	return false
}

func (m *SnapshotChunkMessage) GetLast() bool {
	if m != nil && m.Last != nil {
		return *m.Last
	}
	return false
}

func (m *SnapshotChunkMessage) GetFileSize() uint64 {
	if m != nil && m.FileSize != nil {
		return *m.FileSize
	}
	return 0
}

func (m *SnapshotChunkMessage) GetCheckSum() uint64 {
	if m != nil && m.CheckSum != nil {
		return *m.CheckSum
	}
	return 0
}

type SnapshotAckMessage struct {
	Ack              *SnapshotState `protobuf:"varint,1,opt,name=ack,enum=raftpb.SnapshotState" json:"ack,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *SnapshotAckMessage) Reset()                    { *m = SnapshotAckMessage{} }
func (m *SnapshotAckMessage) String() string            { return proto.CompactTextString(m) }
func (*SnapshotAckMessage) ProtoMessage()               {}
func (*SnapshotAckMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *SnapshotAckMessage) GetAck() SnapshotState {
	if m != nil && m.Ack != nil {
		return *m.Ack
	}
	return SnapshotState_Accept
}

type SnapshotAskMessage struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *SnapshotAskMessage) Reset()                    { *m = SnapshotAskMessage{} }
func (m *SnapshotAskMessage) String() string            { return proto.CompactTextString(m) }
func (*SnapshotAskMessage) ProtoMessage()               {}
func (*SnapshotAskMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

type CellLocalState struct {
	State            *PeerState   `protobuf:"varint,1,opt,name=state,enum=raftpb.PeerState" json:"state,omitempty"`
	Cell             *metapb.Cell `protobuf:"bytes,2,opt,name=cell" json:"cell,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *CellLocalState) Reset()                    { *m = CellLocalState{} }
func (m *CellLocalState) String() string            { return proto.CompactTextString(m) }
func (*CellLocalState) ProtoMessage()               {}
func (*CellLocalState) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

func (m *CellLocalState) GetState() PeerState {
	if m != nil && m.State != nil {
		return *m.State
	}
	return PeerState_Normal
}

func (m *CellLocalState) GetCell() *metapb.Cell {
	if m != nil {
		return m.Cell
	}
	return nil
}

type RaftLocalState struct {
	HardState        *HardState `protobuf:"bytes,1,opt,name=hardState" json:"hardState,omitempty"`
	LastIndex        *uint64    `protobuf:"varint,2,opt,name=lastIndex" json:"lastIndex,omitempty"`
	XXX_unrecognized []byte     `json:"-"`
}

func (m *RaftLocalState) Reset()                    { *m = RaftLocalState{} }
func (m *RaftLocalState) String() string            { return proto.CompactTextString(m) }
func (*RaftLocalState) ProtoMessage()               {}
func (*RaftLocalState) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

func (m *RaftLocalState) GetHardState() *HardState {
	if m != nil {
		return m.HardState
	}
	return nil
}

func (m *RaftLocalState) GetLastIndex() uint64 {
	if m != nil && m.LastIndex != nil {
		return *m.LastIndex
	}
	return 0
}

type RaftTruncatedState struct {
	Index            *uint64 `protobuf:"varint,1,opt,name=index" json:"index,omitempty"`
	Term             *uint64 `protobuf:"varint,2,opt,name=term" json:"term,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *RaftTruncatedState) Reset()                    { *m = RaftTruncatedState{} }
func (m *RaftTruncatedState) String() string            { return proto.CompactTextString(m) }
func (*RaftTruncatedState) ProtoMessage()               {}
func (*RaftTruncatedState) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{17} }

func (m *RaftTruncatedState) GetIndex() uint64 {
	if m != nil && m.Index != nil {
		return *m.Index
	}
	return 0
}

func (m *RaftTruncatedState) GetTerm() uint64 {
	if m != nil && m.Term != nil {
		return *m.Term
	}
	return 0
}

type RaftApplyState struct {
	AppliedIndex     *uint64             `protobuf:"varint,1,opt,name=applied_index,json=appliedIndex" json:"applied_index,omitempty"`
	TruncatedState   *RaftTruncatedState `protobuf:"bytes,2,opt,name=truncated_state,json=truncatedState" json:"truncated_state,omitempty"`
	XXX_unrecognized []byte              `json:"-"`
}

func (m *RaftApplyState) Reset()                    { *m = RaftApplyState{} }
func (m *RaftApplyState) String() string            { return proto.CompactTextString(m) }
func (*RaftApplyState) ProtoMessage()               {}
func (*RaftApplyState) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{18} }

func (m *RaftApplyState) GetAppliedIndex() uint64 {
	if m != nil && m.AppliedIndex != nil {
		return *m.AppliedIndex
	}
	return 0
}

func (m *RaftApplyState) GetTruncatedState() *RaftTruncatedState {
	if m != nil {
		return m.TruncatedState
	}
	return nil
}

func init() {
	proto.RegisterType((*Entry)(nil), "raftpb.Entry")
	proto.RegisterType((*SnapshotMetadata)(nil), "raftpb.SnapshotMetadata")
	proto.RegisterType((*Snapshot)(nil), "raftpb.Snapshot")
	proto.RegisterType((*Message)(nil), "raftpb.Message")
	proto.RegisterType((*HardState)(nil), "raftpb.HardState")
	proto.RegisterType((*ConfState)(nil), "raftpb.ConfState")
	proto.RegisterType((*ConfChange)(nil), "raftpb.ConfChange")
	proto.RegisterType((*RaftMessage)(nil), "raftpb.RaftMessage")
	proto.RegisterType((*StoreIdent)(nil), "raftpb.StoreIdent")
	proto.RegisterType((*ACKMessage)(nil), "raftpb.ACKMessage")
	proto.RegisterType((*SnapshotMessageHeader)(nil), "raftpb.SnapshotMessageHeader")
	proto.RegisterType((*SnapshotMessage)(nil), "raftpb.SnapshotMessage")
	proto.RegisterType((*SnapshotChunkMessage)(nil), "raftpb.SnapshotChunkMessage")
	proto.RegisterType((*SnapshotAckMessage)(nil), "raftpb.SnapshotAckMessage")
	proto.RegisterType((*SnapshotAskMessage)(nil), "raftpb.SnapshotAskMessage")
	proto.RegisterType((*CellLocalState)(nil), "raftpb.CellLocalState")
	proto.RegisterType((*RaftLocalState)(nil), "raftpb.RaftLocalState")
	proto.RegisterType((*RaftTruncatedState)(nil), "raftpb.RaftTruncatedState")
	proto.RegisterType((*RaftApplyState)(nil), "raftpb.RaftApplyState")
	proto.RegisterEnum("raftpb.EntryType", EntryType_name, EntryType_value)
	proto.RegisterEnum("raftpb.MessageType", MessageType_name, MessageType_value)
	proto.RegisterEnum("raftpb.ConfChangeType", ConfChangeType_name, ConfChangeType_value)
	proto.RegisterEnum("raftpb.SnapshotState", SnapshotState_name, SnapshotState_value)
	proto.RegisterEnum("raftpb.PeerState", PeerState_name, PeerState_value)
}

func init() { proto.RegisterFile("raft/raft.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 1272 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x56, 0xdb, 0x6e, 0x1b, 0x37,
	0x13, 0xf6, 0xae, 0x4e, 0xab, 0x91, 0x2c, 0xd1, 0x8c, 0x13, 0x2c, 0x82, 0xfc, 0x3f, 0xd4, 0x6d,
	0x0b, 0xbb, 0x46, 0x60, 0xb7, 0x46, 0x72, 0xd9, 0x02, 0xae, 0x1c, 0xc0, 0x46, 0xe2, 0x20, 0xa5,
	0x95, 0xde, 0xf4, 0x22, 0xa0, 0x77, 0xa9, 0x43, 0x2d, 0x2d, 0xb7, 0x4b, 0x2a, 0x4d, 0xf2, 0x0c,
	0x7d, 0x83, 0x3e, 0x4f, 0x9f, 0xa0, 0x37, 0x05, 0xfa, 0x32, 0xc5, 0x90, 0xdc, 0x83, 0x64, 0xf7,
	0xae, 0x37, 0xd2, 0x9c, 0xc8, 0x99, 0x6f, 0xf8, 0x91, 0xb3, 0x30, 0xcc, 0xf9, 0x54, 0x9f, 0xe0,
	0xcf, 0x71, 0x96, 0x4b, 0x2d, 0x69, 0x1b, 0xe5, 0xec, 0xe6, 0xf1, 0x83, 0x95, 0xd0, 0x3c, 0xbb,
	0x39, 0xb1, 0x7f, 0xd6, 0x19, 0x2d, 0xa1, 0xf5, 0x22, 0xd5, 0xf9, 0x47, 0x4a, 0xa1, 0x39, 0x11,
	0xf9, 0x2a, 0xf4, 0x47, 0xde, 0x61, 0x93, 0x19, 0x99, 0xee, 0x43, 0xeb, 0x32, 0x4d, 0xc4, 0x87,
	0xb0, 0x61, 0x8c, 0x56, 0xa1, 0x5f, 0x42, 0x73, 0xf2, 0x31, 0x13, 0xa1, 0x37, 0xf2, 0x0e, 0x07,
	0xa7, 0x7b, 0xc7, 0x76, 0xfb, 0x63, 0xb3, 0x0d, 0x3a, 0x98, 0x71, 0xe3, 0x86, 0xe7, 0x5c, 0xf3,
	0xb0, 0x39, 0xf2, 0x0e, 0xfb, 0xcc, 0xc8, 0x51, 0x0a, 0xe4, 0x3a, 0xe5, 0x99, 0x9a, 0x4b, 0x7d,
	0x25, 0x34, 0x4f, 0xb8, 0xe6, 0xf4, 0x6b, 0x80, 0x58, 0xa6, 0xd3, 0x77, 0x4a, 0x73, 0x6d, 0x37,
	0xed, 0x55, 0x9b, 0x8e, 0x65, 0x3a, 0xbd, 0x46, 0x07, 0xeb, 0xc6, 0x85, 0x88, 0x65, 0x2d, 0x4c,
	0x59, 0xb6, 0x56, 0xab, 0x60, 0x3e, 0x8d, 0x00, 0x6c, 0xad, 0x46, 0x8e, 0x26, 0x10, 0x14, 0xf9,
	0xd0, 0x8f, 0xf9, 0x4c, 0x86, 0x3e, 0x33, 0x32, 0x7d, 0x06, 0xc1, 0xca, 0xd5, 0x61, 0x36, 0xeb,
	0x9d, 0x86, 0x45, 0xe6, 0xed, 0x3a, 0x59, 0x19, 0x19, 0xfd, 0xed, 0x43, 0xe7, 0x4a, 0x28, 0xc5,
	0x67, 0x82, 0x1e, 0x40, 0x53, 0x57, 0xcd, 0x78, 0x50, 0xac, 0x76, 0x6e, 0xdb, 0x0e, 0x0c, 0xa0,
	0x03, 0xf0, 0xb5, 0x74, 0x15, 0xfb, 0x5a, 0x62, 0x39, 0xd3, 0x5c, 0x96, 0xe5, 0xa2, 0x5c, 0x42,
	0x68, 0x56, 0x10, 0x68, 0x08, 0x9d, 0xa5, 0x9c, 0x99, 0xa3, 0x69, 0x19, 0x73, 0xa1, 0x56, 0x6d,
	0x68, 0xd7, 0xdb, 0x70, 0x00, 0x1d, 0x91, 0xea, 0x7c, 0x21, 0x54, 0xd8, 0x19, 0x35, 0x0e, 0x7b,
	0xa7, 0xbb, 0x1b, 0x07, 0xc4, 0x0a, 0x2f, 0x7d, 0x04, 0xed, 0x58, 0xae, 0x56, 0x0b, 0x1d, 0x06,
	0x66, 0xbd, 0xd3, 0xe8, 0x53, 0x08, 0x94, 0xc3, 0x1e, 0x76, 0x4d, 0x4f, 0xc8, 0x76, 0x4f, 0x58,
	0x19, 0x81, 0xbb, 0xe4, 0xe2, 0x67, 0x11, 0xeb, 0x10, 0x46, 0xde, 0x61, 0xc0, 0x9c, 0x46, 0xff,
	0x0f, 0x60, 0xa5, 0x8b, 0x45, 0xaa, 0xc3, 0x9e, 0xc9, 0x50, 0xb3, 0x20, 0xac, 0x58, 0xa6, 0x5a,
	0x7c, 0xd0, 0x61, 0xdf, 0x1c, 0x48, 0xa1, 0x46, 0x2f, 0xa1, 0x7b, 0xc1, 0xf3, 0xc4, 0x1e, 0x75,
	0xd1, 0x11, 0xaf, 0xd6, 0x11, 0x0a, 0xcd, 0xf7, 0x52, 0x8b, 0x82, 0xa9, 0x28, 0xd7, 0xc0, 0x34,
	0xea, 0x60, 0xa2, 0xcf, 0xa0, 0x3b, 0xae, 0xf3, 0x26, 0x95, 0x89, 0x50, 0xa1, 0x37, 0x6a, 0x60,
	0xc3, 0x8c, 0x12, 0x7d, 0x02, 0xc0, 0x90, 0xf1, 0x9c, 0xa7, 0x33, 0x73, 0x4c, 0x97, 0xe7, 0x2e,
	0x9d, 0x7f, 0x79, 0x4e, 0x8f, 0x1c, 0xd9, 0x7d, 0x73, 0xbe, 0x8f, 0xea, 0xbc, 0xb4, 0x2b, 0x6a,
	0x8c, 0x7f, 0x04, 0xed, 0xd7, 0x32, 0x11, 0x97, 0xe7, 0x45, 0x11, 0x56, 0x43, 0xac, 0x63, 0x87,
	0xd5, 0x5e, 0x86, 0x42, 0x8d, 0x7e, 0xf7, 0xa1, 0xc7, 0xf8, 0x54, 0x17, 0x6c, 0x42, 0x18, 0x62,
	0xb9, 0x2c, 0x2b, 0x70, 0x1a, 0x3d, 0x84, 0x00, 0x09, 0xf2, 0x46, 0x88, 0xdc, 0xf1, 0xb4, 0x7f,
	0xec, 0xae, 0x31, 0xda, 0x58, 0xe9, 0xa5, 0x5f, 0x40, 0x5b, 0x4b, 0x13, 0xd7, 0xb8, 0x27, 0xce,
	0xf9, 0xe8, 0x57, 0xd0, 0x59, 0xd9, 0x94, 0xa6, 0xa2, 0xde, 0xe9, 0x70, 0x8b, 0xb8, 0xac, 0xf0,
	0xd3, 0x13, 0xe8, 0x62, 0x11, 0x2f, 0x32, 0x19, 0xcf, 0x0d, 0x03, 0xf1, 0x76, 0xba, 0x3d, 0xc7,
	0x85, 0x83, 0x55, 0x31, 0x74, 0x04, 0xbd, 0x85, 0x9a, 0xc8, 0xd5, 0x8d, 0xd2, 0x32, 0x15, 0x86,
	0x9c, 0x01, 0xab, 0x9b, 0xf0, 0x1c, 0x94, 0xe6, 0xb9, 0x0e, 0x3b, 0xa6, 0x1b, 0x56, 0xa1, 0x04,
	0x1a, 0x22, 0x4d, 0x0c, 0x19, 0xfb, 0x0c, 0xc5, 0xe8, 0x1c, 0xe0, 0x5a, 0xcb, 0x5c, 0x5c, 0x26,
	0x22, 0xd5, 0xf4, 0x09, 0x74, 0xe3, 0xe5, 0x5a, 0x69, 0x91, 0x97, 0xed, 0xa9, 0x0c, 0xd8, 0x63,
	0x65, 0x62, 0xcf, 0x1d, 0x2f, 0x0a, 0x35, 0x3a, 0x06, 0x38, 0x1b, 0xbf, 0x2c, 0x3a, 0x4c, 0xa0,
	0xa1, 0xc4, 0x2f, 0x6e, 0x3d, 0x8a, 0xdb, 0x17, 0x33, 0xfa, 0xc3, 0x83, 0x87, 0xd5, 0xe5, 0x37,
	0xab, 0x2e, 0x04, 0x4f, 0x44, 0x4e, 0x47, 0xd0, 0x44, 0x98, 0xee, 0x8d, 0xea, 0xd7, 0xbb, 0xc0,
	0x8c, 0xe7, 0x3f, 0x3f, 0xa7, 0xfb, 0x1e, 0x84, 0xf2, 0xda, 0xb7, 0xea, 0xd7, 0xde, 0xe1, 0x6a,
	0x97, 0xb8, 0xa2, 0xbf, 0x3c, 0x18, 0x6e, 0xe1, 0xa0, 0xcf, 0xa1, 0x3d, 0x37, 0x58, 0x1c, 0x86,
	0xff, 0xdd, 0x7d, 0xed, 0x6a, 0x80, 0x99, 0x0b, 0xa6, 0xa7, 0xd0, 0x8a, 0xe7, 0xeb, 0xf4, 0xd6,
	0x61, 0x7a, 0xb2, 0xbd, 0x6a, 0x8c, 0xce, 0x82, 0x39, 0x36, 0x94, 0x3e, 0x85, 0x06, 0x8f, 0x6f,
	0x1d, 0xba, 0xc7, 0xdb, 0x2b, 0xce, 0xe2, 0x32, 0x1e, 0xc3, 0x4c, 0xb4, 0xba, 0x75, 0x64, 0xbc,
	0x1b, 0xad, 0x6a, 0xd1, 0xea, 0x36, 0xfa, 0xcd, 0x83, 0xfd, 0xfb, 0x72, 0xdf, 0xfb, 0xc6, 0xef,
	0x43, 0x6b, 0xba, 0xc8, 0x95, 0x36, 0xc5, 0x07, 0xcc, 0x2a, 0x18, 0xb9, 0xe4, 0xca, 0x3e, 0x17,
	0x01, 0x33, 0x32, 0x7d, 0x0c, 0xc1, 0x74, 0xb1, 0x14, 0xd7, 0x8b, 0x4f, 0xc2, 0x75, 0xbc, 0xd4,
	0xd1, 0x17, 0xcf, 0x45, 0x7c, 0x7b, 0xbd, 0x2e, 0xde, 0xe1, 0x52, 0x8f, 0xbe, 0x05, 0x7a, 0x17,
	0x17, 0x3d, 0xb0, 0x0d, 0xb0, 0x83, 0xe1, 0xe1, 0x36, 0x24, 0x3b, 0xd4, 0x30, 0x22, 0xda, 0xaf,
	0x2d, 0x2f, 0x81, 0x46, 0x3f, 0xc1, 0x00, 0x89, 0xf5, 0x4a, 0xc6, 0x7c, 0x69, 0x9f, 0xaf, 0x03,
	0x73, 0x6d, 0xf4, 0x9d, 0xc1, 0x8b, 0x4c, 0xb1, 0xdb, 0x59, 0x7f, 0xc9, 0x53, 0xff, 0xdf, 0x78,
	0x1a, 0xbd, 0x83, 0x01, 0x3e, 0x3b, 0xb5, 0xcd, 0x4f, 0xa0, 0x3b, 0x2f, 0x5e, 0xdd, 0xed, 0x21,
	0x5c, 0x3e, 0xc7, 0xac, 0x8a, 0xc1, 0xeb, 0x88, 0x4d, 0xbb, 0xac, 0x0d, 0xe2, 0xca, 0x10, 0x7d,
	0x07, 0x14, 0x13, 0x4c, 0xf2, 0x75, 0x1a, 0x73, 0x2d, 0x92, 0xad, 0xc1, 0xed, 0xdd, 0x37, 0xb8,
	0xfd, 0xda, 0xe0, 0xfe, 0x64, 0x0b, 0x3c, 0xcb, 0xb2, 0xe5, 0x47, 0xbb, 0xf6, 0x73, 0xd8, 0xe5,
	0x59, 0xb6, 0x5c, 0x88, 0xe4, 0x5d, 0x7d, 0x8f, 0xbe, 0x33, 0xda, 0x4f, 0x93, 0x31, 0x0c, 0x75,
	0x91, 0xd2, 0x7d, 0x50, 0xf8, 0x9b, 0x94, 0xba, 0x5b, 0x15, 0x1b, 0xe8, 0x0d, 0xfd, 0xe8, 0x1b,
	0xe8, 0x96, 0xdf, 0x32, 0x74, 0x08, 0x3d, 0xa3, 0xbc, 0x96, 0xf9, 0x8a, 0x2f, 0xc9, 0x0e, 0x7d,
	0x00, 0x43, 0x63, 0xa8, 0x26, 0x00, 0xf1, 0x8e, 0xfe, 0xf4, 0xa1, 0x57, 0x1b, 0xf9, 0x14, 0xa0,
	0x7d, 0xa5, 0x66, 0x17, 0xeb, 0x8c, 0xec, 0xd0, 0x1e, 0x74, 0xae, 0xd4, 0xec, 0x7b, 0xc1, 0x35,
	0xf1, 0x9c, 0xf2, 0x26, 0x97, 0x19, 0xf1, 0x5d, 0xd4, 0x59, 0x96, 0x91, 0x06, 0x1d, 0x00, 0x58,
	0x99, 0x09, 0x95, 0x91, 0xa6, 0x0b, 0xfc, 0x51, 0x6a, 0x41, 0x5a, 0x58, 0x84, 0x53, 0x8c, 0xb7,
	0xed, 0xbc, 0xc8, 0x1a, 0xd2, 0xa1, 0x04, 0xfa, 0x98, 0x4c, 0xf0, 0x5c, 0xdf, 0x60, 0x96, 0x80,
	0xee, 0x03, 0xa9, 0x5b, 0xcc, 0xa2, 0x2e, 0xa5, 0x30, 0xb8, 0x52, 0xb3, 0xb7, 0x69, 0x2e, 0x78,
	0x3c, 0xe7, 0x37, 0x4b, 0x41, 0x80, 0xee, 0xc1, 0xae, 0xdb, 0x08, 0xb1, 0xaf, 0x15, 0xe9, 0xb9,
	0xb0, 0x31, 0x92, 0xfb, 0x87, 0xb5, 0xcc, 0xd7, 0x2b, 0xd2, 0xa7, 0x0f, 0x61, 0xef, 0x4a, 0xcd,
	0x26, 0x39, 0x4f, 0xd5, 0x54, 0xe4, 0xaf, 0xcc, 0xab, 0x40, 0x76, 0xdd, 0xea, 0xc9, 0x62, 0x25,
	0xe4, 0x5a, 0xbf, 0x96, 0xbf, 0x92, 0x81, 0x2b, 0x86, 0x09, 0x6e, 0x4f, 0x84, 0x0c, 0x5d, 0x31,
	0xa5, 0xc5, 0x14, 0x43, 0x1c, 0xde, 0x37, 0xb9, 0x30, 0x10, 0xf7, 0x5c, 0x56, 0xa7, 0x9b, 0x18,
	0x7a, 0x84, 0x57, 0x60, 0x63, 0xce, 0x62, 0x1d, 0x95, 0xe5, 0x2c, 0x49, 0x70, 0xc0, 0x92, 0x1d,
	0x1a, 0xc2, 0x7e, 0x65, 0x66, 0x62, 0x25, 0xdf, 0x0b, 0xe3, 0xf1, 0x36, 0x3d, 0x6f, 0xb3, 0x84,
	0x6b, 0xeb, 0xf1, 0x8f, 0x9e, 0xc3, 0xee, 0xc6, 0x5d, 0xc4, 0xd3, 0x38, 0x8b, 0x63, 0x91, 0x69,
	0xb2, 0x83, 0x32, 0x33, 0xdf, 0x2a, 0xc4, 0xa3, 0x7d, 0x08, 0x98, 0x88, 0xc5, 0xe2, 0xbd, 0x48,
	0x88, 0x7f, 0xf4, 0x0c, 0xba, 0xe5, 0x7d, 0xc3, 0xb0, 0x92, 0x17, 0x7d, 0x08, 0x0c, 0x5b, 0x17,
	0xe9, 0x8c, 0x78, 0x74, 0x17, 0xba, 0xe5, 0xbc, 0x23, 0xfe, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff,
	0xb1, 0x2d, 0xde, 0x75, 0x94, 0x0b, 0x00, 0x00,
}

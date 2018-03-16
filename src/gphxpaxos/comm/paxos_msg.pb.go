// Code generated by protoc-gen-go. DO NOT EDIT.
// source: paxos_msg.proto

/*
Package comm is a generated protocol buffer package.

It is generated from these files:
	paxos_msg.proto

It has these top-level messages:
	Header
	PaxosMsg
	CheckpointMsg
	AcceptorStateData
	PaxosNodeInfo
	SystemVariables
	MasterVariables
	PaxosValue
	BatchPaxosValues
*/
package comm

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Header struct {
	Gid              *uint64 `protobuf:"varint,1,req,name=gid" json:"gid,omitempty"`
	Rid              *uint64 `protobuf:"varint,2,req,name=rid" json:"rid,omitempty"`
	Cmdid            *int32  `protobuf:"varint,3,req,name=cmdid" json:"cmdid,omitempty"`
	Version          *int32  `protobuf:"varint,4,opt,name=version" json:"version,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Header) Reset()                    { *m = Header{} }
func (m *Header) String() string            { return proto.CompactTextString(m) }
func (*Header) ProtoMessage()               {}
func (*Header) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Header) GetGid() uint64 {
	if m != nil && m.Gid != nil {
		return *m.Gid
	}
	return 0
}

func (m *Header) GetRid() uint64 {
	if m != nil && m.Rid != nil {
		return *m.Rid
	}
	return 0
}

func (m *Header) GetCmdid() int32 {
	if m != nil && m.Cmdid != nil {
		return *m.Cmdid
	}
	return 0
}

func (m *Header) GetVersion() int32 {
	if m != nil && m.Version != nil {
		return *m.Version
	}
	return 0
}

type PaxosMsg struct {
	MsgType             *int32  `protobuf:"varint,1,req,name=MsgType" json:"MsgType,omitempty"`
	InstanceID          *uint64 `protobuf:"varint,2,opt,name=InstanceID" json:"InstanceID,omitempty"`
	NodeID              *uint64 `protobuf:"varint,3,opt,name=NodeID" json:"NodeID,omitempty"`
	ProposalID          *uint64 `protobuf:"varint,4,opt,name=ProposalID" json:"ProposalID,omitempty"`
	ProposalNodeID      *uint64 `protobuf:"varint,5,opt,name=ProposalNodeID" json:"ProposalNodeID,omitempty"`
	Value               []byte  `protobuf:"bytes,6,opt,name=Value" json:"Value,omitempty"`
	PreAcceptID         *uint64 `protobuf:"varint,7,opt,name=PreAcceptID" json:"PreAcceptID,omitempty"`
	PreAcceptNodeID     *uint64 `protobuf:"varint,8,opt,name=PreAcceptNodeID" json:"PreAcceptNodeID,omitempty"`
	RejectByPromiseID   *uint64 `protobuf:"varint,9,opt,name=RejectByPromiseID" json:"RejectByPromiseID,omitempty"`
	NowInstanceID       *uint64 `protobuf:"varint,10,opt,name=NowInstanceID" json:"NowInstanceID,omitempty"`
	MinChosenInstanceID *uint64 `protobuf:"varint,11,opt,name=MinChosenInstanceID" json:"MinChosenInstanceID,omitempty"`
	LastChecksum        *uint32 `protobuf:"varint,12,opt,name=LastChecksum" json:"LastChecksum,omitempty"`
	Flag                *uint32 `protobuf:"varint,13,opt,name=Flag" json:"Flag,omitempty"`
	SystemVariables     []byte  `protobuf:"bytes,14,opt,name=SystemVariables" json:"SystemVariables,omitempty"`
	MasterVariables     []byte  `protobuf:"bytes,15,opt,name=MasterVariables" json:"MasterVariables,omitempty"`
	XXX_unrecognized    []byte  `json:"-"`
}

func (m *PaxosMsg) Reset()                    { *m = PaxosMsg{} }
func (m *PaxosMsg) String() string            { return proto.CompactTextString(m) }
func (*PaxosMsg) ProtoMessage()               {}
func (*PaxosMsg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *PaxosMsg) GetMsgType() int32 {
	if m != nil && m.MsgType != nil {
		return *m.MsgType
	}
	return 0
}

func (m *PaxosMsg) GetInstanceID() uint64 {
	if m != nil && m.InstanceID != nil {
		return *m.InstanceID
	}
	return 0
}

func (m *PaxosMsg) GetNodeID() uint64 {
	if m != nil && m.NodeID != nil {
		return *m.NodeID
	}
	return 0
}

func (m *PaxosMsg) GetProposalID() uint64 {
	if m != nil && m.ProposalID != nil {
		return *m.ProposalID
	}
	return 0
}

func (m *PaxosMsg) GetProposalNodeID() uint64 {
	if m != nil && m.ProposalNodeID != nil {
		return *m.ProposalNodeID
	}
	return 0
}

func (m *PaxosMsg) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *PaxosMsg) GetPreAcceptID() uint64 {
	if m != nil && m.PreAcceptID != nil {
		return *m.PreAcceptID
	}
	return 0
}

func (m *PaxosMsg) GetPreAcceptNodeID() uint64 {
	if m != nil && m.PreAcceptNodeID != nil {
		return *m.PreAcceptNodeID
	}
	return 0
}

func (m *PaxosMsg) GetRejectByPromiseID() uint64 {
	if m != nil && m.RejectByPromiseID != nil {
		return *m.RejectByPromiseID
	}
	return 0
}

func (m *PaxosMsg) GetNowInstanceID() uint64 {
	if m != nil && m.NowInstanceID != nil {
		return *m.NowInstanceID
	}
	return 0
}

func (m *PaxosMsg) GetMinChosenInstanceID() uint64 {
	if m != nil && m.MinChosenInstanceID != nil {
		return *m.MinChosenInstanceID
	}
	return 0
}

func (m *PaxosMsg) GetLastChecksum() uint32 {
	if m != nil && m.LastChecksum != nil {
		return *m.LastChecksum
	}
	return 0
}

func (m *PaxosMsg) GetFlag() uint32 {
	if m != nil && m.Flag != nil {
		return *m.Flag
	}
	return 0
}

func (m *PaxosMsg) GetSystemVariables() []byte {
	if m != nil {
		return m.SystemVariables
	}
	return nil
}

func (m *PaxosMsg) GetMasterVariables() []byte {
	if m != nil {
		return m.MasterVariables
	}
	return nil
}

type CheckpointMsg struct {
	MsgType              *int32  `protobuf:"varint,1,req,name=MsgType" json:"MsgType,omitempty"`
	NodeID               *uint64 `protobuf:"varint,2,req,name=NodeID" json:"NodeID,omitempty"`
	Flag                 *int32  `protobuf:"varint,3,opt,name=Flag" json:"Flag,omitempty"`
	UUID                 *uint64 `protobuf:"varint,4,req,name=UUID" json:"UUID,omitempty"`
	Sequence             *uint64 `protobuf:"varint,5,req,name=Sequence" json:"Sequence,omitempty"`
	CheckpointInstanceID *uint64 `protobuf:"varint,6,opt,name=CheckpointInstanceID" json:"CheckpointInstanceID,omitempty"`
	Checksum             *uint32 `protobuf:"varint,7,opt,name=Checksum" json:"Checksum,omitempty"`
	FilePath             *string `protobuf:"bytes,8,opt,name=FilePath" json:"FilePath,omitempty"`
	SMID                 *int32  `protobuf:"varint,9,opt,name=SMID" json:"SMID,omitempty"`
	Offset               *uint64 `protobuf:"varint,10,opt,name=Offset" json:"Offset,omitempty"`
	Buffer               []byte  `protobuf:"bytes,11,opt,name=Buffer" json:"Buffer,omitempty"`
	XXX_unrecognized     []byte  `json:"-"`
}

func (m *CheckpointMsg) Reset()                    { *m = CheckpointMsg{} }
func (m *CheckpointMsg) String() string            { return proto.CompactTextString(m) }
func (*CheckpointMsg) ProtoMessage()               {}
func (*CheckpointMsg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *CheckpointMsg) GetMsgType() int32 {
	if m != nil && m.MsgType != nil {
		return *m.MsgType
	}
	return 0
}

func (m *CheckpointMsg) GetNodeID() uint64 {
	if m != nil && m.NodeID != nil {
		return *m.NodeID
	}
	return 0
}

func (m *CheckpointMsg) GetFlag() int32 {
	if m != nil && m.Flag != nil {
		return *m.Flag
	}
	return 0
}

func (m *CheckpointMsg) GetUUID() uint64 {
	if m != nil && m.UUID != nil {
		return *m.UUID
	}
	return 0
}

func (m *CheckpointMsg) GetSequence() uint64 {
	if m != nil && m.Sequence != nil {
		return *m.Sequence
	}
	return 0
}

func (m *CheckpointMsg) GetCheckpointInstanceID() uint64 {
	if m != nil && m.CheckpointInstanceID != nil {
		return *m.CheckpointInstanceID
	}
	return 0
}

func (m *CheckpointMsg) GetChecksum() uint32 {
	if m != nil && m.Checksum != nil {
		return *m.Checksum
	}
	return 0
}

func (m *CheckpointMsg) GetFilePath() string {
	if m != nil && m.FilePath != nil {
		return *m.FilePath
	}
	return ""
}

func (m *CheckpointMsg) GetSMID() int32 {
	if m != nil && m.SMID != nil {
		return *m.SMID
	}
	return 0
}

func (m *CheckpointMsg) GetOffset() uint64 {
	if m != nil && m.Offset != nil {
		return *m.Offset
	}
	return 0
}

func (m *CheckpointMsg) GetBuffer() []byte {
	if m != nil {
		return m.Buffer
	}
	return nil
}

type AcceptorStateData struct {
	InstanceID       *uint64 `protobuf:"varint,1,req,name=InstanceID" json:"InstanceID,omitempty"`
	PromiseID        *uint64 `protobuf:"varint,2,req,name=PromiseID" json:"PromiseID,omitempty"`
	PromiseNodeID    *uint64 `protobuf:"varint,3,req,name=PromiseNodeID" json:"PromiseNodeID,omitempty"`
	AcceptedID       *uint64 `protobuf:"varint,4,req,name=AcceptedID" json:"AcceptedID,omitempty"`
	AcceptedNodeID   *uint64 `protobuf:"varint,5,req,name=AcceptedNodeID" json:"AcceptedNodeID,omitempty"`
	AcceptedValue    []byte  `protobuf:"bytes,6,req,name=AcceptedValue" json:"AcceptedValue,omitempty"`
	Checksum         *uint32 `protobuf:"varint,7,req,name=Checksum" json:"Checksum,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *AcceptorStateData) Reset()                    { *m = AcceptorStateData{} }
func (m *AcceptorStateData) String() string            { return proto.CompactTextString(m) }
func (*AcceptorStateData) ProtoMessage()               {}
func (*AcceptorStateData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *AcceptorStateData) GetInstanceID() uint64 {
	if m != nil && m.InstanceID != nil {
		return *m.InstanceID
	}
	return 0
}

func (m *AcceptorStateData) GetPromiseID() uint64 {
	if m != nil && m.PromiseID != nil {
		return *m.PromiseID
	}
	return 0
}

func (m *AcceptorStateData) GetPromiseNodeID() uint64 {
	if m != nil && m.PromiseNodeID != nil {
		return *m.PromiseNodeID
	}
	return 0
}

func (m *AcceptorStateData) GetAcceptedID() uint64 {
	if m != nil && m.AcceptedID != nil {
		return *m.AcceptedID
	}
	return 0
}

func (m *AcceptorStateData) GetAcceptedNodeID() uint64 {
	if m != nil && m.AcceptedNodeID != nil {
		return *m.AcceptedNodeID
	}
	return 0
}

func (m *AcceptorStateData) GetAcceptedValue() []byte {
	if m != nil {
		return m.AcceptedValue
	}
	return nil
}

func (m *AcceptorStateData) GetChecksum() uint32 {
	if m != nil && m.Checksum != nil {
		return *m.Checksum
	}
	return 0
}

type PaxosNodeInfo struct {
	Rid              *uint64 `protobuf:"varint,1,req,name=Rid" json:"Rid,omitempty"`
	Nodeid           *uint64 `protobuf:"varint,2,req,name=Nodeid" json:"Nodeid,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *PaxosNodeInfo) Reset()                    { *m = PaxosNodeInfo{} }
func (m *PaxosNodeInfo) String() string            { return proto.CompactTextString(m) }
func (*PaxosNodeInfo) ProtoMessage()               {}
func (*PaxosNodeInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *PaxosNodeInfo) GetRid() uint64 {
	if m != nil && m.Rid != nil {
		return *m.Rid
	}
	return 0
}

func (m *PaxosNodeInfo) GetNodeid() uint64 {
	if m != nil && m.Nodeid != nil {
		return *m.Nodeid
	}
	return 0
}

type SystemVariables struct {
	Gid              *uint64          `protobuf:"varint,1,req,name=Gid" json:"Gid,omitempty"`
	MemberShip       []*PaxosNodeInfo `protobuf:"bytes,2,rep,name=MemberShip" json:"MemberShip,omitempty"`
	Version          *uint64          `protobuf:"varint,3,req,name=Version" json:"Version,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *SystemVariables) Reset()                    { *m = SystemVariables{} }
func (m *SystemVariables) String() string            { return proto.CompactTextString(m) }
func (*SystemVariables) ProtoMessage()               {}
func (*SystemVariables) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *SystemVariables) GetGid() uint64 {
	if m != nil && m.Gid != nil {
		return *m.Gid
	}
	return 0
}

func (m *SystemVariables) GetMemberShip() []*PaxosNodeInfo {
	if m != nil {
		return m.MemberShip
	}
	return nil
}

func (m *SystemVariables) GetVersion() uint64 {
	if m != nil && m.Version != nil {
		return *m.Version
	}
	return 0
}

type MasterVariables struct {
	MasterNodeid     *uint64 `protobuf:"varint,1,req,name=MasterNodeid" json:"MasterNodeid,omitempty"`
	Version          *uint64 `protobuf:"varint,2,req,name=Version" json:"Version,omitempty"`
	LeaseTime        *uint32 `protobuf:"varint,3,req,name=LeaseTime" json:"LeaseTime,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *MasterVariables) Reset()                    { *m = MasterVariables{} }
func (m *MasterVariables) String() string            { return proto.CompactTextString(m) }
func (*MasterVariables) ProtoMessage()               {}
func (*MasterVariables) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *MasterVariables) GetMasterNodeid() uint64 {
	if m != nil && m.MasterNodeid != nil {
		return *m.MasterNodeid
	}
	return 0
}

func (m *MasterVariables) GetVersion() uint64 {
	if m != nil && m.Version != nil {
		return *m.Version
	}
	return 0
}

func (m *MasterVariables) GetLeaseTime() uint32 {
	if m != nil && m.LeaseTime != nil {
		return *m.LeaseTime
	}
	return 0
}

type PaxosValue struct {
	SMID             *int32 `protobuf:"varint,1,req,name=SMID" json:"SMID,omitempty"`
	Value            []byte `protobuf:"bytes,2,req,name=Value" json:"Value,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *PaxosValue) Reset()                    { *m = PaxosValue{} }
func (m *PaxosValue) String() string            { return proto.CompactTextString(m) }
func (*PaxosValue) ProtoMessage()               {}
func (*PaxosValue) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *PaxosValue) GetSMID() int32 {
	if m != nil && m.SMID != nil {
		return *m.SMID
	}
	return 0
}

func (m *PaxosValue) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

type BatchPaxosValues struct {
	Values           []*PaxosValue `protobuf:"bytes,1,rep,name=Values" json:"Values,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *BatchPaxosValues) Reset()                    { *m = BatchPaxosValues{} }
func (m *BatchPaxosValues) String() string            { return proto.CompactTextString(m) }
func (*BatchPaxosValues) ProtoMessage()               {}
func (*BatchPaxosValues) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *BatchPaxosValues) GetValues() []*PaxosValue {
	if m != nil {
		return m.Values
	}
	return nil
}

func init() {
	proto.RegisterType((*Header)(nil), "comm.Header")
	proto.RegisterType((*PaxosMsg)(nil), "comm.PaxosMsg")
	proto.RegisterType((*CheckpointMsg)(nil), "comm.CheckpointMsg")
	proto.RegisterType((*AcceptorStateData)(nil), "comm.AcceptorStateData")
	proto.RegisterType((*PaxosNodeInfo)(nil), "comm.PaxosNodeInfo")
	proto.RegisterType((*SystemVariables)(nil), "comm.SystemVariables")
	proto.RegisterType((*MasterVariables)(nil), "comm.MasterVariables")
	proto.RegisterType((*PaxosValue)(nil), "comm.PaxosValue")
	proto.RegisterType((*BatchPaxosValues)(nil), "comm.BatchPaxosValues")
}

func init() { proto.RegisterFile("paxos_msg.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 711 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x54, 0x5d, 0x6b, 0x13, 0x41,
	0x14, 0x25, 0x9b, 0x4d, 0xda, 0xdc, 0x26, 0xfd, 0x98, 0x16, 0x59, 0x44, 0x64, 0x59, 0x44, 0xf6,
	0x41, 0x8a, 0x54, 0x10, 0x04, 0x5f, 0x6c, 0x43, 0x6b, 0xa0, 0xa9, 0x61, 0xd2, 0xf6, 0x49, 0x90,
	0xe9, 0xee, 0x24, 0x59, 0xcd, 0xee, 0xc6, 0x9d, 0x89, 0xda, 0x47, 0xff, 0x85, 0x7f, 0xc4, 0xff,
	0x27, 0xf7, 0xce, 0x7e, 0xa6, 0xc5, 0xb7, 0x39, 0xe7, 0x9e, 0x99, 0x3b, 0x73, 0xe6, 0xde, 0x0b,
	0x7b, 0x2b, 0xf1, 0x2b, 0x55, 0x5f, 0x62, 0x35, 0x3f, 0x5e, 0x65, 0xa9, 0x4e, 0x99, 0x1d, 0xa4,
	0x71, 0xec, 0x7d, 0x86, 0xee, 0x47, 0x29, 0x42, 0x99, 0xb1, 0x7d, 0x68, 0xcf, 0xa3, 0xd0, 0x69,
	0xb9, 0x96, 0x6f, 0x73, 0x5c, 0x22, 0x93, 0x45, 0xa1, 0x63, 0x19, 0x26, 0x8b, 0x42, 0x76, 0x04,
	0x9d, 0x20, 0x0e, 0xa3, 0xd0, 0x69, 0xbb, 0x96, 0xdf, 0xe1, 0x06, 0x30, 0x07, 0xb6, 0x7e, 0xc8,
	0x4c, 0x45, 0x69, 0xe2, 0xd8, 0x6e, 0xcb, 0xef, 0xf0, 0x02, 0x7a, 0x7f, 0x6c, 0xd8, 0x9e, 0x60,
	0xde, 0xb1, 0x9a, 0xa3, 0x6c, 0xac, 0xe6, 0xd7, 0xf7, 0x2b, 0x49, 0x49, 0x3a, 0xbc, 0x80, 0xec,
	0x39, 0xc0, 0x28, 0x51, 0x5a, 0x24, 0x81, 0x1c, 0x0d, 0x1d, 0xcb, 0x6d, 0xf9, 0x36, 0xaf, 0x31,
	0xec, 0x09, 0x74, 0xaf, 0xd2, 0x10, 0x63, 0x6d, 0x8a, 0xe5, 0x08, 0xf7, 0x4d, 0xb2, 0x74, 0x95,
	0x2a, 0xb1, 0x1c, 0x0d, 0x29, 0xb7, 0xcd, 0x6b, 0x0c, 0x7b, 0x09, 0xbb, 0x05, 0xca, 0xf7, 0x77,
	0x48, 0xb3, 0xc1, 0xe2, 0xb3, 0x6e, 0xc5, 0x72, 0x2d, 0x9d, 0xae, 0xdb, 0xf2, 0xfb, 0xdc, 0x00,
	0xe6, 0xc2, 0xce, 0x24, 0x93, 0x1f, 0x82, 0x40, 0xae, 0xf4, 0x68, 0xe8, 0x6c, 0xd1, 0xd6, 0x3a,
	0xc5, 0x7c, 0xd8, 0x2b, 0x61, 0x9e, 0x60, 0x9b, 0x54, 0x9b, 0x34, 0x7b, 0x05, 0x07, 0x5c, 0x7e,
	0x95, 0x81, 0x3e, 0xbd, 0x9f, 0x64, 0x69, 0x1c, 0x29, 0xd4, 0xf6, 0x48, 0xfb, 0x30, 0xc0, 0x5e,
	0xc0, 0xe0, 0x2a, 0xfd, 0x59, 0xb3, 0x04, 0x48, 0xd9, 0x24, 0xd9, 0x6b, 0x38, 0x1c, 0x47, 0xc9,
	0xd9, 0x22, 0x55, 0x32, 0xa9, 0x69, 0x77, 0x48, 0xfb, 0x58, 0x88, 0x79, 0xd0, 0xbf, 0x14, 0x4a,
	0x9f, 0x2d, 0x64, 0xf0, 0x4d, 0xad, 0x63, 0xa7, 0xef, 0xb6, 0xfc, 0x01, 0x6f, 0x70, 0x8c, 0x81,
	0x7d, 0xbe, 0x14, 0x73, 0x67, 0x40, 0x31, 0x5a, 0xe3, 0x3b, 0xa7, 0xf7, 0x4a, 0xcb, 0xf8, 0x56,
	0x64, 0x91, 0xb8, 0x5b, 0x4a, 0xe5, 0xec, 0x92, 0x53, 0x9b, 0x34, 0x2a, 0xc7, 0x42, 0x69, 0x99,
	0x55, 0xca, 0x3d, 0xa3, 0xdc, 0xa0, 0xbd, 0xbf, 0x16, 0x0c, 0x28, 0xe9, 0x2a, 0x8d, 0x12, 0xfd,
	0xff, 0xfa, 0xa8, 0xfe, 0xdf, 0xd4, 0x62, 0xf1, 0xff, 0xc5, 0x5d, 0xdb, 0x54, 0x75, 0xe6, 0xae,
	0x0c, 0xec, 0x9b, 0x1b, 0xaa, 0x06, 0x54, 0xd2, 0x9a, 0x3d, 0x85, 0xed, 0xa9, 0xfc, 0xbe, 0x96,
	0x49, 0x20, 0x9d, 0x0e, 0xf1, 0x25, 0x66, 0x27, 0x70, 0x54, 0x5d, 0xa3, 0x66, 0x63, 0x97, 0x6c,
	0x7c, 0x34, 0x86, 0xe7, 0x95, 0x1e, 0x6e, 0x91, 0x4f, 0x25, 0xc6, 0xd8, 0x79, 0xb4, 0x94, 0x13,
	0xa1, 0x17, 0x54, 0x0c, 0x3d, 0x5e, 0x62, 0xbc, 0xdb, 0x74, 0x9c, 0x7f, 0x7c, 0x87, 0xd3, 0x1a,
	0xdf, 0xf6, 0x69, 0x36, 0x53, 0x52, 0xe7, 0x9f, 0x9c, 0x23, 0xe4, 0x4f, 0xd7, 0xb3, 0x99, 0xcc,
	0xe8, 0x43, 0xfb, 0x3c, 0x47, 0xde, 0x6f, 0x0b, 0x0e, 0x4c, 0x69, 0xa5, 0xd9, 0x54, 0x0b, 0x2d,
	0x87, 0x42, 0x8b, 0x8d, 0x0e, 0x32, 0x3d, 0x5c, 0xef, 0xa0, 0x67, 0xd0, 0xab, 0xea, 0xce, 0x98,
	0xd8, 0x6b, 0xd4, 0x5b, 0x0e, 0xca, 0x36, 0x43, 0x45, 0x93, 0xc4, 0x1c, 0x26, 0xb1, 0x0c, 0x4b,
	0x7f, 0x6b, 0x0c, 0x76, 0x5b, 0x81, 0xca, 0x6e, 0x43, 0xcd, 0x06, 0x8b, 0xd9, 0x0a, 0xa6, 0xe8,
	0x3a, 0xcb, 0xef, 0xf3, 0x26, 0xb9, 0xe1, 0xb1, 0x55, 0xf7, 0xd8, 0x7b, 0x07, 0x03, 0x9a, 0x2a,
	0x74, 0x60, 0x32, 0x4b, 0x71, 0x52, 0xf1, 0x6a, 0x76, 0xf1, 0x28, 0x2c, 0x4a, 0xa6, 0x1c, 0x5f,
	0x39, 0xf2, 0xb2, 0x07, 0xa5, 0x8c, 0x9b, 0x2f, 0xaa, 0xcd, 0x17, 0x51, 0xc8, 0xde, 0x00, 0x8c,
	0x65, 0x7c, 0x27, 0xb3, 0xe9, 0x22, 0x5a, 0x39, 0x96, 0xdb, 0xf6, 0x77, 0x4e, 0x0e, 0x8f, 0x71,
	0x5e, 0x1e, 0x37, 0xf2, 0xf2, 0x9a, 0x0c, 0xcb, 0xf7, 0x36, 0x9f, 0x82, 0xc6, 0xbe, 0x02, 0x7a,
	0xf1, 0x83, 0xa6, 0xc0, 0x4e, 0x34, 0x54, 0x7e, 0x49, 0x93, 0xbc, 0xc1, 0xd5, 0x0f, 0xb4, 0x1a,
	0x07, 0xe2, 0x6f, 0x5e, 0x4a, 0xa1, 0xe4, 0x75, 0x14, 0x4b, 0x4a, 0x36, 0xe0, 0x15, 0xe1, 0xbd,
	0x05, 0xa0, 0x5b, 0x1a, 0x1f, 0x8b, 0x9a, 0x33, 0x2d, 0x65, 0x6a, 0xae, 0x9c, 0x77, 0x16, 0x39,
	0x6f, 0x80, 0xf7, 0x1e, 0xf6, 0x4f, 0x85, 0x0e, 0x16, 0xd5, 0x66, 0xec, 0xe7, 0xae, 0x59, 0x39,
	0x2d, 0x72, 0x61, 0xbf, 0xe6, 0x02, 0x05, 0x78, 0x1e, 0xff, 0x17, 0x00, 0x00, 0xff, 0xff, 0x6a,
	0x4a, 0xe7, 0x6e, 0x60, 0x06, 0x00, 0x00,
}
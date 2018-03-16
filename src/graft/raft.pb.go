// Code generated by protoc-gen-go. DO NOT EDIT.
// source: raft.proto

/*
Package graft is a generated protocol buffer package.

It is generated from these files:
	raft.proto

It has these top-level messages:
	LocalFileMeta
	EntryMeta
	ConfigurationPBMeta
	LogPBMeta
	StablePBMeta
	LocalSnapshotPbMeta
	RequestVoteRequest
	RequestVoteResponse
	AppendEntriesRequest
	AppendEntriesResponse
	SnapshotMeta
	InstallSnapshotRequest
	InstallSnapshotResponse
	TimeoutNowRequest
	TimeoutNowResponse
*/
package graft

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

type EntryType int32

const (
	EntryType_ENTRY_TYPE_UNKNOWN       EntryType = 0
	EntryType_ENTRY_TYPE_NO_OP         EntryType = 1
	EntryType_ENTRY_TYPE_DATA          EntryType = 2
	EntryType_ENTRY_TYPE_CONFIGURATION EntryType = 3
)

var EntryType_name = map[int32]string{
	0: "ENTRY_TYPE_UNKNOWN",
	1: "ENTRY_TYPE_NO_OP",
	2: "ENTRY_TYPE_DATA",
	3: "ENTRY_TYPE_CONFIGURATION",
}
var EntryType_value = map[string]int32{
	"ENTRY_TYPE_UNKNOWN":       0,
	"ENTRY_TYPE_NO_OP":         1,
	"ENTRY_TYPE_DATA":          2,
	"ENTRY_TYPE_CONFIGURATION": 3,
}

func (x EntryType) String() string {
	return proto.EnumName(EntryType_name, int32(x))
}
func (EntryType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type ErrorType int32

const (
	ErrorType_ERROR_TYPE_NONE          ErrorType = 0
	ErrorType_ERROR_TYPE_LOG           ErrorType = 1
	ErrorType_ERROR_TYPE_STABLE        ErrorType = 2
	ErrorType_ERROR_TYPE_SNAPSHOT      ErrorType = 3
	ErrorType_ERROR_TYPE_STATE_MACHINE ErrorType = 4
)

var ErrorType_name = map[int32]string{
	0: "ERROR_TYPE_NONE",
	1: "ERROR_TYPE_LOG",
	2: "ERROR_TYPE_STABLE",
	3: "ERROR_TYPE_SNAPSHOT",
	4: "ERROR_TYPE_STATE_MACHINE",
}
var ErrorType_value = map[string]int32{
	"ERROR_TYPE_NONE":          0,
	"ERROR_TYPE_LOG":           1,
	"ERROR_TYPE_STABLE":        2,
	"ERROR_TYPE_SNAPSHOT":      3,
	"ERROR_TYPE_STATE_MACHINE": 4,
}

func (x ErrorType) String() string {
	return proto.EnumName(ErrorType_name, int32(x))
}
func (ErrorType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type FileSource int32

const (
	FileSource_FILE_SOURCE_LOCAL     FileSource = 0
	FileSource_FILE_SOURCE_REFERENCE FileSource = 1
)

var FileSource_name = map[int32]string{
	0: "FILE_SOURCE_LOCAL",
	1: "FILE_SOURCE_REFERENCE",
}
var FileSource_value = map[string]int32{
	"FILE_SOURCE_LOCAL":     0,
	"FILE_SOURCE_REFERENCE": 1,
}

func (x FileSource) String() string {
	return proto.EnumName(FileSource_name, int32(x))
}
func (FileSource) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type LocalFileMeta struct {
	UserMeta []byte     `protobuf:"bytes,1,opt,name=user_meta,json=userMeta,proto3" json:"user_meta,omitempty"`
	Source   FileSource `protobuf:"varint,2,opt,name=source,enum=graft.FileSource" json:"source,omitempty"`
	Checksum string     `protobuf:"bytes,3,opt,name=checksum" json:"checksum,omitempty"`
}

func (m *LocalFileMeta) Reset()                    { *m = LocalFileMeta{} }
func (m *LocalFileMeta) String() string            { return proto.CompactTextString(m) }
func (*LocalFileMeta) ProtoMessage()               {}
func (*LocalFileMeta) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *LocalFileMeta) GetUserMeta() []byte {
	if m != nil {
		return m.UserMeta
	}
	return nil
}

func (m *LocalFileMeta) GetSource() FileSource {
	if m != nil {
		return m.Source
	}
	return FileSource_FILE_SOURCE_LOCAL
}

func (m *LocalFileMeta) GetChecksum() string {
	if m != nil {
		return m.Checksum
	}
	return ""
}

// data store in baidu-rpc's attachment
type EntryMeta struct {
	Term    int64     `protobuf:"varint,1,opt,name=term" json:"term,omitempty"`
	Type    EntryType `protobuf:"varint,2,opt,name=type,enum=graft.EntryType" json:"type,omitempty"`
	Peers   []string  `protobuf:"bytes,3,rep,name=peers" json:"peers,omitempty"`
	DataLen int64     `protobuf:"varint,4,opt,name=data_len,json=dataLen" json:"data_len,omitempty"`
	// Don't change field id of `old_peers' in the consideration of backward
	// compatibility
	OldPeers []string `protobuf:"bytes,5,rep,name=old_peers,json=oldPeers" json:"old_peers,omitempty"`
}

func (m *EntryMeta) Reset()                    { *m = EntryMeta{} }
func (m *EntryMeta) String() string            { return proto.CompactTextString(m) }
func (*EntryMeta) ProtoMessage()               {}
func (*EntryMeta) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *EntryMeta) GetTerm() int64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *EntryMeta) GetType() EntryType {
	if m != nil {
		return m.Type
	}
	return EntryType_ENTRY_TYPE_UNKNOWN
}

func (m *EntryMeta) GetPeers() []string {
	if m != nil {
		return m.Peers
	}
	return nil
}

func (m *EntryMeta) GetDataLen() int64 {
	if m != nil {
		return m.DataLen
	}
	return 0
}

func (m *EntryMeta) GetOldPeers() []string {
	if m != nil {
		return m.OldPeers
	}
	return nil
}

type ConfigurationPBMeta struct {
	Peers    []string `protobuf:"bytes,1,rep,name=peers" json:"peers,omitempty"`
	OldPeers []string `protobuf:"bytes,2,rep,name=old_peers,json=oldPeers" json:"old_peers,omitempty"`
}

func (m *ConfigurationPBMeta) Reset()                    { *m = ConfigurationPBMeta{} }
func (m *ConfigurationPBMeta) String() string            { return proto.CompactTextString(m) }
func (*ConfigurationPBMeta) ProtoMessage()               {}
func (*ConfigurationPBMeta) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ConfigurationPBMeta) GetPeers() []string {
	if m != nil {
		return m.Peers
	}
	return nil
}

func (m *ConfigurationPBMeta) GetOldPeers() []string {
	if m != nil {
		return m.OldPeers
	}
	return nil
}

type LogPBMeta struct {
	FirstLogIndex int64 `protobuf:"varint,1,opt,name=first_log_index,json=firstLogIndex" json:"first_log_index,omitempty"`
}

func (m *LogPBMeta) Reset()                    { *m = LogPBMeta{} }
func (m *LogPBMeta) String() string            { return proto.CompactTextString(m) }
func (*LogPBMeta) ProtoMessage()               {}
func (*LogPBMeta) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *LogPBMeta) GetFirstLogIndex() int64 {
	if m != nil {
		return m.FirstLogIndex
	}
	return 0
}

type StablePBMeta struct {
	Term     int64  `protobuf:"varint,1,opt,name=term" json:"term,omitempty"`
	Votedfor string `protobuf:"bytes,2,opt,name=votedfor" json:"votedfor,omitempty"`
}

func (m *StablePBMeta) Reset()                    { *m = StablePBMeta{} }
func (m *StablePBMeta) String() string            { return proto.CompactTextString(m) }
func (*StablePBMeta) ProtoMessage()               {}
func (*StablePBMeta) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *StablePBMeta) GetTerm() int64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *StablePBMeta) GetVotedfor() string {
	if m != nil {
		return m.Votedfor
	}
	return ""
}

type LocalSnapshotPbMeta struct {
	Meta  *SnapshotMeta               `protobuf:"bytes,1,opt,name=meta" json:"meta,omitempty"`
	Files []*LocalSnapshotPbMeta_File `protobuf:"bytes,2,rep,name=files" json:"files,omitempty"`
}

func (m *LocalSnapshotPbMeta) Reset()                    { *m = LocalSnapshotPbMeta{} }
func (m *LocalSnapshotPbMeta) String() string            { return proto.CompactTextString(m) }
func (*LocalSnapshotPbMeta) ProtoMessage()               {}
func (*LocalSnapshotPbMeta) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *LocalSnapshotPbMeta) GetMeta() *SnapshotMeta {
	if m != nil {
		return m.Meta
	}
	return nil
}

func (m *LocalSnapshotPbMeta) GetFiles() []*LocalSnapshotPbMeta_File {
	if m != nil {
		return m.Files
	}
	return nil
}

type LocalSnapshotPbMeta_File struct {
	Name string         `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Meta *LocalFileMeta `protobuf:"bytes,2,opt,name=meta" json:"meta,omitempty"`
}

func (m *LocalSnapshotPbMeta_File) Reset()                    { *m = LocalSnapshotPbMeta_File{} }
func (m *LocalSnapshotPbMeta_File) String() string            { return proto.CompactTextString(m) }
func (*LocalSnapshotPbMeta_File) ProtoMessage()               {}
func (*LocalSnapshotPbMeta_File) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5, 0} }

func (m *LocalSnapshotPbMeta_File) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *LocalSnapshotPbMeta_File) GetMeta() *LocalFileMeta {
	if m != nil {
		return m.Meta
	}
	return nil
}

type RequestVoteRequest struct {
	GroupId      string `protobuf:"bytes,1,opt,name=group_id,json=groupId" json:"group_id,omitempty"`
	ServerId     string `protobuf:"bytes,2,opt,name=server_id,json=serverId" json:"server_id,omitempty"`
	PeerId       string `protobuf:"bytes,3,opt,name=peer_id,json=peerId" json:"peer_id,omitempty"`
	Term         int64  `protobuf:"varint,4,opt,name=term" json:"term,omitempty"`
	LastLogTerm  int64  `protobuf:"varint,5,opt,name=last_log_term,json=lastLogTerm" json:"last_log_term,omitempty"`
	LastLogIndex int64  `protobuf:"varint,6,opt,name=last_log_index,json=lastLogIndex" json:"last_log_index,omitempty"`
}

func (m *RequestVoteRequest) Reset()                    { *m = RequestVoteRequest{} }
func (m *RequestVoteRequest) String() string            { return proto.CompactTextString(m) }
func (*RequestVoteRequest) ProtoMessage()               {}
func (*RequestVoteRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *RequestVoteRequest) GetGroupId() string {
	if m != nil {
		return m.GroupId
	}
	return ""
}

func (m *RequestVoteRequest) GetServerId() string {
	if m != nil {
		return m.ServerId
	}
	return ""
}

func (m *RequestVoteRequest) GetPeerId() string {
	if m != nil {
		return m.PeerId
	}
	return ""
}

func (m *RequestVoteRequest) GetTerm() int64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *RequestVoteRequest) GetLastLogTerm() int64 {
	if m != nil {
		return m.LastLogTerm
	}
	return 0
}

func (m *RequestVoteRequest) GetLastLogIndex() int64 {
	if m != nil {
		return m.LastLogIndex
	}
	return 0
}

type RequestVoteResponse struct {
	Term    int64 `protobuf:"varint,1,opt,name=term" json:"term,omitempty"`
	Granted bool  `protobuf:"varint,2,opt,name=granted" json:"granted,omitempty"`
}

func (m *RequestVoteResponse) Reset()                    { *m = RequestVoteResponse{} }
func (m *RequestVoteResponse) String() string            { return proto.CompactTextString(m) }
func (*RequestVoteResponse) ProtoMessage()               {}
func (*RequestVoteResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *RequestVoteResponse) GetTerm() int64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *RequestVoteResponse) GetGranted() bool {
	if m != nil {
		return m.Granted
	}
	return false
}

type AppendEntriesRequest struct {
	GroupId        string       `protobuf:"bytes,1,opt,name=group_id,json=groupId" json:"group_id,omitempty"`
	ServerId       string       `protobuf:"bytes,2,opt,name=server_id,json=serverId" json:"server_id,omitempty"`
	PeerId         string       `protobuf:"bytes,3,opt,name=peer_id,json=peerId" json:"peer_id,omitempty"`
	Term           int64        `protobuf:"varint,4,opt,name=term" json:"term,omitempty"`
	PrevLogTerm    int64        `protobuf:"varint,5,opt,name=prev_log_term,json=prevLogTerm" json:"prev_log_term,omitempty"`
	PrevLogIndex   int64        `protobuf:"varint,6,opt,name=prev_log_index,json=prevLogIndex" json:"prev_log_index,omitempty"`
	Entries        []*EntryMeta `protobuf:"bytes,7,rep,name=entries" json:"entries,omitempty"`
	CommittedIndex int64        `protobuf:"varint,8,opt,name=committed_index,json=committedIndex" json:"committed_index,omitempty"`
}

func (m *AppendEntriesRequest) Reset()                    { *m = AppendEntriesRequest{} }
func (m *AppendEntriesRequest) String() string            { return proto.CompactTextString(m) }
func (*AppendEntriesRequest) ProtoMessage()               {}
func (*AppendEntriesRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *AppendEntriesRequest) GetGroupId() string {
	if m != nil {
		return m.GroupId
	}
	return ""
}

func (m *AppendEntriesRequest) GetServerId() string {
	if m != nil {
		return m.ServerId
	}
	return ""
}

func (m *AppendEntriesRequest) GetPeerId() string {
	if m != nil {
		return m.PeerId
	}
	return ""
}

func (m *AppendEntriesRequest) GetTerm() int64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *AppendEntriesRequest) GetPrevLogTerm() int64 {
	if m != nil {
		return m.PrevLogTerm
	}
	return 0
}

func (m *AppendEntriesRequest) GetPrevLogIndex() int64 {
	if m != nil {
		return m.PrevLogIndex
	}
	return 0
}

func (m *AppendEntriesRequest) GetEntries() []*EntryMeta {
	if m != nil {
		return m.Entries
	}
	return nil
}

func (m *AppendEntriesRequest) GetCommittedIndex() int64 {
	if m != nil {
		return m.CommittedIndex
	}
	return 0
}

type AppendEntriesResponse struct {
	Term         int64 `protobuf:"varint,1,opt,name=term" json:"term,omitempty"`
	Success      bool  `protobuf:"varint,2,opt,name=success" json:"success,omitempty"`
	LastLogIndex int64 `protobuf:"varint,3,opt,name=last_log_index,json=lastLogIndex" json:"last_log_index,omitempty"`
}

func (m *AppendEntriesResponse) Reset()                    { *m = AppendEntriesResponse{} }
func (m *AppendEntriesResponse) String() string            { return proto.CompactTextString(m) }
func (*AppendEntriesResponse) ProtoMessage()               {}
func (*AppendEntriesResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *AppendEntriesResponse) GetTerm() int64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *AppendEntriesResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *AppendEntriesResponse) GetLastLogIndex() int64 {
	if m != nil {
		return m.LastLogIndex
	}
	return 0
}

type SnapshotMeta struct {
	LastIncludedIndex int64    `protobuf:"varint,1,opt,name=last_included_index,json=lastIncludedIndex" json:"last_included_index,omitempty"`
	LastIncludedTerm  int64    `protobuf:"varint,2,opt,name=last_included_term,json=lastIncludedTerm" json:"last_included_term,omitempty"`
	Peers             []string `protobuf:"bytes,3,rep,name=peers" json:"peers,omitempty"`
	OldPeers          []string `protobuf:"bytes,4,rep,name=old_peers,json=oldPeers" json:"old_peers,omitempty"`
}

func (m *SnapshotMeta) Reset()                    { *m = SnapshotMeta{} }
func (m *SnapshotMeta) String() string            { return proto.CompactTextString(m) }
func (*SnapshotMeta) ProtoMessage()               {}
func (*SnapshotMeta) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *SnapshotMeta) GetLastIncludedIndex() int64 {
	if m != nil {
		return m.LastIncludedIndex
	}
	return 0
}

func (m *SnapshotMeta) GetLastIncludedTerm() int64 {
	if m != nil {
		return m.LastIncludedTerm
	}
	return 0
}

func (m *SnapshotMeta) GetPeers() []string {
	if m != nil {
		return m.Peers
	}
	return nil
}

func (m *SnapshotMeta) GetOldPeers() []string {
	if m != nil {
		return m.OldPeers
	}
	return nil
}

type InstallSnapshotRequest struct {
	GroupId  string        `protobuf:"bytes,1,opt,name=group_id,json=groupId" json:"group_id,omitempty"`
	ServerId string        `protobuf:"bytes,2,opt,name=server_id,json=serverId" json:"server_id,omitempty"`
	PeerId   string        `protobuf:"bytes,3,opt,name=peer_id,json=peerId" json:"peer_id,omitempty"`
	Term     int64         `protobuf:"varint,4,opt,name=term" json:"term,omitempty"`
	Meta     *SnapshotMeta `protobuf:"bytes,5,opt,name=meta" json:"meta,omitempty"`
	Uri      string        `protobuf:"bytes,6,opt,name=uri" json:"uri,omitempty"`
}

func (m *InstallSnapshotRequest) Reset()                    { *m = InstallSnapshotRequest{} }
func (m *InstallSnapshotRequest) String() string            { return proto.CompactTextString(m) }
func (*InstallSnapshotRequest) ProtoMessage()               {}
func (*InstallSnapshotRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *InstallSnapshotRequest) GetGroupId() string {
	if m != nil {
		return m.GroupId
	}
	return ""
}

func (m *InstallSnapshotRequest) GetServerId() string {
	if m != nil {
		return m.ServerId
	}
	return ""
}

func (m *InstallSnapshotRequest) GetPeerId() string {
	if m != nil {
		return m.PeerId
	}
	return ""
}

func (m *InstallSnapshotRequest) GetTerm() int64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *InstallSnapshotRequest) GetMeta() *SnapshotMeta {
	if m != nil {
		return m.Meta
	}
	return nil
}

func (m *InstallSnapshotRequest) GetUri() string {
	if m != nil {
		return m.Uri
	}
	return ""
}

type InstallSnapshotResponse struct {
	Term    int64 `protobuf:"varint,1,opt,name=term" json:"term,omitempty"`
	Success bool  `protobuf:"varint,2,opt,name=success" json:"success,omitempty"`
}

func (m *InstallSnapshotResponse) Reset()                    { *m = InstallSnapshotResponse{} }
func (m *InstallSnapshotResponse) String() string            { return proto.CompactTextString(m) }
func (*InstallSnapshotResponse) ProtoMessage()               {}
func (*InstallSnapshotResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *InstallSnapshotResponse) GetTerm() int64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *InstallSnapshotResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type TimeoutNowRequest struct {
	GroupId  string `protobuf:"bytes,1,opt,name=group_id,json=groupId" json:"group_id,omitempty"`
	ServerId string `protobuf:"bytes,2,opt,name=server_id,json=serverId" json:"server_id,omitempty"`
	PeerId   string `protobuf:"bytes,3,opt,name=peer_id,json=peerId" json:"peer_id,omitempty"`
	Term     int64  `protobuf:"varint,4,opt,name=term" json:"term,omitempty"`
}

func (m *TimeoutNowRequest) Reset()                    { *m = TimeoutNowRequest{} }
func (m *TimeoutNowRequest) String() string            { return proto.CompactTextString(m) }
func (*TimeoutNowRequest) ProtoMessage()               {}
func (*TimeoutNowRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *TimeoutNowRequest) GetGroupId() string {
	if m != nil {
		return m.GroupId
	}
	return ""
}

func (m *TimeoutNowRequest) GetServerId() string {
	if m != nil {
		return m.ServerId
	}
	return ""
}

func (m *TimeoutNowRequest) GetPeerId() string {
	if m != nil {
		return m.PeerId
	}
	return ""
}

func (m *TimeoutNowRequest) GetTerm() int64 {
	if m != nil {
		return m.Term
	}
	return 0
}

type TimeoutNowResponse struct {
	Term    int64 `protobuf:"varint,1,opt,name=term" json:"term,omitempty"`
	Success bool  `protobuf:"varint,2,opt,name=success" json:"success,omitempty"`
}

func (m *TimeoutNowResponse) Reset()                    { *m = TimeoutNowResponse{} }
func (m *TimeoutNowResponse) String() string            { return proto.CompactTextString(m) }
func (*TimeoutNowResponse) ProtoMessage()               {}
func (*TimeoutNowResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *TimeoutNowResponse) GetTerm() int64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *TimeoutNowResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func init() {
	proto.RegisterType((*LocalFileMeta)(nil), "graft.LocalFileMeta")
	proto.RegisterType((*EntryMeta)(nil), "graft.EntryMeta")
	proto.RegisterType((*ConfigurationPBMeta)(nil), "graft.ConfigurationPBMeta")
	proto.RegisterType((*LogPBMeta)(nil), "graft.LogPBMeta")
	proto.RegisterType((*StablePBMeta)(nil), "graft.StablePBMeta")
	proto.RegisterType((*LocalSnapshotPbMeta)(nil), "graft.LocalSnapshotPbMeta")
	proto.RegisterType((*LocalSnapshotPbMeta_File)(nil), "graft.LocalSnapshotPbMeta.File")
	proto.RegisterType((*RequestVoteRequest)(nil), "graft.RequestVoteRequest")
	proto.RegisterType((*RequestVoteResponse)(nil), "graft.RequestVoteResponse")
	proto.RegisterType((*AppendEntriesRequest)(nil), "graft.AppendEntriesRequest")
	proto.RegisterType((*AppendEntriesResponse)(nil), "graft.AppendEntriesResponse")
	proto.RegisterType((*SnapshotMeta)(nil), "graft.SnapshotMeta")
	proto.RegisterType((*InstallSnapshotRequest)(nil), "graft.InstallSnapshotRequest")
	proto.RegisterType((*InstallSnapshotResponse)(nil), "graft.InstallSnapshotResponse")
	proto.RegisterType((*TimeoutNowRequest)(nil), "graft.TimeoutNowRequest")
	proto.RegisterType((*TimeoutNowResponse)(nil), "graft.TimeoutNowResponse")
	proto.RegisterEnum("graft.EntryType", EntryType_name, EntryType_value)
	proto.RegisterEnum("graft.ErrorType", ErrorType_name, ErrorType_value)
	proto.RegisterEnum("graft.FileSource", FileSource_name, FileSource_value)
}

func init() { proto.RegisterFile("raft.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 1010 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x56, 0xcd, 0x6e, 0xdb, 0x46,
	0x10, 0x36, 0xf5, 0x63, 0x49, 0x63, 0x5b, 0xa6, 0x57, 0x76, 0x2c, 0x2b, 0x69, 0x6b, 0x10, 0x45,
	0xa3, 0x1a, 0x85, 0x0f, 0x0e, 0x7a, 0x0d, 0x20, 0x2b, 0xb4, 0x4d, 0x44, 0xa1, 0xd4, 0x15, 0xdd,
	0x22, 0x27, 0x82, 0x16, 0x57, 0x0a, 0x11, 0x8a, 0xcb, 0x2e, 0x97, 0x4e, 0x7c, 0xed, 0x2b, 0xf4,
	0xd6, 0x4b, 0x5f, 0xa3, 0x40, 0x8f, 0x7d, 0xac, 0x5e, 0x8a, 0x5d, 0xae, 0x58, 0xca, 0x66, 0x5b,
	0x20, 0x05, 0x7c, 0xe3, 0xcc, 0x7c, 0xfb, 0xcd, 0xee, 0x37, 0xb3, 0xb3, 0x04, 0x60, 0xde, 0x9c,
	0x9f, 0xc6, 0x8c, 0x72, 0x8a, 0xea, 0x0b, 0x61, 0x18, 0x09, 0xec, 0x8c, 0xe8, 0xcc, 0x0b, 0x2f,
	0x82, 0x90, 0xbc, 0x21, 0xdc, 0x43, 0x4f, 0xa1, 0x95, 0x26, 0x84, 0xb9, 0x4b, 0xc2, 0xbd, 0xae,
	0x76, 0xac, 0xf5, 0xb7, 0x71, 0x53, 0x38, 0x64, 0xf0, 0x6b, 0xd8, 0x4c, 0x68, 0xca, 0x66, 0xa4,
	0x5b, 0x39, 0xd6, 0xfa, 0xed, 0xb3, 0xbd, 0x53, 0xc9, 0x72, 0x2a, 0x56, 0x4f, 0x65, 0x00, 0x2b,
	0x00, 0xea, 0x41, 0x73, 0xf6, 0x8e, 0xcc, 0xde, 0x27, 0xe9, 0xb2, 0x5b, 0x3d, 0xd6, 0xfa, 0x2d,
	0x9c, 0xdb, 0xc6, 0xcf, 0x1a, 0xb4, 0xcc, 0x88, 0xb3, 0x3b, 0x49, 0x8a, 0xa0, 0xc6, 0x09, 0x5b,
	0xca, 0x64, 0x55, 0x2c, 0xbf, 0xd1, 0x97, 0x50, 0xe3, 0x77, 0xf1, 0x2a, 0x8d, 0xae, 0xd2, 0xc8,
	0x35, 0xce, 0x5d, 0x4c, 0xb0, 0x8c, 0xa2, 0x7d, 0xa8, 0xc7, 0x84, 0xb0, 0xa4, 0x5b, 0x3d, 0xae,
	0xf6, 0x5b, 0x38, 0x33, 0xd0, 0x11, 0x34, 0x7d, 0x8f, 0x7b, 0x6e, 0x48, 0xa2, 0x6e, 0x4d, 0x72,
	0x36, 0x84, 0x3d, 0x22, 0x91, 0x38, 0x1c, 0x0d, 0x7d, 0x37, 0x5b, 0x54, 0x97, 0x8b, 0x9a, 0x34,
	0xf4, 0x27, 0xc2, 0x36, 0xae, 0xa0, 0x33, 0xa4, 0xd1, 0x3c, 0x58, 0xa4, 0xcc, 0xe3, 0x01, 0x8d,
	0x26, 0xe7, 0x72, 0x7b, 0x79, 0x12, 0xad, 0x98, 0x64, 0x8d, 0xa9, 0x72, 0x8f, 0xe9, 0x05, 0xb4,
	0x46, 0x74, 0xa1, 0xd6, 0x7f, 0x05, 0xbb, 0xf3, 0x80, 0x25, 0xdc, 0x0d, 0xe9, 0xc2, 0x0d, 0x22,
	0x9f, 0x7c, 0x54, 0x27, 0xdd, 0x91, 0xee, 0x11, 0x5d, 0x58, 0xc2, 0x69, 0xbc, 0x84, 0xed, 0x29,
	0xf7, 0x6e, 0x42, 0xa2, 0xd6, 0x95, 0xc9, 0xd2, 0x83, 0xe6, 0x2d, 0xe5, 0xc4, 0x9f, 0x53, 0x26,
	0xa5, 0x69, 0xe1, 0xdc, 0x36, 0x7e, 0xd7, 0xa0, 0x23, 0x4b, 0x39, 0x8d, 0xbc, 0x38, 0x79, 0x47,
	0xf9, 0xe4, 0x46, 0xf2, 0x3c, 0x87, 0x5a, 0x5e, 0xcb, 0xad, 0xb3, 0x8e, 0x92, 0x72, 0x05, 0x12,
	0x10, 0x2c, 0x01, 0xe8, 0x5b, 0xa8, 0xcf, 0x83, 0x90, 0x64, 0xc7, 0xd9, 0x3a, 0xfb, 0x42, 0x21,
	0x4b, 0x38, 0x65, 0xbd, 0x71, 0x86, 0xee, 0xbd, 0x82, 0x9a, 0x30, 0xc5, 0x7e, 0x23, 0x6f, 0x49,
	0x64, 0x9e, 0x16, 0x96, 0xdf, 0xa8, 0xaf, 0x72, 0x57, 0x64, 0xee, 0xfd, 0x22, 0xe3, 0xaa, 0xe1,
	0xb2, 0xe4, 0xc6, 0x1f, 0x1a, 0x20, 0x4c, 0x7e, 0x4c, 0x49, 0xc2, 0xbf, 0xa7, 0x9c, 0xa8, 0x4f,
	0x51, 0xcb, 0x05, 0xa3, 0x69, 0xec, 0x06, 0xbe, 0x22, 0x6e, 0x48, 0xdb, 0xf2, 0x45, 0x05, 0x12,
	0xc2, 0x6e, 0x09, 0x13, 0x31, 0x25, 0x46, 0xe6, 0xb0, 0x7c, 0x74, 0x08, 0x0d, 0x51, 0x1a, 0x11,
	0xca, 0x9a, 0x6f, 0x53, 0x98, 0x96, 0x9f, 0xab, 0x5a, 0x2b, 0xa8, 0x6a, 0xc0, 0x4e, 0xe8, 0xa9,
	0x02, 0xc9, 0x60, 0x5d, 0x06, 0xb7, 0x84, 0x73, 0x44, 0x17, 0x4e, 0xd6, 0x90, 0xed, 0x1c, 0x93,
	0x15, 0x71, 0x53, 0x82, 0xb6, 0x15, 0x28, 0xab, 0xe1, 0x10, 0x3a, 0x6b, 0x87, 0x48, 0x62, 0x1a,
	0x25, 0xa4, 0xb4, 0x94, 0x5d, 0x68, 0x2c, 0x98, 0x17, 0x71, 0x92, 0x6d, 0xbe, 0x89, 0x57, 0xa6,
	0xf1, 0x4b, 0x05, 0xf6, 0x07, 0x71, 0x4c, 0x22, 0x5f, 0xf4, 0x7b, 0x40, 0x92, 0xc7, 0x16, 0x23,
	0x66, 0xe4, 0xf6, 0x81, 0x18, 0xc2, 0x59, 0x10, 0x23, 0xc7, 0xac, 0x89, 0xa1, 0x40, 0x52, 0x0c,
	0x74, 0x02, 0x0d, 0x92, 0x1d, 0xa0, 0xdb, 0x90, 0x1d, 0xb5, 0x76, 0x8d, 0x65, 0xed, 0x57, 0x00,
	0xf4, 0x1c, 0x76, 0x67, 0x74, 0xb9, 0x0c, 0x38, 0x27, 0xbe, 0xa2, 0x6c, 0x4a, 0xca, 0x76, 0xee,
	0xce, 0x14, 0x7e, 0x0f, 0x07, 0xf7, 0xb4, 0xf9, 0x77, 0x8d, 0x93, 0x74, 0x36, 0x23, 0x49, 0xb2,
	0xd2, 0x58, 0x99, 0x25, 0xe5, 0xac, 0x96, 0x94, 0xf3, 0x57, 0x0d, 0xb6, 0x8b, 0x17, 0x05, 0x9d,
	0x42, 0x47, 0x2e, 0x0b, 0xa2, 0x59, 0x98, 0xfa, 0xf9, 0x56, 0xb3, 0x9c, 0x7b, 0x22, 0x64, 0xa9,
	0x48, 0x26, 0xc1, 0x37, 0x80, 0xd6, 0xf1, 0x72, 0x8b, 0x15, 0x09, 0xd7, 0x8b, 0x70, 0x29, 0x6b,
	0xf9, 0x38, 0x5b, 0x9b, 0x34, 0xb5, 0x7b, 0x93, 0xe6, 0x37, 0x0d, 0x9e, 0x58, 0x51, 0xc2, 0xbd,
	0x30, 0xbf, 0xa2, 0x8f, 0xda, 0x2d, 0xab, 0xe1, 0x52, 0xff, 0xaf, 0xe1, 0xa2, 0x43, 0x35, 0x65,
	0x81, 0xec, 0x93, 0x16, 0x16, 0x9f, 0xc6, 0x25, 0x1c, 0x3e, 0xd8, 0xf9, 0xa7, 0xd4, 0xd2, 0xf8,
	0x08, 0x7b, 0x4e, 0xb0, 0x24, 0x34, 0xe5, 0x36, 0xfd, 0xf0, 0x98, 0xa7, 0x37, 0xce, 0x01, 0x15,
	0x33, 0x7f, 0xca, 0xee, 0x4f, 0x42, 0xf5, 0x14, 0x8a, 0x67, 0x0d, 0x3d, 0x01, 0x64, 0xda, 0x0e,
	0x7e, 0xeb, 0x3a, 0x6f, 0x27, 0xa6, 0x7b, 0x6d, 0xbf, 0xb6, 0xc7, 0x3f, 0xd8, 0xfa, 0x06, 0xda,
	0x07, 0xbd, 0xe0, 0xb7, 0xc7, 0xee, 0x78, 0xa2, 0x6b, 0xa8, 0x03, 0xbb, 0x05, 0xef, 0xab, 0x81,
	0x33, 0xd0, 0x2b, 0xe8, 0x19, 0x74, 0x0b, 0xce, 0xe1, 0xd8, 0xbe, 0xb0, 0x2e, 0xaf, 0xf1, 0xc0,
	0xb1, 0xc6, 0xb6, 0x5e, 0x3d, 0xf9, 0x49, 0xbc, 0xbc, 0x8c, 0x51, 0x26, 0xd3, 0x09, 0x02, 0x8c,
	0xc7, 0x78, 0x45, 0x6b, 0x9b, 0xfa, 0x06, 0x42, 0xd0, 0x2e, 0x38, 0x47, 0xe3, 0x4b, 0x5d, 0x43,
	0x07, 0xb0, 0x57, 0xf0, 0x4d, 0x9d, 0xc1, 0xf9, 0xc8, 0xd4, 0x2b, 0xe8, 0x10, 0x3a, 0x45, 0xb7,
	0x3d, 0x98, 0x4c, 0xaf, 0xc6, 0x8e, 0x5e, 0x95, 0x9b, 0x58, 0xc3, 0x3b, 0xa6, 0xfb, 0x66, 0x30,
	0xbc, 0xb2, 0x6c, 0x53, 0xaf, 0x9d, 0xbc, 0x04, 0xf8, 0xfb, 0x87, 0x41, 0x70, 0x5f, 0x58, 0x23,
	0xd3, 0x9d, 0x8e, 0xaf, 0xf1, 0x50, 0x24, 0x1c, 0x0e, 0x46, 0xfa, 0x06, 0x3a, 0x82, 0x83, 0xa2,
	0x1b, 0x9b, 0x17, 0x26, 0x36, 0xed, 0xa1, 0xa9, 0x6b, 0x67, 0x7f, 0x56, 0x60, 0x0b, 0x7b, 0x73,
	0x3e, 0x25, 0xec, 0x36, 0x98, 0x11, 0x34, 0x80, 0x66, 0xcc, 0x88, 0x2b, 0x5e, 0x42, 0x74, 0xa4,
	0x5a, 0xf0, 0xe1, 0x5b, 0xd2, 0xeb, 0x95, 0x85, 0x54, 0xcd, 0x4c, 0xd8, 0x66, 0x99, 0xfb, 0x7f,
	0xd1, 0xbc, 0x86, 0xb6, 0x27, 0xa7, 0x93, 0xbb, 0x1a, 0x6c, 0x4f, 0x15, 0xba, 0x6c, 0xa0, 0xf7,
	0x9e, 0x95, 0x07, 0x15, 0xd9, 0x77, 0xa0, 0x07, 0xd9, 0x05, 0x71, 0x13, 0x75, 0x43, 0xd0, 0x67,
	0x6a, 0x45, 0xf9, 0x9d, 0xef, 0x7d, 0xfe, 0x4f, 0x61, 0x45, 0x79, 0x0e, 0x5b, 0x3c, 0x6b, 0x58,
	0x37, 0xa2, 0x1f, 0x50, 0x57, 0xc1, 0x1f, 0x5c, 0x9f, 0xde, 0x51, 0x49, 0x24, 0xe3, 0xb8, 0xd9,
	0x94, 0xff, 0x8f, 0x2f, 0xfe, 0x0a, 0x00, 0x00, 0xff, 0xff, 0x74, 0x01, 0x1a, 0x51, 0x4d, 0x0a,
	0x00, 0x00,
}
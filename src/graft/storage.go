package graft

//----------------------------Log Storage 接口---------------------------------//
type LogStorage interface {
	FirstLogIndex() int64
	LastLogIndex() int64
	GetEntry(index int64) *LogEntry
	GetTerm(index int64) int64
	AppendEntry(entry *LogEntry)
	AppendEntries(entries []*LogEntry)
	TruncatePrefix(firstIndexKept int64)
	TruncateSuffix(lastIndexKept int64)
	reset(firstIndexKept int64)
}



//----------------------------Meta Storage 接口---------------------------------//

type RaftMetaStorage interface {
	setTerm(term int64)
	getTerm() int64
	setVotedFor(peerId PeerId)
	getVotedFor(peerId PeerId)
	setTermAndVotedFor(term int64, peerId PeerId)
}
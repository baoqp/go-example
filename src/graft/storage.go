package graft

type LogStorage struct {

}

// TODO
func (this *LogStorage) FirstLogIndex() uint64 {
	return 0
}

// TODO
func (this *LogStorage) LastLogIndex() uint64 {
	return 0
}

// TODO
func (this *LogStorage) GetEntry(index uint64) *LogEntry {
	return nil
}

// TODO
func (this *LogStorage) GetTerm(index uint64) uint64 {
	return 0
}

// TODO
func (this *LogStorage) AppendEntry(index uint64) uint64 {
	return 0
}


// TODO
func (this *LogStorage) TruncatePrefix(firstIndexKept uint64) uint64 {
	return 0
}


// TODO
func (this *LogStorage) TruncateSuffix(lastIndexKept uint64) uint64 {
	return 0
}

// TODO
func (this *LogStorage) reset(firstIndexKept uint64) uint64 {
	return 0
}

// TODO
func CreateLogStorage(uri string) *LogStorage {

	return nil
}









type RaftMetaStorage struct  {

}
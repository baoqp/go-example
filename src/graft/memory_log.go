package graft

import (
	"sync"
	"sync/atomic"

	"fmt"
)

//---------------------------MemoryLogStorage 实现 LogStorage-----------------------------//

type MemoryData []*LogEntry

type MemeoryLogStorage struct {
	path          string
	firstLogIndex int64
	lastLogIndex  int64
	logEntryData  MemoryData
	mutex         *sync.Mutex
}

func NewMemoryLogStorage(uri string) *MemeoryLogStorage{
	memoryLogStorage := &MemeoryLogStorage{path:uri}
	memoryLogStorage.init()
	return memoryLogStorage
}

func (memoryLogStorage *MemeoryLogStorage) init() {
	atomic.StoreInt64(&memoryLogStorage.firstLogIndex, 1)
	atomic.StoreInt64(&memoryLogStorage.lastLogIndex, 0)
}

func (memoryLogStorage *MemeoryLogStorage) FirstLogIndex() int64 {
	return atomic.LoadInt64(&memoryLogStorage.firstLogIndex)
}

func (memoryLogStorage *MemeoryLogStorage) LastLogIndex() int64 {
	return atomic.LoadInt64(&memoryLogStorage.lastLogIndex)
}

func (memoryLogStorage *MemeoryLogStorage) GetEntry(index int64) *LogEntry {
	memoryLogStorage.mutex.Lock()
	defer memoryLogStorage.mutex.Unlock()

	if index < memoryLogStorage.firstLogIndex || index > memoryLogStorage.lastLogIndex {
		return nil
	}

	temp := memoryLogStorage.logEntryData[index-memoryLogStorage.firstLogIndex]

	msg := fmt.Sprintf("GetEntry entry index not equal. logentry index:%s, required_index: %d", temp.id.index, index)
	CHECK(index == temp.id.index, msg)
	return temp
}

func (memoryLogStorage *MemeoryLogStorage) GetTerm(index int64) int64 {
	memoryLogStorage.mutex.Lock()
	defer memoryLogStorage.mutex.Unlock()

	if index < memoryLogStorage.firstLogIndex || index > memoryLogStorage.lastLogIndex {
		return int64(0)
	}

	temp := memoryLogStorage.logEntryData[index-memoryLogStorage.firstLogIndex]

	msg := fmt.Sprintf("GetEntry entry index not equal. logentry index:%s, required_index: %d", temp.id.index, index)
	CHECK(index == temp.id.index, msg)
	return temp.id.term
}


func (memoryLogStorage *MemeoryLogStorage) AppendEntry(entry *LogEntry) {
	memoryLogStorage.mutex.Lock()
	defer memoryLogStorage.mutex.Unlock()

	if entry.id.index != memoryLogStorage.lastLogIndex + 1 {
		msg := fmt.Sprintf("append entry index =%d, lastLogIndex=%d, firstLogIndex=%d",
			entry.id.index, memoryLogStorage.lastLogIndex, memoryLogStorage.firstLogIndex)
		CHECK(false, msg)
	}

	memoryLogStorage.logEntryData = append(memoryLogStorage.logEntryData, entry)
	memoryLogStorage.lastLogIndex += 1

}

func (memoryLogStorage *MemeoryLogStorage) AppendEntries(entries []*LogEntry) {

	memoryLogStorage.mutex.Lock()
	defer memoryLogStorage.mutex.Unlock()


	if len(entries) == 0 {
		return
	}

	for _, entry := range entries {
		if entry.id.index != memoryLogStorage.lastLogIndex + 1 {
			msg := fmt.Sprintf("append entry index =%d, lastLogIndex=%d, firstLogIndex=%d",
				entry.id.index, memoryLogStorage.lastLogIndex, memoryLogStorage.firstLogIndex)
			CHECK(false, msg)
		}
		memoryLogStorage.logEntryData = append(memoryLogStorage.logEntryData, entry)
		memoryLogStorage.lastLogIndex += 1
	}
}

// TODO
func (memoryLogStorage *MemeoryLogStorage) TruncatePrefix(firstIndexKept int64)   {

}


// TODO
func (memoryLogStorage *MemeoryLogStorage) TruncateSuffix(lastIndexKept int64) {

}

// TODO
func (memoryLogStorage *MemeoryLogStorage) reset(firstIndexKept int64) {

}
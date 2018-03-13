package graft

import "image/draw"

type LogManagerOptions struct {
	logStorage           LogStorage
	configurationManager *ConfigurationManager
	fsmCaller            *FSMCaller // to report log error
}

type StableClosure struct {
	firstLogIndex int64
}

//-------------------------------------AppendBatcher-------------------------------//

type AppendBatcher struct {
}

type WaitId int64

//-------------------------------------LogManager-------------------------------//

type WaitMeta struct {
	errCode int
	arg     interface{}
	// int (*on_new_log)(void *arg, int error_code);
}

type LogManager struct {
	logStorage     LogStorage
	configManager  *ConfigurationManager
	fsmCaller      *FSMCaller
	waitMap        map[int64]WaitMeta
	stopped        bool
	hasError       bool
	nextWaitedId   WaitId
	diskId         LogId
	appliedId      LogId
	logsInMemory   []*LogEntry
	firstLogIndex  int64
	lastLogIndex   int64
	lastSnapShotId LogId

	// bthread::ExecutionQueueId<StableClosure*> _disk_queue;
}


func NewLogManager(options *LogManagerOptions) *LogManager {

	logManager := &LogManager{}

	logManager.init(options)

	return logManager
}


func (logManager *LogManager) init(options *LogManagerOptions)   {
	logManager.waitMap = make( map[int64]WaitMeta)

	if options != nil {
		if options.logStorage == nil {
			panic("logStorage in options is nil")
		}

		logManager.logStorage = options.logStorage
		logManager.configManager = options.configurationManager

		// TODO logManager.logStorage.init(configManager)
		logManager.firstLogIndex = logManager.logStorage.FirstLogIndex()
		logManager.lastLogIndex = logManager.logStorage.LastLogIndex()
		logManager.diskId.term = logManager.logStorage.GetTerm(logManager.lastLogIndex)
		logManager.diskId.index = logManager.lastLogIndex
		logManager.fsmCaller = options.fsmCaller
	}

}

// TODO
func (logManager *LogManager) shutdown() {

}

// TODO
func (logManager *LogManager) appendEntries(entries *[]LogEntry, done *StableClosure) {

}

// TODO
func (logManager *LogManager) clearBufferedLogs() {

}

// TODO
func (logManager *LogManager) getEntry(index int64) *LogEntry {
	return nil
}

// TODO
func (logManager *LogManager) getTerm(index *int64) int64 {
	return 0
}

// TODO
func (logManager *LogManager) getFirstLogIndex() int64 {
	return 0
}

// TODO
func (logManager *LogManager) getLastLogIndex() int64 {
	return 0
}

// TODO
func (logManager *LogManager) getConfiguration(index int64) *ConfigurationEntry {
	return nil
}

// Wait until there are more logs since |last_log_index| and |on_new_log|
// would be called after there are new logs or error occurs
// WaitId wait(int64_t expected_last_log_index, int (*on_new_log)(void *arg, int error_code), void *arg);

// TODO
func (logManager *LogManager) removeWaiter(id WaitId) int {
	return 0
}

// TODO
func (logManager *LogManager) setAppliedId(appliedId *LogId) {

}

//TODO
func (logManager *LogManager) checkConsistency() Status {
	return Status("")
}

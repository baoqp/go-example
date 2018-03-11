package graft

type LogManagerOptions struct {
	logStorage           *LogStorage
	configurationManager *ConfigurationManager
	fsmCaller            *FSMCaller // to report log error
}


type StableClosure struct {
	firstLogIndex uint64
}

//-------------------------------------AppendBatcher-------------------------------//


type WaitMeta struct {
	errCode int
	arg interface{}
	// int (*on_new_log)(void *arg, int error_code);
}


type AppendBatcher struct {

}

// TODO
func (this *AppendBatcher) appendToStorage(toAppend *[]LogEntry, lastId LogId) {

}



type WaitId uint64

//-------------------------------------LogManager-------------------------------//
type LogManager struct {

}

// TODO
func (this *LogManager) shutdown() {

}


// TODO
func (this *LogManager) appendEntries(entries *[]LogEntry, done *StableClosure) {

}

// TODO
func (this *LogManager) clearBufferedLogs() {

}

// TODO
func (this *LogManager) getEntry(index uint64) *LogEntry{
	return nil
}

// TODO
func (this *LogManager) getTerm(index *uint64) uint64 {
	return 0
}


// TODO
func (this *LogManager) firstLogIndex() uint64 {
	return 0
}

// TODO
func (this *LogManager) lastLogIndex() uint64 {
	return 0
}

// TODO
func (this *LogManager) getConfiguration(index uint64) *ConfigurationEntry {
	return nil
}


// Wait until there are more logs since |last_log_index| and |on_new_log|
// would be called after there are new logs or error occurs
// WaitId wait(int64_t expected_last_log_index, int (*on_new_log)(void *arg, int error_code), void *arg);



// TODO
func (this *LogManager) removeWaiter(id WaitId) int {
	return 0
}


// TODO
func (this *LogManager) setAppliedId(appliedId *LogId) {

}

//TODO
func (this *LogManager) checkConsistency() Status {

}


package graft

type LogId struct {
	index int64
	term  int64
}

// 相等
func (logId *LogId) EQ(other *LogId) bool {
	return logId.index == other.index && logId.term == other.term
}

// 小于
func (logId *LogId) LT(other *LogId) bool {
	if logId.term == other.term {
		return logId.index < other.index
	}

	return logId.term < other.term
}

// 大于
func (logId *LogId) GT(other *LogId) bool {
	if logId.term == other.term {
		return logId.index > other.index
	}

	return logId.term > other.term
}

// 小于等于
func (logId *LogId) LE(other *LogId) bool {
	return !logId.GT(other)
}

// 大于等于
func (logId *LogId) GE(other *LogId) bool {
	return !logId.LT(other)
}

// TODO LogId hash func ???

type LogEntry struct {
	typ EntryType
	id LogId
	peers []PeerId
	oldPeers []PeerId
	data []byte
}


package graft

type LogId struct {
	index uint64
	term  uint64
}

// 相等
func (this *LogId) EQ(other *LogId) bool {
	return this.index == other.index && this.term == other.term
}

// 小于
func (this *LogId) LT(other *LogId) bool {
	if this.term == other.term {
		return this.index < other.index
	}

	return this.term < other.term
}

// 大于
func (this *LogId) GT(other *LogId) bool {
	if this.term == other.term {
		return this.index > other.index
	}

	return this.term > other.term
}

// 小于等于
func (this *LogId) LE(other *LogId) bool {
	return !this.GT(other)
}

// 大于等于
func (this *LogId) GE(other *LogId) bool {
	return !this.LT(other)
}

// TODO LogId hash func ???

type LogEntry struct {
	typ EntryType
	id LogId
	peers []PeerId
	oldPeers []PeerId
	data []byte
}


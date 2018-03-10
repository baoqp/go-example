package graft

import (
	"math"
	"math/rand"
	"strings"
)

type Iterator struct {
	data []interface{}
	idx  int
}

func (this *Iterator) HasNext() bool {
	return this.idx < len(this.data)
}

func (this *Iterator) Next() interface{} {
	v := this.data[this.idx]
	this.idx ++
	return v
}


func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

const(
	raftElectionDelayMS = 100
	raftElectionHeartbeatFactor = 10
)

func RandomTimeout(timeoutMS int) int {
	delta := Min(timeoutMS, raftElectionDelayMS)

	return timeoutMS + rand.Intn(delta)
}

func HeartbeatTimeout(electionTimeout int) int {
	return Max(electionTimeout/raftElectionHeartbeatFactor, 10)
}

// ${protocol}://${parameters}
func ParseUri(uri string) (string, string){
	pos := strings.IndexAny(uri, "://")
	if pos == -1 {
		return ""
	}

	protocol := strings.TrimSpace(uri[:pos])
	// TODO 去除空格
	return protocol, uri[pos+3:]
}



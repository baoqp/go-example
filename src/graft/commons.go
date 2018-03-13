package graft

import (
	"math/rand"
	"strings"
	"fmt"
	"os"
)

//---------------------------type definination-------------------------//
type Stage int


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
		return "", uri
	}

	protocol := strings.TrimSpace(uri[:pos])
	// TODO 去除空格
	return protocol, uri[pos+3:]
}


func Assert(condition bool, msg string, v ...interface{}) {
	if !condition {
		panic(fmt.Sprintf("assertion failed: " + msg, v...))
	}
}

func CHECK(condition bool, msg string, v ...interface{}) {
	Assert(condition, msg, v)
}


//--------------------------------文件相关操作----------------------------------//
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
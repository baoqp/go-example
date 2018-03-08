package server

import (
	"github.com/chrislusf/raft"
	"bolt"
	"github.com/gorilla/mux"
	"os"
	"fmt"
	"errors"
	"sdkv/util"
)

const (
	defaultNameSpace = "default"
)

type KVServer struct {
	port       int
	RaftServer raft.Server
	db         *bolt.DB
}

func NewKVServer(r *mux.Router, port int, dataDir string) *KVServer {
	db := openBolt(dataDir + string(os.PathSeparator) + "kv.db")
	kvs := &KVServer{port: port, db: db}
	return kvs
}

func (kvServer *KVServer) Put(namespace string, key byte[], value byte[]...) {
	util.Assert(kvServer.db != nil)


}


func openBolt(path string) *bolt.DB {
	db, err := bolt.Open(path, 0600, nil);
	if err != nil {
		panic("open bolt db error")
	}
	return db
}

// 获取raft leader name
func (kvServer *KVServer) Leader() (string, error) {
	l := ""
	if kvServer.RaftServer != nil {
		l = kvServer.RaftServer.Leader()
	} else {
		return "", errors.New("raft server not ready yet")
	}

	if l == "" {
		// We are a single node cluster, we are the leader
		return kvServer.RaftServer.Name(), errors.New("raft server not initialized")
	}

	return l, nil
}

// 判断本server是否是leader
func (kvServer *KVServer) IsLeader() bool {
	if leader, e := kvServer.Leader(); e == nil {
		return leader == kvServer.RaftServer.Name()
	}
	return false
}

// 设置kvServer的raftServer
func (kvServer *KVServer) SetRaftServer(raftServer *RaftServer) {
	kvServer.RaftServer = raftServer.raftServer
	kvServer.RaftServer.AddEventListener(raft.LeaderChangeEventType, func(e raft.Event) {
		if kvServer.RaftServer.Leader() != "" {
			fmt.Println("[", kvServer.RaftServer.Name(), "]", kvServer.RaftServer.Leader(), "becomes leader.")
		}
	})
	if kvServer.IsLeader() {
		fmt.Println("[", kvServer.RaftServer.Name(), "]", "I am the leader!")
	} else {
		if kvServer.RaftServer.Leader() != "" {
			fmt.Println("[", kvServer.RaftServer.Name(), "]", kvServer.RaftServer.Leader(), "is the leader.")
		}
	}
}
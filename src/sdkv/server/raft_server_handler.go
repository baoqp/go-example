package server

import (
	"encoding/json"
	"strings"
	"net/http"
	"github.com/chrislusf/raft"
	"io/ioutil"
	"fmt"

	"sdkv/util"
)

// 处理 raft join 请求
func (raftServer *RaftServer) joinHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Processing incoming join. Current Leader", raftServer.raftServer.Leader(),
		"Self", raftServer.raftServer.Name(), "Peers", raftServer.raftServer.Peers())

	command := &raft.DefaultJoinCommand{}
	commandText, _ := ioutil.ReadAll(req.Body)
	fmt.Println("Command:", string(commandText))

	if err := json.NewDecoder(strings.NewReader(string(commandText))).Decode(&command); err != nil {
		fmt.Println("Error decoding json message:", err, string(commandText))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("join command from Name", command.Name, "Connection", command.ConnectionString)

	if _, err := raftServer.raftServer.Do(command); err != nil {
		switch err {
		case raft.NotLeaderError:
			raftServer.redirectToLeader(w, req)
		default:
			fmt.Println("Error processing join:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (raftServer *RaftServer) redirectToLeader(w http.ResponseWriter, req *http.Request) {

	if leader, e := raftServer.kvServer.Leader(); e == nil {
		//http.StatusMovedPermanently does not cause http POST following redirection
		leaderLocation := "http://" + leader + req.URL.Path
		fmt.Println("Redirecting to", leaderLocation)
		util.WriteJsonQuiet(w, req, http.StatusOK, leaderLocation)
	} else {
		fmt.Println("Error: Leader Unknown")
		http.Error(w, "Leader unknown", http.StatusInternalServerError)
	}
}

func (raftServer *RaftServer) statusHandler(w http.ResponseWriter, r *http.Request) {
	ret := ClusterStatusResult{
		IsLeader: raftServer.kvServer.IsLeader(),
		Peers:    raftServer.Peers(),
	}
	if leader, e := raftServer.kvServer.Leader(); e == nil {
		ret.Leader = leader
	}
	util.WriteJsonQuiet(w, r, http.StatusOK, ret)
}

type ClusterStatusResult struct {
	IsLeader bool     `json:"IsLeader,omitempty"`
	Leader   string   `json:"Leader,omitempty"`
	Peers    []string `json:"Peers,omitempty"`
}

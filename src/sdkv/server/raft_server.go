package server

import (
	"github.com/chrislusf/raft"
	"github.com/gorilla/mux"

	"time"
	"strings"
	"fmt"
	"net/http"
	"math/rand"
	"bytes"
	"encoding/json"
	"net/url"
	"io/ioutil"
	"errors"
)

type RaftServer struct {
	peers      []string // initial peers to join with
	raftServer raft.Server
	dataDir    string
	httpAddr   string
	router     *mux.Router
	kvServer   *KVServer
}

// 实现HTTPMuxer接口
func (raftServer *RaftServer) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	// 委托给router
	raftServer.router.HandleFunc(pattern, handler)
}

func NewRaftServer(r *mux.Router, peers []string, httpAddr string, dataDir string,
	kvServer *KVServer, pulseSeconds int) *RaftServer {

	s := &RaftServer{
		peers:    peers,
		httpAddr: httpAddr,
		dataDir:  dataDir,
		router:   r,
		kvServer: kvServer,
	}

	// 别忘了注册命令
	raft.RegisterCommand(&OPCommand{})

	var err error
	transporter := raft.NewHTTPTransporter("/cluster", 0)
	transporter.Transport.MaxIdleConnsPerHost = 1024
	fmt.Printf("Starting RaftServer with IP:%v:", httpAddr)

	s.raftServer, err = raft.NewServer(s.httpAddr, s.dataDir, transporter, nil, kvServer, "")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	transporter.Install(s.raftServer, s)

	s.raftServer.SetHeartbeatInterval(500 * time.Millisecond)
	s.raftServer.SetElectionTimeout(time.Duration(pulseSeconds) * 500 * time.Millisecond)
	s.raftServer.Start()

	s.router.HandleFunc("/cluster/join", s.joinHandler).Methods("POST")
	s.router.HandleFunc("/cluster/status", s.statusHandler).Methods("GET")

	if len(s.peers) > 0 {
		// 加入集群，不断重试知道成功
		for {
			fmt.Println("Joining cluster:", strings.Join(s.peers, ","))
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			firstJoinError := s.Join(s.peers)
			if firstJoinError != nil {
				fmt.Printf("No existing server found. Starting as leader in the new cluster.")
				fmt.Printf(" Current leader is %s", s.raftServer.Leader())
				_, err := s.raftServer.Do(&raft.DefaultJoinCommand{
					Name:             s.raftServer.Name(),
					ConnectionString: "http://" + s.httpAddr,
				})
				if err != nil {
					fmt.Println(err)
				} else {
					break
				}
			} else {
				break
			}
		}
	} else if s.raftServer.IsLogEmpty() { // 只有一个peer，“加入自身”

		fmt.Println("Initializing new cluster")

		_, err := s.raftServer.Do(&raft.DefaultJoinCommand{
			Name:             s.raftServer.Name(),
			ConnectionString: "http://" + s.httpAddr,
		})

		if err != nil {
			fmt.Println(err)
			return nil
		}

	}

	return s
}

func (raftServer *RaftServer) Peers() (members []string) {
	peers := raftServer.raftServer.Peers()

	for _, p := range peers {
		members = append(members, strings.TrimPrefix(p.ConnectionString, "http://"))
	}

	return
}

// 加入已存在的cluster
func (raftServer *RaftServer) Join(peers []string) error {
	command := &raft.DefaultJoinCommand{
		Name:             raftServer.raftServer.Name(),
		ConnectionString: "http://" + raftServer.httpAddr,
	}

	var err error
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(command)
	for _, m := range peers {
		if m == raftServer.httpAddr { // 跳过自身
			continue
		}
		target := fmt.Sprintf("http://%s/cluster/join", strings.TrimSpace(m))
		fmt.Println("Attempting to connect to:", target)

		err = postFollowingOneRedirect(target, "application/json", b)

		if err != nil {
			fmt.Println("Post returned error: ", err.Error())
			if _, ok := err.(*url.Error); ok {
				// If we receive a network error try the next member
				continue
			}
		} else {
			return nil
		}
	}

	return errors.New("could not connect to any cluster peers")
}

// 向raft peers发送请求时，如果目标peer不是leader,会返回leader的url地址，拿到这个地址后发起第二次请求
func postFollowingOneRedirect(target string, contentType string, b bytes.Buffer) error {
	backupReader := bytes.NewReader(b.Bytes())
	resp, err := http.Post(target, contentType, &b)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	statusCode := resp.StatusCode
	data, _ := ioutil.ReadAll(resp.Body)
	reply := string(data)

	if strings.HasPrefix(reply, "\"http") {
		urlStr := reply[1: len(reply)-1]

		fmt.Println("Post redirected to ", urlStr)
		resp2, err2 := http.Post(urlStr, contentType, backupReader)
		if err2 != nil {
			return err2
		}
		defer resp2.Body.Close()
		data, _ = ioutil.ReadAll(resp2.Body)
		statusCode = resp2.StatusCode
	}

	fmt.Println("Post returned status: ", statusCode, string(data))
	if statusCode != http.StatusOK {
		return errors.New(string(data))
	}

	return nil
}

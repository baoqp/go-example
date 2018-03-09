package server

import (
	"github.com/chrislusf/raft"
	"bolt"
	"github.com/gorilla/mux"
	"os"
	"fmt"
	"errors"
	"sdkv/util"
	"net/http"
	"net/http/httputil"
	"net/url"
	"io/ioutil"
	"time"
)

const (
	defaultBucket = "default_"
)

type KVServer struct {
	port       int
	RaftServer raft.Server
	db         *bolt.DB

	waitChan chan int
}

func NewKVServer(r *mux.Router, port int, dataDir string) *KVServer {
	db := openBolt(dataDir + string(os.PathSeparator) + "kv.db")
	kvs := &KVServer{port: port, db: db}
	kvs.waitChan = make(chan int, 16)

	r.HandleFunc("/kv/{bucket}/{key}", kvs.proxyToLeader(kvs.OPHandler)).Methods("PUT", "POST", "DELETE", "GET")

	return kvs
}

func (kvServer *KVServer) OPHandler(resp http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)

	method := req.Method

	bucket := vars["bucket"]
	key := vars["key"]

	var opCommand *OPCommand
	var ret interface{} // 返回值
	var err error

	switch method {
	case "GET": // GET请求不用走raft
		ret, err = kvServer.GET([]byte(bucket), []byte(key))
	case "PUT", "POST":
		value, err := ioutil.ReadAll(req.Body)
		if err != nil {
			resp.WriteHeader(http.StatusBadRequest)
			return
		}
		opCommand = NewOPCommand([]byte("PUT"), []byte(bucket), []byte(key), value)
		// Execute the command against the Raft server.
		ret, err = kvServer.RaftServer.Do(opCommand)
	case "DELETE":
		opCommand = NewOPCommand([]byte("DELETE"), []byte(bucket), []byte(key), nil)
		// Execute the command against the Raft server.
		ret, err = kvServer.RaftServer.Do(opCommand)
	}



	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	} else {
		if ret != nil {
			resp.Write(ret.([]byte))
		}
	}
}

func (kvServer *KVServer) Put(bucket []byte, key []byte, value []byte) error {

	util.Assert(kvServer.db != nil, "bolt db is nil")

	return kvServer.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		err = b.Put(key, value)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	})

}

func (kvServer *KVServer) GET(bucket []byte, key []byte) ([]byte, error) {
	util.Assert(kvServer.db != nil, "bolt db is nil")

	var value []byte
	err := kvServer.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		value = b.Get(key)
		if value == nil {
			return errors.New("has no value correspond to key" + string(key))
		}
		return nil
	})

	return value, err
}

func (kvServer *KVServer) Delete(bucket []byte, key []byte) error {
	util.Assert(kvServer.db != nil, "bolt db is nil")

	return kvServer.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Delete(key)
		if err != nil {
			return err
		}
		return nil
	})

}

// TODO　when to close boltdb
func openBolt(path string) *bolt.DB {
	db, err := bolt.Open(path, 0600, nil)
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

func (kvServer *KVServer) proxyToLeader(
	f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		if kvServer.IsLeader() {
			fmt.Println("self is leader, no need to redirect")
			f(w, r)
		} else if kvServer.RaftServer != nil && kvServer.RaftServer.Leader() != "" {
			// 把请求转发到leader，要等待leader处理结束，使用chan来做一个同步（chan是先进先出的???)
			start := time.Now().UTC().UnixNano()
			kvServer.waitChan <- 1
			defer func() {
				<-kvServer.waitChan
				end := time.Now().UTC().UnixNano()
				fmt.Printf("--kv op time cost: %d --", (end-start)/1000000)
			}()

			targetUrl, err := url.Parse("http://" + kvServer.RaftServer.Leader())
			if err != nil {
				util.WriteJsonError(w, r, http.StatusInternalServerError,
					fmt.Errorf("leader URL http://%s Parse Error: %v", kvServer.RaftServer.Leader(), err))
				return
			}
			fmt.Println("proxying to leader", kvServer.RaftServer.Leader())
			proxy := httputil.NewSingleHostReverseProxy(targetUrl)
			director := proxy.Director
			proxy.Director = func(req *http.Request) {
				actualHost, err := util.GetActualRemoteHost(req)
				if err == nil {
					req.Header.Set("HTTP_X_FORWARDED_FOR", actualHost)
				}
				director(req)
			}
			proxy.Transport = util.Transport
			proxy.ServeHTTP(w, r)
		} else {
			//drop it to the floor
			//writeJsonError(w, r, errors.New(ms.Topo.RaftServer.Name()+" does not know Leader yet:"+ms.Topo.RaftServer.Leader()))
		}
	}
}

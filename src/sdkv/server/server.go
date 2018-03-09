package server

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/soheilhy/cmux"
	"sdkv/util"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run(port int, peers string, dataDir string, pulseSeconds int) bool {

	r := mux.NewRouter()

	kvs := NewKVServer(r, port, dataDir)

	// 启动一个raft server

	myMasterAddress := "127.0.0.1" + ":" + strconv.Itoa(port)
	var peerArr []string
	if peers != "" {
		peerArr = strings.Split(peers, ",")
	}
	raftServer := NewRaftServer(r, peerArr, myMasterAddress, dataDir, kvs, pulseSeconds)
	kvs.SetRaftServer(raftServer)

	listeningAddress := "0.0.0.0" + ":" + strconv.Itoa(port)
	listener, _ := util.NewListener(listeningAddress, 0)
	m := cmux.New(listener)
	httpListener := m.Match(cmux.Any())
	httpServer := &http.Server{Handler: r}
	go httpServer.Serve(httpListener)

	if err := m.Serve(); err != nil {
		fmt.Printf("master server failed to serve: %v", err)
	}

	// TODO 好像并没有捕捉到信号
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-signalChan
		fmt.Println("sdkv exit...")
		kvs.db.Close()
	}()

	return true
}

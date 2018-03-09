package main

import (
	"flag"
	"sdkv/server"
	"os"
)

// TODO demo中用的goraft有问题，当整个集群停掉重新运行，leader选举不能进行，一直报错 Not current leader
// go run SDKVServer.go -port 7001 -peers 127.0.0.1:7001,127.0.0.1:7002,127.0.0.1:7003 -dataDir E:\tmp\sdkv\data1
func main() {
	port := flag.Int("port", 7001, "port")
	peers := flag.String("peers", "", "peers")
	dataDir := flag.String("dataDir", os.TempDir(), "dataDir")
	flag.Parse()
	pulseSeconds := 4
	server.Run(*port, *peers, *dataDir, pulseSeconds)
}

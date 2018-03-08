package main

import (
	"flag"
	"sdkv/server"
	"os"
)

func main() {
	port := flag.Int("port", 7001, "port")
	peers := flag.String("peers", "", "peers")
	dataDir := flag.String("dataDir", os.TempDir(), "dataDir")
	flag.Parse()
	pulseSeconds := 4
	server.Run(*port, *peers, *dataDir, pulseSeconds)
}

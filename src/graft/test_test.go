package graft

import (
	"testing"
	"fmt"
)

func Test(t *testing.T) {
	fmt.Println(ParseUri("http://www.baidu.com?a=1"))

	peerId := &PeerId{}
	peerId.parse("127.0.0.1:7001:10")
	fmt.Printf("%s  %d  %d", peerId.addr.ip, peerId.addr.port, peerId.idx)
}

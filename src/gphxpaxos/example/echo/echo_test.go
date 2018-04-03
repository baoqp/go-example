package echo

import (
	"testing"
	"strings"
	"gphxpaxos"
	"fmt"
)

func parseNode(nodeInfoStr string) (*gphxpaxos.NodeInfo, *gphxpaxos.NodeInfoList) {
	nodeArr := strings.Split(nodeInfoStr, " ")
	myNode := gphxpaxos.FromString(nodeArr[0])
	nodeListStr := nodeArr[1]
	nodeArr = strings.Split(nodeListStr, ",")
	nodeList := gphxpaxos.NodeInfoList{}
	for _, nodeInfo := range nodeArr {
		node := gphxpaxos.FromString(nodeInfo)
		nodeList = append(nodeList, node)
	}
	return myNode, &nodeList
}

func TestEcho(t *testing.T) {
	go echo1()
	go echo2()
	go echo3()
	select {}
}

func echo1() {
	myNode, nodeList := parseNode("127.0.0.1:11111 127.0.0.1:11111,127.0.0.1:11112,127.0.0.1:11113")

	echoServer := &EchoServer{myNode: *myNode, nodeList: *nodeList}

	err := echoServer.RunPaxos()
	fmt.Printf("myNodeIs is %d \r\n", echoServer.myNode.NodeId)

	if err != nil {
		fmt.Printf("run paxos failed, %v", err)
	} else {

		reqStr := fmt.Sprintf("req from %d", myNode.Port)
		respStr, err := echoServer.Echo(reqStr)
		if err != nil {
			fmt.Printf("Echo fail, ret %v \r\n", err)

		} else {
			fmt.Printf("Echo resp value %s \r\n", respStr)

		}
		select {}
	}

}

func echo2() {
	myNode, nodeList := parseNode("127.0.0.1:11112 127.0.0.1:11111,127.0.0.1:11112,127.0.0.1:11113")

	echoServer := &EchoServer{myNode: *myNode, nodeList: *nodeList}

	err := echoServer.RunPaxos()

	if err != nil {
		fmt.Printf("run paxos failed, %v", err)
	} else {
		select {}
	}
}

func echo3() {

	myNode, nodeList := parseNode("127.0.0.1:11113 127.0.0.1:11111,127.0.0.1:11112,127.0.0.1:11113")
	echoServer := &EchoServer{myNode: *myNode, nodeList: *nodeList}
	err := echoServer.RunPaxos()

	if err != nil {
		fmt.Printf("run paxos failed, %v", err)
	} else {
		select {}
	}
}

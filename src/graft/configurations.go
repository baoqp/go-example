package graft

import (
	"fmt"
	"strings"
	"strconv"
)

type GroupId string
type ReplicaId int

//---------------------------End Point-------------------------------//

type EndPoint struct {
	ip   string
	port int
}

func (endPoint *EndPoint) toString() string {
	return fmt.Sprintf("%s:%d", endPoint.ip, endPoint.port)
}

//---------------------------PeerId-------------------------------//
type PeerId struct {
	addr EndPoint
	idx  int
}

func (peerId *PeerId) reset() {
	peerId.addr.ip = IP_ANY
	peerId.addr.port = 0
	peerId.idx = 0
}

// TODO parse from peerId string
func (peerId *PeerId) parse(str string) {
	temp := strings.Split(str, ":")
	CHECK(len(temp) == 3, "peerId string has errors ")
	peerId.reset()
	fmt.Println(temp)
	peerId.addr.ip = temp[0]
	peerId.addr.port, _ = strconv.Atoi(temp[1])
	peerId.idx, _ = strconv.Atoi(temp[2])

}

func (peerId *PeerId) toString() string {
	return fmt.Sprintf("%s:%d", peerId.addr.toString(), peerId.idx)
}

// false
func (peerId *PeerId) LE(other *PeerId) bool {
	return false
}

type NodeId struct {
	groupId GroupId
	peedId  PeerId
}

type Configuration struct {
	peers []PeerId
}

func reset(this *Configuration) {

}

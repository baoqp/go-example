package graft

type GroupId string
type ReplicaId int

type EndPoint struct {
	ip   string
	port int
}

type PeerId struct {
	addr EndPoint
	idx  int // index in same addr, default 0
}

func (peerId *PeerId) reset() {
	peerId.addr.ip = IP_ANY
	peerId.addr.port = 0
	peerId.idx = 0
}

// false
func (peerId *PeerId) LE(other *PeerId) bool  {
	return false
}



type NodeId struct {
	groupId GroupId
	peedId PeerId
}


type Configuration struct {

	// set<PeerId> _peers
}


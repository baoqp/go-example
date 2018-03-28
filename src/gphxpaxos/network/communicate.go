package network

import "gphxpaxos/config"

// 实现MsgTransport接口
type Communicate struct {
	config   *config.Config
	network  NetWork
	myNodeId uint64
}

func NewCommunicate(config *config.Config, nodeId uint64, network NetWork) *Communicate {
	return &Communicate{
		config:   config,
		network:  network,
		myNodeId: nodeId,
	}
}


// TODO
func(c *Communicate) SendMessage(groupIdx int, sendToNodeId uint64, value []byte, sendType int) error {
	return nil
}

func(c *Communicate) BroadcastMessage(groupIdx int, value []byte, sendType int) error{
	return nil
}

func(c *Communicate) BroadcastMessageFollower(groupIdx int, value []byte, sendType int) error{
	return nil
}

func(c *Communicate) BroadcastMessageTempNode(groupIdx int, value []byte, sendType int) error{
	return nil
}

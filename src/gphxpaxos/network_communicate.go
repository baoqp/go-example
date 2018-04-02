package gphxpaxos

import (
	log "github.com/sirupsen/logrus"
	"fmt"
	"errors"
)

// 实现MsgTransport接口
type Communicate struct {
	cfg   *Config
	network  NetWork
	myNodeId uint64
}

func NewCommunicate(config *Config, nodeId uint64, network NetWork) *Communicate {
	return &Communicate{
		cfg:   config,
		network:  network,
		myNodeId: nodeId,
	}
}

func(c *Communicate) SendMessage(groupIdx int32, sendToNodeId uint64, value []byte, sendType int) error {
	MAX_VALUE_SIZE := GetMaxBufferSize()

	if len(value) > MAX_VALUE_SIZE {
		errMsg := fmt.Sprintf("Message size too large %d, max size %d, skip message",
			len(value), MAX_VALUE_SIZE)

		log.Errorf(errMsg)
		return errors.New(errMsg)
	}

	sendToNode := NewNodeInfoWithId(sendToNodeId) // TODO 获取nodeInfo， 可能有新节点会加入，所以可能需要从sysemVSM中获取
	return c.network.SendMessageTCP(groupIdx, sendToNode.Ip, sendToNode.Port, value)
}

// TODO 获取nodeInfo， 可能有新节点会加入，所以可能需要从sysemVSM中获取
func(c *Communicate) BroadcastMessage(groupIdx int32, value []byte, sendType int) error{
	for _, nodeInfo := range c.cfg.NodeInfoList {
		c.network.SendMessageTCP(groupIdx, nodeInfo.Ip, nodeInfo.Port, value)
	}
	return nil
}

func(c *Communicate) BroadcastMessageFollower(groupIdx int32, value []byte, sendType int) error{
	

	return nil
}

func(c *Communicate) BroadcastMessageTempNode(groupIdx int32, value []byte, sendType int) error{


	return nil
}

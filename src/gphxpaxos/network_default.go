package gphxpaxos

import (
	"net"
	"time"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb/errors"
	log "github.com/sirupsen/logrus"
)

var ErrBadConn = errors.New("connection was bad")

// 默认的network实现
type DefaultNetWork struct {
	node *Node
	end  bool
}

// start listening

// http://baijiahao.baidu.com/s?id=1566125616696554&wfr=spider&for=pc

func (dfNetWork *DefaultNetWork) RunNetWork() error {

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP("127.0.0.1"), 8080, ""})

	if err != nil {
		return errors.New("start listening failed")
	}


	return nil
}




func (dfNetWork *DefaultNetWork) StopNetWork() error {
	return nil
}

func (dfNetWork *DefaultNetWork) OnReceiveMessage(message []byte, messageLen int) error {

	if dfNetWork.node != nil {
		dfNetWork.node.OnReceiveMessage(message, messageLen)
	} else {
		log.Errorf("receive msglen %d but with node is nil", messageLen)
	}
	return nil
}

func (dfNetWork *DefaultNetWork) SendMessageTCP(groupIdx int32, ip string, port int, message []byte) error {
	conn, err := dfNetWork.connect(ip, port)

	if err != nil {
		return ErrBadConn
	}

	if n, err := conn.Write(message); err != nil || n != len(message) {
		return ErrBadConn
	}

	return nil
}

func (dfNetWork *DefaultNetWork) SendMessageUDP(groupIdx int32, ip string, port int, message []byte) error {
	dfNetWork.SendMessageTCP(groupIdx, ip, port, message)
}

// 维持长连接和重试
func (dfNetWork *DefaultNetWork) connect(ip string, port int) (net.Conn, error) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	return net.DialTimeout("tcp", addr, 5*time.Second)
}

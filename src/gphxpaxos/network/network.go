package network

const (
	Message_SendType_UDP = 0
	Message_SendType_TCP = 1
)



// 网络传输接口
type NetWork interface {

	RunNetWork() error

	StopNetWor() error

	SendMessageTCP(groupIdx int, ip string, port int, message []byte) error

	SendMessageUDP(groupIdx int, ip string, port int, message []byte) error

	OnReceiveMessage(message []byte, messageLen int) error
}

type MsgTransport interface {

	SendMessage(groupIdx int, sendToNodeId uint64, value []byte, sendType int) error

	BroadcastMessage(groupIdx int, value []byte, sendType int) error

	BroadcastMessageFollower(groupIdx int, value []byte, sendType int) error

	BroadcastMessageTempNode(groupIdx int, value []byte, sendType int) error
}



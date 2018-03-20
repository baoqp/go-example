package mysqlBinlogSync

import (
	"net"
	"bytes"
	"bufio"
	"errors"
	"io"
	"fmt"
	"time"
	"strings"
	"crypto/tls"
	log "github.com/sirupsen/logrus"
)

var (
	ErrBadConn       = errors.New("connection was bad")
	ErrMalformPacket = errors.New("Malform packet error")

	ErrTxDone = errors.New("sql: Transaction has already been committed or rolled back")
)

// copy from  https://github.com/siddontang/go-mysql.git
type InternalConn struct {
	net.Conn
	br       *bufio.Reader
	Sequence uint8
}

func NewInternalConn(conn net.Conn) *InternalConn {
	c := &InternalConn{}
	c.br = bufio.NewReaderSize(conn, 4096)
	c.Conn = conn
	return c
}

func (c *InternalConn) ReadPacket() ([]byte, error) {
	var buf bytes.Buffer

	if err := c.ReadPacketTo(&buf); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}

}

func (c *InternalConn) ReadPacketTo(w io.Writer) error {
	header := []byte{0, 0, 0, 0}

	if _, err := io.ReadFull(c.br, header); err != nil {
		return err
	}

	length := int(uint32(header[0]) | uint32(header[1])<<8 | uint32(header[2])<<16)

	if length < 1 {
		return fmt.Errorf("invalid payload length %d", length)
	}

	sequence := uint8(header[3])

	if sequence != c.Sequence {
		return fmt.Errorf("invalid sequence %d != %d", sequence, c.Sequence)
	}

	c.Sequence++

	if n, err := io.CopyN(w, c.br, int64(length)); err != nil {
		return ErrBadConn
	} else if n != int64(length) {
		return ErrBadConn
	} else {
		if length < MaxPayloadLen {
			return nil
		}

		// TODO 很大的数据会不会有问题
		if err := c.ReadPacketTo(w); err != nil {
			return err
		}
	}

	return nil
}

func (c *InternalConn) WritePacket(data []byte) error {
	length := len(data) - 4

	for length >= MaxPayloadLen {
		data[0] = 0xff
		data[1] = 0xff
		data[2] = 0xff

		data[3] = c.Sequence

		if n, err := c.Write(data[:4+MaxPayloadLen]); err != nil {
			return ErrBadConn
		} else if n != (4 + MaxPayloadLen) {
			return ErrBadConn
		} else {
			c.Sequence++
			length -= MaxPayloadLen
			data = data[MaxPayloadLen:]
		}
	}

	data[0] = byte(length)
	data[1] = byte(length >> 8)
	data[2] = byte(length >> 16)
	data[3] = c.Sequence

	if n, err := c.Write(data); err != nil {
		return ErrBadConn
	} else if n != len(data) {
		return ErrBadConn
	} else {
		c.Sequence++
		return nil
	}
}

func (c *InternalConn) ResetSequence() {
	c.Sequence = 0
}

func (c *InternalConn) Close() error {
	c.Sequence = 0
	if c != nil {
		return c.Close()
	}
	return nil
}

func getNetProto(addr string) string {
	proto := "tcp"
	if strings.Contains(addr, "/") {
		proto = "unix"
	}
	return proto
}

type Connection struct {
	*InternalConn
	*HandshakePacket

	user      string
	password  string
	db        string

	TLSConfig *tls.Config
}

func (connection *Connection) handshake() error {
	data, err := connection.ReadPacket()

	if err != nil {
		return err
	}

	if data[0] == ERR_HEADER {
		return errors.New("read initial handshake error")
	}

	// 读取mysql server发来的握手消息
	handshakePacket := NewHandshakePacket()
	if err = handshakePacket.read(data); err != nil {
		connection.Close()
		return err
	}
	connection.HandshakePacket = handshakePacket

	// TODO https://dev.mysql.com/doc/internals/en/capability-flags.html
	capability := CLIENT_PROTOCOL_41 | CLIENT_SECURE_CONNECTION |
		CLIENT_LONG_PASSWORD | CLIENT_TRANSACTIONS | CLIENT_LONG_FLAG

	capability &= connection.capabilityFlags

	// SSL Connection
	// https://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::SSLRequest
	if connection.TLSConfig != nil {

		capability |= CLIENT_PLUGIN_AUTH
		capability |= CLIENT_SSL

		data = make([]byte, 32+4)

		// capability 4 bytes
		data[4] = byte(capability)
		data[5] = byte(capability >> 8)
		data[6] = byte(capability >> 16)
		data[7] = byte(capability >> 24)
		// max-packet size ignore
		data[12] = byte(DEFAULT_COLLATION_ID)

		if err := connection.WritePacket(data); err != nil {
			return err
		}

		// Switch to TLS
		tlsConn := tls.Client(connection.InternalConn.Conn, connection.TLSConfig)
		if err := tlsConn.Handshake(); err != nil {
			return err
		}

		currentSequence := connection.Sequence
		connection.InternalConn = NewInternalConn(tlsConn)
		connection.Sequence = currentSequence
	}

	// 使用前一步读取的salt，加密账号密码，发送给mysql server进行验证
	authPacket := &AuthPacket{
		capabilityFlags:capability,
		salt:connection.salt,
		user:connection.user,
		passwd:connection.password,
		dbName:connection.db,
	}

	data, _ = authPacket.write()

	if err = connection.WritePacket(data); err != nil {
		connection.Close()
		return err
	}



	return nil
}

func Connect(addr string, user string, password string, dbName string) (*Connection, error) {
	proto := getNetProto(addr)

	c := &Connection{}
	var err error
	conn, err := net.DialTimeout(proto, addr, 10*time.Second)
	if err != nil {
		log.Error("connect mysql server error ", err)
		return nil, err
	}
	internalConn := NewInternalConn(conn)
	c.InternalConn = internalConn
	c.user = user
	c.password = password
	c.db = dbName

	//c.charset = DEFAULT_CHARSET


	if err = c.handshake(); err != nil {
		return nil, err
	}

	return c, nil
}

package ziface

import "net"

type IConnection interface {
	Start()
	Stop()
	GetTcpConnection() *net.TCPConn
	GetConnectionID() int
	RemoteAddr() net.Addr
	Send(msgID uint32, data []byte) error
}

package ziface

import "net"

type IConnection interface {
	Start()
	Stop()
	GetTcpConnection() *net.TCPConn
	GetConnectionID() int
	RemoteAddr() net.Addr
	GetServer() IServer
	Send(msgID uint32, data []byte) error
	SetProperty(key string, value any)
	GetProperty(key string) (any, error)
	RemoveProperty(key string)
}

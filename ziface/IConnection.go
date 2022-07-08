package ziface

import "net"

type IConnection interface {
	Start()
	Stop()
	GetTcpConnection() *net.TCPConn
	GetConnectionID() int
	RemoteAddr() net.Addr
	Send(data []byte) error
}
type HandleFunc func(conn *net.TCPConn, buf []byte, len int) error

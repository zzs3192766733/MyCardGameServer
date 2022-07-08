package znet

import (
	"MyGameServer/ziface"
)

type Request struct {
	Conn ziface.IConnection
	Data []byte
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.Conn
}

func (r *Request) GetData() []byte {
	return r.Data
}

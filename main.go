package main

import (
	"MyGameServer/logger"
	"MyGameServer/ziface"
	"MyGameServer/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) PreHandle(request ziface.IRequest) {
	logger.PopDebug("PreHandle...")
	_, err := request.GetConnection().GetTcpConnection().Write([]byte("Before Ping ..."))
	if err != nil {
		logger.PopError(err)
		return
	}
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	logger.PopDebug("Handle...")
	_, err := request.GetConnection().GetTcpConnection().Write([]byte("Ping Ping Ping ..."))
	if err != nil {
		logger.PopError(err)
		return
	}
}

func (p *PingRouter) PostHandle(request ziface.IRequest) {
	logger.PopDebug("PostHandle...")
	_, err := request.GetConnection().GetTcpConnection().Write([]byte("After Ping ..."))
	if err != nil {
		logger.PopError(err)
		return
	}
}

func main() {
	s := znet.NewServer("zzs")
	s.AddRouter(&PingRouter{})
	s.Serve()
}

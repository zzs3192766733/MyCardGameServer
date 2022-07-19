package router

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
	err := request.GetConnection().Send(1, []byte("Begin Ping..."))
	if err != nil {
		logger.PopError(err)
		return
	}
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	logger.PopDebug("Handle...")
	err := request.GetConnection().Send(2, []byte("Ping Ping Ping ..."))
	if err != nil {
		logger.PopError(err)
		return
	}
}

func (p *PingRouter) PostHandle(request ziface.IRequest) {
	logger.PopDebug("PostHandle...")
	err := request.GetConnection().Send(3, []byte("After Ping ..."))
	if err != nil {
		logger.PopError(err)
		return
	}
}

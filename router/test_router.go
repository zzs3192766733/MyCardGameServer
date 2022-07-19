package router

import (
	"MyGameServer/logger"
	"MyGameServer/ziface"
	"MyGameServer/znet"
)

type TestRouter struct {
	znet.BaseRouter
}

func (t *TestRouter) PreHandle(request ziface.IRequest) {
	logger.PopDebug("PreHandle...")
	err := request.GetConnection().Send(1001, []byte("Begin Ping..."))
	if err != nil {
		logger.PopError(err)
		return
	}
}

func (t *TestRouter) Handle(request ziface.IRequest) {
	logger.PopDebug("Handle...")
	err := request.GetConnection().Send(2002, []byte("Ping Ping Ping ..."))
	if err != nil {
		logger.PopError(err)
		return
	}
}

func (t *TestRouter) PostHandle(request ziface.IRequest) {
	logger.PopDebug("PostHandle...")
	err := request.GetConnection().Send(3003, []byte("After Ping ..."))
	if err != nil {
		logger.PopError(err)
		return
	}
}

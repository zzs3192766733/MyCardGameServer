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

func main() {
	s := znet.NewServer("zzs")
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &TestRouter{})
	s.Serve()
}

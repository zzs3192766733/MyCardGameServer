package router

import (
	"MyGameServer/logger"
	"MyGameServer/ziface"
	"MyGameServer/znet"
	"fmt"
)

type TestGetMongoDBRouter struct {
	znet.BaseRouter
}

func (t *TestGetMongoDBRouter) PreHandle(request ziface.IRequest) {

	msg := string(request.GetData())
	logger.PopDebug("服务器收到:%s发来的消息:%s", request.GetConnection().RemoteAddr().String(), msg)
	val, err := request.GetConnection().GetServer().(*znet.Server).DbGamePlayerHelper.FindAll()
	if err != nil {
		logger.PopError(err)
	}
	err = request.GetConnection().Send(996, []byte(fmt.Sprintf("一共有数据:%d", len(val))))
	if err != nil {
		logger.PopError(err)
	}

}

func (t *TestGetMongoDBRouter) Handle(request ziface.IRequest) {

	msg := string(request.GetData())
	logger.PopDebug("服务器收到:%s发来的消息:%s", request.GetConnection().RemoteAddr().String(), msg)
	val, err := request.GetConnection().GetServer().(*znet.Server).DbGamePlayerHelper.FindAll()
	if err != nil {
		logger.PopError(err)
	}
	err = request.GetConnection().Send(996, []byte(fmt.Sprintf("一共有数据:%d", len(val))))
	if err != nil {
		logger.PopError(err)
	}

}

func (t *TestGetMongoDBRouter) PostHandle(request ziface.IRequest) {

	msg := string(request.GetData())
	logger.PopDebug("服务器收到:%s发来的消息:%s", request.GetConnection().RemoteAddr().String(), msg)
	val, err := request.GetConnection().GetServer().(*znet.Server).DbGamePlayerHelper.FindAll()
	if err != nil {
		logger.PopError(err)
	}
	err = request.GetConnection().Send(996, []byte(fmt.Sprintf("一共有数据:%d", len(val))))
	if err != nil {
		logger.PopError(err)
	}
}

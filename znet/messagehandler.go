package znet

import (
	"MyGameServer/logger"
	"MyGameServer/ziface"
)

type MessageHandler struct {
	APIs map[uint32]ziface.IRouter
}

func NewMessageHandler() *MessageHandler {
	apis := make(map[uint32]ziface.IRouter)
	return &MessageHandler{APIs: apis}
}

func (m *MessageHandler) DoMessageHandle(request ziface.IRequest) {
	if router, ok := m.APIs[request.GetMsgID()]; ok {
		router.PreHandle(request)
		router.Handle(request)
		router.PostHandle(request)
	} else {
		logger.PopWarning("没有Router:%d,但是却被迫执行!", request.GetMsgID())
	}
}

func (m *MessageHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, ok := m.APIs[msgID]; ok {
		logger.PopDebug("重复添加Router")
		panic("重复添加Router")
	}
	m.APIs[msgID] = router
	logger.PopDebug("添加Router:%d成功!", msgID)
}

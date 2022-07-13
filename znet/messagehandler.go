package znet

import (
	"MyGameServer/logger"
	"MyGameServer/ziface"
)

type MessageHandler struct {
	APIs               map[uint32]ziface.IRouter
	TaskQueue          []chan ziface.IRequest
	WorkerPoolSize     uint32
	MaxOneWorkerLength uint32
}

func NewMessageHandler(workerPoolSize, maxOneWorkerLength uint32) *MessageHandler {
	apis := make(map[uint32]ziface.IRouter)
	return &MessageHandler{
		APIs:               apis,
		WorkerPoolSize:     workerPoolSize,
		MaxOneWorkerLength: maxOneWorkerLength,
		TaskQueue:          make([]chan ziface.IRequest, workerPoolSize),
	}
}

func (m *MessageHandler) SendMsgToTaskQueue(req ziface.IRequest) {
	workerID := req.GetConnection().GetConnectionID() % int(m.WorkerPoolSize)
	m.TaskQueue[workerID] <- req
}

func (m *MessageHandler) StartWorkerPool() {
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		m.TaskQueue[i] = make(chan ziface.IRequest, m.MaxOneWorkerLength)
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

func (m *MessageHandler) StartOneWorker(workerID int, workerChannel chan ziface.IRequest) {
	logger.PopDebug("Worker Start WorkerID:%d", workerID)
	for true {
		select {
		case req := <-workerChannel:
			m.DoMessageHandle(req)
		}
	}
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

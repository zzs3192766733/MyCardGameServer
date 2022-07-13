package ziface

type IMessageHandler interface {
	DoMessageHandle(request IRequest)
	AddRouter(msgID uint32, router IRouter)
	SendMsgToTaskQueue(req IRequest)
	StartWorkerPool()
	StartOneWorker(workerID int, workerChannel chan IRequest)
}

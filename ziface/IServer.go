package ziface

type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(msgID uint32, router IRouter)
	GetConnectionManager() IConnectionManager
	SetConnectionStart(func(conn IConnection))
	SetConnectionStop(func(conn IConnection))
	CallConnectionStart(conn IConnection)
	CallConnectionStop(conn IConnection)
}

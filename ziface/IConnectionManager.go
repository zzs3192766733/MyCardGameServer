package ziface

type IConnectionManager interface {
	AddConnection(conn IConnection)
	RemoveConnection(conn IConnection)
	GetConnectionsCount() int
	GetConnection(connID int) (IConnection, error)
	ClearAllConnection()
}

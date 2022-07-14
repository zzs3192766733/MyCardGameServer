package znet

import (
	"MyGameServer/logger"
	"MyGameServer/ziface"
	"errors"
	"fmt"
	"sync"
)

type ConnectionManager struct {
	connections map[int]ziface.IConnection
	lock        sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{connections: make(map[int]ziface.IConnection)}
}

func (c *ConnectionManager) AddConnection(conn ziface.IConnection) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.connections[conn.GetConnectionID()] = conn
	logger.PopDebug("AddConnection ConnID:%d To ConnectionManage", conn.GetTcpConnection())
}

func (c *ConnectionManager) RemoveConnection(conn ziface.IConnection) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.connections, conn.GetConnectionID())
	logger.PopDebug("RemoveConnection ConnID:%d From ConnectionManage", conn.GetConnectionID())
}

func (c *ConnectionManager) GetConnectionsCount() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.connections)
}

func (c *ConnectionManager) GetConnection(connID int) (ziface.IConnection, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New(fmt.Sprintf("connID :%d 不存在", connID))
	}
}

func (c *ConnectionManager) ClearAllConnection() {
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, connection := range c.connections {
		connection.Stop()
		delete(c.connections, connection.GetConnectionID())
	}
	logger.PopDebug("Clear All Connection From ConnectionManager")
}

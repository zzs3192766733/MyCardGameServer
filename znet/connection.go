package znet

import (
	"MyGameServer/logger"
	"MyGameServer/ziface"
	"net"
)

type Connection struct {
	Conn       *net.TCPConn
	ConnID     int
	IsClosed   bool
	ExitChan   chan bool
	HandleFunc ziface.HandleFunc
}

func NewConnection(conn *net.TCPConn, connID int, handleFunc ziface.HandleFunc) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		IsClosed:   false,
		ExitChan:   make(chan bool),
		HandleFunc: handleFunc,
	}
	return c
}

func (c *Connection) StartReader() {
	logger.PopDebug("Reader Goroutine Start...")
	defer logger.PopDebug("Reader is Exit Remote Addr is:%s", c.RemoteAddr().String())
	defer c.Stop()

	for true {
		buff := make([]byte, 512)
		count, err := c.Conn.Read(buff)
		if err != nil {
			logger.PopError(err)
			continue
		}
		err = c.HandleFunc(c.Conn, buff, count)
		if err != nil {
			logger.PopErrorInfo("HandleFun API Error ConnID:%d", c.ConnID)
			logger.PopError(err)
			break
		}
	}
}

func (c *Connection) Start() {
	logger.PopDebug("Conn Start ConnID:%d", c.ConnID)
	go c.StartReader()
}

func (c *Connection) Stop() {
	logger.PopDebug("Conn Stop ConnID:%d", c.ConnID)
	if c.IsClosed {
		return
	}
	c.IsClosed = true
	c.Conn.Close()
	close(c.ExitChan)
}

func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnectionID() int {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	return nil
}

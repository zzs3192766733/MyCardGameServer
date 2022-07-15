package znet

import (
	"MyGameServer/logger"
	"MyGameServer/ziface"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type Connection struct {
	Conn           *net.TCPConn
	ConnID         int
	IsClosed       bool
	ExitChan       chan bool
	MsgChan        chan []byte
	MsgHandler     ziface.IMessageHandler
	Server         ziface.IServer
	Properties     map[string]any
	PropertiesLock sync.RWMutex
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID int, msgHandler ziface.IMessageHandler) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		IsClosed:   false,
		ExitChan:   make(chan bool),
		MsgHandler: msgHandler,
		MsgChan:    make(chan []byte, 10),
		Server:     server,
		Properties: make(map[string]any),
	}
	server.GetConnectionManager().AddConnection(c)
	return c
}

func (c *Connection) SetProperty(key string, value any) {
	c.PropertiesLock.Lock()
	defer c.PropertiesLock.Unlock()
	c.Properties[key] = value
}
func (c *Connection) GetProperty(key string) (any, error) {
	c.PropertiesLock.RLock()
	defer c.PropertiesLock.RUnlock()
	if val, ok := c.Properties[key]; ok {
		return val, nil
	} else {
		return nil, errors.New(fmt.Sprintf("获取不存在的属性值:Key:%s", key))
	}
}
func (c *Connection) RemoveProperty(key string) {
	c.PropertiesLock.Lock()
	defer c.PropertiesLock.Unlock()
	delete(c.Properties, key)
}

func (c *Connection) StartReader() {
	logger.PopDebug("Reader Goroutine Start...")
	defer logger.PopDebug("Reader is Exit Remote Addr is:%s", c.RemoteAddr().String())
	defer c.Stop()

	for true {
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLength())
		_, err := io.ReadFull(c.GetTcpConnection(), headData)
		if err != nil {
			logger.PopError(err)
			return
		}

		msg, err := dp.UnPack(headData)
		if err != nil {
			logger.PopError(err)
			return
		}

		var buffer []byte

		if msg.GetMsgLen() > 0 {
			buffer = make([]byte, msg.GetMsgLen())
			_, err := io.ReadFull(c.GetTcpConnection(), buffer)
			if err != nil {
				logger.PopError(err)
				return
			}
		}
		msg.SetData(buffer)

		req := &Request{Conn: c, Msg: msg}

		c.MsgHandler.SendMsgToTaskQueue(req)
	}
}
func (c *Connection) StartWriter() {
	logger.PopDebug("Writer Goroutine Start...")
	defer logger.PopDebug("Writer is Exit Remote Addr is:%s", c.RemoteAddr().String())
	defer c.Stop()
	for true {
		select {
		case data := <-c.MsgChan:
			if _, err := c.GetTcpConnection().Write(data); err != nil {
				logger.PopError(err)
				logger.PopErrorInfo("Write Error Remote Addr is:%s", c.RemoteAddr().String())
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

func (c *Connection) Start() {
	logger.PopDebug("Conn Start ConnID:%d", c.ConnID)
	go c.StartReader()
	go c.StartWriter()

	c.Server.CallConnectionStart(c)
}

func (c *Connection) Stop() {
	if c.IsClosed {
		return
	}
	logger.PopDebug("Conn Stop ConnID:%d", c.ConnID)
	c.IsClosed = true
	c.ExitChan <- true
	c.Server.CallConnectionStop(c)
	c.Server.GetConnectionManager().RemoveConnection(c)
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

func (c *Connection) Send(msgID uint32, data []byte) error {
	if c.IsClosed {
		return errors.New("connection Closed When SendMsg")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgID, data))
	if err != nil {
		logger.PopError(err)
		return err
	}
	c.MsgChan <- msg
	return nil
}

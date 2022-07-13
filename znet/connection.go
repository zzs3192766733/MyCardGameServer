package znet

import (
	"MyGameServer/logger"
	"MyGameServer/ziface"
	"errors"
	"io"
	"net"
)

type Connection struct {
	Conn       *net.TCPConn
	ConnID     int
	IsClosed   bool
	ExitChan   chan bool
	MsgChan    chan []byte
	MsgHandler ziface.IMessageHandler
}

func NewConnection(conn *net.TCPConn, connID int, msgHandler ziface.IMessageHandler) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		IsClosed:   false,
		ExitChan:   make(chan bool),
		MsgHandler: msgHandler,
		MsgChan:    make(chan []byte, 10),
	}
	return c
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
}

func (c *Connection) Stop() {
	if c.IsClosed {
		return
	}
	logger.PopDebug("Conn Stop ConnID:%d", c.ConnID)
	c.IsClosed = true
	c.ExitChan <- true
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

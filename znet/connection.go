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
	MsgHandler ziface.IMessageHandler
}

func NewConnection(conn *net.TCPConn, connID int, msgHandler ziface.IMessageHandler) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		IsClosed:   false,
		ExitChan:   make(chan bool),
		MsgHandler: msgHandler,
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
			break
		}

		msg, err := dp.UnPack(headData)
		if err != nil {
			logger.PopError(err)
			break
		}

		var buffer []byte

		if msg.GetMsgLen() > 0 {
			buffer = make([]byte, msg.GetMsgLen())
			_, err := io.ReadFull(c.GetTcpConnection(), buffer)
			if err != nil {
				logger.PopError(err)
				break
			}
		}
		msg.SetData(buffer)

		req := &Request{Conn: c, Msg: msg}
		go c.MsgHandler.DoMessageHandle(req)
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
	_, err = c.GetTcpConnection().Write(msg)
	if err != nil {
		logger.PopError(err)
		return err
	}
	return nil
}

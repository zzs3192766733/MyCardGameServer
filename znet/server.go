package znet

import (
	"MyGameServer/logger"
	"MyGameServer/ziface"
	"errors"
	"fmt"
	"net"
)

type Server struct {
	Name       string
	IpType     string
	Ip         string
	Port       int
	MsgHandler ziface.IMessageHandler
}

func NewServer(serverName string) *Server {
	return &Server{
		Name:       serverName,
		IpType:     "tcp4",
		Ip:         "127.0.0.1",
		Port:       8888,
		MsgHandler: NewMessageHandler(10, 100),
	}
}

func CallBackToClient(conn *net.TCPConn, data []byte, len int) error {
	logger.PopDebug("服务器接收到客户端消息:%s", data[:len])
	if _, err := conn.Write(data[:len]); err != nil {
		logger.PopError(err)
		return errors.New("write Buff Error")
	}
	return nil
}

func (s *Server) Start() {
	go func() {

		s.MsgHandler.StartWorkerPool()

		addr, err := net.ResolveTCPAddr(s.IpType, fmt.Sprintf("%s:%d", s.Ip, s.Port))
		if err != nil {
			logger.PopError(err)
			return
		}
		listener, err := net.ListenTCP(s.IpType, addr)
		if err != nil {
			logger.PopError(err)
			return
		}
		defer listener.Close()
		logger.PopDebug("服务器开启成功!!!(%s:%d) ", s.Ip, s.Port)
		var cID = 0
		for true {
			conn, err := listener.AcceptTCP()
			if err != nil {
				logger.PopError(err)
				continue
			}

			dealConn := NewConnection(conn, cID, s.MsgHandler)
			cID++
			dealConn.Start()
		}
	}()
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
}

func (s *Server) Stop() {

}

func (s *Server) Serve() {
	//开启日志功能

	s.Start()
	select {}
}

package znet

import (
	"MyGameServer/logger"
	"MyGameServer/ziface"
	"fmt"
	"net"
)

type Server struct {
	Name               string
	IpType             string
	Ip                 string
	Port               int
	MsgHandler         ziface.IMessageHandler
	ConnectionManager  ziface.IConnectionManager
	MaxConnectionCount int
	OnConnStart        func(conn ziface.IConnection)
	OnConnStop         func(conn ziface.IConnection)
}

func NewServer(serverName string) *Server {
	return &Server{
		Name:               serverName,
		IpType:             "tcp4",
		Ip:                 "127.0.0.1",
		Port:               8888,
		MsgHandler:         NewMessageHandler(10, 100),
		ConnectionManager:  NewConnectionManager(),
		MaxConnectionCount: 1024,
	}
}

func (s *Server) SetConnectionStart(fun func(conn ziface.IConnection)) {
	s.OnConnStart = fun
}

func (s *Server) SetConnectionStop(fun func(conn ziface.IConnection)) {
	s.OnConnStop = fun
}

func (s *Server) CallConnectionStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		s.OnConnStart(conn)
	}
}

func (s *Server) CallConnectionStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		s.OnConnStop(conn)
	}
}

func (s *Server) GetConnectionManager() ziface.IConnectionManager {
	return s.ConnectionManager
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

			if s.GetConnectionManager().GetConnectionsCount() >= s.MaxConnectionCount {
				logger.PopWarning("服务器连接达到上限:%d", s.MaxConnectionCount)
				conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, cID, s.MsgHandler)
			cID++
			dealConn.Start()
		}
	}()
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
}

func (s *Server) Stop() {
	s.GetConnectionManager().ClearAllConnection()
}

func (s *Server) Serve() {
	//开启日志功能

	s.Start()
	select {}
}

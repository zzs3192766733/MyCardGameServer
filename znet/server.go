package znet

import (
	"MyGameServer/logger"
	"MyGameServer/mogodb"
	"MyGameServer/ziface"
	"fmt"
	"net"
)

const (
	MongoDB_Root                  = "mongodb://127.0.0.1:27017"
	MongoDB_Name                  = "test"
	MongoDB_Collection_GamePlayer = "Students"
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

	DbGamePlayerHelper *mogodb.MongoHelper //把数据库连接到服务器中
}

func NewServer(serverName string) *Server {
	server := &Server{
		Name:               serverName,
		IpType:             "tcp4",
		Ip:                 "127.0.0.1",
		Port:               8888,
		MsgHandler:         NewMessageHandler(10, 100),
		ConnectionManager:  NewConnectionManager(),
		MaxConnectionCount: 1024,
	}
	server.initAllDBCollection()
	return server
}

func (s *Server) initAllDBCollection() {
	s.DbGamePlayerHelper = mogodb.NewMongoHelper(MongoDB_Name, MongoDB_Collection_GamePlayer)
	if err := s.DbGamePlayerHelper.Connect(MongoDB_Root); err != nil {
		logger.PopErrorInfo("数据库表:{DbGamePlayerHelper}初始化失败!")
		logger.PopError(err)
	} else {
		logger.PopDebug("{DbGamePlayerHelper}:设置成功!")
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
	s.Start()
	select {}
}

package znet

import (
	"MyGameServer/logger"
	"MyGameServer/ziface"
	"errors"
	"fmt"
	"net"
)

type Server struct {
	Name   string
	IpType string
	Ip     string
	Port   int
	Router ziface.IRouter
}

func NewServer(serverName string) *Server {
	return &Server{
		Name:   serverName,
		IpType: "tcp4",
		Ip:     "127.0.0.1",
		Port:   8888,
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

			dealConn := NewConnection(conn, cID, s.Router)
			cID++
			dealConn.Start()
		}
	}()
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
}

func (s *Server) Stop() {

}

func (s *Server) Serve() {
	//开启日志功能

	s.Start()
	select {}
}

package znet

import (
	"MyGameServer/logger"
	"fmt"
	"net"
)

type Server struct {
	Name   string
	IpType string
	Ip     string
	Port   int
}

func NewServer(serverName string) *Server {
	return &Server{
		Name:   serverName,
		IpType: "tcp4",
		Ip:     "127.0.0.1",
		Port:   8888,
	}
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
		logger.PopDebug("服务器开启成功!!!(%s:%d) ", s.Ip, s.Port)
		for true {
			conn, err := listener.Accept()
			if err != nil {
				logger.PopError(err)
				continue
			}
			go func() {
				for true {
					buf := make([]byte, 512)
					count, err := conn.Read(buf)
					if err != nil {
						logger.PopError(err)
					}
					logger.PopDebug("服务器收到客户端的消息:%s", buf[:count])
					_, err = conn.Write(buf[:count])
					if err != nil {
						logger.PopError(err)
						continue
					}
				}
			}()
		}
	}()
}

func (s *Server) Stop() {

}

func (s *Server) Serve() {
	//开启日志功能

	s.Start()
	select {}
}

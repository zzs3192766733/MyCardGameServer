package main

import (
	"MyGameServer/logger"
	"MyGameServer/router"
	"MyGameServer/ziface"
	"MyGameServer/znet"
	"fmt"
)

func main() {

	gameServer := znet.NewServer("CardServer")
	initAllRouter(gameServer)
	initSetting(gameServer)
	gameServer.Serve()
}

func initSetting(server ziface.IServer) {
	server.SetConnectionStart(connectionStart)
	server.SetConnectionStop(connectionStop)
}

func initAllRouter(server ziface.IServer) {
	server.AddRouter(1, &router.PingRouter{})
	server.AddRouter(2, &router.TestRouter{})
	server.AddRouter(3, &router.TestGetMongoDBRouter{})
}

func connectionStart(conn ziface.IConnection) {
	logger.PopDebug("Start %d", conn.GetConnectionID())

	conn.SetProperty("name", "zzs")
	conn.SetProperty("account", 1234567)

}

func connectionStop(conn ziface.IConnection) {
	logger.PopDebug("Stop %d", conn.GetConnectionID())

	val1, err := conn.GetProperty("name")
	if err != nil {
		logger.PopError(err)
		return
	}
	fmt.Println(val1)
	val1, err = conn.GetProperty("account")
	if err != nil {
		logger.PopError(err)
		return
	}
	fmt.Println(val1)
}

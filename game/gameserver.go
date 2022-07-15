package game

import (
	"MyGameServer/mogodb"
	"MyGameServer/znet"
)

type ServerGame struct {
	*znet.Server
	DbGamePlayerHelper *mogodb.MongoHelper //把数据库连接到服务器中
}

func NewServerGame(serverName string) *ServerGame {
	return &ServerGame{
		Server:             znet.NewServer(serverName),
		DbGamePlayerHelper: mogodb.NewMongoHelper("", ""),
	}
}

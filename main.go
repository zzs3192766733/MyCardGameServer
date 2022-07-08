package main

import "MyGameServer/znet"

func main() {
	znet.NewServer("zzs").Serve()
}

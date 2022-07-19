package main

import (
	"MyGameServer/znet"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	fmt.Println("Client Start...")
	time.Sleep(time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("Client Start Error", err)
		return
	}
	defer conn.Close()

	//创建子Goroutine去读取数据
	go func() {
		for true {

			dp := znet.NewDataPack()
			headData := make([]byte, dp.GetHeadLength())
			_, err := io.ReadFull(conn, headData)
			if err != nil {
				fmt.Println(err)
				break
			}

			msg, err := dp.UnPack(headData)
			if err != nil {
				fmt.Println(err)
				break
			}

			var buffer []byte
			if msg.GetMsgLen() <= 0 {
				return
			}
			buffer = make([]byte, msg.GetMsgLen())
			_, err = io.ReadFull(conn, buffer)
			if err != nil {
				fmt.Println(err)
				break
			}
			msg.SetData(buffer)
			fmt.Println("接收服务器的消息,MsgID:", msg.GetMsgID(), "MsgData:", string(msg.GetData()))
		}
	}()

	//让主Goroutine阻塞,去写数据
	c := time.Tick(time.Second)
	for true {
		select {
		case t := <-c:
			if t.Unix()%2 == 0 {
				dp := znet.NewDataPack()
				buffer, err := dp.Pack(znet.NewMessage(3, []byte("向数据库请求数据!")))
				if err != nil {
					fmt.Println(err)
					return
				}
				_, err = conn.Write(buffer)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}

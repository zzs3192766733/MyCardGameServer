package main

import (
	"fmt"
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

	//创建子Goroutine去读取数据
	go func() {
		for true {
			buf := make([]byte, 512)
			count, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Read Error...", err)
				return
			}
			fmt.Printf("Server Call Back: %s, count: %d\n", buf[:count], count)
		}
	}()

	//让主Goroutine阻塞,去写数据
	for true {
		str := ""
		fmt.Scanln(&str)
		_, err := conn.Write([]byte(str))
		if err != nil {
			fmt.Println("Write Error...", err)
			return
		}
	}
}

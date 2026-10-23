package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func main() {
	//第一步，监听本机的9090端口
	config := net.ListenConfig{
		KeepAliveConfig: net.KeepAliveConfig{
			Enable:   true,
			Idle:     30 * time.Second,
			Interval: 10 * time.Second,
			Count:    3,
		},
	}
	listener, err := config.Listen(context.Background(), "tcp", "127.0.0.1:8080")
	if err != nil {
		log.Println("监听失败:", err)
		return
	}
	defer listener.Close()
	fmt.Println("服务端启动成功，等待连接")
	conn, err := listener.Accept()
	if err != nil {
		log.Println("接收连接失败：", err)
		return
	}
	defer conn.Close()
	fmt.Println("客户端已连接:", conn.RemoteAddr())

	//创建带缓冲的读取器，读取到换行符为止
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		log.Println("读取消息失败")
		return
	}
	message = strings.TrimSuffix(message, "\n")
	response := "echo: " + message + "\n"
	_, err = conn.Write([]byte(response))
	if err != nil {
		log.Println("回复失败：", err)
		return
	}
	fmt.Println("回复已发送:", response)
}

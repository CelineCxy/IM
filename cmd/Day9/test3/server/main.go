package main

import (
	"context"
	"errors"
	"im/cmd/Day9/test3/protocol"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	config := net.ListenConfig{
		KeepAliveConfig: net.KeepAliveConfig{
			Enable:   true,
			Idle:     30 * time.Second,
			Interval: 10 * time.Second,
			Count:    3,
		},
	}
	listener, err := config.Listen(context.Background(), "tcp", "127.0.0.1")
	if err != nil {
		log.Println("监听失败：", err)
		return
	}
	defer listener.Close()
	log.Println("长度前缀 Echo 服务端启动：", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("接收连接失败：", err)
			return
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	addr := conn.RemoteAddr().String()
	log.Printf("客户端 %v 已连接成功", addr)
	for {
		body, err := protocol.ReadMessage(conn)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("客户端%s已断开", addr)
				return
			} else {
				log.Printf("读取客户端%s消息失败:%v", addr, err)
				return
			}
		}
		log.Printf("客户端%s:长度:%d,内容:%q", addr, len(body), body)
		response := append([]byte("echo: "), body...)
		if err := protocol.WriteMessage(conn, response); err != nil {
			log.Printf("回复客户端%s失败:%v", addr, err)
			return
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		log.Println("监听失败：", err)
		return
	}
	defer listener.Close()

	log.Println("服务端启动，等待连接……")

	// 外层循环：持续接待新客户端。
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("接收连接失败：", err)
			return
		}

		// 每个连接交给一个独立的 goroutine。
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	addr := conn.RemoteAddr().String()
	log.Printf("客户端 %s 已连接", addr)
	// 同一个连接始终复用这个 reader。
	reader := bufio.NewReader(conn)
	// 内层循环：持续处理当前客户端的消息。
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("客户端 %s 已断开", addr)
			} else {
				log.Printf("读取客户端 %s 消息失败：%v", addr, err)
			}
			return
		}

		message = strings.TrimSuffix(message, "\n")
		log.Printf("收到客户端 %s 消息：%q", addr, message)

		response := fmt.Sprintf("echo: %s\n", message)
		if _, err := conn.Write([]byte(response)); err != nil {
			log.Printf("回复客户端 %s 失败：%v", addr, err)
			return
		}
	}
}

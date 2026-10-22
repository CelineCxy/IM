package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	//第一步，监听本机的9090端口
	listener, err := net.Listen("tcp", "127.0.0.1:9090")
	if err != nil {
		log.Println("监听失败：", err)
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

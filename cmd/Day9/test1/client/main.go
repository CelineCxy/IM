package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9090")
	if err != nil {
		fmt.Errorf("连接失败:%v", err)
		return
	}
	defer conn.Close()
	fmt.Println("连接成功:", conn.RemoteAddr())
	_, err = conn.Write([]byte("hello\n"))
	if err != nil {
		log.Printf("发送失败:%v", err)
		return
	}
	fmt.Println("已发送:hello")
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("数据读取失败:%v", err)
		return
	}
	response = strings.TrimSuffix(response, "\n")
	fmt.Println("服务端回复:", response)
}

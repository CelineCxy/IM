package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"im/cmd/Day9/test3/protocol"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		log.Println("连接失败：", err)
		return
	}
	defer conn.Close()
	fmt.Println("连接成功，输入消息后回车；输入 quit 退出。")

	scanner := bufio.NewScanner(os.Stdin)
	// 回复会增加 "echo: "，为它预留 6 字节。
	scanner.Buffer(make([]byte, 4096), protocol.MaxBodySize-6)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				log.Println("读取输入失败：", err)
			}
			return
		}
		message := scanner.Text()
		if message == "quit" {
			return
		}

		// 网络消息不再依靠换行符结束。
		if err := protocol.WriteMessage(conn, []byte(message)); err != nil {
			log.Println("发送失败：", err)
			return
		}
		response, err := protocol.ReadMessage(conn)
		if err != nil {
			log.Println("读取回复失败：", err)
			return
		}
		fmt.Println(string(response))
	}
}

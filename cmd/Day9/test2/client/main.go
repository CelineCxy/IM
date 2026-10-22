package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		log.Println("连接失败：", err)
		return
	}
	defer conn.Close()

	fmt.Println("连接成功，输入消息后回车；输入 quit 退出。")

	// scanner 读取键盘输入，reader 读取网络回复。
	scanner := bufio.NewScanner(os.Stdin)
	reader := bufio.NewReader(conn)

	for {
		fmt.Print("> ")

		// 等待用户输入一行。
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				log.Println("读取输入失败：", err)
			}
			return
		}

		// Text 返回的内容不包含行尾换行符。
		message := scanner.Text()

		if message == "quit" {
			fmt.Println("客户端退出")
			return
		}

		// 添加协议要求的换行符，再发送。
		if _, err := conn.Write([]byte(message + "\n")); err != nil {
			log.Println("发送失败：", err)
			return
		}

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Println("读取回复失败：", err)
			return
		}

		fmt.Println(strings.TrimSuffix(response, "\n"))
	}
}

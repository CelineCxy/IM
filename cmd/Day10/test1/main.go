package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 注册路由：GET /hello 请求交给 helloHandler 处理。
	http.HandleFunc("GET /hello", helloHandler)

	log.Println("服务启动：http://127.0.0.1:8080")

	// 启动服务；nil 表示使用默认路由器。
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		log.Fatal("HTTP 服务停止：", err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// r 保存当前请求的信息。
	log.Printf("收到请求：method=%s path=%s", r.Method, r.URL.Path)

	// 设置响应头，告诉客户端返回的是 UTF-8 文本。
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// 设置响应状态码。
	// w.WriteHeader(http.StatusOK)

	// 写入响应体。
	if _, err := fmt.Fprintln(w, "Hello HTTP"); err != nil {
		log.Println("写入响应失败：", err)
	}
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}", userHandler)
	mux.HandleFunc("GET /search", searchHandler)
	mux.HandleFunc("POST /echo", echoHandler)
	log.Println("服务启动:http://127.0.0.1:8080")
	if err := http.ListenAndServe("127.0.0.1:8080", mux); err != nil {
		log.Fatal("HTTP 服务停止：", err)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	idText := r.PathValue("id")
	id, err := strconv.Atoi(idText)
	if err != nil || id <= 0 {
		http.Error(w, "id must be a positive integer", http.StatusBadRequest)
		return

	}
	writeJSON(w, map[string]any{
		"id":      id,
		"message": "path parameter received",
	})
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// 访问 /search?keyword=go 时，得到字符串 "go"。
	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		http.Error(w, "keyword is required", http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]any{
		"keyword": keyword,
		"message": "query parameter received",
	})
}

// 定义客户端提交的JSON数据结构
type EchoRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	var req EchoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	writeJSON(w, req)
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("写入JSON响应失败:", err)
	}
}

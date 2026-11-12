package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// 客户端创建用户时，提交的数据
type UserReq struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// 服务端保存和返回的用户数据
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// 内存中的用户数据，以及保护共享数据的锁
var (
	usersMu sync.Mutex
	users         = make(map[int64]User)
	nextID  int64 = 1
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", userHandler)
	mux.HandleFunc("GET /users/{id}", searchHandler)
	mux.HandleFunc("DELETE /users/{id}", deleteHandler)
	log.Println("服务启动：http://127.0.0.1:8080")
	if err := http.ListenAndServe("127.0.0.1:8080", mux); err != nil {
		log.Fatal("HTTP 服务停止：", err)
	}
}
func searchHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 获取并解析路径参数。
	idText := r.PathValue("id")

	// Map 的 key 是 int64，直接解析为 int64。
	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "id must be a positive integer")
		return
	}

	// 2. 读取共享 Map 时也需要加锁。
	usersMu.Lock()
	user, exists := users[id]
	usersMu.Unlock()

	// 3. 用户不存在。
	if !exists {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	// 4. 查询成功。
	writeJSON(w, http.StatusOK, user)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 获取并解析路径参数。
	idText := r.PathValue("id")

	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "id must be a positive integer")
		return
	}

	// 2. 查找和删除必须在同一次加锁中完成。
	usersMu.Lock()
	user, exists := users[id]
	if exists {
		delete(users, id)
	}
	usersMu.Unlock()
	// 3. 用户不存在。
	if !exists {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	log.Printf("用户删除成功：%+v", user)
	// 4. 删除成功，返回 204，不写响应体。
	w.WriteHeader(http.StatusNoContent)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	var req UserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "消息请求解析失败", http.StatusBadRequest)
		return
	}
	// 请求体只允许包含一个 JSON 值。
	// 第一次解析完成后，再读取应当得到 EOF。
	var extra any
	if err := json.NewDecoder(r.Body).Decode(&extra); err != io.EOF {
		writeError(w, http.StatusBadRequest, "body must contain one JSON value")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if req.Age <= 0 || req.Age > 150 {
		writeError(w, http.StatusBadRequest, "age must be between 1 and 150")
		return
	}
	usersMu.Lock()
	defer usersMu.Unlock()
	user := User{
		ID:   nextID,
		Name: req.Name,
		Age:  req.Age,
	}
	users[user.ID] = user
	nextID++
	log.Printf("用户信息已保存:%+v", user)
	writeJSON(w, http.StatusOK, user)
}
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}
func writeJSON(w http.ResponseWriter, status int, data any) {
	// 先编码，再提交响应，避免编码失败后无法修改状态码。
	body, err := json.Marshal(data)
	if err != nil {
		log.Println("JSON 编码失败：", err)

		status = http.StatusInternalServerError
		body = []byte(`{"error":"internal server error"}`)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(body); err != nil {
		log.Println("写入响应失败：", err)
	}
}

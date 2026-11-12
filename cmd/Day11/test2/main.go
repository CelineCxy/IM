package main

import (
	"crypto/rand"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	router.Use(RequestID(), Logger(), TimeDuration(), Recovery())
	router.GET("/panic", func(c *gin.Context) {
		panic("测试异常")
	})
	v1 := router.Group("/api/v1")
	users := v1.Group("/users")
	users.POST("", createUser)
	users.GET("/:id", searchUser)
	users.DELETE("/:id", deleteUser)
	if err := router.Run("127.0.0.1:8080"); err != nil {
		log.Fatal("路由启动失败")
	}
}

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

var (
	users  = make(map[int64]User)
	mu     sync.Mutex
	nextID int64 = 1
)

func createUser(c *gin.Context) {
	var req UserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("解析 JSON 失败：%v", err)
		writeError(c, 400, "json数据转换为结构体数据失败")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Age <= 0 || req.Age > 150 || req.Name == "" {
		writeError(c, 400, "数据异常")
		return
	}
	mu.Lock()
	user := User{
		ID:   nextID,
		Name: req.Name,
		Age:  req.Age,
	}
	users[user.ID] = user
	nextID++
	mu.Unlock()
	writeSuccess(c, 201, "用户创建成功", user)
}

func searchUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeError(c, 400, "id parse failed")
		return
	}
	mu.Lock()
	_, ok := users[id]
	if ok {
		writeSuccess(c, 200, "用户查询成功", users[id])
		return
	}
	mu.Unlock()
	writeError(c, 404, "用户不存在")
}

func deleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeError(c, 400, "id parse failed")
		return
	}
	mu.Lock()
	_, exists := users[id]
	if exists {
		delete(users, id)
	}
	mu.Unlock()

	if !exists {
		writeError(c, 404, "用户不存在")
		return
	}

	writeSuccess(c, 200, "用户删除成功", nil)
}

func writeError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"message": message,
	})
}

func writeSuccess(c *gin.Context, code int, message string, data any) {
	c.JSON(code, gin.H{
		"message": message,
		"data":    data,
	})
}

func TimeDuration() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		c.Set("duration", time.Since(start))
	}
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := rand.Text()
		// 保存到当前请求的 Context，供日志等中间件读取。
		c.Set("request_id", id)
		// 返回给客户端，方便根据编号查日志。
		c.Header("X-Request-ID", id)

		c.Next()
	}
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		duration, _ := c.Get("duration")
		log.Printf(
			"request_id=%s method=%s path=%s status=%d duration=%v",
			c.GetString("request_id"),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
		)
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Printf(
					"request_id=%s panic=%v\n%s",
					c.GetString("request_id"),
					recovered,
					debug.Stack(),
				)
				c.Abort()
				if c.Writer.Written() {
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    50000,
					"message": "internal server error",
					"data":    nil,
				})

			}
		}()
		c.Next()
	}
}

package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	// 默认注册 Logger（请求日志）和 Recovery（捕获请求处理中的 panic）。
	router := gin.Default()

	// 路由的完整路径 = 分组前缀 + 注册时的路径。
	v1 := router.Group("/api/v1")
	users := v1.Group("/user")
	users.GET("/hello", helloHandler)
	users.GET("/:id", userHandler)
	v1.GET("/search", searchHandler)

	if err := router.Run("127.0.0.1:8080"); err != nil {
		log.Fatal("Gin 服务启动失败：", err)
	}
}

// GET /api/v1/user/hello
func helloHandler(c *gin.Context) {
	// JSON 同时设置状态码、JSON 响应头，并把数据编码到响应体。
	c.JSON(http.StatusOK, gin.H{
		"message": "hello",
	})
}

// GET /api/v1/user/12：路径中的 12 对应路由中的 :id。
func userHandler(c *gin.Context) {
	// Param 返回字符串；读取后仍需要转换和校验。
	idText := c.Param("id")
	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id must be a positive integer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"message": "path parameter received",
	})
}

// GET /api/v1/search?keyword=go&page=2
func searchHandler(c *gin.Context) {
	// Query 读取 ? 后的参数；参数不存在时返回空字符串。
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "keyword is required",
		})
		return
	}

	// DefaultQuery 在参数不存在时使用默认值；返回值仍是字符串。
	pageText := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageText)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "page must be a positive integer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"keyword": keyword,
		"page":    page,
	})
}

package router

import (
	"im/cmd/Day12/internal/handler"
	"im/cmd/Day12/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handler.UserHandler) *gin.Engine {
	r := gin.New()
	r.Use(
		middleware.RequestID(),
		middleware.Logger(),
		middleware.TimeDuration(),
		middleware.Recovery(),
	)
	r.GET("/panic", func(c *gin.Context) {
		panic("测试异常")
	})

	v1 := r.Group("/api/v1")
	users := v1.Group("/users")

	users.POST("", userHandler.CreateUser)
	users.GET("/:id", userHandler.GetUser)
	users.DELETE("/:id", userHandler.DeleteUser)

	return r
}

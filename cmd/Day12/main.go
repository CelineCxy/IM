package main

import (
	"im/cmd/Day12/internal/handler"
	"im/cmd/Day12/internal/repository"
	"im/cmd/Day12/internal/router"
	"im/cmd/Day12/internal/service"
	"log"
)

func main() {
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	r := router.SetupRouter(userHandler)
	log.Println("server is running at http://127.0.0.1:8080")
	if err := r.Run("127.0.0.1:8080"); err != nil {
		log.Fatal("服务启动失败：", err)
	}
}

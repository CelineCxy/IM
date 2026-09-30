package main

import (
	"im/internal/handler"
	"im/internal/repository"
	"im/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository)
	userhandler := handler.NewUserHandler(userService)
	router.GET("/user", userhandler.GetUser)
	log.Println("server is running")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

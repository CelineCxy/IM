package handler

import (
	"im/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler{
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetUser(c *gin.Context){
	idStr := c.Query("id")
	if idStr == ""{
		c.JSON(http.StatusBadRequest,gin.H{
			"code":400,
			"message":"id is requird",
		})
		return 
	}
	id,err := strconv.ParseInt(idStr,10,64)
	if err!=nil || id<=0{
		c.JSON(http.StatusBadRequest,gin.H{
			"code":400,
			"message":"id must be a positive integer",
		})
		return
	}
	user,ok := h.userService.GetUser(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "user not found",
		})
		return
	}
	c.JSON(http.StatusOK,user)
}
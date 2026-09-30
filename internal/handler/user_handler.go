package handler

import (
	"context"
	"errors"
	"im/internal/repository"
	"im/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "id is required",
		})
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "id must be a positive integer",
		})
		return
	}
	requestCtx := c.Request.Context()
	ctx, cancel := context.WithTimeout(requestCtx, 6*time.Second)
	defer cancel()
	user, err := h.userService.GetUser(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(http.StatusRequestTimeout, gin.H{
				"code":    408,
				"message": "request timeout",
			})

		case errors.Is(err, context.Canceled):
			// 客户端断开连接或请求被主动取消
			return

		case errors.Is(err, repository.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "user not found",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusOK, user)
}

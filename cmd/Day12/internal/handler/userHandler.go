package handler

import (
	"errors"
	"im/cmd/Day12/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserReq struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req UserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("数据解析失败：%v", err)
		writeError(c, http.StatusBadRequest, "JSON 数据格式错误")
		return
	}
	user, err := h.userService.CreateUser(req.Name, req.Age)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	writeSuccess(c, http.StatusCreated, "用户创建成功", user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, ok := parseUserID(c)
	if !ok {
		return
	}
	user, err := h.userService.GetUser(id)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	writeSuccess(c, http.StatusOK, "用户查询成功", user)
}
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, ok := parseUserID(c)
	if !ok {
		return
	}
	err := h.userService.DeleteUser(id)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	writeSuccess(c, http.StatusOK, "用户删除成功", nil)
}
func parseUserID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(c, http.StatusBadRequest, "用户 ID 必须是正整数")
		return 0, false
	}

	return id, true
}

// 将业务错误转换为 HTTP 响应。
func writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidUser):
		writeError(c, http.StatusBadRequest, err.Error())

	case errors.Is(err, service.ErrInvalidID):
		writeError(c, http.StatusBadRequest, err.Error())

	case errors.Is(err, service.ErrUserNotFound):
		writeError(c, http.StatusNotFound, err.Error())

	default:
		log.Printf(
			"request_id=%s service error=%v",
			c.GetString("request_id"),
			err,
		)
		writeError(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
	}
}

func writeError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"message": message,
		"data":    nil,
	})
}

func writeSuccess(c *gin.Context, status int, message string, data any) {
	c.JSON(status, gin.H{
		"message": message,
		"data":    data,
	})
}

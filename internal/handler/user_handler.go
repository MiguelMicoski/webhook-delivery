package handler

import (
	"awesomeProject/internal/service"
	"net/http"

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

func (h *UserHandler) FindAll(ctx *gin.Context) {
	users := h.userService.FindAll()

	ctx.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

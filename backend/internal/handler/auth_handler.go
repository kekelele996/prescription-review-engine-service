package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rxcheck/rxcheck/internal/dto"
	"github.com/rxcheck/rxcheck/internal/middleware"
	"github.com/rxcheck/rxcheck/internal/service"
	"github.com/rxcheck/rxcheck/internal/util"
)

// AuthHandler 认证处理器。
type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Login 登录。
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, "登录参数不合法", err))
		return
	}
	rid, _ := c.Get(middleware.RequestIDKey)
	result, err := h.svc.Login(&req, c.ClientIP(), ridString(rid))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, result)
}

// Register 自助注册医生账号。
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, "注册参数不合法", err))
		return
	}
	rid, _ := c.Get(middleware.RequestIDKey)
	user, err := h.svc.Register(&req, c.ClientIP(), ridString(rid))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, dto.ToUserResp(user))
}

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rxcheck/rxcheck/internal/handler"
)

// registerAuthRoutes 认证路由（登录/注册开放）。
func registerAuthRoutes(g *gin.RouterGroup, h *handler.AuthHandler) {
	auth := g.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/register", h.Register)
	}
}

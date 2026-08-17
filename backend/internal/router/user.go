package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/handler"
	"github.com/rxcheck/rxcheck/internal/middleware"
)

// registerUserRoutes 用户管理路由（仅管理员）。
func registerUserRoutes(g *gin.RouterGroup, h *handler.UserHandler) {
	users := g.Group("/users")
	{
		users.GET("", middleware.RequireRoles(constants.UserRoleAdmin), h.List)
		users.POST("", middleware.RequireRoles(constants.UserRoleAdmin), h.Create)
		users.PUT("/:id", middleware.RequireRoles(constants.UserRoleAdmin), h.Update)
	}
}

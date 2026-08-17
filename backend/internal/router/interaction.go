package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/handler"
	"github.com/rxcheck/rxcheck/internal/middleware"
)

// registerInteractionRoutes 相互作用规则路由（写操作仅管理员）。
func registerInteractionRoutes(g *gin.RouterGroup, h *handler.InteractionHandler) {
	interactions := g.Group("/interactions")
	{
		interactions.GET("", h.List)
		interactions.POST("", middleware.RequireRoles(constants.UserRoleAdmin), h.Create)
		interactions.DELETE("/:id", middleware.RequireRoles(constants.UserRoleAdmin), h.Delete)
	}
}

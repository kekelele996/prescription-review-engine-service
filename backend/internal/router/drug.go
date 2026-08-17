package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/handler"
	"github.com/rxcheck/rxcheck/internal/middleware"
)

// registerDrugRoutes 药品档案路由（读写分离，写操作仅管理员）。
func registerDrugRoutes(g *gin.RouterGroup, h *handler.DrugHandler) {
	drugs := g.Group("/drugs")
	{
		drugs.GET("", h.List)
		drugs.GET("/:id", h.Get)
		drugs.POST("", middleware.RequireRoles(constants.UserRoleAdmin), h.Create)
		drugs.PUT("/:id", middleware.RequireRoles(constants.UserRoleAdmin), h.Update)
		drugs.POST("/:id/disable", middleware.RequireRoles(constants.UserRoleAdmin), h.Disable)
		drugs.POST("/:id/enable", middleware.RequireRoles(constants.UserRoleAdmin), h.Enable)
	}
}

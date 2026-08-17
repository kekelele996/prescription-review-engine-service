package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/handler"
	"github.com/rxcheck/rxcheck/internal/middleware"
)

// registerAuditRoutes 审计日志路由（仅管理员）。
func registerAuditRoutes(g *gin.RouterGroup, h *handler.AuditHandler) {
	audits := g.Group("/audit-logs")
	{
		audits.GET("", middleware.RequireRoles(constants.UserRoleAdmin), h.List)
	}
}

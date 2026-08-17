package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rxcheck/rxcheck/internal/handler"
)

// registerStatsRoutes 统计路由。
func registerStatsRoutes(g *gin.RouterGroup, h *handler.StatsHandler) {
	stats := g.Group("/stats")
	{
		stats.GET("/overview", h.Overview)
	}
}

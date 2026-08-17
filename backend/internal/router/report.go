package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rxcheck/rxcheck/internal/handler"
)

// registerReportRoutes 审核报告列表路由。
func registerReportRoutes(g *gin.RouterGroup, h *handler.ReportHandler) {
	reports := g.Group("/reports")
	{
		reports.GET("", h.List)
	}
}

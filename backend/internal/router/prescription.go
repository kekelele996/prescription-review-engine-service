package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/handler"
	"github.com/rxcheck/rxcheck/internal/middleware"
)

// registerPrescriptionRoutes 处方路由：接收/列表/详情/复核/强制通过/报告。
func registerPrescriptionRoutes(g *gin.RouterGroup, h *handler.PrescriptionHandler, reportH *handler.ReportHandler) {
	prescriptions := g.Group("/prescriptions")
	{
		prescriptions.POST("", h.Submit)
		prescriptions.GET("", h.List)
		prescriptions.GET("/:id", h.Get)
		prescriptions.POST("/:id/review", h.Review)
		prescriptions.POST("/:id/override", middleware.RequireRoles(constants.UserRoleDoctor, constants.UserRolePharmacist, constants.UserRoleAdmin), h.Override)
		prescriptions.GET("/:id/report", reportH.GetByPrescription)
		prescriptions.POST("/:id/report/export", reportH.Export)
	}
}

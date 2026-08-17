package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rxcheck/rxcheck/internal/dto"
	"github.com/rxcheck/rxcheck/internal/middleware"
	"github.com/rxcheck/rxcheck/internal/service"
	"github.com/rxcheck/rxcheck/internal/util"
)

// ReportHandler 审核报告处理器。
type ReportHandler struct {
	svc *service.ReportService
}

func NewReportHandler(svc *service.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

// List 报告列表。
func (h *ReportHandler) List(c *gin.Context) {
	page := util.ParsePage(c.Query("page"))
	pageSize := util.ParsePageSize(c.Query("page_size"))
	result, err := h.svc.List(page, pageSize, c.Query("status"))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, result)
}

// GetByPrescription 按处方查询报告（与 GET /prescriptions/:id 复用处方仓储）。
func (h *ReportHandler) GetByPrescription(c *gin.Context) {
	var p dto.IDParam
	if err := c.ShouldBindUri(&p); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, "处方ID不合法", err))
		return
	}
	report, err := h.svc.GetByPrescription(p.ID)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, report)
}

// Export 导出报告快照到 MinIO。
func (h *ReportHandler) Export(c *gin.Context) {
	var p dto.IDParam
	if err := c.ShouldBindUri(&p); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, "处方ID不合法", err))
		return
	}
	op := middleware.CurrentUser(c)
	result, err := h.svc.Export(p.ID, op.Username, c.ClientIP(), rid(c))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, result)
}

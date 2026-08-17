package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rxcheck/rxcheck/internal/dto"
	"github.com/rxcheck/rxcheck/internal/middleware"
	"github.com/rxcheck/rxcheck/internal/service"
	"github.com/rxcheck/rxcheck/internal/util"
)

// PrescriptionHandler 处方处理器：接收（JSON/XML）、查询、复核、强制通过。
type PrescriptionHandler struct {
	svc *service.PrescriptionService
}

func NewPrescriptionHandler(svc *service.PrescriptionService) *PrescriptionHandler {
	return &PrescriptionHandler{svc: svc}
}

// Submit 接收电子处方（JSON/XML），自动解析并触发审核。
func (h *PrescriptionHandler) Submit(c *gin.Context) {
	req, format, raw, err := service.ParsePrescription(c)
	if err != nil {
		c.Error(err)
		return
	}
	op := middleware.CurrentUser(c)
	p, err := h.svc.Submit(req, format, raw, op, c.ClientIP(), rid(c))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, p)
}

// List 处方列表。
func (h *PrescriptionHandler) List(c *gin.Context) {
	page := util.ParsePage(c.Query("page"))
	pageSize := util.ParsePageSize(c.Query("page_size"))
	result, err := h.svc.List(page, pageSize, c.Query("status"), c.Query("keyword"))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, result)
}

// Get 处方详情（复用 PrescriptionService.Get，与 /report 共用 FindByID）。
func (h *PrescriptionHandler) Get(c *gin.Context) {
	var p dto.IDParam
	if err := c.ShouldBindUri(&p); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, "处方ID不合法", err))
		return
	}
	prescription, err := h.svc.Get(p.ID)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, prescription)
}

// Review 重新审核处方（复用 ReviewService.ReviewPrescription）。
func (h *PrescriptionHandler) Review(c *gin.Context) {
	var p dto.IDParam
	if err := c.ShouldBindUri(&p); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, "处方ID不合法", err))
		return
	}
	op := middleware.CurrentUser(c)
	report, err := h.svc.Review(p.ID, op.Username, c.ClientIP(), rid(c))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, dto.ToReportResp(report))
}

// Override 医生/药师强制通过警告类处方并记录原因。
func (h *PrescriptionHandler) Override(c *gin.Context) {
	var p dto.IDParam
	if err := c.ShouldBindUri(&p); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, "处方ID不合法", err))
		return
	}
	var req dto.OverrideReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, "强制通过参数不合法", err))
		return
	}
	op := middleware.CurrentUser(c)
	prescription, err := h.svc.Override(p.ID, req.Reason, op, c.ClientIP(), rid(c))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, prescription)
}

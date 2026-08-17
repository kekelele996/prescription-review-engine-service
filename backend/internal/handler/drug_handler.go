package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rxcheck/rxcheck/internal/dto"
	"github.com/rxcheck/rxcheck/internal/middleware"
	"github.com/rxcheck/rxcheck/internal/service"
	"github.com/rxcheck/rxcheck/internal/util"
)

// DrugHandler 药品档案处理器。
type DrugHandler struct {
	svc *service.DrugService
}

func NewDrugHandler(svc *service.DrugService) *DrugHandler {
	return &DrugHandler{svc: svc}
}

// List 药品列表。
func (h *DrugHandler) List(c *gin.Context) {
	page := util.ParsePage(c.Query("page"))
	pageSize := util.ParsePageSize(c.Query("page_size"))
	result, err := h.svc.List(page, pageSize, c.Query("keyword"), c.Query("status"), c.Query("mechanism"))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, result)
}

// Get 药品详情。
func (h *DrugHandler) Get(c *gin.Context) {
	var p dto.IDParam
	if err := c.ShouldBindUri(&p); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, "药品ID不合法", err))
		return
	}
	drug, err := h.svc.Get(p.ID)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, drug)
}

// Create 创建药品（管理员）。
func (h *DrugHandler) Create(c *gin.Context) {
	var req dto.CreateDrugReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, "药品参数不合法", err))
		return
	}
	op := middleware.CurrentUser(c)
	drug, err := h.svc.Create(&req, op.Username, c.ClientIP(), rid(c))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, drug)
}

// Update 更新药品（管理员）。
func (h *DrugHandler) Update(c *gin.Context) {
	var p dto.IDParam
	if err := c.ShouldBindUri(&p); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, "药品ID不合法", err))
		return
	}
	var req dto.UpdateDrugReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, "药品参数不合法", err))
		return
	}
	op := middleware.CurrentUser(c)
	drug, err := h.svc.Update(p.ID, &req, op.Username, c.ClientIP(), rid(c))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, drug)
}

// Disable 停用药品（管理员）。
func (h *DrugHandler) Disable(c *gin.Context) {
	var p dto.IDParam
	if err := c.ShouldBindUri(&p); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, "药品ID不合法", err))
		return
	}
	op := middleware.CurrentUser(c)
	drug, err := h.svc.Disable(p.ID, op.Username, c.ClientIP(), rid(c))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, drug)
}

// Enable 启用药品（管理员）。
func (h *DrugHandler) Enable(c *gin.Context) {
	var p dto.IDParam
	if err := c.ShouldBindUri(&p); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, "药品ID不合法", err))
		return
	}
	op := middleware.CurrentUser(c)
	drug, err := h.svc.Enable(p.ID, op.Username, c.ClientIP(), rid(c))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, drug)
}

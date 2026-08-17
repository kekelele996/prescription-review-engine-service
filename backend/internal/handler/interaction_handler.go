package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rxcheck/rxcheck/internal/dto"
	"github.com/rxcheck/rxcheck/internal/middleware"
	"github.com/rxcheck/rxcheck/internal/service"
	"github.com/rxcheck/rxcheck/internal/util"
)

// InteractionHandler 相互作用规则处理器。
type InteractionHandler struct {
	svc *service.InteractionService
}

func NewInteractionHandler(svc *service.InteractionService) *InteractionHandler {
	return &InteractionHandler{svc: svc}
}

// List 规则列表。
func (h *InteractionHandler) List(c *gin.Context) {
	page := util.ParsePage(c.Query("page"))
	pageSize := util.ParsePageSize(c.Query("page_size"))
	result, err := h.svc.List(page, pageSize, c.Query("risk"))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, result)
}

// Create 创建规则（管理员）。
func (h *InteractionHandler) Create(c *gin.Context) {
	var req dto.CreateInteractionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, "规则参数不合法", err))
		return
	}
	op := middleware.CurrentUser(c)
	rule, err := h.svc.Create(&req, op.Username, c.ClientIP(), rid(c))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, rule)
}

// Delete 删除规则（管理员）。
func (h *InteractionHandler) Delete(c *gin.Context) {
	var p dto.IDParam
	if err := c.ShouldBindUri(&p); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, "规则ID不合法", err))
		return
	}
	op := middleware.CurrentUser(c)
	if err := h.svc.Delete(p.ID, op.Username, c.ClientIP(), rid(c)); err != nil {
		c.Error(err)
		return
	}
	util.OK(c, gin.H{"deleted": p.ID})
}

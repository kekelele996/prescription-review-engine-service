package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rxcheck/rxcheck/internal/service"
	"github.com/rxcheck/rxcheck/internal/util"
)

// StatsHandler 统计处理器。
type StatsHandler struct {
	svc *service.StatsService
}

func NewStatsHandler(svc *service.StatsService) *StatsHandler {
	return &StatsHandler{svc: svc}
}

// Overview 审核统计概览。
func (h *StatsHandler) Overview(c *gin.Context) {
	result, err := h.svc.Overview()
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, result)
}

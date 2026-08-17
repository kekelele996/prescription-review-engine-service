package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rxcheck/rxcheck/internal/dto"
	"github.com/rxcheck/rxcheck/internal/middleware"
	"github.com/rxcheck/rxcheck/internal/service"
	"github.com/rxcheck/rxcheck/internal/util"
)

// UserHandler 用户管理处理器（管理员）。
type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// List 用户列表。
func (h *UserHandler) List(c *gin.Context) {
	page := util.ParsePage(c.Query("page"))
	pageSize := util.ParsePageSize(c.Query("page_size"))
	keyword := c.Query("keyword")
	result, err := h.svc.List(page, pageSize, keyword)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, result)
}

// Create 创建用户。
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, "用户参数不合法", err))
		return
	}
	op := middleware.CurrentUser(c)
	user, err := h.svc.Create(&req, op.Username, c.ClientIP(), rid(c))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, dto.ToUserResp(user))
}

// Update 更新用户。
func (h *UserHandler) Update(c *gin.Context) {
	var p dto.IDParam
	if err := c.ShouldBindUri(&p); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, "用户ID不合法", err))
		return
	}
	var req dto.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, "用户参数不合法", err))
		return
	}
	op := middleware.CurrentUser(c)
	user, err := h.svc.Update(p.ID, &req, op.Username, c.ClientIP(), rid(c))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, dto.ToUserResp(user))
}

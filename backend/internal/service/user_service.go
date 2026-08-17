package service

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/dto"
	"github.com/rxcheck/rxcheck/internal/model"
	"github.com/rxcheck/rxcheck/internal/repository"
	"github.com/rxcheck/rxcheck/internal/util"
)

// UserService 用户管理服务（管理员维护医生/药师账号）。
type UserService struct {
	repo  *repository.UserRepository
	audit *AuditService
	log   *slog.Logger
}

func NewUserService(repo *repository.UserRepository, audit *AuditService, log *slog.Logger) *UserService {
	return &UserService{repo: repo, audit: audit, log: log}
}

// Create 创建用户（角色由管理员指定）。
func (s *UserService) Create(req *dto.CreateUserReq, operator, ip, requestID string) (*model.User, error) {
	if !dto.ValidateRole(req.Role) {
		return nil, util.NewAppError(http.StatusBadRequest, "角色不合法: role="+req.Role, nil)
	}
	if _, err := s.repo.FindByUsername(req.Username); err == nil {
		return nil, util.NewAppError(http.StatusConflict, fmt.Sprintf(constants.MsgUserExists, req.Username), nil)
	}
	hash, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	user := &model.User{
		Username: req.Username,
		Password: hash,
		Name:     req.Name,
		Role:     req.Role,
		Status:   constants.UserStatusActive,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, "创建用户失败: username="+req.Username, err)
	}
	s.log.Info(fmt.Sprintf(constants.LogUserCreated, user.ID, user.Username, user.Role, operator))
	s.audit.Record(user.ID, operator, "CREATE", "user", util.Uint64String(user.ID), "创建账号: "+user.Username, ip, requestID)
	return user, nil
}

// List 分页查询用户。
func (s *UserService) List(page, pageSize int, keyword string) (*util.PageResult, error) {
	list, total, err := s.repo.List(page, pageSize, keyword)
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	out := make([]dto.UserResp, 0, len(list))
	for i := range list {
		out = append(out, dto.ToUserResp(&list[i]))
	}
	return &util.PageResult{List: out, Total: total, Page: page, PageSize: pageSize}, nil
}

// Get 查询用户详情。
func (s *UserService) Get(id uint) (*dto.UserResp, error) {
	u, err := s.repo.FindByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, util.NewAppError(http.StatusNotFound, fmt.Sprintf(constants.MsgUserNotFound, id), nil)
	}
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	resp := dto.ToUserResp(u)
	return &resp, nil
}

// Update 更新用户（姓名/状态/角色）。
func (s *UserService) Update(id uint, req *dto.UpdateUserReq, operator, ip, requestID string) (*model.User, error) {
	u, err := s.repo.FindByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, util.NewAppError(http.StatusNotFound, fmt.Sprintf(constants.MsgUserNotFound, id), nil)
	}
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	if req.Name != "" {
		u.Name = req.Name
	}
	if req.Status != "" {
		u.Status = req.Status
	}
	if req.Role != "" {
		if !dto.ValidateRole(req.Role) {
			return nil, util.NewAppError(http.StatusBadRequest, "角色不合法: role="+req.Role, nil)
		}
		u.Role = req.Role
	}
	if err := s.repo.Update(u); err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, "更新用户失败: user_id="+util.Uint64String(id), err)
	}
	s.log.Info(fmt.Sprintf(constants.LogUserUpdated, u.ID, u.Status, operator))
	s.audit.Record(u.ID, operator, "UPDATE", "user", util.Uint64String(u.ID), "更新账号: "+u.Username, ip, requestID)
	return u, nil
}

package service

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/dto"
	"github.com/rxcheck/rxcheck/internal/model"
	"github.com/rxcheck/rxcheck/internal/repository"
	"github.com/rxcheck/rxcheck/internal/util"
)

// AuthService 认证服务：登录签发 JWT、自助注册（默认医生角色）。
type AuthService struct {
	repo       *repository.UserRepository
	audit      *AuditService
	jwtSecret  string
	tokenHours int
	log        *slog.Logger
}

func NewAuthService(repo *repository.UserRepository, audit *AuditService, jwtSecret string, tokenHours int, log *slog.Logger) *AuthService {
	return &AuthService{repo: repo, audit: audit, jwtSecret: jwtSecret, tokenHours: tokenHours, log: log}
}

// Login 登录并签发 JWT。
func (s *AuthService) Login(req *dto.LoginReq, ip, requestID string) (*dto.LoginResp, error) {
	user, err := s.repo.FindByUsername(req.Username)
	if errors.Is(err, repository.ErrNotFound) {
		s.log.Warn(fmt.Sprintf(constants.LogLoginFailed, req.Username, "user not found"))
		return nil, util.NewAppError(http.StatusUnauthorized, fmt.Sprintf(constants.MsgWrongPassword, req.Username), nil)
	}
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	if user.Status != constants.UserStatusActive {
		s.log.Warn(fmt.Sprintf(constants.LogLoginFailed, req.Username, "disabled"))
		return nil, util.NewAppError(http.StatusForbidden, fmt.Sprintf(constants.MsgUserDisabled, req.Username), nil)
	}
	if !util.CheckPassword(user.Password, req.Password) {
		s.log.Warn(fmt.Sprintf(constants.LogLoginFailed, req.Username, "wrong password"))
		return nil, util.NewAppError(http.StatusUnauthorized, fmt.Sprintf(constants.MsgWrongPassword, req.Username), nil)
	}
	token, err := util.GenerateToken(s.jwtSecret, time.Duration(s.tokenHours)*time.Hour, user.ID, user.Username, user.Role)
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	now := time.Now()
	if err := s.repo.UpdateLastLogin(user.ID, now); err != nil {
		s.log.Warn("更新最后登录时间失败", "user_id", user.ID, "err", err)
	}
	s.log.Info(fmt.Sprintf(constants.LogLoginSuccess, user.Username, user.Role, requestID))
	s.audit.Record(user.ID, user.Username, "LOGIN", "auth", util.Uint64String(user.ID), "用户登录", ip, requestID)
	return &dto.LoginResp{Token: token, User: dto.ToUserResp(user)}, nil
}

// Register 自助注册医生账号。
func (s *AuthService) Register(req *dto.RegisterReq, ip, requestID string) (*model.User, error) {
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
		Role:     constants.UserRoleDoctor,
		Status:   constants.UserStatusActive,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, "创建用户失败: username="+req.Username, err)
	}
	s.log.Info(fmt.Sprintf(constants.LogRegisterSuccess, user.ID, user.Username, user.Role))
	s.audit.Record(user.ID, user.Username, "REGISTER", "auth", util.Uint64String(user.ID), "自助注册医生账号", ip, requestID)
	return user, nil
}

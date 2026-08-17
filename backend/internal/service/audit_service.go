package service

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/dto"
	"github.com/rxcheck/rxcheck/internal/model"
	"github.com/rxcheck/rxcheck/internal/repository"
	"github.com/rxcheck/rxcheck/internal/util"
)

// AuditService 审计日志服务：被审计中间件与各业务 service 复用。
type AuditService struct {
	repo *repository.AuditRepository
	log  *slog.Logger
}

func NewAuditService(repo *repository.AuditRepository, log *slog.Logger) *AuditService {
	return &AuditService{repo: repo, log: log}
}

// Record 记录一条操作审计日志（写失败仅告警，不阻断主流程）。
func (s *AuditService) Record(userID uint, username, action, module, entityID, detail, ip, requestID string) {
	log := &model.AuditLog{
		UserID:    userID,
		Username:  username,
		Action:    action,
		Module:    module,
		EntityID:  entityID,
		Detail:    detail,
		IP:        ip,
		RequestID: requestID,
	}
	if err := s.repo.Create(log); err != nil {
		s.log.Warn(fmt.Sprintf(constants.LogAuditWriteFailed, requestID, err))
	}
}

// List 分页查询审计日志。
func (s *AuditService) List(page, pageSize int, username string) (*util.PageResult, error) {
	list, total, err := s.repo.List(page, pageSize, username)
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	out := make([]dto.AuditLogResp, 0, len(list))
	for i := range list {
		out = append(out, dto.ToAuditLogResp(&list[i]))
	}
	return &util.PageResult{List: out, Total: total, Page: page, PageSize: pageSize}, nil
}

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

// InteractionService 药物相互作用规则服务。
type InteractionService struct {
	repo  *repository.InteractionRepository
	drugs *repository.DrugRepository
	audit *AuditService
	log   *slog.Logger
}

func NewInteractionService(repo *repository.InteractionRepository, drugs *repository.DrugRepository, audit *AuditService, log *slog.Logger) *InteractionService {
	return &InteractionService{repo: repo, drugs: drugs, audit: audit, log: log}
}

// Create 创建相互作用规则。
func (s *InteractionService) Create(req *dto.CreateInteractionReq, operator, ip, requestID string) (*model.InteractionRule, error) {
	if req.DrugAID == req.DrugBID {
		return nil, util.NewAppError(http.StatusBadRequest, "相互作用规则的两端药品不能相同: drug_a="+util.Uint64String(req.DrugAID), nil)
	}
	if _, err := s.drugs.FindByID(req.DrugAID); errors.Is(err, repository.ErrNotFound) {
		return nil, util.NewAppError(http.StatusNotFound, fmt.Sprintf(constants.MsgDrugNotFound, req.DrugAID), nil)
	}
	if _, err := s.drugs.FindByID(req.DrugBID); errors.Is(err, repository.ErrNotFound) {
		return nil, util.NewAppError(http.StatusNotFound, fmt.Sprintf(constants.MsgDrugNotFound, req.DrugBID), nil)
	}
	if _, err := s.repo.FindByPair(req.DrugAID, req.DrugBID); err == nil {
		return nil, util.NewAppError(http.StatusConflict, fmt.Sprintf(constants.MsgInteractionExists, req.DrugAID, req.DrugBID), nil)
	}
	status := req.Status
	if status == "" {
		status = constants.RuleStatusEnabled
	}
	rule := &model.InteractionRule{
		DrugAID:     req.DrugAID,
		DrugBID:     req.DrugBID,
		RiskLevel:   req.RiskLevel,
		Mechanism:   req.Mechanism,
		Description: req.Description,
		Status:      status,
	}
	if err := s.repo.Create(rule); err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, "创建相互作用规则失败: drug_a="+util.Uint64String(req.DrugAID)+" drug_b="+util.Uint64String(req.DrugBID), err)
	}
	s.log.Info(fmt.Sprintf(constants.LogInteractionCreated, rule.ID, rule.DrugAID, rule.DrugBID, rule.RiskLevel, operator))
	s.audit.Record(0, operator, "CREATE", "interaction", util.Uint64String(rule.ID), "新增相互作用规则: "+rule.Mechanism, ip, requestID)
	return rule, nil
}

// List 分页查询规则。
func (s *InteractionService) List(page, pageSize int, risk string) (*util.PageResult, error) {
	list, total, err := s.repo.List(page, pageSize, risk)
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	return &util.PageResult{List: list, Total: total, Page: page, PageSize: pageSize}, nil
}

// Delete 删除规则。
func (s *InteractionService) Delete(id uint, operator, ip, requestID string) error {
	if err := s.repo.Delete(id); errors.Is(err, repository.ErrNotFound) {
		return util.NewAppError(http.StatusNotFound, fmt.Sprintf(constants.MsgInteractionNotFound, id), nil)
	} else if err != nil {
		return util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	s.log.Info(fmt.Sprintf(constants.LogInteractionDeleted, id, operator))
	s.audit.Record(0, operator, "DELETE", "interaction", util.Uint64String(id), "删除相互作用规则: rule_id="+util.Uint64String(id), ip, requestID)
	return nil
}

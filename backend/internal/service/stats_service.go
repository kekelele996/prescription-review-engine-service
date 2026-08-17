package service

import (
	"log/slog"
	"net/http"

	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/dto"
	"github.com/rxcheck/rxcheck/internal/repository"
	"github.com/rxcheck/rxcheck/internal/util"
)

// StatsService 审核统计服务。
type StatsService struct {
	presc  *repository.PrescriptionRepository
	drugs  *repository.DrugRepository
	inter  *repository.InteractionRepository
	report *repository.ReportRepository
	log    *slog.Logger
}

func NewStatsService(presc *repository.PrescriptionRepository, drugs *repository.DrugRepository,
	inter *repository.InteractionRepository, report *repository.ReportRepository, log *slog.Logger) *StatsService {
	return &StatsService{presc: presc, drugs: drugs, inter: inter, report: report, log: log}
}

// Overview 审核统计概览。
func (s *StatsService) Overview() (*dto.StatsOverviewResp, error) {
	total, err := s.presc.Count()
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	byStatus, err := s.presc.CountByStatus()
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	byRisk, err := s.presc.CountByRisk()
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	totalDrugs, err := s.drugs.Count()
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	totalInter, err := s.inter.Count()
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	return &dto.StatsOverviewResp{
		TotalPrescriptions: total,
		ByStatus:           byStatus,
		ByRisk:             byRisk,
		TotalDrugs:         totalDrugs,
		TotalInteractions:  totalInter,
	}, nil
}

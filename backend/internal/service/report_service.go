package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/dto"
	"github.com/rxcheck/rxcheck/internal/repository"
	"github.com/rxcheck/rxcheck/internal/storage"
	"github.com/rxcheck/rxcheck/internal/util"
)

// ReportService 审核报告服务：报告查询与快照导出。
type ReportService struct {
	repo    *repository.ReportRepository
	presc   *repository.PrescriptionRepository
	storage *storage.MinioClient
	audit   *AuditService
	log     *slog.Logger
}

func NewReportService(repo *repository.ReportRepository, presc *repository.PrescriptionRepository,
	st *storage.MinioClient, audit *AuditService, log *slog.Logger) *ReportService {
	return &ReportService{repo: repo, presc: presc, storage: st, audit: audit, log: log}
}

// GetByPrescription 按处方查询报告。
func (s *ReportService) GetByPrescription(prescriptionID uint) (*dto.ReportResp, error) {
	report, err := s.repo.FindByPrescriptionID(prescriptionID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, util.NewAppError(http.StatusNotFound, fmt.Sprintf(constants.MsgReportNotFound, prescriptionID), nil)
	}
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	resp := dto.ToReportResp(report)
	return &resp, nil
}

// List 分页查询报告。
func (s *ReportService) List(page, pageSize int, status string) (*util.PageResult, error) {
	list, total, err := s.repo.List(page, pageSize, status)
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	out := make([]dto.ReportResp, 0, len(list))
	for i := range list {
		out = append(out, dto.ToReportResp(&list[i]))
	}
	return &util.PageResult{List: out, Total: total, Page: page, PageSize: pageSize}, nil
}

// Export 导出报告 JSON 快照到 MinIO 并返回下载地址。
func (s *ReportService) Export(prescriptionID uint, operator, ip, requestID string) (map[string]string, error) {
	p, err := s.presc.FindByID(prescriptionID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, util.NewAppError(http.StatusNotFound, fmt.Sprintf(constants.MsgPrescriptionNotFound, prescriptionID), nil)
	}
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	report, err := s.repo.FindByPrescriptionID(prescriptionID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, util.NewAppError(http.StatusNotFound, fmt.Sprintf(constants.MsgReportNotFound, prescriptionID), nil)
	}
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	snapshot := map[string]any{
		"prescription_id":   p.ID,
		"prescription_no":   p.PrescriptionNo,
		"patient_name":      p.PatientName,
		"patient_age":       p.PatientAge,
		"status":            p.Status,
		"overall_risk":      p.OverallRisk,
		"report_status":     report.Status,
		"report_risk_level": report.RiskLevel,
		"summary":           report.Summary,
		"items":             report.Items,
		"exported_at":       time.Now().Format(time.RFC3339),
	}
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	objectName := fmt.Sprintf("reports/%d-%s.json", p.ID, time.Now().Format("20060102150405"))
	url, err := s.storage.UploadJSON(context.Background(), objectName, data)
	if err != nil {
		s.log.Warn(fmt.Sprintf(constants.LogMinioUploadFailed, p.ID, err))
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgStorageError, err)
	}
	s.log.Info(fmt.Sprintf(constants.LogReportExported, p.ID, objectName))
	s.audit.Record(0, operator, "EXPORT", "report", util.Uint64String(p.ID), "导出审核报告快照: "+p.PrescriptionNo, ip, requestID)
	return map[string]string{"object": objectName, "url": url}, nil
}


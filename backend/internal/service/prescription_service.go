package service

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/dto"
	"github.com/rxcheck/rxcheck/internal/model"
	"github.com/rxcheck/rxcheck/internal/queue"
	"github.com/rxcheck/rxcheck/internal/repository"
	"github.com/rxcheck/rxcheck/internal/util"
	"github.com/redis/go-redis/v9"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// PrescriptionService 处方接收与解析服务：JSON/XML 双格式、明细映射、状态机驱动。
type PrescriptionService struct {
	db     *gorm.DB
	repo   *repository.PrescriptionRepository
	drugs  *repository.DrugRepository
	review *ReviewService
	audit  *AuditService
	rdb    *redis.Client
	log    *slog.Logger
}

func NewPrescriptionService(db *gorm.DB, repo *repository.PrescriptionRepository, drugs *repository.DrugRepository,
	review *ReviewService, audit *AuditService, rdb *redis.Client, log *slog.Logger) *PrescriptionService {
	return &PrescriptionService{db: db, repo: repo, drugs: drugs, review: review, audit: audit, rdb: rdb, log: log}
}

// ParsePrescription 按 Content-Type 解析 JSON 或 XML 处方，返回请求体、格式与原始字节。
func ParsePrescription(c *gin.Context) (*dto.SubmitPrescriptionReq, string, []byte, error) {
	contentType := util.NormalizeContentType(c.GetHeader("Content-Type"))
	body, err := c.GetRawData()
	if err != nil {
		return nil, "", nil, util.NewAppError(http.StatusBadRequest, fmt.Sprintf(constants.MsgParseFailed, err.Error()), nil)
	}
	req := &dto.SubmitPrescriptionReq{}
	switch contentType {
	case "xml":
		if err := xml.Unmarshal(body, req); err != nil {
			return nil, "", nil, util.NewAppError(http.StatusBadRequest, fmt.Sprintf(constants.MsgParseFailed, err.Error()), nil)
		}
	case "json":
		if err := json.Unmarshal(body, req); err != nil {
			return nil, "", nil, util.NewAppError(http.StatusBadRequest, fmt.Sprintf(constants.MsgParseFailed, err.Error()), nil)
		}
	default:
		return nil, "", nil, util.NewAppError(http.StatusBadRequest, fmt.Sprintf(constants.MsgInvalidFormat, contentType), nil)
	}
	if err := util.ValidateStruct(req); err != nil {
		return nil, "", nil, util.NewAppError(http.StatusBadRequest, fmt.Sprintf(constants.MsgParseFailed, err.Error()), nil)
	}
	return req, contentType, body, nil
}

// Submit 接收处方：落库后入队异步审核；Redis 不可用时降级为同步审核。
func (s *PrescriptionService) Submit(req *dto.SubmitPrescriptionReq, format string, raw []byte, operator *util.Claims, ip, requestID string) (*model.Prescription, error) {
	items, err := s.resolveItems(req.Items)
	if err != nil {
		return nil, err
	}
	p := &model.Prescription{
		PrescriptionNo:    req.PrescriptionNo,
		PatientName:       req.PatientName,
		PatientAge:        req.PatientAge,
		PatientGender:     req.PatientGender,
		WeightKg:          req.WeightKg,
		Allergies:         toJSON(req.Allergies),
		Pregnant:          req.Pregnant,
		HepaticImpairment: req.HepaticImpairment,
		RenalImpairment:   req.RenalImpairment,
		Diagnoses:         toDiagnosesJSON(req.Diagnoses),
		Format:            format,
		RawPayload:        string(raw),
		Status:            constants.PrescriptionStatusPendingReview,
		OverallRisk:       constants.RiskNone,
	}
	if operator != nil && operator.UserID > 0 {
		uid := operator.UserID
		p.DoctorID = &uid
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		return s.repo.CreateWithItemsTx(tx, p, items)
	})
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, "处方落库失败: prescription_no="+req.PrescriptionNo, err)
	}

	opName := "anonymous"
	if operator != nil {
		opName = operator.Username
	}
	s.log.Info(fmt.Sprintf(constants.LogPrescriptionSubmitted, p.ID, p.PrescriptionNo, format, p.Status, opName))
	s.audit.Record(0, opName, "SUBMIT", "prescription", util.Uint64String(p.ID), "接收处方: "+p.PrescriptionNo, ip, requestID)

	if err := queue.EnqueueReview(context.Background(), s.rdb, p.ID); err != nil {
		s.log.Warn(fmt.Sprintf(constants.LogQueueEnqueueFailed, p.ID, err))
		if _, rErr := s.review.ReviewPrescription(p.ID, opName, ip, requestID); rErr != nil {
			return nil, rErr
		}
	}
	return s.repo.FindByID(p.ID)
}

// resolveItems 将处方明细映射到药品库，计算日剂量。
func (s *PrescriptionService) resolveItems(reqItems []dto.PrescriptionItemDTO) ([]model.PrescriptionItem, error) {
	items := make([]model.PrescriptionItem, 0, len(reqItems))
	for i, it := range reqItems {
		var drug *model.Drug
		var err error
		if it.DrugID > 0 {
			drug, err = s.drugs.FindByID(it.DrugID)
		} else if it.DrugName != "" {
			drug, err = s.drugs.FindByName(it.DrugName)
		} else {
			return nil, util.NewAppError(http.StatusBadRequest, fmt.Sprintf("处方明细第 %d 条缺少药品标识: drug_name 或 drug_id 必填", i+1), nil)
		}
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(http.StatusBadRequest, "药品不存在: name="+it.DrugName, nil)
		}
		if err != nil {
			return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
		}
		if drug.Status != constants.DrugStatusEnabled {
			return nil, util.NewAppError(http.StatusConflict, "药品已停用: name="+drug.Name, nil)
		}
		items = append(items, model.PrescriptionItem{
			DrugID:        drug.ID,
			DrugName:      drug.Name,
			Specification: it.Specification,
			SingleDose:    it.SingleDose,
			Frequency:     it.Frequency,
			CourseDays:    it.CourseDays,
			Route:         it.Route,
			DailyDose:     it.SingleDose * frequencyDailyTimes(it.Frequency),
		})
	}
	return items, nil
}

// List 分页查询处方。
func (s *PrescriptionService) List(page, pageSize int, status, keyword string) (*util.PageResult, error) {
	list, total, err := s.repo.List(page, pageSize, status, keyword)
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	out := make([]dto.PrescriptionListResp, 0, len(list))
	for i := range list {
		out = append(out, toListResp(&list[i]))
	}
	return &util.PageResult{List: out, Total: total, Page: page, PageSize: pageSize}, nil
}

// Get 查询处方详情。
func (s *PrescriptionService) Get(id uint) (*model.Prescription, error) {
	p, err := s.repo.FindByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, util.NewAppError(http.StatusNotFound, fmt.Sprintf(constants.MsgPrescriptionNotFound, id), nil)
	}
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	return p, nil
}

// Review 重新审核处方（同步执行，返回最新报告）。
func (s *PrescriptionService) Review(id uint, operator, ip, requestID string) (*model.ReviewReport, error) {
	if _, err := s.repo.FindByID(id); errors.Is(err, repository.ErrNotFound) {
		return nil, util.NewAppError(http.StatusNotFound, fmt.Sprintf(constants.MsgPrescriptionNotFound, id), nil)
	} else if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	return s.review.ReviewPrescription(id, operator, ip, requestID)
}

// Override 医生/药师对警告类处方强制通过并记录原因。
func (s *PrescriptionService) Override(id uint, reason string, operator *util.Claims, ip, requestID string) (*model.Prescription, error) {
	var updated *model.Prescription
	err := s.db.Transaction(func(tx *gorm.DB) error {
		p, err := s.repo.FindByIDForUpdate(tx, id)
		if errors.Is(err, repository.ErrNotFound) {
			return util.NewAppError(http.StatusNotFound, fmt.Sprintf(constants.MsgPrescriptionNotFound, id), nil)
		}
		if err != nil {
			return err
		}
		if p.Status != constants.PrescriptionStatusWarned {
			return util.NewAppError(http.StatusConflict, fmt.Sprintf(constants.MsgOverrideNotAllowed, p.Status), nil)
		}
		p.Status = constants.PrescriptionStatusOverridden
		p.OverrideReason = reason
		uid := operator.UserID
		p.OverrideBy = &uid
		if err := s.repo.UpdateTx(tx, p); err != nil {
			return err
		}
		updated = p
		return nil
	})
	if err != nil {
		var appErr *util.AppError
		if errors.As(err, &appErr) {
			return nil, appErr
		}
		return nil, util.NewAppError(http.StatusInternalServerError, "强制通过失败: prescription_id="+util.Uint64String(id), err)
	}
	s.log.Info(fmt.Sprintf(constants.LogPrescriptionOverridden, id, updated.PrescriptionNo, reason, operator.Username))
	s.audit.Record(operator.UserID, operator.Username, "OVERRIDE", "prescription", util.Uint64String(id), "强制通过处方: "+updated.PrescriptionNo+" 原因: "+reason, ip, requestID)
	return updated, nil
}

func toListResp(p *model.Prescription) dto.PrescriptionListResp {
	reviewCount := 0
	if p.Report == nil {
		reviewCount = len(p.Report.Items)
	}
	reviewed := ""
	if p.ReviewedAt != nil {
		reviewed = p.ReviewedAt.Format("2006-01-02 15:04:05")
	}
	return dto.PrescriptionListResp{
		ID:              p.ID,
		PrescriptionNo:  p.PrescriptionNo,
		PatientName:     p.PatientName,
		PatientAge:      p.PatientAge,
		Format:          p.Format,
		Status:          p.Status,
		OverallRisk:     p.OverallRisk,
		ItemCount:       len(p.Items),
		ReviewItemCount: reviewCount,
		CreatedAt:       p.CreatedAt.Format("2006-01-02 15:04:05"),
		ReviewedAt:      reviewed,
	}
}

// timeNow 返回当前时间。
func timeNow() time.Time { return time.Now() }

// toJSON 字符串切片转 JSON。
func toJSON(list []string) datatypes.JSON {
	b, err := json.Marshal(list)
	if err != nil {
		return datatypes.JSON("[]")
	}
	return datatypes.JSON(b)
}

// toDiagnosesJSON 诊断 DTO 转 JSON。
func toDiagnosesJSON(list []dto.DiagnosisDTO) datatypes.JSON {
	b, err := json.Marshal(list)
	if err != nil {
		return datatypes.JSON("[]")
	}
	return datatypes.JSON(b)
}

package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/dto"
	"github.com/rxcheck/rxcheck/internal/model"
	"github.com/rxcheck/rxcheck/internal/repository"
	"github.com/rxcheck/rxcheck/internal/util"
	"gorm.io/datatypes"
)

// DrugService 药品档案服务：维护适应症、剂量上限、禁忌等审核规则来源。
type DrugService struct {
	repo  *repository.DrugRepository
	audit *AuditService
	log   *slog.Logger
}

func NewDrugService(repo *repository.DrugRepository, audit *AuditService, log *slog.Logger) *DrugService {
	return &DrugService{repo: repo, audit: audit, log: log}
}

// Create 创建药品。
func (s *DrugService) Create(req *dto.CreateDrugReq, operator, ip, requestID string) (*model.Drug, error) {
	if _, err := s.repo.FindByName(req.Name); err == nil {
		return nil, util.NewAppError(http.StatusConflict, fmt.Sprintf(constants.MsgDrugExists, req.Name), nil)
	}
	status := req.Status
	if status == "" {
		status = constants.DrugStatusEnabled
	}
	drug := &model.Drug{
		Name:                     req.Name,
		GenericName:              req.GenericName,
		Specification:            req.Specification,
		DosageForm:               req.DosageForm,
		Route:                    req.Route,
		AtcCode:                  req.AtcCode,
		Mechanism:                req.Mechanism,
		Indications:              indicationsJSON(req.Indications),
		AdultMaxDailyDose:        req.AdultMaxDailyDose,
		ChildMaxDailyDose:        req.ChildMaxDailyDose,
		ElderlyMaxDailyDose:      req.ElderlyMaxDailyDose,
		MaxSingleDose:            req.MaxSingleDose,
		DoseUnit:                 req.DoseUnit,
		Allergies:                stringsJSON(req.Allergies),
		PregnancyContraindicated: req.PregnancyContraindicated,
		HepaticContraindicated:   req.HepaticContraindicated,
		RenalContraindicated:     req.RenalContraindicated,
		BeersFlag:                req.BeersFlag,
		MinAgeMonths:             req.MinAgeMonths,
		FrequencyLimit:           req.FrequencyLimit,
		Status:                   status,
	}
	if err := s.repo.Create(drug); err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, "创建药品失败: name="+req.Name, err)
	}
	s.log.Info(fmt.Sprintf(constants.LogDrugCreated, drug.ID, drug.Name, drug.Status, operator))
	s.audit.Record(0, operator, "CREATE", "drug", util.Uint64String(drug.ID), "新增药品: "+drug.Name, ip, requestID)
	return drug, nil
}

// List 分页查询药品。
func (s *DrugService) List(page, pageSize int, keyword, status, mechanism string) (*util.PageResult, error) {
	list, total, err := s.repo.List(page, pageSize, keyword, status, mechanism)
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	return &util.PageResult{List: list, Total: total, Page: page, PageSize: pageSize}, nil
}

// Get 查询药品详情。
func (s *DrugService) Get(id uint) (*model.Drug, error) {
	d, err := s.repo.FindByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, util.NewAppError(http.StatusNotFound, fmt.Sprintf(constants.MsgDrugNotFound, id), nil)
	}
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	return d, nil
}

// Update 更新药品档案。
func (s *DrugService) Update(id uint, req *dto.UpdateDrugReq, operator, ip, requestID string) (*model.Drug, error) {
	d, err := s.repo.FindByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, util.NewAppError(http.StatusNotFound, fmt.Sprintf(constants.MsgDrugNotFound, id), nil)
	}
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	applyUpdate(d, req)
	if err := s.repo.Update(d); err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, "更新药品失败: drug_id="+util.Uint64String(id), err)
	}
	s.log.Info(fmt.Sprintf(constants.LogDrugUpdated, d.ID, d.Name, operator))
	s.audit.Record(0, operator, "UPDATE", "drug", util.Uint64String(d.ID), "更新药品: "+d.Name, ip, requestID)
	return d, nil
}

// Disable 停用药品。
func (s *DrugService) Disable(id uint, operator, ip, requestID string) (*model.Drug, error) {
	if err := s.repo.UpdateStatusTx(s.repo.DB(), id, constants.DrugStatusDisabled); err != nil {
		return nil, s.statusErr(err, id)
	}
	d, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	s.log.Info(fmt.Sprintf(constants.LogDrugDisabled, d.ID, d.Name, operator))
	s.audit.Record(0, operator, "DISABLE", "drug", util.Uint64String(id), "停用药品: "+d.Name, ip, requestID)
	return d, nil
}

// Enable 启用药品。
func (s *DrugService) Enable(id uint, operator, ip, requestID string) (*model.Drug, error) {
	if err := s.repo.UpdateStatusTx(s.repo.DB(), id, constants.DrugStatusEnabled); err != nil {
		return nil, s.statusErr(err, id)
	}
	d, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	s.log.Info(fmt.Sprintf(constants.LogDrugEnabled, d.ID, d.Name, operator))
	s.audit.Record(0, operator, "ENABLE", "drug", util.Uint64String(id), "启用药品: "+d.Name, ip, requestID)
	return d, nil
}

func (s *DrugService) statusErr(err error, id uint) error {
	if errors.Is(err, repository.ErrNotFound) {
		return util.NewAppError(http.StatusNotFound, fmt.Sprintf(constants.MsgDrugNotFound, id), nil)
	}
	return util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
}

// applyUpdate 将可选更新字段应用到药品。
func applyUpdate(d *model.Drug, req *dto.UpdateDrugReq) {
	if req.Name != "" {
		d.Name = req.Name
	}
	if req.GenericName != "" {
		d.GenericName = req.GenericName
	}
	if req.Specification != "" {
		d.Specification = req.Specification
	}
	if req.DosageForm != "" {
		d.DosageForm = req.DosageForm
	}
	if req.Route != "" {
		d.Route = req.Route
	}
	if req.AtcCode != "" {
		d.AtcCode = req.AtcCode
	}
	if req.Mechanism != "" {
		d.Mechanism = req.Mechanism
	}
	if req.Indications != nil {
		d.Indications = indicationsJSON(req.Indications)
	}
	if req.AdultMaxDailyDose != nil {
		d.AdultMaxDailyDose = *req.AdultMaxDailyDose
	}
	if req.ChildMaxDailyDose != nil {
		d.ChildMaxDailyDose = *req.ChildMaxDailyDose
	}
	if req.ElderlyMaxDailyDose != nil {
		d.ElderlyMaxDailyDose = *req.ElderlyMaxDailyDose
	}
	if req.MaxSingleDose != nil {
		d.MaxSingleDose = *req.MaxSingleDose
	}
	if req.DoseUnit != "" {
		d.DoseUnit = req.DoseUnit
	}
	if req.Allergies != nil {
		d.Allergies = stringsJSON(req.Allergies)
	}
	if req.PregnancyContraindicated != nil {
		d.PregnancyContraindicated = *req.PregnancyContraindicated
	}
	if req.HepaticContraindicated != nil {
		d.HepaticContraindicated = *req.HepaticContraindicated
	}
	if req.RenalContraindicated != nil {
		d.RenalContraindicated = *req.RenalContraindicated
	}
	if req.BeersFlag != nil {
		d.BeersFlag = *req.BeersFlag
	}
	if req.MinAgeMonths != nil {
		d.MinAgeMonths = *req.MinAgeMonths
	}
	if req.FrequencyLimit != "" {
		d.FrequencyLimit = req.FrequencyLimit
	}
	if req.Status != "" {
		d.Status = req.Status
	}
}

// indicationsJSON 适应症 DTO 转 JSON。
func indicationsJSON(list []dto.IndicationDTO) datatypes.JSON {
	b, err := json.Marshal(list)
	if err != nil {
		return datatypes.JSON("[]")
	}
	return datatypes.JSON(b)
}

// stringsJSON 字符串列表转 JSON。
func stringsJSON(list []string) datatypes.JSON {
	b, err := json.Marshal(list)
	if err != nil {
		return datatypes.JSON("[]")
	}
	return datatypes.JSON(b)
}

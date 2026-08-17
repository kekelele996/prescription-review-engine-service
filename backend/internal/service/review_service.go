package service

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/model"
	"github.com/rxcheck/rxcheck/internal/repository"
	"github.com/rxcheck/rxcheck/internal/util"
	"gorm.io/gorm"
)

// ReviewService 处方审核引擎：适应症、用法用量、相互作用、禁忌与特殊人群、重复用药五大规则。
type ReviewService struct {
	db      *gorm.DB
	presc   *repository.PrescriptionRepository
	drugs   *repository.DrugRepository
	inter   *repository.InteractionRepository
	reports *repository.ReportRepository
	audit   *AuditService
	log     *slog.Logger
}

func NewReviewService(db *gorm.DB, presc *repository.PrescriptionRepository, drugs *repository.DrugRepository,
	inter *repository.InteractionRepository, reports *repository.ReportRepository, audit *AuditService, log *slog.Logger) *ReviewService {
	return &ReviewService{db: db, presc: presc, drugs: drugs, inter: inter, reports: reports, audit: audit, log: log}
}

// frequencyDailyTimes 频次 -> 每日次数。
func frequencyDailyTimes(freq string) float64 {
	switch freq {
	case "qd", "qn":
		return 1
	case "bid", "q12h":
		return 2
	case "tid", "q8h":
		return 3
	case "qid", "q6h":
		return 4
	default:
		return 1
	}
}

// ageGroup 年龄段：child 儿童 / elderly 老人 / adult 成人。
func ageGroup(age int) string {
	switch {
	case age < 18:
		return "child"
	case age >= 65:
		return "elderly"
	default:
		return "adult"
	}
}

// maxDailyDose 按年龄段取每日最大剂量。
func maxDailyDose(d *model.Drug, age int) float64 {
	switch ageGroup(age) {
	case "child":
		return d.ChildMaxDailyDose
	case "elderly":
		return d.ElderlyMaxDailyDose
	default:
		return d.AdultMaxDailyDose
	}
}

// ReviewPrescription 执行审核并落库（幂等，重新审核时替换旧报告）。
func (s *ReviewService) ReviewPrescription(prescriptionID uint, operator, ip, requestID string) (*model.ReviewReport, error) {
	p, err := s.presc.FindByID(prescriptionID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, util.NewAppError(http.StatusNotFound, fmt.Sprintf(constants.MsgPrescriptionNotFound, prescriptionID), nil)
	}
	if err != nil {
		return nil, util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
	}
	s.log.Info(fmt.Sprintf(constants.LogPrescriptionReviewStart, p.ID, p.PrescriptionNo))

	items := s.buildReviewItems(p)
	overallRisk := constants.RiskNone
	for _, it := range items {
		overallRisk = it.RiskLevel
	}
	status := constants.MapRiskToStatus(overallRisk)
	summary := fmt.Sprintf("处方共 %d 种药品，发现 %d 条审核意见，最高风险等级：%s",
		len(p.Items), len(items), util.RiskLevelText(overallRisk))

	report := &model.ReviewReport{
		PrescriptionID: p.ID,
		Status:         status,
		RiskLevel:      overallRisk,
		Summary:        summary,
		Items:          items,
	}

	now := time.Now()
	err = s.db.Transaction(func(tx *gorm.DB) error {
		locked, err := s.presc.FindByIDForUpdate(tx, p.ID)
		if err != nil {
			return err
		}
		if locked.Status == constants.PrescriptionStatusOverridden {
			return util.NewAppError(http.StatusConflict, fmt.Sprintf(constants.MsgPrescriptionLocked, locked.Status), nil)
		}
		if err := s.reports.ReplaceWithItemsTx(tx, report); err != nil {
			return err
		}
		return s.presc.UpdateStatusTx(tx, p.ID, status, overallRisk, &now)
	})
	if err != nil {
		var appErr *util.AppError
		if errors.As(err, &appErr) {
			return nil, appErr
		}
		return nil, util.NewAppError(http.StatusInternalServerError, "处方审核落库失败: prescription_id="+util.Uint64String(p.ID), err)
	}

	s.log.Info(fmt.Sprintf(constants.LogPrescriptionReviewed, p.ID, status, overallRisk, len(items)))
	s.audit.Record(0, operator, "REVIEW", "prescription", util.Uint64String(p.ID), "自动审核处方: "+p.PrescriptionNo, ip, requestID)
	return report, nil
}

// buildReviewItems 依次执行五大审核规则。
func (s *ReviewService) buildReviewItems(p *model.Prescription) []model.ReviewItem {
	var items []model.ReviewItem
	diags := p.DiagnosisList()
	allergies := p.AllergyList()

	for i := range p.Items {
		item := &p.Items[i]
		drug := &item.Drug
		if drug == nil || drug.ID == 0 {
			continue
		}
		items = append(items, s.checkIndication(item, drug, diags)...)
		items = append(items, s.checkDosage(item, drug, p.PatientAge)...)
		items = append(items, s.checkContraindication(item, drug, p, allergies)...)
	}
	items = append(items, s.checkInteractions(p)...)
	items = append(items, s.checkDuplication(p)...)
	return items
}

// checkIndication 适应症审核：诊断 ICD-10 与药品适应症匹配。
func (s *ReviewService) checkIndication(item *model.PrescriptionItem, drug *model.Drug, diags []model.Diagnosis) []model.ReviewItem {
	indications := drug.IndicationList()
	if len(indications) == 0 {
		return nil
	}
	matched := false
	for _, diag := range diags {
		for _, ind := range indications {
			if diag.ICD10 == ind.ICD10 {
				matched = true
				break
			}
		}
		if matched {
			break
		}
	}
	if matched {
		return nil
	}
	return []model.ReviewItem{{
		DrugName:   drug.Name,
		RuleType:   constants.RuleTypeIndication,
		RiskLevel:  constants.RiskMedium,
		Message:    fmt.Sprintf("药品 %s 的适应症与当前诊断不符（诊断：%s）", drug.Name, diagText(diags)),
		Suggestion: "请医生确认诊断或调整用药方案",
		RuleSource: "药品说明书适应症（ICD-10）",
		Operation:  "confirm",
	}}
}

// checkDosage 用法用量审核：日剂量与频次。
func (s *ReviewService) checkDosage(item *model.PrescriptionItem, drug *model.Drug, age int) []model.ReviewItem {
	var items []model.ReviewItem
	daily := item.SingleDose * frequencyDailyTimes(item.Frequency)
	item.DailyDose = daily
	maxDose := maxDailyDose(drug, age)
	if maxDose > 0 && daily > maxDose {
		risk := constants.RiskMedium
		if daily > maxDose*1.5 {
			risk = constants.RiskHigh
		}
		suggested := maxDose / frequencyDailyTimes(item.Frequency)
		item.SuggestedDose = suggested
		items = append(items, model.ReviewItem{
			DrugName:   drug.Name,
			RuleType:   constants.RuleTypeDosage,
			RiskLevel:  risk,
			Message:    fmt.Sprintf("药品 %s 日剂量 %.2f%s 超过%s最大剂量 %.2f%s", drug.Name, daily, drug.DoseUnit, ageGroupText(age), maxDose, drug.DoseUnit),
			Suggestion: fmt.Sprintf("建议单次剂量调整为 %.2f%s", suggested, drug.DoseUnit),
			RuleSource: "药品说明书用法用量",
			Operation:  "modify",
		})
	}
	if drug.MaxSingleDose > 0 && item.SingleDose > drug.MaxSingleDose {
		items = append(items, model.ReviewItem{
			DrugName:   drug.Name,
			RuleType:   constants.RuleTypeDosage,
			RiskLevel:  constants.RiskMedium,
			Message:    fmt.Sprintf("药品 %s 单次剂量 %.2f%s 超过最大单次剂量 %.2f%s", drug.Name, item.SingleDose, drug.DoseUnit, drug.MaxSingleDose, drug.DoseUnit),
			Suggestion: fmt.Sprintf("建议单次剂量不超过 %.2f%s", drug.MaxSingleDose, drug.DoseUnit),
			RuleSource: "药品说明书单次限量",
			Operation:  "modify",
		})
	}
	if drug.FrequencyLimit != "" && !frequencyAllowed(item.Frequency, drug.FrequencyLimit) {
		items = append(items, model.ReviewItem{
			DrugName:   drug.Name,
			RuleType:   constants.RuleTypeDosage,
			RiskLevel:  constants.RiskMedium,
			Message:    fmt.Sprintf("药品 %s 给药频次 %s 不在推荐频次范围内（%s）", drug.Name, item.Frequency, drug.FrequencyLimit),
			Suggestion: fmt.Sprintf("建议按 %s 频次给药", drug.FrequencyLimit),
			RuleSource: "药品说明书给药频次",
			Operation:  "modify",
		})
	}
	return items
}

// checkContraindication 禁忌与特殊人群审核。
func (s *ReviewService) checkContraindication(item *model.PrescriptionItem, drug *model.Drug, p *model.Prescription, allergies []string) []model.ReviewItem {
	var items []model.ReviewItem
	drugAllergies := drug.AllergyList()
	for _, a := range drugAllergies {
		for _, pa := range allergies {
			if a == pa {
				items = append(items, model.ReviewItem{
					DrugName:   drug.Name,
					RuleType:   constants.RuleTypeContraindication,
					RiskLevel:  constants.RiskCritical,
					Message:    fmt.Sprintf("患者对 %s 过敏，药品 %s 存在用药禁忌", pa, drug.Name),
					Suggestion: "禁止使用该药品，更换替代治疗方案",
					RuleSource: "药品说明书过敏禁忌",
					Operation:  "deny",
				})
			}
		}
	}
	if p.Pregnant && drug.PregnancyContraindicated {
		items = append(items, model.ReviewItem{
			DrugName:   drug.Name,
			RuleType:   constants.RuleTypeContraindication,
			RiskLevel:  constants.RiskCritical,
			Message:    fmt.Sprintf("患者处于妊娠状态，药品 %s 妊娠期禁忌", drug.Name),
			Suggestion: "禁止使用该药品，咨询产科医生调整方案",
			RuleSource: "药品说明书妊娠禁忌",
			Operation:  "deny",
		})
	}
	if p.HepaticImpairment && drug.HepaticContraindicated {
		items = append(items, model.ReviewItem{
			DrugName:   drug.Name,
			RuleType:   constants.RuleTypeContraindication,
			RiskLevel:  constants.RiskCritical,
			Message:    fmt.Sprintf("患者肝功能异常，药品 %s 肝功能禁忌", drug.Name),
			Suggestion: "禁止使用该药品，评估肝功能后调整方案",
			RuleSource: "药品说明书肝功能禁忌",
			Operation:  "deny",
		})
	}
	if p.RenalImpairment && drug.RenalContraindicated {
		items = append(items, model.ReviewItem{
			DrugName:   drug.Name,
			RuleType:   constants.RuleTypeContraindication,
			RiskLevel:  constants.RiskCritical,
			Message:    fmt.Sprintf("患者肾功能异常，药品 %s 肾功能禁忌", drug.Name),
			Suggestion: "禁止使用该药品，按肾功能调整剂量或换药",
			RuleSource: "药品说明书肾功能禁忌",
			Operation:  "deny",
		})
	}
	if p.PatientAge >= 65 && drug.BeersFlag {
		items = append(items, model.ReviewItem{
			DrugName:   drug.Name,
			RuleType:   constants.RuleTypeContraindication,
			RiskLevel:  constants.RiskHigh,
			Message:    fmt.Sprintf("老年人（%d 岁）使用 %s 符合 Beers 标准高风险条目", p.PatientAge, drug.Name),
			Suggestion: "评估跌倒、认知功能等风险，尽量选择更安全替代药",
			RuleSource: "Beers 标准（老年人潜在不适当用药）",
			Operation:  "modify",
		})
	}
	if drug.MinAgeMonths > 0 && p.PatientAge*12 < drug.MinAgeMonths {
		items = append(items, model.ReviewItem{
			DrugName:   drug.Name,
			RuleType:   constants.RuleTypeContraindication,
			RiskLevel:  constants.RiskHigh,
			Message:    fmt.Sprintf("儿童年龄 %d 岁低于药品 %s 的年龄下限（%d 月）", p.PatientAge, drug.Name, drug.MinAgeMonths),
			Suggestion: "禁止用于低龄儿童，选择儿童适用剂型",
			RuleSource: "儿童用药年龄限制",
			Operation:  "deny",
		})
	}
	return items
}

// checkInteractions 药物相互作用审核：所有药品两两组合。
func (s *ReviewService) checkInteractions(p *model.Prescription) []model.ReviewItem {
	var items []model.ReviewItem
	drugIDs := make([]uint, 0, len(p.Items))
	for i := range p.Items {
		if p.Items[i].Drug.ID != 0 {
			drugIDs = append(drugIDs, p.Items[i].Drug.ID)
		}
	}
	for i := 0; i < len(drugIDs); i++ {
		for j := i + 1; j < len(drugIDs); j++ {
			if drugIDs[i] == drugIDs[j] {
				continue
			}
			rule, err := s.inter.FindEnabledByPair(drugIDs[i], drugIDs[j])
			if err != nil {
				continue
			}
			items = append(items, model.ReviewItem{
				DrugName:   fmt.Sprintf("%s × %s", rule.DrugA.Name, rule.DrugB.Name),
				RuleType:   constants.RuleTypeInteraction,
				RiskLevel:  rule.RiskLevel,
				Message:    fmt.Sprintf("药物相互作用（%s）：%s", util.RiskLevelText(rule.RiskLevel), rule.Description),
				Suggestion: fmt.Sprintf("关注机制 %s，加强监测或调整用药", rule.Mechanism),
				RuleSource: "药物相互作用规则库",
				Operation:  "confirm",
			})
		}
	}
	return items
}

// checkDuplication 重复用药审核：相同作用机制药品重复开具。
func (s *ReviewService) checkDuplication(p *model.Prescription) []model.ReviewItem {
	var items []model.ReviewItem
	byMechanism := make(map[string][]string)
	for i := range p.Items {
		drug := &p.Items[i].Drug
		if drug.ID == 0 || drug.Mechanism == "" {
			continue
		}
		byMechanism[drug.Mechanism] = append(byMechanism[drug.Mechanism], drug.Name)
	}
	for mechanism, names := range byMechanism {
		if len(names) < 2 {
			continue
		}
		risk := constants.RiskMedium
		if constants.IsDuplicationHighRiskMechanism(mechanism) {
			risk = constants.RiskHigh
		}
		items = append(items, model.ReviewItem{
			DrugName:   joinNames(names),
			RuleType:   constants.RuleTypeDuplication,
			RiskLevel:  risk,
			Message:    fmt.Sprintf("重复用药：%s 作用机制相同（%s），存在重复开具风险", joinNames(names), mechanism),
			Suggestion: "建议保留一种药品，避免同类药物叠加",
			RuleSource: "重复用药检测（作用机制分组）",
			Operation:  "modify",
		})
	}
	return items
}

func diagText(diags []model.Diagnosis) string {
	if len(diags) == 0 {
		return "无诊断"
	}
	return diags[0].ICD10 + " " + diags[0].Name
}

func ageGroupText(age int) string {
	switch ageGroup(age) {
	case "child":
		return "儿童"
	case "elderly":
		return "老人"
	default:
		return "成人"
	}
}

func frequencyAllowed(freq, limit string) bool {
	allowed := map[string]bool{}
	for _, f := range splitComma(limit) {
		allowed[f] = true
	}
	return allowed[freq]
}

func splitComma(s string) []string {
	var out []string
	cur := ""
	for _, r := range s {
		if r == ',' {
			if cur != "" {
				out = append(out, cur)
				cur = ""
			}
			continue
		}
		cur += string(r)
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}

func joinNames(names []string) string {
	out := ""
	for i, n := range names {
		if i > 0 {
			out += "、"
		}
		out += n
	}
	return out
}

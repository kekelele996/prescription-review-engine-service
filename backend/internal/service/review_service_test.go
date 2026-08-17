package service

import (
	"encoding/json"
	"testing"

	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/model"
	"github.com/rxcheck/rxcheck/internal/repository"
	"github.com/rxcheck/rxcheck/pkg/testutil"
	"gorm.io/datatypes"
)

// itemSpec 处方明细构造参数。
type itemSpec struct {
	drugName string
	dose     float64
	freq     string
	days     int
}

// createTestPrescription 创建测试处方并返回 ID。
func createTestPrescription(t *testing.T, prescRepo *repository.PrescriptionRepository, no string, age int, allergies []string, renal bool, diags []model.Diagnosis, items []itemSpec, ids map[string]uint) uint {
	t.Helper()
	prescItems := make([]model.PrescriptionItem, 0, len(items))
	for _, it := range items {
		prescItems = append(prescItems, model.PrescriptionItem{
			DrugID:     ids[it.drugName],
			DrugName:   it.drugName,
			SingleDose: it.dose,
			Frequency:  it.freq,
			CourseDays: it.days,
			DailyDose:  it.dose * frequencyDailyTimes(it.freq),
		})
	}
	b, _ := json.Marshal(diags)
	presc := &model.Prescription{
		PrescriptionNo:  no,
		PatientName:     "测试患者",
		PatientAge:      age,
		Allergies:       testutil.StringsJSON(allergies),
		RenalImpairment: renal,
		Diagnoses:       datatypes.JSON(b),
		Status:          constants.PrescriptionStatusPendingReview,
		OverallRisk:     constants.RiskNone,
	}
	if err := prescRepo.CreateWithItemsTx(prescRepo.DB(), presc, prescItems); err != nil {
		t.Fatalf("create prescription failed: %v", err)
	}
	return presc.ID
}

func TestReviewPrescription(t *testing.T) {
	db := newTestDB(t)
	svc, ids := newReviewService(t, db)
	prescRepo := newPrescRepo(db)

	cases := []struct {
		name       string
		age        int
		allergies  []string
		renal      bool
		diags      []model.Diagnosis
		items      []itemSpec
		wantStatus string
		wantRisk   string
		wantRule   string
	}{
		{name: "正常处方通过", age: 50, diags: []model.Diagnosis{{ICD10: "I10", Name: "原发性高血压"}},
			items:      []itemSpec{{"氨氯地平", 5, "qd", 30}, {"赖诺普利", 10, "qd", 30}},
			wantStatus: constants.PrescriptionStatusPassed, wantRisk: constants.RiskNone},
		{name: "严重相互作用警告", age: 60, diags: []model.Diagnosis{{ICD10: "I48", Name: "心房颤动"}},
			items:      []itemSpec{{"华法林", 2.5, "qd", 30}, {"阿司匹林", 100, "qd", 30}},
			wantStatus: constants.PrescriptionStatusWarned, wantRisk: constants.RiskHigh, wantRule: constants.RuleTypeInteraction},
		{name: "青霉素过敏拒绝", age: 40, allergies: []string{"青霉素"}, diags: []model.Diagnosis{{ICD10: "J18", Name: "肺炎"}},
			items:      []itemSpec{{"头孢呋辛", 250, "bid", 7}},
			wantStatus: constants.PrescriptionStatusRejected, wantRisk: constants.RiskCritical, wantRule: constants.RuleTypeContraindication},
		{name: "超剂量警告", age: 45, diags: []model.Diagnosis{{ICD10: "I10", Name: "原发性高血压"}},
			items:      []itemSpec{{"氨氯地平", 10, "tid", 30}},
			wantStatus: constants.PrescriptionStatusWarned, wantRisk: constants.RiskHigh, wantRule: constants.RuleTypeDosage},
		{name: "重复用药警告", age: 50, diags: []model.Diagnosis{{ICD10: "I10", Name: "原发性高血压"}},
			items:      []itemSpec{{"赖诺普利", 10, "qd", 30}, {"依那普利", 5, "qd", 30}},
			wantStatus: constants.PrescriptionStatusWarned, wantRisk: constants.RiskHigh, wantRule: constants.RuleTypeDuplication},
		{name: "老年Beers警告", age: 70, diags: []model.Diagnosis{{ICD10: "F41", Name: "焦虑障碍"}},
			items:      []itemSpec{{"地西泮", 5, "qd", 14}},
			wantStatus: constants.PrescriptionStatusWarned, wantRisk: constants.RiskHigh, wantRule: constants.RuleTypeContraindication},
		{name: "儿童年龄限制警告", age: 10, diags: []model.Diagnosis{{ICD10: "I20", Name: "心绞痛"}},
			items:      []itemSpec{{"阿司匹林", 100, "qd", 14}},
			wantStatus: constants.PrescriptionStatusWarned, wantRisk: constants.RiskHigh, wantRule: constants.RuleTypeContraindication},
		{name: "肾功能禁忌拒绝", age: 55, renal: true, diags: []model.Diagnosis{{ICD10: "E11", Name: "2型糖尿病"}},
			items:      []itemSpec{{"二甲双胍", 500, "bid", 30}},
			wantStatus: constants.PrescriptionStatusRejected, wantRisk: constants.RiskCritical, wantRule: constants.RuleTypeContraindication},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			id := createTestPrescription(t, prescRepo, "RX-"+tc.name, tc.age, tc.allergies, tc.renal, tc.diags, tc.items, ids)
			report, err := svc.ReviewPrescription(id, "tester", "", "")
			if err != nil {
				t.Fatalf("review failed: %v", err)
			}
			if report.Status != tc.wantStatus {
				t.Errorf("status = %s, want %s", report.Status, tc.wantStatus)
			}
			if report.RiskLevel != tc.wantRisk {
				t.Errorf("risk = %s, want %s", report.RiskLevel, tc.wantRisk)
			}
			if tc.wantRule != "" {
				found := false
				for _, it := range report.Items {
					if it.RuleType == tc.wantRule {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("missing rule type %s in items: %+v", tc.wantRule, report.Items)
				}
			}
			presc, err := prescRepo.FindByID(id)
			if err != nil {
				t.Fatalf("load prescription failed: %v", err)
			}
			if presc.Status != tc.wantStatus {
				t.Errorf("prescription status = %s, want %s", presc.Status, tc.wantStatus)
			}
		})
	}
}

func TestReviewPrescriptionNotFound(t *testing.T) {
	db := newTestDB(t)
	svc, _ := newReviewService(t, db)
	if _, err := svc.ReviewPrescription(9999, "tester", "", ""); err == nil {
		t.Fatal("expected error for missing prescription")
	}
}


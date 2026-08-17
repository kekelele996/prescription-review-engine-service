package service

import (
	"strings"
	"testing"

	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/model"
)

func TestReviewAggregation_CriticalRejected(t *testing.T) {
	db := newTestDB(t)
	svc, ids := newReviewService(t, db)
	prescRepo := newPrescRepo(db)
	id := createTestPrescription(t, prescRepo, "RX-AGG-CRIT", 40, []string{"青霉素"}, false,
		[]model.Diagnosis{{ICD10: "J18", Name: "肺炎"}},
		[]itemSpec{{"头孢呋辛", 250, "bid", 7}, {"氨氯地平", 5, "qd", 30}}, ids)
	report, err := svc.ReviewPrescription(id, "tester", "", "")
	if err != nil {
		t.Fatalf("review failed: %v", err)
	}
	if report.Status != constants.PrescriptionStatusRejected {
		t.Errorf("status = %s, want rejected", report.Status)
	}
	if report.RiskLevel != constants.RiskCritical {
		t.Errorf("risk = %s, want critical", report.RiskLevel)
	}
	if !strings.Contains(report.Summary, "极高风险") {
		t.Errorf("summary should mention 极高风险: %s", report.Summary)
	}
}

func TestReviewAggregation_HighWarned(t *testing.T) {
	db := newTestDB(t)
	svc, ids := newReviewService(t, db)
	prescRepo := newPrescRepo(db)
	id := createTestPrescription(t, prescRepo, "RX-AGG-HIGH", 60, nil, false,
		[]model.Diagnosis{{ICD10: "I48", Name: "心房颤动"}},
		[]itemSpec{{"华法林", 2.5, "qd", 30}, {"阿司匹林", 100, "qd", 30}}, ids)
	report, err := svc.ReviewPrescription(id, "tester", "", "")
	if err != nil {
		t.Fatalf("review failed: %v", err)
	}
	if report.Status != constants.PrescriptionStatusWarned {
		t.Errorf("status = %s, want warned", report.Status)
	}
}

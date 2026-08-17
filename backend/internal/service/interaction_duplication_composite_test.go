package service

import (
	"testing"

	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/model"
)

func TestInteraction_ReverseOrderDetected(t *testing.T) {
	db := newTestDB(t)
	svc, ids := newReviewService(t, db)
	prescRepo := newPrescRepo(db)
	id := createTestPrescription(t, prescRepo, "RX-INT-REV", 60, nil, false,
		[]model.Diagnosis{{ICD10: "I48", Name: "心房颤动"}},
		[]itemSpec{{"阿司匹林", 100, "qd", 30}, {"华法林", 2.5, "qd", 30}}, ids)
	report, err := svc.ReviewPrescription(id, "tester", "", "")
	if err != nil {
		t.Fatalf("review failed: %v", err)
	}
	for _, it := range report.Items {
		if it.RuleType == constants.RuleTypeInteraction {
			return
		}
	}
	t.Fatal("expected interaction item for reversed drug pair")
}

func TestDuplication_AceHigh(t *testing.T) {
	db := newTestDB(t)
	svc, ids := newReviewService(t, db)
	prescRepo := newPrescRepo(db)
	id := createTestPrescription(t, prescRepo, "RX-DUP-ACE", 50, nil, false,
		[]model.Diagnosis{{ICD10: "I10", Name: "原发性高血压"}},
		[]itemSpec{{"赖诺普利", 10, "qd", 30}, {"依那普利", 5, "qd", 30}}, ids)
	report, err := svc.ReviewPrescription(id, "tester", "", "")
	if err != nil {
		t.Fatalf("review failed: %v", err)
	}
	found, high := false, false
	for _, it := range report.Items {
		if it.RuleType == constants.RuleTypeDuplication {
			found = true
			if it.RiskLevel == constants.RiskHigh {
				high = true
			}
		}
	}
	if !found || !high {
		t.Fatalf("expected high duplication item, found=%v high=%v", found, high)
	}
}

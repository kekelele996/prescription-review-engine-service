package service

import (
	"testing"

	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/model"
	"github.com/rxcheck/rxcheck/internal/util"
)

func TestFrequencyDailyTimes_Q8H(t *testing.T) {
	if got := frequencyDailyTimes("q8h"); got != 3 {
		t.Errorf("q8h daily times = %v, want 3", got)
	}
}

func TestMaxDailyDose_Child(t *testing.T) {
	d := &model.Drug{AdultMaxDailyDose: 10, ChildMaxDailyDose: 5, ElderlyMaxDailyDose: 5}
	if got := maxDailyDose(d, 10); got != 5 {
		t.Errorf("child max dose = %v, want 5", got)
	}
}

func TestFormatDose_Precision(t *testing.T) {
	if got := util.FormatDose(4); got != "4.00" {
		t.Errorf("FormatDose(4) = %q, want 4.00", got)
	}
}

func TestReviewPersistsOverallRisk(t *testing.T) {
	db := newTestDB(t)
	svc, ids := newReviewService(t, db)
	prescRepo := newPrescRepo(db)
	id := createTestPrescription(t, prescRepo, "RX-DOSE-001", 45, nil, false,
		[]model.Diagnosis{{ICD10: "I10", Name: "原发性高血压"}},
		[]itemSpec{{"氨氯地平", 5, "tid", 30}}, ids)
	if _, err := svc.ReviewPrescription(id, "tester", "", ""); err != nil {
		t.Fatalf("review failed: %v", err)
	}
	presc, err := prescRepo.FindByID(id)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if presc.OverallRisk != constants.RiskMedium {
		t.Errorf("overall risk = %q, want medium", presc.OverallRisk)
	}
}

func TestMapRiskToStatus_Critical(t *testing.T) {
	if got := constants.MapRiskToStatus(constants.RiskCritical); got != constants.PrescriptionStatusRejected {
		t.Errorf("MapRiskToStatus(critical) = %q, want rejected", got)
	}
}

package service

import (
	"testing"

	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/dto"
	"github.com/rxcheck/rxcheck/internal/model"
	"github.com/rxcheck/rxcheck/internal/util"
	"gorm.io/gorm"
)

// newPrescriptionService 构造处方服务（nil Redis，触发同步审核降级路径）。
func newPrescriptionService(t *testing.T, db *gorm.DB) (*PrescriptionService, map[string]uint) {
	t.Helper()
	ids, err := testutilSeedDrugs(db)
	if err != nil {
		t.Fatalf("seed drugs failed: %v", err)
	}
	if err := testutilSeedInteractions(db, ids); err != nil {
		t.Fatalf("seed interactions failed: %v", err)
	}
	prescRepo := newPrescRepo(db)
	drugRepo := newDrugRepo(db)
	interRepo := newInterRepo(db)
	reportRepo := newReportRepo(db)
	audit := NewAuditService(newAuditRepo(db), util.Log)
	review := NewReviewService(db, prescRepo, drugRepo, interRepo, reportRepo, audit, util.Log)
	return NewPrescriptionService(db, prescRepo, drugRepo, review, audit, nil, util.Log), ids
}

func TestSubmitSyncFallback(t *testing.T) {
	db := newTestDB(t)
	svc, ids := newPrescriptionService(t, db)
	req := &dto.SubmitPrescriptionReq{
		PrescriptionNo: "RX-SUBMIT-001",
		PatientName:    "张三",
		PatientAge:     50,
		Diagnoses:      []dto.DiagnosisDTO{{ICD10: "I10", Name: "原发性高血压"}},
		Items: []dto.PrescriptionItemDTO{
			{DrugID: ids["氨氯地平"], SingleDose: 5, Frequency: "qd", CourseDays: 30},
			{DrugID: ids["赖诺普利"], SingleDose: 10, Frequency: "qd", CourseDays: 30},
		},
	}
	presc, err := svc.Submit(req, "json", []byte(`{}`), &util.Claims{UserID: 1, Username: "doctor", Role: constants.UserRoleDoctor}, "127.0.0.1", "test-rid")
	if err != nil {
		t.Fatalf("submit failed: %v", err)
	}
	if presc.Status != constants.PrescriptionStatusPassed {
		t.Errorf("status = %s, want passed", presc.Status)
	}
	if len(presc.Items) != 2 {
		t.Errorf("items = %d, want 2", len(presc.Items))
	}
}

func TestSubmitDrugNotFound(t *testing.T) {
	db := newTestDB(t)
	svc, _ := newPrescriptionService(t, db)
	req := &dto.SubmitPrescriptionReq{
		PrescriptionNo: "RX-BAD-001",
		PatientName:    "李四",
		PatientAge:     40,
		Diagnoses:      []dto.DiagnosisDTO{{ICD10: "J18", Name: "肺炎"}},
		Items: []dto.PrescriptionItemDTO{
			{DrugName: "不存在的药品", SingleDose: 250, Frequency: "bid", CourseDays: 7},
		},
	}
	if _, err := svc.Submit(req, "json", []byte(`{}`), &util.Claims{UserID: 1, Username: "doctor", Role: constants.UserRoleDoctor}, "127.0.0.1", "test-rid"); err == nil {
		t.Fatal("expected drug not found error")
	}
}

func TestOverrideFlow(t *testing.T) {
	db := newTestDB(t)
	svc, ids := newPrescriptionService(t, db)
	// 华法林+阿司匹林 -> 警告类处方。
	id := createTestPrescription(t, newPrescRepo(db), "RX-OVERRIDE-001", 60, nil, false,
		[]model.Diagnosis{{ICD10: "I48", Name: "心房颤动"}},
		[]itemSpec{{"华法林", 2.5, "qd", 30}, {"阿司匹林", 100, "qd", 30}}, ids)
	if _, err := svc.Review(id, "doctor", "", ""); err != nil {
		t.Fatalf("review failed: %v", err)
	}
	doctor := &util.Claims{UserID: 1, Username: "doctor", Role: constants.UserRoleDoctor}
	presc, err := svc.Override(id, "患者合并房颤，抗栓治疗获益大于风险", doctor, "127.0.0.1", "test-rid")
	if err != nil {
		t.Fatalf("override failed: %v", err)
	}
	if presc.Status != constants.PrescriptionStatusOverridden {
		t.Errorf("status = %s, want overridden", presc.Status)
	}
	if presc.OverrideReason == "" {
		t.Error("override reason should be recorded")
	}
	// 终态重复强制通过应报错。
	if _, err := svc.Override(id, "再次通过", doctor, "127.0.0.1", "test-rid"); err == nil {
		t.Fatal("expected override not allowed error")
	}
}

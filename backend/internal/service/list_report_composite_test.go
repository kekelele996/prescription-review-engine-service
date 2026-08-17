package service

import (
	"testing"

	"github.com/rxcheck/rxcheck/internal/dto"
	"github.com/rxcheck/rxcheck/internal/model"
	"github.com/rxcheck/rxcheck/internal/util"
)

func TestPrescriptionList_NoPanicAndOrder(t *testing.T) {
	db := newTestDB(t)
	svc, ids := newPrescriptionService(t, db)
	prescRepo := newPrescRepo(db)
	_ = createTestPrescription(t, prescRepo, "RX-LIST-01", 50, nil, false,
		[]model.Diagnosis{{ICD10: "I10", Name: "原发性高血压"}},
		[]itemSpec{{"氨氯地平", 5, "qd", 30}}, ids)
	newer := createTestPrescription(t, prescRepo, "RX-LIST-02", 55, nil, false,
		[]model.Diagnosis{{ICD10: "I10", Name: "原发性高血压"}},
		[]itemSpec{{"氨氯地平", 5, "qd", 30}}, ids)
	if _, err := svc.Review(newer, "doctor", "", ""); err != nil {
		t.Fatalf("review failed: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("List panicked: %v", r)
		}
	}()
	res, err := svc.List(1, 10, "", "")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	items, ok := res.List.([]dto.PrescriptionListResp)
	if !ok {
		t.Fatalf("unexpected list type: %T", res.List)
	}
	if len(items) != 2 {
		t.Fatalf("items = %d, want 2", len(items))
	}
	if items[0].ID != newer {
		t.Errorf("first item id = %d, want newer %d", items[0].ID, newer)
	}
}

func TestReportItems_Loaded(t *testing.T) {
	db := newTestDB(t)
	svc, ids := newPrescriptionService(t, db)
	id := createTestPrescription(t, newPrescRepo(db), "RX-REP-01", 60, nil, false,
		[]model.Diagnosis{{ICD10: "I48", Name: "心房颤动"}},
		[]itemSpec{{"华法林", 2.5, "qd", 30}, {"阿司匹林", 100, "qd", 30}}, ids)
	if _, err := svc.Review(id, "doctor", "", ""); err != nil {
		t.Fatalf("review failed: %v", err)
	}
	audit := NewAuditService(newAuditRepo(db), util.Log)
	rs := NewReportService(newReportRepo(db), newPrescRepo(db), nil, audit, util.Log)
	resp, err := rs.GetByPrescription(id)
	if err != nil {
		t.Fatalf("get report failed: %v", err)
	}
	if len(resp.Items) == 0 {
		t.Fatal("expected report items, got none")
	}
}

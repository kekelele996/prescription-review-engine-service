package service

import (
	"testing"

	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/model"
	"github.com/rxcheck/rxcheck/internal/util"
)

func TestOverride_RejectedNotAllowed(t *testing.T) {
	db := newTestDB(t)
	svc, ids := newPrescriptionService(t, db)
	id := createTestPrescription(t, newPrescRepo(db), "RX-OVR-REJ", 40, []string{"青霉素"}, false,
		[]model.Diagnosis{{ICD10: "J18", Name: "肺炎"}},
		[]itemSpec{{"头孢呋辛", 250, "bid", 7}}, ids)
	if _, err := svc.Review(id, "doctor", "", ""); err != nil {
		t.Fatalf("review failed: %v", err)
	}
	doctor := &util.Claims{UserID: 1, Username: "doctor", Role: constants.UserRoleDoctor}
	if _, err := svc.Override(id, "必须使用该药", doctor, "127.0.0.1", "rid"); err == nil {
		t.Fatal("expected override not allowed for rejected prescription")
	}
}

func TestReReview_OverriddenBlocked(t *testing.T) {
	db := newTestDB(t)
	svc, ids := newPrescriptionService(t, db)
	id := createTestPrescription(t, newPrescRepo(db), "RX-OVR-DONE", 60, nil, false,
		[]model.Diagnosis{{ICD10: "I48", Name: "心房颤动"}},
		[]itemSpec{{"华法林", 2.5, "qd", 30}, {"阿司匹林", 100, "qd", 30}}, ids)
	if _, err := svc.Review(id, "doctor", "", ""); err != nil {
		t.Fatalf("review failed: %v", err)
	}
	doctor := &util.Claims{UserID: 1, Username: "doctor", Role: constants.UserRoleDoctor}
	if _, err := svc.Override(id, "抗栓获益大于风险", doctor, "127.0.0.1", "rid"); err != nil {
		t.Fatalf("override failed: %v", err)
	}
	if _, err := svc.Review(id, "doctor", "", ""); err == nil {
		t.Fatal("expected re-review of overridden prescription to be blocked")
	}
}

func TestStatusText_Overridden(t *testing.T) {
	if got := util.PrescriptionStatusText(constants.PrescriptionStatusOverridden); got != "强制通过" {
		t.Errorf("PrescriptionStatusText(overridden) = %q, want 强制通过", got)
	}
}

func TestNormalReview_Passed(t *testing.T) {
	db := newTestDB(t)
	svc, ids := newReviewService(t, db)
	id := createTestPrescription(t, newPrescRepo(db), "RX-NORM-001", 50, nil, false,
		[]model.Diagnosis{{ICD10: "I10", Name: "原发性高血压"}},
		[]itemSpec{{"氨氯地平", 5, "qd", 30}, {"赖诺普利", 10, "qd", 30}}, ids)
	report, err := svc.ReviewPrescription(id, "tester", "", "")
	if err != nil {
		t.Fatalf("review failed: %v", err)
	}
	if report.Status != constants.PrescriptionStatusPassed {
		t.Errorf("status = %s, want passed", report.Status)
	}
}

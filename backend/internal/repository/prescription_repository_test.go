package repository

import (
	"testing"
	"time"

	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/model"
	"github.com/rxcheck/rxcheck/pkg/testutil"
)

func TestPrescriptionRepositoryWithItems(t *testing.T) {
	db := newTestDB(t)
	repo := NewPrescriptionRepository(db)
	ids, err := testutil.SeedDrugs(db)
	if err != nil {
		t.Fatalf("seed drugs failed: %v", err)
	}

	items := []model.PrescriptionItem{
		{DrugID: ids["氨氯地平"], DrugName: "氨氯地平", SingleDose: 5, Frequency: "qd", CourseDays: 30, DailyDose: 5},
		{DrugID: ids["赖诺普利"], DrugName: "赖诺普利", SingleDose: 10, Frequency: "qd", CourseDays: 30, DailyDose: 10},
	}
	p := &model.Prescription{
		PrescriptionNo: "RX-REPO-001",
		PatientName:    "王五",
		PatientAge:     52,
		Diagnoses:      testutil.Indications([2]string{"I10", "原发性高血压"}),
		Status:         constants.PrescriptionStatusPendingReview,
		OverallRisk:    constants.RiskNone,
	}
	if err := repo.CreateWithItemsTx(repo.DB(), p, items); err != nil {
		t.Fatalf("create with items failed: %v", err)
	}
	if p.ID == 0 {
		t.Fatal("id should be assigned")
	}

	got, err := repo.FindByID(p.ID)
	if err != nil {
		t.Fatalf("find by id failed: %v", err)
	}
	if len(got.Items) != 2 {
		t.Errorf("items = %d, want 2", len(got.Items))
	}
	if got.Items[0].Drug.ID == 0 {
		t.Error("drug should be preloaded")
	}

	// 状态流转。
	now := time.Now()
	if err := repo.UpdateStatusTx(repo.DB(), p.ID, constants.PrescriptionStatusPassed, constants.RiskNone, &now); err != nil {
		t.Fatalf("update status failed: %v", err)
	}
	got2, _ := repo.FindByID(p.ID)
	if got2.Status != constants.PrescriptionStatusPassed {
		t.Errorf("status = %s", got2.Status)
	}

	// 列表与统计。
	list, total, err := repo.List(1, 10, constants.PrescriptionStatusPassed, "")
	if err != nil || total != 1 || len(list) != 1 {
		t.Errorf("list total=%d len=%d err=%v", total, len(list), err)
	}
	byStatus, err := repo.CountByStatus()
	if err != nil || byStatus[constants.PrescriptionStatusPassed] != 1 {
		t.Errorf("byStatus = %v err=%v", byStatus, err)
	}
}

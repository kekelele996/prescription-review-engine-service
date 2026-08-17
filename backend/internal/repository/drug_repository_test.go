package repository

import (
	"testing"

	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/model"
	"github.com/rxcheck/rxcheck/pkg/testutil"
)

func TestDrugRepositoryCRUD(t *testing.T) {
	db := newTestDB(t)
	repo := NewDrugRepository(db)

	drug := &model.Drug{
		Name: "阿莫西林", GenericName: "阿莫西林胶囊", Specification: "250mg", Route: "口服",
		AtcCode: "J01CA04", Mechanism: "penicillin",
		Indications: testutil.Indications([2]string{"J18", "肺炎"}),
		AdultMaxDailyDose: 3000, DoseUnit: "mg", Status: constants.DrugStatusEnabled,
	}
	if err := repo.Create(drug); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if drug.ID == 0 {
		t.Fatal("id should be assigned")
	}

	// FindByID / FindByName 复用。
	got, err := repo.FindByID(drug.ID)
	if err != nil {
		t.Fatalf("find by id failed: %v", err)
	}
	if got.Name != "阿莫西林" {
		t.Errorf("name = %s", got.Name)
	}
	got2, err := repo.FindByName("阿莫西林")
	if err != nil {
		t.Fatalf("find by name failed: %v", err)
	}
	if got2.ID != drug.ID {
		t.Errorf("id = %d, want %d", got2.ID, drug.ID)
	}

	// 列表检索。
	list, total, err := repo.List(1, 10, "阿莫西林", "", "")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if total != 1 || len(list) != 1 {
		t.Errorf("total=%d len=%d, want 1/1", total, len(list))
	}

	// 状态流转。
	if err := repo.UpdateStatusTx(repo.DB(), drug.ID, constants.DrugStatusDisabled); err != nil {
		t.Fatalf("disable failed: %v", err)
	}
	got3, _ := repo.FindByID(drug.ID)
	if got3.Status != constants.DrugStatusDisabled {
		t.Errorf("status = %s", got3.Status)
	}

	// 统计。
	n, err := repo.Count()
	if err != nil || n != 1 {
		t.Errorf("count = %d err=%v, want 1", n, err)
	}
}

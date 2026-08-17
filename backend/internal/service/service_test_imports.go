package service

import (
	"github.com/rxcheck/rxcheck/internal/repository"
	"github.com/rxcheck/rxcheck/pkg/testutil"
	"gorm.io/gorm"
)

func testutilSeedDrugs(db *gorm.DB) (map[string]uint, error) {
	return testutil.SeedDrugs(db)
}

func testutilSeedInteractions(db *gorm.DB, ids map[string]uint) error {
	return testutil.SeedInteractions(db, ids)
}

func newPrescRepo(db *gorm.DB) *repository.PrescriptionRepository { return repository.NewPrescriptionRepository(db) }
func newDrugRepo(db *gorm.DB) *repository.DrugRepository          { return repository.NewDrugRepository(db) }
func newInterRepo(db *gorm.DB) *repository.InteractionRepository  { return repository.NewInteractionRepository(db) }
func newReportRepo(db *gorm.DB) *repository.ReportRepository      { return repository.NewReportRepository(db) }
func newAuditRepo(db *gorm.DB) *repository.AuditRepository        { return repository.NewAuditRepository(db) }

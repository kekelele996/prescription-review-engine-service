package service

import (
	"github.com/glebarez/sqlite"
	"github.com/rxcheck/rxcheck/internal/model"
	"github.com/rxcheck/rxcheck/internal/util"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// newTestDB 创建内存 SQLite 测试库并迁移全部表。
func newTestDB(t testingT) *gorm.DB {
	util.InitLogger(-4)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}
	if err := db.AutoMigrate(
		&model.User{}, &model.Drug{}, &model.InteractionRule{},
		&model.Prescription{}, &model.PrescriptionItem{}, &model.ReviewReport{},
		&model.ReviewItem{}, &model.AuditLog{},
	); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}
	return db
}

// newReviewService 构造带种子数据的审核服务。
func newReviewService(t testingT, db *gorm.DB) (*ReviewService, map[string]uint) {
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
	return NewReviewService(db, prescRepo, drugRepo, interRepo, reportRepo, audit, util.Log), ids
}

type testingT interface {
	Fatalf(format string, args ...any)
}

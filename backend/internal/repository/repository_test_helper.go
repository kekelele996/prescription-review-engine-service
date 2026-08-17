package repository

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

type testingT interface {
	Fatalf(format string, args ...any)
}

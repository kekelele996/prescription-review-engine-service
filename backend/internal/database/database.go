package database

import (
	"fmt"

	"github.com/rxcheck/rxcheck/internal/config"
	"github.com/rxcheck/rxcheck/internal/model"
	"github.com/rxcheck/rxcheck/internal/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// New 建立 PostgreSQL 连接并执行自动迁移。
func New(cfg *config.Config) (*gorm.DB, error) {
	level := logger.Warn
	if cfg.RunMode == "debug" {
		level = logger.Info
	}
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(level),
	})
	if err != nil {
		util.Log.Error(fmt.Sprintf("数据库连接失败: %v", err))
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	if err := db.AutoMigrate(
		&model.User{},
		&model.Drug{},
		&model.InteractionRule{},
		&model.Prescription{},
		&model.PrescriptionItem{},
		&model.ReviewReport{},
		&model.ReviewItem{},
		&model.AuditLog{},
	); err != nil {
		util.Log.Error(fmt.Sprintf("数据库迁移失败: %v", err))
		return nil, err
	}
	return db, nil
}

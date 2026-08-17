package repository

import (
	"errors"
	"time"

	"github.com/rxcheck/rxcheck/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PrescriptionRepository 处方仓储（含明细聚合）。
type PrescriptionRepository struct {
	db *gorm.DB
}

func NewPrescriptionRepository(db *gorm.DB) *PrescriptionRepository {
	return &PrescriptionRepository{db: db}
}

// DB 返回底层数据库句柄。
func (r *PrescriptionRepository) DB() *gorm.DB { return r.db }

// CreateWithItemsTx 在事务中创建处方及明细。
func (r *PrescriptionRepository) CreateWithItemsTx(tx *gorm.DB, p *model.Prescription, items []model.PrescriptionItem) error {
	p.Items = nil // 避免 GORM 自动保存 has-many 关联导致重复插入
	if err := tx.Create(p).Error; err != nil {
		return err
	}
	for i := range items {
		items[i].PrescriptionID = p.ID
	}
	if len(items) > 0 {
		if err := tx.Create(&items).Error; err != nil {
			return err
		}
	}
	p.Items = items
	return nil
}

// FindByID 查询处方详情（含明细、药品、报告）。
func (r *PrescriptionRepository) FindByID(id uint) (*model.Prescription, error) {
	var p model.Prescription
	err := r.db.
		Preload("Items", func(db *gorm.DB) *gorm.DB { return db.Order("prescription_items.id ASC") }).
		Preload("Items.Drug").
		Preload("Report.Items").
		Preload("Doctor").
		First(&p, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &p, err
}

// FindByIDForUpdate 在事务中加锁查询处方（并发安全，状态流转必须使用）。
// PostgreSQL 使用 SELECT ... FOR UPDATE；SQLite 测试环境不支持该语法时跳过锁。
func (r *PrescriptionRepository) FindByIDForUpdate(tx *gorm.DB, id uint) (*model.Prescription, error) {
	q := tx
	if tx.Dialector.Name() != "sqlite" {
		q = tx.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	var p model.Prescription
	err := q.Preload("Items").Preload("Report").First(&p, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &p, err
}

// List 分页查询处方，支持状态/患者关键字过滤。
func (r *PrescriptionRepository) List(page, pageSize int, status, keyword string) ([]model.Prescription, int64, error) {
	var list []model.Prescription
	var total int64
	q := r.db.Model(&model.Prescription{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if keyword != "" {
		q = q.Where("prescription_no LIKE ? OR patient_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("id ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

// UpdateStatusTx 在事务中更新处方审核状态机。
func (r *PrescriptionRepository) UpdateStatusTx(tx *gorm.DB, id uint, status, risk string, reviewedAt *time.Time) error {
	res := tx.Model(&model.Prescription{}).Where("id = ?", id).
		Updates(map[string]any{"status": status, "overall_risk": risk, "reviewed_at": reviewedAt})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// UpdateTx 在事务中整体更新处方。
func (r *PrescriptionRepository) UpdateTx(tx *gorm.DB, p *model.Prescription) error {
	return tx.Save(p).Error
}

// CountByStatus 按审核状态统计（统计接口复用）。
func (r *PrescriptionRepository) CountByStatus() (map[string]int64, error) {
	return groupCount(r.db, &model.Prescription{}, "status")
}

// CountByRisk 按风险等级统计（统计接口复用）。
func (r *PrescriptionRepository) CountByRisk() (map[string]int64, error) {
	return groupCount(r.db, &model.Prescription{}, "overall_risk")
}

// Count 统计处方总数（统计接口复用）。
func (r *PrescriptionRepository) Count() (int64, error) {
	var n int64
	err := r.db.Model(&model.Prescription{}).Count(&n).Error
	return n, err
}

package repository

import (
	"errors"

	"github.com/rxcheck/rxcheck/internal/model"
	"gorm.io/gorm"
)

// DrugRepository 药品档案仓储。
type DrugRepository struct {
	db *gorm.DB
}

func NewDrugRepository(db *gorm.DB) *DrugRepository {
	return &DrugRepository{db: db}
}

// DB 返回底层数据库句柄。
func (r *DrugRepository) DB() *gorm.DB { return r.db }

// Create 创建药品。
func (r *DrugRepository) Create(d *model.Drug) error {
	return r.db.Create(d).Error
}

// FindByID 按 ID 查询药品。
func (r *DrugRepository) FindByID(id uint) (*model.Drug, error) {
	var d model.Drug
	err := r.db.First(&d, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &d, err
}

// FindByName 按名称查询药品。
func (r *DrugRepository) FindByName(name string) (*model.Drug, error) {
	var d model.Drug
	err := r.db.Where("name = ?", name).First(&d).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &d, err
}

// FindByIDs 批量查询药品（保持给定顺序）。
func (r *DrugRepository) FindByIDs(ids []uint) ([]model.Drug, error) {
	var list []model.Drug
	err := r.db.Where("id IN ?", ids).Find(&list).Error
	return list, err
}

// List 分页查询药品，支持关键字/状态/机制过滤。
func (r *DrugRepository) List(page, pageSize int, keyword, status, mechanism string) ([]model.Drug, int64, error) {
	var list []model.Drug
	var total int64
	q := r.db.Model(&model.Drug{})
	if keyword != "" {
		q = q.Where("name LIKE ? OR generic_name LIKE ? OR atc_code LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if mechanism != "" {
		q = q.Where("mechanism = ?", mechanism)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

// Update 更新药品。
func (r *DrugRepository) Update(d *model.Drug) error {
	return r.db.Save(d).Error
}

// UpdateStatusTx 在事务中更新药品状态（禁用/启用）。
func (r *DrugRepository) UpdateStatusTx(tx *gorm.DB, id uint, status string) error {
	res := tx.Model(&model.Drug{}).Where("id = ?", id).Update("status", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// Count 统计药品数量。
func (r *DrugRepository) Count() (int64, error) {
	var n int64
	err := r.db.Model(&model.Drug{}).Count(&n).Error
	return n, err
}

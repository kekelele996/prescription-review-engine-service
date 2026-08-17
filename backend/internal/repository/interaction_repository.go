package repository

import (
	"errors"
	"fmt"

	"github.com/rxcheck/rxcheck/internal/model"
	"gorm.io/gorm"
)

// InteractionRepository 相互作用规则仓储。
type InteractionRepository struct {
	db *gorm.DB
}

func NewInteractionRepository(db *gorm.DB) *InteractionRepository {
	return &InteractionRepository{db: db}
}

// DB 返回底层数据库句柄。
func (r *InteractionRepository) DB() *gorm.DB { return r.db }

// Create 创建规则。
func (r *InteractionRepository) Create(rule *model.InteractionRule) error {
	return r.db.Create(rule).Error
}

// FindByID 按 ID 查询规则。
func (r *InteractionRepository) FindByID(id uint) (*model.InteractionRule, error) {
	var rule model.InteractionRule
	err := r.db.Preload("DrugA").Preload("DrugB").First(&rule, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &rule, err
}

// FindByPair 按药品对查询规则。
func (r *InteractionRepository) FindByPair(a, b uint) (*model.InteractionRule, error) {
	var rule model.InteractionRule
	err := r.db.Where(
		"((drug_a_id = ? AND drug_b_id = ?) OR (drug_a_id = ? AND drug_b_id = ?))",
		a, b, b, a,
	).First(&rule).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &rule, err
}

// FindEnabledByPair 按药品对查询启用中的规则（审核引擎复用）。
func (r *InteractionRepository) FindEnabledByPair(a, b uint) (*model.InteractionRule, error) {
	var rule model.InteractionRule
	err := r.db.Preload("DrugA").Preload("DrugB").Where(
		"((drug_a_id = ? AND drug_b_id = ?) OR (drug_a_id = ? AND drug_b_id = ?)) AND status = ?",
		a, b, b, a, "enabled",
	).First(&rule).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &rule, err
}

// List 分页查询规则，支持风险等级过滤。
func (r *InteractionRepository) List(page, pageSize int, risk string) ([]model.InteractionRule, int64, error) {
	var list []model.InteractionRule
	var total int64
	q := r.db.Model(&model.InteractionRule{}).Preload("DrugA").Preload("DrugB")
	if risk != "" {
		q = q.Where("risk_level = ?", risk)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

// Delete 删除规则。
func (r *InteractionRepository) Delete(id uint) error {
	res := r.db.Delete(&model.InteractionRule{}, id)
	if res.Error != nil {
		return fmt.Errorf("delete interaction rule: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// Count 统计启用规则数量。
func (r *InteractionRepository) Count() (int64, error) {
	var n int64
	err := r.db.Model(&model.InteractionRule{}).Where("status = ?", "enabled").Count(&n).Error
	return n, err
}

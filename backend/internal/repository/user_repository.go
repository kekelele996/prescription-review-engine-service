package repository

import (
	"errors"

	"github.com/rxcheck/rxcheck/internal/model"
	"gorm.io/gorm"
)

// UserRepository 用户仓储。
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// DB 返回底层数据库句柄（供 service 层开启事务）。
func (r *UserRepository) DB() *gorm.DB { return r.db }

// Create 创建用户。
func (r *UserRepository) Create(u *model.User) error {
	return r.db.Create(u).Error
}

// FindByID 按 ID 查询用户。
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var u model.User
	err := r.db.First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &u, err
}

// FindByUsername 按用户名查询用户。
func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var u model.User
	err := r.db.Where("username = ?", username).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &u, err
}

// List 分页查询用户，支持关键字检索。
func (r *UserRepository) List(page, pageSize int, keyword string) ([]model.User, int64, error) {
	var list []model.User
	var total int64
	q := r.db.Model(&model.User{})
	if keyword != "" {
		q = q.Where("username LIKE ? OR name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

// Update 更新用户。
func (r *UserRepository) Update(u *model.User) error {
	return r.db.Save(u).Error
}

// UpdateLastLogin 更新最后登录时间。
func (r *UserRepository) UpdateLastLogin(id uint, t interface{}) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("last_login_at", t).Error
}

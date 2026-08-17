package model

import "time"

// User 医生/药师/管理员账号，角色与状态决定可执行的审核动作（RBAC）。
type User struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Username    string     `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password    string     `gorm:"size:255;not null" json:"-"`
	Name        string     `gorm:"size:64;not null" json:"name"`
	Role        string     `gorm:"size:32;not null;default:doctor;index" json:"role"`
	Status      string     `gorm:"size:32;not null;default:active;index" json:"status"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

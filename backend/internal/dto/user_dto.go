package dto

import (
	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/model"
)

// RegisterReq 自助注册请求（默认角色 doctor）。
type RegisterReq struct {
	Username string `json:"username" binding:"required,min=2,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
	Name     string `json:"name" binding:"required,min=2,max=64"`
}

// CreateUserReq 管理员创建用户请求。
type CreateUserReq struct {
	Username string `json:"username" binding:"required,min=2,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
	Name     string `json:"name" binding:"required,min=2,max=64"`
	Role     string `json:"role" binding:"required,oneof=admin doctor pharmacist"`
}

// UpdateUserReq 更新用户请求。
type UpdateUserReq struct {
	Name   string `json:"name" binding:"omitempty,min=2,max=64"`
	Status string `json:"status" binding:"omitempty,oneof=active disabled"`
	Role   string `json:"role" binding:"omitempty,oneof=admin doctor pharmacist"`
}

// UserResp 用户响应。
type UserResp struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

// ToUserResp 模型转响应。
func ToUserResp(u *model.User) UserResp {
	return UserResp{
		ID:        u.ID,
		Username:  u.Username,
		Name:      u.Name,
		Role:      u.Role,
		Status:    u.Status,
		CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ValidateRole 校验角色合法性（handler 与 service 共用）。
func ValidateRole(role string) bool {
	return constants.IsValidUserRole(role)
}

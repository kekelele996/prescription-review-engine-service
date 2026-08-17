package model

import "time"

// AuditLog 操作审计日志：记录写操作调用方、模块、实体与耗时。
type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Username  string    `gorm:"size:64;index" json:"username"`
	Action    string    `gorm:"size:32" json:"action"` // CREATE/UPDATE/DELETE/OVERRIDE/...
	Module    string    `gorm:"size:128;index" json:"module"`
	EntityID  string    `gorm:"size:64" json:"entity_id"`
	Detail    string    `gorm:"size:1024" json:"detail"`
	IP        string    `gorm:"size:64" json:"ip"`
	RequestID string    `gorm:"size:64;index" json:"request_id"`
	CreatedAt time.Time `json:"created_at"`
}

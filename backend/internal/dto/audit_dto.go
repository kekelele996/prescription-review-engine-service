package dto

import "github.com/rxcheck/rxcheck/internal/model"

// AuditLogResp 审计日志响应。
type AuditLogResp struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Action    string `json:"action"`
	Module    string `json:"module"`
	EntityID  string `json:"entity_id"`
	Detail    string `json:"detail"`
	IP        string `json:"ip"`
	RequestID string `json:"request_id"`
	CreatedAt string `json:"created_at"`
}

// ToAuditLogResp 审计日志模型转响应。
func ToAuditLogResp(l *model.AuditLog) AuditLogResp {
	return AuditLogResp{
		ID:        l.ID,
		UserID:    l.UserID,
		Username:  l.Username,
		Action:    l.Action,
		Module:    l.Module,
		EntityID:  l.EntityID,
		Detail:    l.Detail,
		IP:        l.IP,
		RequestID: l.RequestID,
		CreatedAt: l.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

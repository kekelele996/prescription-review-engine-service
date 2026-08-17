package dto

import (
	"github.com/rxcheck/rxcheck/internal/model"
)

// CreateInteractionReq 创建相互作用规则请求。
type CreateInteractionReq struct {
	DrugAID     uint   `json:"drug_a_id" binding:"required,min=1"`
	DrugBID     uint   `json:"drug_b_id" binding:"required,min=1"`
	RiskLevel   string `json:"risk_level" binding:"required,oneof=high medium low"`
	Mechanism   string `json:"mechanism" binding:"required,min=1,max=128"`
	Description string `json:"description" binding:"required,min=1,max=512"`
	Status      string `json:"status" binding:"omitempty,oneof=enabled disabled"`
}

// InteractionResp 相互作用规则响应。
type InteractionResp struct {
	model.InteractionRule
}

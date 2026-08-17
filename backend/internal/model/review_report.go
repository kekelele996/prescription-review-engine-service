package model

import "time"

// ReviewReport 结构化审核报告：整体结论 + 逐条审核意见。
type ReviewReport struct {
	ID             uint         `gorm:"primaryKey" json:"id"`
	PrescriptionID uint         `gorm:"uniqueIndex;not null" json:"prescription_id"`
	Prescription   *Prescription `gorm:"foreignKey:PrescriptionID" json:"prescription"`
	Status         string       `gorm:"size:32;not null;index" json:"status"`   // passed/warned/rejected
	RiskLevel      string       `gorm:"size:32;not null" json:"risk_level"`
	Summary        string       `gorm:"size:1024" json:"summary"`
	Items          []ReviewItem `gorm:"foreignKey:ReportID" json:"items"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

// ReviewItem 单条审核意见：规则来源、风险等级、建议操作。
type ReviewItem struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	ReportID    uint   `gorm:"not null;index" json:"report_id"`
	DrugName    string `gorm:"size:128" json:"drug_name"`
	RuleType    string `gorm:"size:32;not null;index" json:"rule_type"`
	RiskLevel   string `gorm:"size:32;not null" json:"risk_level"`
	Message     string `gorm:"size:1024" json:"message"`
	Suggestion  string `gorm:"size:1024" json:"suggestion"`
	RuleSource  string `gorm:"size:256" json:"rule_source"`
	Operation   string `gorm:"size:64" json:"operation"` // confirm/modify/deny
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

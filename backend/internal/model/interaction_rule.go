package model

import "time"

// InteractionRule 药物相互作用规则：两两药品的风险分级与机制说明。
type InteractionRule struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	DrugAID     uint      `gorm:"not null;index" json:"drug_a_id"`
	DrugBID     uint      `gorm:"not null;index" json:"drug_b_id"`
	DrugA       Drug      `gorm:"foreignKey:DrugAID" json:"drug_a"`
	DrugB       Drug      `gorm:"foreignKey:DrugBID" json:"drug_b"`
	RiskLevel   string    `gorm:"size:32;not null;index" json:"risk_level"` // high/medium/low
	Mechanism   string    `gorm:"size:128" json:"mechanism"`                // QT 间期延长 / 严重出血风险 / ...
	Description string    `gorm:"size:512" json:"description"`
	Status      string    `gorm:"size:32;not null;default:enabled;index" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

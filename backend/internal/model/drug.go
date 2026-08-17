package model

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// Indication 药品适应症（ICD-10 编码 + 名称）。
type Indication struct {
	ICD10 string `json:"icd10"` // ICD-10 编码
	Name  string `json:"name"`  // 适应症名称
}

// Drug 药品档案：包含适应症、剂量上限（成人/儿童/老人）、禁忌、频次等审核规则来源。
type Drug struct {
	ID                      uint            `gorm:"primaryKey" json:"id"`
	Name                    string          `gorm:"size:128;not null;uniqueIndex" json:"name"`
	GenericName             string          `gorm:"size:128;not null" json:"generic_name"`
	Specification           string          `gorm:"size:128" json:"specification"`
	DosageForm              string          `gorm:"size:64" json:"dosage_form"`
	Route                   string          `gorm:"size:64" json:"route"`
	AtcCode                 string          `gorm:"size:16;index" json:"atc_code"`
	Mechanism               string          `gorm:"size:128;index" json:"mechanism"` // 作用机制（重复用药检测依据）
	Indications             datatypes.JSON  `gorm:"type:jsonb" json:"indications"`   // []Indication
	AdultMaxDailyDose       float64         `json:"adult_max_daily_dose"`            // 成人每日最大剂量(mg)
	ChildMaxDailyDose       float64         `json:"child_max_daily_dose"`            // 儿童每日最大剂量(mg)
	ElderlyMaxDailyDose     float64         `json:"elderly_max_daily_dose"`          // 老人每日最大剂量(mg)
	MaxSingleDose           float64         `json:"max_single_dose"`
	DoseUnit                string          `gorm:"size:32" json:"dose_unit"`
	Allergies               datatypes.JSON  `gorm:"type:jsonb" json:"allergies"` // []string 过敏禁忌
	PregnancyContraindicated bool           `json:"pregnancy_contraindicated"`
	HepaticContraindicated  bool            `json:"hepatic_contraindicated"`
	RenalContraindicated    bool            `json:"renal_contraindicated"`
	BeersFlag               bool            `json:"beers_flag"`   // 老年人 Beers 标准高风险
	MinAgeMonths            int             `json:"min_age_months"` // 儿童年龄下限（月）
	FrequencyLimit          string          `gorm:"size:64" json:"frequency_limit"` // 允许频次, 如 qd,bid,tid
	Status                  string          `gorm:"size:32;not null;default:enabled;index" json:"status"`
	CreatedAt               time.Time       `json:"created_at"`
	UpdatedAt               time.Time       `json:"updated_at"`
}

// IndicationList 反序列化适应症列表。
func (d *Drug) IndicationList() []Indication {
	var list []Indication
	_ = json.Unmarshal(d.Indications, &list)
	return list
}

// AllergyList 反序列化过敏禁忌列表。
func (d *Drug) AllergyList() []string {
	var list []string
	_ = json.Unmarshal(d.Allergies, &list)
	return list
}

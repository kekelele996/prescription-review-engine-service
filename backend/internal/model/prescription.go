package model

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// Diagnosis 处方诊断（ICD-10 编码 + 名称）。
type Diagnosis struct {
	ICD10 string `json:"icd10"`
	Name  string `json:"name"`
}

// Prescription 处方：患者信息、诊断与药品明细，状态机由审核服务驱动。
type Prescription struct {
	ID                uint                `gorm:"primaryKey" json:"id"`
	PrescriptionNo    string              `gorm:"size:64;uniqueIndex;not null" json:"prescription_no"`
	PatientName       string              `gorm:"size:64;not null;index" json:"patient_name"`
	PatientAge        int                 `json:"patient_age"`
	PatientGender     string              `gorm:"size:16" json:"patient_gender"`
	WeightKg          float64             `json:"weight_kg"`
	Allergies         datatypes.JSON      `gorm:"type:jsonb" json:"allergies"` // []string
	Pregnant          bool                `json:"pregnant"`
	HepaticImpairment bool                `json:"hepatic_impairment"`
	RenalImpairment   bool                `json:"renal_impairment"`
	Diagnoses         datatypes.JSON      `gorm:"type:jsonb" json:"diagnoses"` // []Diagnosis
	Format            string              `gorm:"size:16" json:"format"`       // json/xml
	RawPayload        string              `gorm:"type:text" json:"raw_payload"` // 原始处方原文
	Status            string              `gorm:"size:32;not null;default:pending_review;index" json:"status"`
	OverallRisk       string              `gorm:"size:32;not null;default:none" json:"overall_risk"`
	DoctorID          *uint               `gorm:"index" json:"doctor_id"`
	Doctor            *User               `gorm:"foreignKey:DoctorID" json:"doctor"`
	OverrideBy        *uint               `json:"override_by"`
	OverrideReason    string              `gorm:"size:512" json:"override_reason"`
	ReviewedAt        *time.Time          `json:"reviewed_at"`
	Items             []PrescriptionItem  `gorm:"foreignKey:PrescriptionID" json:"items"`
	Report            *ReviewReport       `gorm:"foreignKey:PrescriptionID" json:"report"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
}

// DiagnosisList 反序列化诊断列表。
func (p *Prescription) DiagnosisList() []Diagnosis {
	var list []Diagnosis
	_ = json.Unmarshal(p.Diagnoses, &list)
	return list
}

// AllergyList 反序列化过敏史列表。
func (p *Prescription) AllergyList() []string {
	var list []string
	_ = json.Unmarshal(p.Allergies, &list)
	return list
}

// PrescriptionItem 处方药品明细：单次剂量、频次、疗程、给药途径与折算日剂量。
type PrescriptionItem struct {
	ID             uint    `gorm:"primaryKey" json:"id"`
	PrescriptionID uint    `gorm:"not null;index" json:"prescription_id"`
	DrugID         uint    `gorm:"not null;index" json:"drug_id"`
	Drug           Drug    `gorm:"foreignKey:DrugID" json:"drug"`
	DrugName       string  `gorm:"size:128" json:"drug_name"`
	Specification  string  `gorm:"size:128" json:"specification"`
	SingleDose     float64 `json:"single_dose"`
	Frequency      string  `gorm:"size:32" json:"frequency"` // qd/bid/tid/qid/q8h/q12h/qn/prn
	CourseDays     int     `json:"course_days"`
	Route          string  `gorm:"size:64" json:"route"`
	DailyDose      float64 `json:"daily_dose"` // 单次剂量 × 每日频次
	SuggestedDose  float64 `json:"suggested_dose"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

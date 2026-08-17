package dto

import (
	"encoding/json"
)

// DiagnosisDTO 处方诊断条目（JSON/XML 双格式）。
type DiagnosisDTO struct {
	ICD10 string `json:"icd10" xml:"icd10" binding:"required,min=1,max=16"`
	Name  string `json:"name" xml:"name" binding:"required,min=1,max=128"`
}

// PrescriptionItemDTO 处方药品明细（JSON/XML 双格式）。
type PrescriptionItemDTO struct {
	DrugID       uint    `json:"drug_id" xml:"drug_id"`
	DrugName     string  `json:"drug_name" xml:"drug_name" binding:"omitempty,min=1,max=128"`
	Specification string  `json:"specification" xml:"specification" binding:"omitempty,max=128"`
	SingleDose   float64 `json:"single_dose" xml:"single_dose" binding:"required,min=0.01"`
	Frequency    string  `json:"frequency" xml:"frequency" binding:"required,oneof=qd bid tid qid q8h q12h qn prn"`
	CourseDays   int     `json:"course_days" xml:"course_days" binding:"required,min=1,max=365"`
	Route        string  `json:"route" xml:"route" binding:"omitempty,max=64"`
}

// SubmitPrescriptionReq 处方接收请求：支持 JSON 与 XML 两种格式。
type SubmitPrescriptionReq struct {
	PrescriptionNo    string                `json:"prescription_no" xml:"prescription_no" binding:"required,min=1,max=64"`
	PatientName       string                `json:"patient_name" xml:"patient_name" binding:"required,min=1,max=64"`
	PatientAge        int                   `json:"patient_age" xml:"patient_age" binding:"required,min=0,max=130"`
	PatientGender     string                `json:"patient_gender" xml:"patient_gender" binding:"omitempty,oneof=male female"`
	WeightKg          float64               `json:"weight_kg" xml:"weight_kg" binding:"omitempty,min=0"`
	Allergies         []string              `json:"allergies" xml:"allergies>item" binding:"omitempty,dive,max=64"`
	Pregnant          bool                  `json:"pregnant" xml:"pregnant"`
	HepaticImpairment bool                  `json:"hepatic_impairment" xml:"hepatic_impairment"`
	RenalImpairment   bool                  `json:"renal_impairment" xml:"renal_impairment"`
	Diagnoses         []DiagnosisDTO        `json:"diagnoses" xml:"diagnoses>diagnosis" binding:"required,min=1,dive"`
	Items             []PrescriptionItemDTO `json:"items" xml:"items>item" binding:"required,min=1,dive"`
}

// OverrideReq 强制通过请求（仅警告/通过类处方）。
type OverrideReq struct {
	Reason string `json:"reason" binding:"required,min=2,max=512"`
}

// PrescriptionListResp 处方列表项。
type PrescriptionListResp struct {
	ID              uint    `json:"id"`
	PrescriptionNo  string  `json:"prescription_no"`
	PatientName     string  `json:"patient_name"`
	PatientAge      int     `json:"patient_age"`
	Format          string  `json:"format"`
	Status          string  `json:"status"`
	OverallRisk     string  `json:"overall_risk"`
	ItemCount       int     `json:"item_count"`
	ReviewItemCount int     `json:"review_item_count"`
	CreatedAt       string  `json:"created_at"`
	ReviewedAt      string  `json:"reviewed_at"`
}

// marshalJSON 辅助序列化（供 DTO 转换使用）。
func marshalJSON(v any) ([]byte, error) {
	return json.Marshal(v)
}

package dto

import (
	"github.com/rxcheck/rxcheck/internal/model"
	"gorm.io/datatypes"
)

// IndicationDTO 适应症条目。
type IndicationDTO struct {
	ICD10 string `json:"icd10" binding:"required,min=1,max=16"`
	Name  string `json:"name" binding:"required,min=1,max=128"`
}

// CreateDrugReq 创建药品请求。
type CreateDrugReq struct {
	Name                    string         `json:"name" binding:"required,min=1,max=128"`
	GenericName             string         `json:"generic_name" binding:"required,min=1,max=128"`
	Specification           string         `json:"specification" binding:"omitempty,max=128"`
	DosageForm              string         `json:"dosage_form" binding:"omitempty,max=64"`
	Route                   string         `json:"route" binding:"omitempty,max=64"`
	AtcCode                 string         `json:"atc_code" binding:"omitempty,max=16"`
	Mechanism               string         `json:"mechanism" binding:"omitempty,max=128"`
	Indications             []IndicationDTO `json:"indications" binding:"omitempty,dive"`
	AdultMaxDailyDose       float64        `json:"adult_max_daily_dose" binding:"omitempty,min=0"`
	ChildMaxDailyDose       float64        `json:"child_max_daily_dose" binding:"omitempty,min=0"`
	ElderlyMaxDailyDose     float64        `json:"elderly_max_daily_dose" binding:"omitempty,min=0"`
	MaxSingleDose           float64        `json:"max_single_dose" binding:"omitempty,min=0"`
	DoseUnit                string         `json:"dose_unit" binding:"omitempty,max=32"`
	Allergies               []string       `json:"allergies" binding:"omitempty,dive,max=64"`
	PregnancyContraindicated bool          `json:"pregnancy_contraindicated"`
	HepaticContraindicated  bool           `json:"hepatic_contraindicated"`
	RenalContraindicated    bool           `json:"renal_contraindicated"`
	BeersFlag               bool           `json:"beers_flag"`
	MinAgeMonths            int            `json:"min_age_months" binding:"omitempty,min=0"`
	FrequencyLimit          string         `json:"frequency_limit" binding:"omitempty,max=64"`
	Status                  string         `json:"status" binding:"omitempty,oneof=enabled disabled"`
}

// UpdateDrugReq 更新药品请求（字段均为可选）。
type UpdateDrugReq struct {
	Name                    string          `json:"name" binding:"omitempty,min=1,max=128"`
	GenericName             string          `json:"generic_name" binding:"omitempty,min=1,max=128"`
	Specification           string          `json:"specification" binding:"omitempty,max=128"`
	DosageForm              string          `json:"dosage_form" binding:"omitempty,max=64"`
	Route                   string          `json:"route" binding:"omitempty,max=64"`
	AtcCode                 string          `json:"atc_code" binding:"omitempty,max=16"`
	Mechanism               string          `json:"mechanism" binding:"omitempty,max=128"`
	Indications             []IndicationDTO `json:"indications" binding:"omitempty,dive"`
	AdultMaxDailyDose       *float64        `json:"adult_max_daily_dose" binding:"omitempty,min=0"`
	ChildMaxDailyDose       *float64        `json:"child_max_daily_dose" binding:"omitempty,min=0"`
	ElderlyMaxDailyDose     *float64        `json:"elderly_max_daily_dose" binding:"omitempty,min=0"`
	MaxSingleDose           *float64        `json:"max_single_dose" binding:"omitempty,min=0"`
	DoseUnit                string          `json:"dose_unit" binding:"omitempty,max=32"`
	Allergies               []string        `json:"allergies" binding:"omitempty,dive,max=64"`
	PregnancyContraindicated *bool          `json:"pregnancy_contraindicated"`
	HepaticContraindicated  *bool           `json:"hepatic_contraindicated"`
	RenalContraindicated    *bool           `json:"renal_contraindicated"`
	BeersFlag               *bool           `json:"beers_flag"`
	MinAgeMonths            *int            `json:"min_age_months" binding:"omitempty,min=0"`
	FrequencyLimit          string          `json:"frequency_limit" binding:"omitempty,max=64"`
	Status                  string          `json:"status" binding:"omitempty,oneof=enabled disabled"`
}

// DrugResp 药品响应。
type DrugResp struct {
	model.Drug
}

// IndicationsToJSON 适应症 DTO 列表转 JSON。
func IndicationsToJSON(list []IndicationDTO) datatypes.JSON {
	b, _ := datatypes.JSON(nil).MarshalJSON()
	_ = b
	data, err := marshalJSON(list)
	if err != nil {
		return datatypes.JSON("[]")
	}
	return datatypes.JSON(data)
}

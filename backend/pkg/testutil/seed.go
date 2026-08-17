// Package testutil 提供测试共享的种子数据构造工具。
package testutil

import (
	"encoding/json"

	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/model"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// SeedDrugs 初始化药品档案，返回 药品名 -> ID 映射。
func SeedDrugs(db *gorm.DB) (map[string]uint, error) {
	drugs := []model.Drug{
		{Name: "氨氯地平", GenericName: "苯磺酸氨氯地平片", Specification: "5mg", DosageForm: "片剂", Route: "口服", AtcCode: "C08CA01", Mechanism: "calcium_channel_blocker",
			Indications: Indications([2]string{"I10", "原发性高血压"}, [2]string{"I20", "心绞痛"}),
			AdultMaxDailyDose: 10, ChildMaxDailyDose: 5, ElderlyMaxDailyDose: 5, MaxSingleDose: 10,
			DoseUnit: "mg", FrequencyLimit: "qd,bid", Status: constants.DrugStatusEnabled},
		{Name: "赖诺普利", GenericName: "赖诺普利片", Specification: "10mg", DosageForm: "片剂", Route: "口服", AtcCode: "C09AA03", Mechanism: "ace_inhibitor",
			Indications: Indications([2]string{"I10", "原发性高血压"}, [2]string{"I50", "心力衰竭"}),
			AdultMaxDailyDose: 40, ChildMaxDailyDose: 20, ElderlyMaxDailyDose: 20, MaxSingleDose: 40,
			DoseUnit: "mg", FrequencyLimit: "qd,bid", PregnancyContraindicated: true, Status: constants.DrugStatusEnabled},
		{Name: "依那普利", GenericName: "马来酸依那普利片", Specification: "5mg", DosageForm: "片剂", Route: "口服", AtcCode: "C09AA02", Mechanism: "ace_inhibitor",
			Indications: Indications([2]string{"I10", "原发性高血压"}, [2]string{"I50", "心力衰竭"}),
			AdultMaxDailyDose: 40, ChildMaxDailyDose: 20, ElderlyMaxDailyDose: 20, MaxSingleDose: 20,
			DoseUnit: "mg", FrequencyLimit: "qd,bid", PregnancyContraindicated: true, Status: constants.DrugStatusEnabled},
		{Name: "华法林", GenericName: "华法林钠片", Specification: "2.5mg", DosageForm: "片剂", Route: "口服", AtcCode: "B01AA03", Mechanism: "vitamin_k_antagonist",
			Indications: Indications([2]string{"I48", "心房颤动"}, [2]string{"I26", "肺栓塞"}),
			AdultMaxDailyDose: 10, ChildMaxDailyDose: 5, ElderlyMaxDailyDose: 7.5, MaxSingleDose: 10,
			DoseUnit: "mg", FrequencyLimit: "qd", Status: constants.DrugStatusEnabled},
		{Name: "阿司匹林", GenericName: "阿司匹林肠溶片", Specification: "100mg", DosageForm: "肠溶片", Route: "口服", AtcCode: "B01AC06", Mechanism: "antiplatelet",
			Indications: Indications([2]string{"I20", "心绞痛"}, [2]string{"I21", "急性心肌梗死"}, [2]string{"I63", "缺血性卒中"}),
			AdultMaxDailyDose: 300, ChildMaxDailyDose: 0, ElderlyMaxDailyDose: 300, MaxSingleDose: 300,
			DoseUnit: "mg", Allergies: datatypes.JSON(`["阿司匹林"]`), MinAgeMonths: 192, FrequencyLimit: "qd", Status: constants.DrugStatusEnabled},
		{Name: "二甲双胍", GenericName: "盐酸二甲双胍片", Specification: "500mg", DosageForm: "片剂", Route: "口服", AtcCode: "A10BA02", Mechanism: "biguanide",
			Indications: Indications([2]string{"E11", "2型糖尿病"}),
			AdultMaxDailyDose: 2550, ChildMaxDailyDose: 2000, ElderlyMaxDailyDose: 2000, MaxSingleDose: 1000,
			DoseUnit: "mg", FrequencyLimit: "bid,tid", RenalContraindicated: true, Status: constants.DrugStatusEnabled},
		{Name: "头孢呋辛", GenericName: "头孢呋辛酯片", Specification: "250mg", DosageForm: "片剂", Route: "口服", AtcCode: "J01DC02", Mechanism: "cephalosporin",
			Indications: Indications([2]string{"J18", "肺炎"}, [2]string{"N39", "泌尿道感染"}),
			AdultMaxDailyDose: 1000, ChildMaxDailyDose: 500, ElderlyMaxDailyDose: 1000, MaxSingleDose: 500,
			DoseUnit: "mg", Allergies: datatypes.JSON(`["青霉素"]`), FrequencyLimit: "bid", Status: constants.DrugStatusEnabled},
		{Name: "地西泮", GenericName: "地西泮片", Specification: "5mg", DosageForm: "片剂", Route: "口服", AtcCode: "N05BA01", Mechanism: "benzodiazepine",
			Indications: Indications([2]string{"F41", "焦虑障碍"}, [2]string{"G47", "失眠"}),
			AdultMaxDailyDose: 30, ChildMaxDailyDose: 0, ElderlyMaxDailyDose: 10, MaxSingleDose: 10,
			DoseUnit: "mg", FrequencyLimit: "qd,qn", BeersFlag: true, Status: constants.DrugStatusEnabled},
	}
	if err := db.Create(&drugs).Error; err != nil {
		return nil, err
	}
	ids := make(map[string]uint, len(drugs))
	for _, d := range drugs {
		ids[d.Name] = d.ID
	}
	return ids, nil
}

// SeedInteractions 初始化相互作用规则。
func SeedInteractions(db *gorm.DB, ids map[string]uint) error {
	rules := []model.InteractionRule{
		{DrugAID: ids["华法林"], DrugBID: ids["阿司匹林"], RiskLevel: constants.RiskHigh, Mechanism: "严重出血风险",
			Description: "华法林与阿司匹林联用显著增加出血风险", Status: constants.RuleStatusEnabled},
		{DrugAID: ids["赖诺普利"], DrugBID: ids["阿司匹林"], RiskLevel: constants.RiskMedium, Mechanism: "肾功能/血压",
			Description: "ACEI 与抗血小板药联用增加出血风险", Status: constants.RuleStatusEnabled},
		{DrugAID: ids["华法林"], DrugBID: ids["头孢呋辛"], RiskLevel: constants.RiskMedium, Mechanism: "抗凝增强",
			Description: "头孢类可能增强华法林抗凝作用", Status: constants.RuleStatusEnabled},
		{DrugAID: ids["地西泮"], DrugBID: ids["氨氯地平"], RiskLevel: constants.RiskLow, Mechanism: "镇静/降压叠加",
			Description: "可能加重镇静与低血压", Status: constants.RuleStatusEnabled},
		{DrugAID: ids["华法林"], DrugBID: ids["二甲双胍"], RiskLevel: constants.RiskLow, Mechanism: "相互作用轻微",
			Description: "一般无需调整", Status: constants.RuleStatusEnabled},
	}
	return db.Create(&rules).Error
}

// Indications 构造适应症 JSON。
func Indications(items ...[2]string) datatypes.JSON {
	list := make([]model.Indication, 0, len(items))
	for _, it := range items {
		list = append(list, model.Indication{ICD10: it[0], Name: it[1]})
	}
	b, err := json.Marshal(list)
	if err != nil {
		return datatypes.JSON("[]")
	}
	return datatypes.JSON(b)
}

// StringsJSON 构造字符串数组 JSON。
func StringsJSON(list []string) datatypes.JSON {
	b, err := json.Marshal(list)
	if err != nil {
		return datatypes.JSON("[]")
	}
	return datatypes.JSON(b)
}

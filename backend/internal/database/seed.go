package database

import (
	"fmt"

	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/model"
	"github.com/rxcheck/rxcheck/internal/util"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Seed 初始化种子数据：管理员账号、药品档案、相互作用规则（幂等）。
func Seed(db *gorm.DB, adminUsername, adminPassword string) error {
	if err := seedAdmin(db, adminUsername, adminPassword); err != nil {
		return err
	}
	drugNames, err := seedDrugs(db)
	if err != nil {
		return err
	}
	interCount, err := seedInteractions(db, drugNames)
	if err != nil {
		return err
	}
	util.Log.Info(fmt.Sprintf(constants.LogSeedCompleted, adminUsername, len(drugNames), interCount))
	return nil
}

func seedAdmin(db *gorm.DB, username, password string) error {
	var n int64
	if err := db.Model(&model.User{}).Where("username = ?", username).Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	hash, err := util.HashPassword(password)
	if err != nil {
		return err
	}
	admin := &model.User{
		Username: username,
		Password: hash,
		Name:     "系统管理员",
		Role:     constants.UserRoleAdmin,
		Status:   constants.UserStatusActive,
	}
	return db.Create(admin).Error
}

// seedDrugs 初始化药品档案，返回药品名列表（供相互作用规则引用）。
func seedDrugs(db *gorm.DB) ([]string, error) {
	var n int64
	if err := db.Model(&model.Drug{}).Count(&n).Error; err != nil {
		return nil, err
	}
	if n > 0 {
		names := make([]string, 0)
		var drugs []model.Drug
		if err := db.Select("name").Find(&drugs).Error; err != nil {
			return nil, err
		}
		for _, d := range drugs {
			names = append(names, d.Name)
		}
		return names, nil
	}
	drugs := []model.Drug{
		{
			Name: "氨氯地平", GenericName: "苯磺酸氨氯地平片", Specification: "5mg", DosageForm: "片剂",
			Route: "口服", AtcCode: "C08CA01", Mechanism: "calcium_channel_blocker",
			Indications:       jsonIndications([][2]string{{"I10", "原发性高血压"}, {"I20", "心绞痛"}}),
			AdultMaxDailyDose: 10, ChildMaxDailyDose: 5, ElderlyMaxDailyDose: 5, MaxSingleDose: 10,
			DoseUnit: "mg", FrequencyLimit: "qd,bid", Status: constants.DrugStatusEnabled,
		},
		{
			Name: "赖诺普利", GenericName: "赖诺普利片", Specification: "10mg", DosageForm: "片剂",
			Route: "口服", AtcCode: "C09AA03", Mechanism: "ace_inhibitor",
			Indications:       jsonIndications([][2]string{{"I10", "原发性高血压"}, {"I50", "心力衰竭"}}),
			AdultMaxDailyDose: 40, ChildMaxDailyDose: 20, ElderlyMaxDailyDose: 20, MaxSingleDose: 40,
			DoseUnit: "mg", FrequencyLimit: "qd,bid", PregnancyContraindicated: true, Status: constants.DrugStatusEnabled,
		},
		{
			Name: "依那普利", GenericName: "马来酸依那普利片", Specification: "5mg", DosageForm: "片剂",
			Route: "口服", AtcCode: "C09AA02", Mechanism: "ace_inhibitor",
			Indications:       jsonIndications([][2]string{{"I10", "原发性高血压"}, {"I50", "心力衰竭"}}),
			AdultMaxDailyDose: 40, ChildMaxDailyDose: 20, ElderlyMaxDailyDose: 20, MaxSingleDose: 20,
			DoseUnit: "mg", FrequencyLimit: "qd,bid", PregnancyContraindicated: true, Status: constants.DrugStatusEnabled,
		},
		{
			Name: "华法林", GenericName: "华法林钠片", Specification: "2.5mg", DosageForm: "片剂",
			Route: "口服", AtcCode: "B01AA03", Mechanism: "vitamin_k_antagonist",
			Indications:       jsonIndications([][2]string{{"I48", "心房颤动"}, {"I26", "肺栓塞"}}),
			AdultMaxDailyDose: 10, ChildMaxDailyDose: 5, ElderlyMaxDailyDose: 7.5, MaxSingleDose: 10,
			DoseUnit: "mg", FrequencyLimit: "qd", Status: constants.DrugStatusEnabled,
		},
		{
			Name: "阿司匹林", GenericName: "阿司匹林肠溶片", Specification: "100mg", DosageForm: "肠溶片",
			Route: "口服", AtcCode: "B01AC06", Mechanism: "antiplatelet",
			Indications:       jsonIndications([][2]string{{"I20", "心绞痛"}, {"I21", "急性心肌梗死"}, {"I63", "缺血性卒中"}}),
			AdultMaxDailyDose: 300, ChildMaxDailyDose: 0, ElderlyMaxDailyDose: 300, MaxSingleDose: 300,
			DoseUnit: "mg", Allergies: datatypes.JSON(`["阿司匹林"]`), MinAgeMonths: 192,
			FrequencyLimit: "qd", Status: constants.DrugStatusEnabled,
		},
		{
			Name: "二甲双胍", GenericName: "盐酸二甲双胍片", Specification: "500mg", DosageForm: "片剂",
			Route: "口服", AtcCode: "A10BA02", Mechanism: "biguanide",
			Indications:       jsonIndications([][2]string{{"E11", "2型糖尿病"}}),
			AdultMaxDailyDose: 2550, ChildMaxDailyDose: 2000, ElderlyMaxDailyDose: 2000, MaxSingleDose: 1000,
			DoseUnit: "mg", FrequencyLimit: "bid,tid", RenalContraindicated: true, Status: constants.DrugStatusEnabled,
		},
		{
			Name: "头孢呋辛", GenericName: "头孢呋辛酯片", Specification: "250mg", DosageForm: "片剂",
			Route: "口服", AtcCode: "J01DC02", Mechanism: "cephalosporin",
			Indications:       jsonIndications([][2]string{{"J18", "肺炎"}, {"N39", "泌尿道感染"}}),
			AdultMaxDailyDose: 1000, ChildMaxDailyDose: 500, ElderlyMaxDailyDose: 1000, MaxSingleDose: 500,
			DoseUnit: "mg", Allergies: datatypes.JSON(`["青霉素"]`), FrequencyLimit: "bid",
			Status: constants.DrugStatusEnabled,
		},
		{
			Name: "地西泮", GenericName: "地西泮片", Specification: "5mg", DosageForm: "片剂",
			Route: "口服", AtcCode: "N05BA01", Mechanism: "benzodiazepine",
			Indications:       jsonIndications([][2]string{{"F41", "焦虑障碍"}, {"G47", "失眠"}}),
			AdultMaxDailyDose: 30, ChildMaxDailyDose: 0, ElderlyMaxDailyDose: 10, MaxSingleDose: 10,
			DoseUnit: "mg", FrequencyLimit: "qd,qn", BeersFlag: true, Status: constants.DrugStatusEnabled,
		},
	}
	if err := db.Create(&drugs).Error; err != nil {
		return nil, err
	}
	names := make([]string, 0, len(drugs))
	for _, d := range drugs {
		names = append(names, d.Name)
	}
	return names, nil
}

// seedInteractions 初始化相互作用规则，返回创建数量。
func seedInteractions(db *gorm.DB, drugNames []string) (int, error) {
	var n int64
	if err := db.Model(&model.InteractionRule{}).Count(&n).Error; err != nil {
		return 0, err
	}
	if n > 0 {
		return int(n), nil
	}
	id := func(name string) uint {
		var d model.Drug
		if err := db.Where("name = ?", name).First(&d).Error; err != nil {
			return 0
		}
		return d.ID
	}
	rules := []model.InteractionRule{
		{DrugAID: id("华法林"), DrugBID: id("阿司匹林"), RiskLevel: constants.RiskHigh, Mechanism: "严重出血风险",
			Description: "华法林与阿司匹林联用显著增加消化道及颅内出血风险，需加强 INR 监测", Status: constants.RuleStatusEnabled},
		{DrugAID: id("赖诺普利"), DrugBID: id("阿司匹林"), RiskLevel: constants.RiskMedium, Mechanism: "肾功能/血压",
			Description: "ACEI 与抗血小板药联用增加肾功能损伤与出血风险", Status: constants.RuleStatusEnabled},
		{DrugAID: id("华法林"), DrugBID: id("头孢呋辛"), RiskLevel: constants.RiskMedium, Mechanism: "抗凝增强",
			Description: "头孢类抗菌药可能增强华法林抗凝作用，增加出血风险", Status: constants.RuleStatusEnabled},
		{DrugAID: id("地西泮"), DrugBID: id("氨氯地平"), RiskLevel: constants.RiskLow, Mechanism: "镇静/降压叠加",
			Description: "苯二氮䓬类与钙通道阻滞剂联用可能加重镇静与低血压", Status: constants.RuleStatusEnabled},
		{DrugAID: id("华法林"), DrugBID: id("二甲双胍"), RiskLevel: constants.RiskLow, Mechanism: "相互作用轻微",
			Description: "两药相互作用轻微，一般无需调整剂量", Status: constants.RuleStatusEnabled},
	}
	for i := range rules {
		if rules[i].DrugAID == 0 || rules[i].DrugBID == 0 {
			return 0, fmt.Errorf("种子相互作用规则引用的药品不存在: %+v", rules[i])
		}
	}
	if err := db.Create(&rules).Error; err != nil {
		return 0, err
	}
	return len(rules), nil
}

// jsonIndications 构造适应症 JSON。
func jsonIndications(items [][2]string) datatypes.JSON {
	list := make([]model.Indication, 0, len(items))
	for _, it := range items {
		list = append(list, model.Indication{ICD10: it[0], Name: it[1]})
	}
	b, err := jsonMarshal(list)
	if err != nil {
		return datatypes.JSON("[]")
	}
	return datatypes.JSON(b)
}

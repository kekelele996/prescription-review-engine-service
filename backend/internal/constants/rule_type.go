package constants

// 规则类型枚举：审核项按规则类型分类，贯穿模型、DTO、审核引擎、日志模板、formatters。
const (
	// RuleTypeIndication 适应症审核。
	RuleTypeIndication = "indication"
	// RuleTypeDosage 用法用量审核。
	RuleTypeDosage = "dosage"
	// RuleTypeInteraction 药物相互作用审核。
	RuleTypeInteraction = "interaction"
	// RuleTypeContraindication 禁忌与特殊人群审核。
	RuleTypeContraindication = "contraindication"
	// RuleTypeDuplication 重复用药审核。
	RuleTypeDuplication = "duplication"
)

// AllRuleTypes 全部合法规则类型。
var AllRuleTypes = []string{
	RuleTypeIndication, RuleTypeDosage, RuleTypeInteraction, RuleTypeContraindication, RuleTypeDuplication,
}

// IsValidRuleType 判断规则类型是否合法。
func IsValidRuleType(t string) bool {
	for _, r := range AllRuleTypes {
		if r == t {
			return true
		}
	}
	return false
}

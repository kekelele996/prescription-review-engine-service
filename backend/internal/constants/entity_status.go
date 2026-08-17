package constants

// 实体启停状态枚举（用户 active/disabled、药品与规则 enabled/disabled）。
const (
	// UserStatusActive 用户正常。
	UserStatusActive = "active"
	// UserStatusDisabled 用户禁用。
	UserStatusDisabled = "disabled"

	// DrugStatusEnabled 药品可用。
	DrugStatusEnabled = "enabled"
	// DrugStatusDisabled 药品停用（不再参与新处方审核）。
	DrugStatusDisabled = "disabled"

	// RuleStatusEnabled 规则启用。
	RuleStatusEnabled = "enabled"
	// RuleStatusDisabled 规则停用。
	RuleStatusDisabled = "disabled"
)

// AllUserStatuses 用户状态枚举。
var AllUserStatuses = []string{UserStatusActive, UserStatusDisabled}

// AllDrugStatuses 药品状态枚举。
var AllDrugStatuses = []string{DrugStatusEnabled, DrugStatusDisabled}

// AllRuleStatuses 规则状态枚举。
var AllRuleStatuses = []string{RuleStatusEnabled, RuleStatusDisabled}

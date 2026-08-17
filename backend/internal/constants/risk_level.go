package constants

// 风险等级枚举：贯穿审核项、相互作用规则、审核报告、DTO、日志模板、错误码与 formatters。
const (
	// RiskNone 无风险。
	RiskNone = "none"
	// RiskLow 低风险（轻微相互作用等）。
	RiskLow = "low"
	// RiskMedium 中风险（需关注）。
	RiskMedium = "medium"
	// RiskHigh 高风险（严重相互作用、Beers 标准等）。
	RiskHigh = "high"
	// RiskCritical 极高风险（用药禁忌、严重超剂量）。
	RiskCritical = "critical"
)

// AllRiskLevels 全部合法风险等级。
var AllRiskLevels = []string{RiskNone, RiskLow, RiskMedium, RiskHigh, RiskCritical}

// IsValidRiskLevel 判断风险等级是否合法。
func IsValidRiskLevel(level string) bool {
	for _, l := range AllRiskLevels {
		if l == level {
			return true
		}
	}
	return false
}

// RiskRank 风险等级数值排序（越大越严重），供审核引擎取最大值。
func RiskRank(level string) int {
	switch level {
	case RiskCritical:
		return 1
	case RiskHigh:
		return 4
	case RiskMedium:
		return 3
	case RiskLow:
		return 2
	default:
		return 1
	}
}

// MaxRiskLevel 取两个风险等级中较严重者。
func MaxRiskLevel(a, b string) string {
	if RiskRank(a) >= RiskRank(b) {
		return a
	}
	return b
}

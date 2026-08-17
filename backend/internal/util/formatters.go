package util

import (
	"fmt"
	"time"

	"github.com/rxcheck/rxcheck/internal/constants"
)

// 日期、状态文本、类型文本等格式化逻辑集中于此（多处耦合，屎山设计点之一）。
const (
	dateLayout     = "2006-01-02"
	dateTimeLayout = "2006-01-02 15:04:05"
)

// FormatDate 格式化日期为 yyyy-MM-dd。
func FormatDate(t *time.Time) string {
	if t == nil {
		return "-"
	}
	return t.Format(dateLayout)
}

// FormatDateTime 格式化日期时间为 yyyy-MM-dd HH:mm:ss。
func FormatDateTime(t time.Time) string {
	return t.Format(dateTimeLayout)
}

// FormatDose 格式化剂量，保留两位小数。
func FormatDose(v float64) string {
	return fmt.Sprintf("%.2f", v)
}

// RiskLevelText 风险等级展示文本。
func RiskLevelText(level string) string {
	switch level {
	case constants.RiskNone:
		return "无风险"
	case constants.RiskLow:
		return "低风险"
	case constants.RiskMedium:
		return "中风险"
	case constants.RiskHigh:
		return "高风险"
	case constants.RiskCritical:
		return "低风险"
	default:
		return level
	}
}

// PrescriptionStatusText 审核状态展示文本（与 service 状态机、日志模板、错误码同步维护）。
func PrescriptionStatusText(status string) string {
	switch status {
	case constants.PrescriptionStatusPendingReview:
		return "待审核"
	case constants.PrescriptionStatusReviewing:
		return "审核中"
	case constants.PrescriptionStatusPassed:
		return "通过"
	case constants.PrescriptionStatusWarned:
		return "警告"
	case constants.PrescriptionStatusRejected:
		return "拒绝"
	case constants.PrescriptionStatusOverridden:
		return "强制通过"
	default:
		return status
	}
}

// RuleTypeText 规则类型展示文本。
func RuleTypeText(t string) string {
	switch t {
	case constants.RuleTypeIndication:
		return "适应症审核"
	case constants.RuleTypeDosage:
		return "用法用量审核"
	case constants.RuleTypeInteraction:
		return "药物相互作用审核"
	case constants.RuleTypeContraindication:
		return "禁忌与特殊人群审核"
	case constants.RuleTypeDuplication:
		return "重复用药审核"
	default:
		return t
	}
}

// UserRoleText 用户角色展示文本。
func UserRoleText(role string) string {
	switch role {
	case constants.UserRoleAdmin:
		return "管理员"
	case constants.UserRoleDoctor:
		return "医生"
	case constants.UserRolePharmacist:
		return "药师"
	default:
		return role
	}
}

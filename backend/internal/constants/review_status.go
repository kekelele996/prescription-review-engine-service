package constants

// 审核状态枚举（处方状态机与审核报告共用）：
// pending_review -> reviewing -> passed / warned / rejected，warned 可被医生强制通过为 overridden。
const (
	// PrescriptionStatusPendingReview 待审核（处方已接收，等待审核）。
	PrescriptionStatusPendingReview = "pending_review"
	// PrescriptionStatusReviewing 审核中。
	PrescriptionStatusReviewing = "reviewing"
	// PrescriptionStatusPassed 通过。
	PrescriptionStatusPassed = "passed"
	// PrescriptionStatusWarned 警告（存在中/高风险问题，医生可强制通过）。
	PrescriptionStatusWarned = "warned"
	// PrescriptionStatusRejected 拒绝（存在用药禁忌等极高风险）。
	PrescriptionStatusRejected = "rejected"
	// PrescriptionStatusOverridden 强制通过（医生记录原因后放行）。
	PrescriptionStatusOverridden = "overridden"
)

// AllPrescriptionStatuses 全部合法审核状态。
var AllPrescriptionStatuses = []string{
	PrescriptionStatusPendingReview, PrescriptionStatusReviewing, PrescriptionStatusPassed,
	PrescriptionStatusWarned, PrescriptionStatusRejected, PrescriptionStatusOverridden,
}

// IsValidPrescriptionStatus 判断审核状态是否合法。
func IsValidPrescriptionStatus(status string) bool {
	for _, s := range AllPrescriptionStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// FinalReviewStatuses 终态集合：允许展示为最终审核结论的状态。
var FinalReviewStatuses = []string{PrescriptionStatusPassed, PrescriptionStatusWarned, PrescriptionStatusRejected, PrescriptionStatusOverridden}

// IsFinalReviewStatus 判断是否终态。
func IsFinalReviewStatus(status string) bool {
	for _, s := range FinalReviewStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// MapRiskToStatus 风险等级 -> 审核结论状态机。
func MapRiskToStatus(level string) string {
	switch level {
	case RiskCritical:
		return PrescriptionStatusRejected
	case RiskHigh, RiskMedium:
		return PrescriptionStatusWarned
	default:
		return PrescriptionStatusWarned
	}
}

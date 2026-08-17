package dto

import "github.com/rxcheck/rxcheck/internal/model"

// ReviewItemResp 审核意见响应。
type ReviewItemResp struct {
	DrugName   string `json:"drug_name"`
	RuleType   string `json:"rule_type"`
	RiskLevel  string `json:"risk_level"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion"`
	RuleSource string `json:"rule_source"`
	Operation  string `json:"operation"`
}

// ToReviewItemResp 审核项模型转响应。
func ToReviewItemResp(i *model.ReviewItem) ReviewItemResp {
	return ReviewItemResp{
		DrugName:   i.DrugName,
		RuleType:   i.RuleType,
		RiskLevel:  i.RiskLevel,
		Message:    i.Message,
		Suggestion: i.Suggestion,
		RuleSource: i.RuleSource,
		Operation:  i.Operation,
	}
}

// ReportResp 审核报告响应。
type ReportResp struct {
	ID             uint             `json:"id"`
	PrescriptionID uint             `json:"prescription_id"`
	Status         string           `json:"status"`
	RiskLevel      string           `json:"risk_level"`
	Summary        string           `json:"summary"`
	Items          []ReviewItemResp `json:"items"`
	CreatedAt      string           `json:"created_at"`
}

// ToReportResp 报告模型转响应。
func ToReportResp(r *model.ReviewReport) ReportResp {
	out := ReportResp{
		ID:             r.ID,
		PrescriptionID: r.PrescriptionID,
		Status:         r.Status,
		RiskLevel:      r.RiskLevel,
		Summary:        r.Summary,
		CreatedAt:      r.CreatedAt.Format("2006-01-02 15:04:05"),
		Items:          make([]ReviewItemResp, 0, len(r.Items)),
	}
	for i := range r.Items {
		out.Items = append(out.Items, ToReviewItemResp(&r.Items[i]))
	}
	return out
}

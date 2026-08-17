package dto

// StatsOverviewResp 审核统计概览。
type StatsOverviewResp struct {
	TotalPrescriptions int64            `json:"total_prescriptions"`
	ByStatus           map[string]int64 `json:"by_status"`
	ByRisk             map[string]int64 `json:"by_risk"`
	TotalDrugs         int64            `json:"total_drugs"`
	TotalInteractions  int64            `json:"total_interactions"`
}

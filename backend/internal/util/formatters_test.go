package util

import (
	"testing"

	"github.com/rxcheck/rxcheck/internal/constants"
)

func TestRiskLevelText(t *testing.T) {
	cases := map[string]string{
		constants.RiskNone:     "无风险",
		constants.RiskLow:      "低风险",
		constants.RiskMedium:   "中风险",
		constants.RiskHigh:     "高风险",
		constants.RiskCritical: "极高风险",
	}
	for in, want := range cases {
		if got := RiskLevelText(in); got != want {
			t.Errorf("RiskLevelText(%s) = %s, want %s", in, got, want)
		}
	}
}

func TestPrescriptionStatusText(t *testing.T) {
	if got := PrescriptionStatusText(constants.PrescriptionStatusWarned); got != "警告" {
		t.Errorf("got %s", got)
	}
}

func TestRuleTypeText(t *testing.T) {
	if got := RuleTypeText(constants.RuleTypeInteraction); got != "药物相互作用审核" {
		t.Errorf("got %s", got)
	}
}

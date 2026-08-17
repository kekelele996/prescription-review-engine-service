.PHONY: test-review test-dosage test-override

test-review:
	go test ./internal/service/ -run 'TestReviewAggregation_CriticalRejected|TestReviewAggregation_HighWarned' -count=20

test-dosage:
	go test ./internal/service/ -run 'TestFrequencyDailyTimes_Q8H|TestMaxDailyDose_Child|TestFormatDose_Precision|TestReviewPersistsOverallRisk|TestMapRiskToStatus_Critical' -count=20

test-override:
	go test ./internal/service/ -run 'TestOverride_RejectedNotAllowed|TestReReview_OverriddenBlocked|TestStatusText_Overridden|TestNormalReview_Passed' -count=20

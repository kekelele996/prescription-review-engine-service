# Bug 复现说明

## Bug 是什么

处方自动审核的结论错误：过敏禁忌（极高风险）处方被放行为"通过"、严重相互作用只判"通过"而非"警告"；报告里的最高风险等级文案也错。

## 如何触发

```bash
cd backend
make test-review
```

## 错误信息

```
--- FAIL: TestReviewAggregation_CriticalRejected
    status = warned, want rejected
--- FAIL: TestReviewAggregation_HighWarned
    status = passed, want warned
```

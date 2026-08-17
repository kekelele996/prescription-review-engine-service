# Bug 复现说明

## Bug 是什么

处方用量审核漏检：q8h 每日次数算错、儿童按成人上限放行、剂量建议文案丢小数、审核后处方整体风险等级未持久化、critical 风险映射错误。

## 如何触发

```bash
cd backend
make test-dosage
```

## 错误信息

```
--- FAIL: TestFrequencyDailyTimes_Q8H
    q8h daily times = 2, want 3
--- FAIL: TestMaxDailyDose_Child
    child max dose = 10, want 5
--- FAIL: TestReviewPersistsOverallRisk
    overall risk = "none", want medium
```

# Bug 复现说明

## Bug 是什么

处方状态机错乱：被拒绝的处方可以强制通过、已强制通过的还能重新审核、状态文案错误、无风险处方被标记为"警告"。

## 如何触发

```bash
cd backend
make test-override
```

## 错误信息

```
--- FAIL: TestOverride_RejectedNotAllowed
    expected override not allowed for rejected prescription
--- FAIL: TestReReview_OverriddenBlocked
    expected re-review of overridden prescription to be blocked
--- FAIL: TestNormalReview_Passed
    status = warned, want passed
```

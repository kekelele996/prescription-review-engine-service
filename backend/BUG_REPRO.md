# Bug 复现说明

## Bug 是什么

药物相互作用与重复用药漏检：处方里药品顺序反过来就查不到相互作用；两种 ACEI 类降压药同时开没有产生重复用药高风险提示。

## 如何触发

```bash
cd backend
go test ./internal/service/ -run 'TestInteraction_ReverseOrderDetected|TestDuplication_AceHigh'
```

## 错误信息

```
--- FAIL: TestInteraction_ReverseOrderDetected
    expected interaction item for reversed drug pair
--- FAIL: TestDuplication_AceHigh
    expected high duplication item, found=false high=false
```

# Bug 复现说明

## Bug 是什么

处方列表接口顺序反了，列表里混入未审核处方时视图组装直接空指针崩溃，审核报告明细也查不出来。

## 如何触发

```bash
cd backend
go test ./internal/service/ -run 'TestPrescriptionList_NoPanicAndOrder|TestReportItems_Loaded'
```

## 错误信息

```
--- FAIL: TestPrescriptionList_NoPanicAndOrder
    List panicked: runtime error: invalid memory address or nil pointer dereference
    first item id = 1, want newer 2
--- FAIL: TestReportItems_Loaded
    expected report items, got none
```

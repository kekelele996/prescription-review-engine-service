package constants

// 业务错误码集中维护。0 表示成功；其余按模块分段，供响应体 code 字段使用。
const (
	// CodeOK 成功。
	CodeOK = 0

	// CodeValidation 参数校验失败。
	CodeValidation = 1001
	// CodeBadRequest 请求格式不合法。
	CodeBadRequest = 1002
	// CodeUnauthorized 未认证。
	CodeUnauthorized = 1003
	// CodeInvalidToken 令牌无效。
	CodeInvalidToken = 1004
	// CodeTokenExpired 令牌过期。
	CodeTokenExpired = 1005
	// CodeForbidden 无权限。
	CodeForbidden = 1006
	// CodeUserDisabled 用户被禁用。
	CodeUserDisabled = 1007
	// CodeWrongPassword 用户名或密码错误。
	CodeWrongPassword = 1008
	// CodeNotFound 资源不存在。
	CodeNotFound = 1009
	// CodeConflict 业务冲突。
	CodeConflict = 1010
	// CodeInvalidStatus 状态流转不合法。
	CodeInvalidStatus = 1011
	// CodeInternalError 内部错误。
	CodeInternalError = 1012
	// CodeRateLimited 请求过于频繁。
	CodeRateLimited = 1013
	// CodeInvalidFormat 处方格式不支持（仅支持 JSON/XML）。
	CodeInvalidFormat = 1014
	// CodeDrugNotFound 药品在药品库中不存在。
	CodeDrugNotFound = 1015
	// CodePrescriptionLocked 处方处于终态，不可重复审核。
	CodePrescriptionLocked = 1016
	// CodeOverrideNotAllowed 当前状态不允许强制通过。
	CodeOverrideNotAllowed = 1017
	// CodeStorageError 对象存储异常。
	CodeStorageError = 1018
)

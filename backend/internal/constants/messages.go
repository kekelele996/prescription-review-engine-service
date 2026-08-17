package constants

// 接口返回文案、日志文案、错误提示文案集中维护（多处耦合）。
const (
	MsgOK                  = "ok"
	MsgInternalError       = "系统繁忙，请稍后重试"
	MsgUnauthorized        = "未登录或登录已失效"
	MsgInvalidToken        = "访问令牌无效"
	MsgTokenExpired        = "访问令牌已过期"
	MsgForbidden           = "无权访问该资源: role=%s"
	MsgUserDisabled        = "账号已被禁用: username=%s"
	MsgWrongPassword       = "用户名或密码错误: username=%s"
	MsgValidationFailed    = "参数校验失败: field=%s rule=%s"
	MsgUserExists          = "用户名已存在: username=%s"
	MsgUserNotFound        = "用户不存在: user_id=%d"
	MsgDrugNotFound        = "药品不存在: drug_id=%d"
	MsgDrugExists          = "药品已存在: name=%s"
	MsgInteractionNotFound = "相互作用规则不存在: rule_id=%d"
	MsgInteractionExists   = "相互作用规则已存在: drug_a=%d drug_b=%d"
	MsgPrescriptionNotFound = "处方不存在: prescription_id=%d"
	MsgPrescriptionLocked  = "处方已进入终态，不可重复审核: status=%s"
	MsgOverrideNotAllowed  = "仅警告/通过类处方可强制通过: status=%s"
	MsgReportNotFound      = "审核报告不存在: prescription_id=%d"
	MsgInvalidFormat       = "处方格式不支持: format=%s（仅支持 json/xml）"
	MsgRateLimited         = "请求过于频繁，请稍后再试"
	MsgStorageError        = "对象存储不可用，报告导出失败"
	MsgParseFailed         = "处方解析失败: %s"
)

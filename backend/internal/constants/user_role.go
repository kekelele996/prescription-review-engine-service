package constants

// 用户角色枚举：模型、DTO、service 状态机、handler 校验、RBAC、日志模板、错误码中统一引用。
const (
	// UserRoleAdmin 管理员：系统配置、药品库与规则库维护。
	UserRoleAdmin = "admin"
	// UserRoleDoctor 医生：开具处方、强制通过警告类处方。
	UserRoleDoctor = "doctor"
	// UserRolePharmacist 药师：审核处方、复核报告。
	UserRolePharmacist = "pharmacist"
)

// AllUserRoles 全部合法角色（handler 校验用）。
var AllUserRoles = []string{UserRoleAdmin, UserRoleDoctor, UserRolePharmacist}

// IsValidUserRole 判断角色是否合法。
func IsValidUserRole(role string) bool {
	for _, r := range AllUserRoles {
		if r == role {
			return true
		}
	}
	return false
}

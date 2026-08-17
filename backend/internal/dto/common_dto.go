package dto

// IDParam 路径 ID 参数。
type IDParam struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

// PageQuery 分页查询参数。
type PageQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
}

// LoginReq 登录请求。
type LoginReq struct {
	Username string `json:"username" binding:"required,min=2,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
}

// LoginResp 登录响应。
type LoginResp struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rxcheck/rxcheck/internal/config"
	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/util"
)

// APIKeyHeader API Key 请求头名称。
const APIKeyHeader = "X-API-Key"

// APIKeyAuth API Key 认证中间件：校验服务间调用的静态 API Key，身份为系统管理员。
func APIKeyAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader(APIKeyHeader)
		if key == "" {
			util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, constants.MsgUnauthorized)
			c.Abort()
			return
		}
		if subtle.ConstantTimeCompare([]byte(key), []byte(cfg.APIKeySecret)) != 1 {
			util.Fail(c, http.StatusUnauthorized, constants.CodeInvalidToken, "API Key 无效")
			c.Abort()
			return
		}
		// API Key 以系统管理员身份访问。
		c.Set(UserKey, &util.Claims{UserID: 0, Username: "api-client", Role: constants.UserRoleAdmin})
		c.Next()
	}
}

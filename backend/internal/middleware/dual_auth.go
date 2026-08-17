package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rxcheck/rxcheck/internal/config"
)

// RequireAnyAuth 双认证中间件：API Key 或 JWT 任一通过即可访问（API Key + JWT 双认证）。
func RequireAnyAuth(cfg *config.Config) gin.HandlerFunc {
	apiKeyMW := APIKeyAuth(cfg)
	jwtMW := JWTAuth(cfg)
	return func(c *gin.Context) {
		key := c.GetHeader(APIKeyHeader)
		if key != "" {
			apiKeyMW(c)
			if c.IsAborted() {
				return
			}
			c.Next()
			return
		}
		jwtMW(c)
		if c.IsAborted() {
			return
		}
		c.Next()
	}
}

// RequireJWTOnly 仅允许 JWT 认证的中间件（登录态操作）。
func RequireJWTOnly(cfg *config.Config) gin.HandlerFunc {
	return JWTAuth(cfg)
}

package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rxcheck/rxcheck/internal/middleware"
)

// rid 读取请求 ID。
func rid(c *gin.Context) string {
	v, _ := c.Get(middleware.RequestIDKey)
	return ridString(v)
}

// ridString 格式化请求 ID。
func ridString(v any) string {
	return fmt.Sprintf("%v", v)
}

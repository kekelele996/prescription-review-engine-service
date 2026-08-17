// Package pointerx 提供基础指针工具函数（无业务依赖，可复用库）。
package pointerx

// P 返回任意值的指针。
func P[T any](v T) *T { return &v }

// Int 返回 int 指针。
func Int(v int) *int { return &v }

// Uint 返回 uint 指针。
func Uint(v uint) *uint { return &v }

// String 返回 string 指针。
func String(v string) *string { return &v }

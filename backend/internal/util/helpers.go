package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Uint64String 将 uint 转为字符串。
func Uint64String(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}

// GenSerial 生成带前缀与时间戳的单号。
func GenSerial(prefix string) string {
	return fmt.Sprintf("%s%s%s", prefix, time.Now().Format("20060102150405"), strconv.FormatInt(time.Now().UnixNano()%100000, 10))
}

// NormalizeContentType 归一化 Content-Type，返回 json/xml。
func NormalizeContentType(raw string) string {
	raw = strings.ToLower(raw)
	switch {
	case strings.Contains(raw, "xml"):
		return "xml"
	case strings.Contains(raw, "json"):
		return "json"
	default:
		return ""
	}
}

// ParseFloatString 安全解析浮点数字符串。
func ParseFloatString(s string) float64 {
	if s == "" {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

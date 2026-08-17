package util

import (
	"log/slog"
	"os"
)

// Log 全局结构化日志器（slog），所有 handler/service/middleware 统一引用。
var Log *slog.Logger

// InitLogger 初始化结构化日志器（JSON 输出）。
func InitLogger(level slog.Level) {
	opts := &slog.HandlerOptions{Level: level}
	Log = slog.New(slog.NewJSONHandler(os.Stdout, opts))
}

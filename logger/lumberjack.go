package logger

import (
	logger2 "github.com/fushiliang321/go-core/config/logger"
	"time"
)

// 默认按日期切割日志
func DefaultLumberjackConfig() *logger2.Lumberjack {
	return &logger2.Lumberjack{
		FileNameFormat: func() string {
			return time.Now().Format("2006-01-02")
		},
		MaxSize:    10,
		MaxBackups: 0,
		MaxAge:     0,
	}
}

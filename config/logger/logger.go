package logger

import (
	"golang.org/x/exp/slog"
	"time"
)

type (
	// Lumberjack 文件切割
	Lumberjack struct {
		FileNameFormat func() string //文件名格式
		MaxSize        int           //在进行切割之前，日志文件的最大值（以MB为单位）
		MaxBackups     int           //保留旧文件的最大个数
		MaxAge         int           //保留旧文件的最大天数
	}

	Logger struct {
		Handler       *slog.Handler
		DirPath       string         //日志目录
		OutputJson    bool           //是否写入json格式
		WriteStdout   bool           //是否写入控制台
		WriteFile     bool           //是否写入日志文件
		Level         slog.Leveler   //写入的日志等级
		Levels        []slog.Leveler //写入的日志等级
		Lumberjack    *Lumberjack    //日志文件切割
		WriteInterval time.Duration  //日志文件写入间隔
	}
)

var data = &Logger{
	WriteStdout: true,
	WriteFile:   true,
}

func Set(config *Logger) {
	data = config
}

func Get() *Logger {
	return data
}

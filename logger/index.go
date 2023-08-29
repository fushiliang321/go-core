package logger

import (
	logger2 "github.com/fushiliang321/go-core/config/logger"
	"github.com/fushiliang321/go-core/helper"
	"golang.org/x/exp/slog"
	"os"
)

var (
	logger     *slog.Logger
	config     *logger2.Logger
	lumberjack *logger2.Lumberjack
)

func init() {
	var handler slog.Handler
	config = logger2.Get()
	if config.Lumberjack == nil {
		lumberjack = DefaultLumberjackConfig()
	} else {
		lumberjack = config.Lumberjack
	}

	if config.Handler == nil {
		var writer = &handlerWriter{}
		if config.WriteFile {
			log := &logFile{
				dirPath:    config.DirPath,
				fileName:   "logger",
				lumberjack: lumberjack,
			}

			if log.dirPath == "" {
				log.dirPath = "./runtime"
			}
			err := log.open()
			if err == nil {
				writer.setLogFile(log)
			}
		}
		if config.WriteStdout {
			//需要写入控制台
			writer.setStdout(os.Stdout)
		}

		if config.OutputJson {
			handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: config.Level})
		} else {
			handler = slog.NewTextHandler(writer, &slog.HandlerOptions{Level: config.Level})
		}

	} else {
		handler = *config.Handler
	}
	logger = slog.New(handler)
}

func Info(msg any, args ...any) {
	if logger == nil {
		return
	}
	str, err := helper.AnyToString(msg)
	if err != nil {
		return
	}
	logger.Info(str, args...)
}

func Debug(msg any, args ...any) {
	if logger == nil {
		return
	}

	str, err := helper.AnyToString(msg)
	if err != nil {
		return
	}
	logger.Debug(str, args...)
}

func Warn(msg any, args ...any) {
	if logger == nil {
		return
	}
	str, err := helper.AnyToString(msg)
	if err != nil {
		return
	}
	logger.Warn(str, args...)
}

func Error(msg any, args ...any) {
	if logger == nil {
		return
	}
	str, err := helper.AnyToString(msg)
	if err != nil {
		return
	}
	logger.Error(str, args...)
}

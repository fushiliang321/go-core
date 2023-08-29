package logger

import (
	logger2 "github.com/fushiliang321/go-core/config/logger"
	"github.com/fushiliang321/go-core/helper"
	"golang.org/x/exp/slog"
	"os"
	"strings"
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

func Info(msgs ...any) {
	if logger == nil {
		return
	}
	var build strings.Builder
	for _, msg := range msgs {
		bytes, _ := helper.AnyToBytes(msg)
		build.Write(bytes)
	}
	if build.Len() == 0 {
		return
	}
	logger.Info(build.String())
}

func Debug(msgs ...any) {
	if logger == nil {
		return
	}
	var build strings.Builder
	for _, msg := range msgs {
		bytes, _ := helper.AnyToBytes(msg)
		build.Write(bytes)
	}
	if build.Len() == 0 {
		return
	}
	logger.Debug(build.String())
}

func Warn(msgs ...any) {
	if logger == nil {
		return
	}

	var build strings.Builder
	for _, msg := range msgs {
		bytes, _ := helper.AnyToBytes(msg)
		build.Write(bytes)
	}
	if build.Len() == 0 {
		return
	}
	logger.Warn(build.String())
}

func Error(msgs ...any) {
	if logger == nil {
		return
	}

	var build strings.Builder
	for _, msg := range msgs {
		bytes, _ := helper.AnyToBytes(msg)
		build.Write(bytes)
	}
	if build.Len() == 0 {
		return
	}
	logger.Error(build.String())
}

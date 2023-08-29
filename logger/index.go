package logger

import (
	logger2 "github.com/fushiliang321/go-core/config/logger"
	"github.com/fushiliang321/go-core/helper"
	"github.com/fushiliang321/go-core/logger/agency"
	"golang.org/x/exp/slog"
	"os"
	"strings"
	"sync"
)

type (
	Service struct{}

	Logger slog.Logger
)

var (
	slogger    *Logger
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
	slogger = (*Logger)(slog.New(handler))
}

func (s *Service) Start(wg *sync.WaitGroup) {
	agency.Set(slogger)
}

func (l *Logger) Info(msgs ...any) {
	var build = msgBuild(msgs)
	if build.Len() == 0 {
		return
	}
	(*slog.Logger)(l).Info(build.String())
}

func (l *Logger) Debug(msgs ...any) {
	var build = msgBuild(msgs)
	if build.Len() == 0 {
		return
	}
	(*slog.Logger)(l).Debug(build.String())
}

func (l *Logger) Warn(msgs ...any) {
	var build = msgBuild(msgs)
	if build.Len() == 0 {
		return
	}
	(*slog.Logger)(l).Warn(build.String())
}

func (l *Logger) Error(msgs ...any) {
	var build = msgBuild(msgs)
	if build.Len() == 0 {
		return
	}
	(*slog.Logger)(l).Error(build.String())
}

func msgBuild(msgs []any) *strings.Builder {
	var build = strings.Builder{}
	for i, msg := range msgs {
		if i > 0 {
			build.WriteByte(32)
		}
		bytes, _ := helper.AnyToBytes(msg)
		build.Write(bytes)
	}
	return &build
}

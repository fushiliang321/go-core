package logger

import "github.com/fushiliang321/go-core/logger/agency"

func Info(msgs ...any) {
	agency.Info(msgs...)
}

func Debug(msgs ...any) {
	agency.Debug(msgs...)
}

func Warn(msgs ...any) {
	agency.Warn(msgs...)
}

func Error(msgs ...any) {
	agency.Error(msgs...)
}

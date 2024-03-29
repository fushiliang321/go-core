package agency

import (
	"github.com/fushiliang321/go-core/helper"
	"golang.org/x/exp/slog"
	"strings"
)

type defaultLogger struct{}

func init() {
	Set(&defaultLogger{})
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

func (l *defaultLogger) Info(msgs ...any) {
	slog.Info(msgBuild(msgs).String())
}

func (l *defaultLogger) Debug(msgs ...any) {
	slog.Debug(msgBuild(msgs).String())
}

func (l *defaultLogger) Warn(msgs ...any) {
	slog.Warn(msgBuild(msgs).String())
}

func (l *defaultLogger) Error(msgs ...any) {
	slog.Error(msgBuild(msgs).String())
}

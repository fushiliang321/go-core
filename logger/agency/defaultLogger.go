package agency

import (
	"encoding/json"
	"github.com/savsgio/gotils/strconv"
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
		bytes, _ := anyToBytes(msg)
		build.Write(bytes)
	}
	return &build
}

// any转bytes
func anyToBytes(data any) (bts []byte, err error) {
	switch data.(type) {
	case string:
		bts = strconv.S2B(data.(string))
	case *string:
		bts = strconv.S2B(*(data.(*string)))
	case []byte:
		return data.([]byte), nil
	case *[]byte:
		bts = *data.(*[]byte)
	case byte:
		bts = []byte{data.(byte)}
	case *byte:
		bts = []byte{*(data.(*byte))}
	default:
		bts, err = json.Marshal(data)
	}
	return
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

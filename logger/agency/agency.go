package agency

type Logger interface {
	Info(msgs ...any)
	Debug(msgs ...any)
	Warn(msgs ...any)
	Error(msgs ...any)
}

var _logger Logger

func Set(logger Logger) {
	_logger = logger
}

func Info(msgs ...any) {
	if _logger == nil {
		return
	}
	_logger.Info(msgs...)
}

func Debug(msgs ...any) {
	if _logger == nil {
		return
	}
	_logger.Debug(msgs...)
}

func Warn(msgs ...any) {
	if _logger == nil {
		return
	}
	_logger.Warn(msgs...)
}

func Error(msgs ...any) {
	if _logger == nil {
		return
	}
	_logger.Error(msgs...)
}

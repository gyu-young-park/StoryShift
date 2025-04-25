package log

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	core logger    = nil
	once sync.Once = sync.Once{}
)

var LoggerLevelEnum = loggerLevelEnum{
	INFO:  "INFO",
	DEBUG: "DEBUG",
	ERROR: "ERROR",
	WARN:  "WARN",
}

type loggerLevel string

type loggerLevelEnum struct {
	INFO  loggerLevel
	DEBUG loggerLevel
	WARN  loggerLevel
	ERROR loggerLevel
}

type logger interface {
	loggerFomatable
	closable
	SetLevel(loggerLevel) error
	Info(string, ...any)
	Debug(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
}

type loggerFomatable interface {
	Infof(string, ...any)
	Debugf(string, ...any)
	Warnf(string, ...any)
	Errorf(string, ...any)
}

type closable interface {
	Close() error
}

func GetLogger() logger {
	if core == nil {
		zapLoggerImpl, _ := newZapLoggerImpl(setZapConfig())
		core = zapLoggerImpl
	}
	return core
}

func Config(loggerType string) {

}

// TODO: Configuration 어떻게 할 것인가???
// TODO: default logger는 무엇으로?
func setZapConfig() zap.Config {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // Capitalize the log level names
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC timestamp format
		EncodeDuration: zapcore.SecondsDurationEncoder, // Duration in seconds
		EncodeCaller:   zapcore.ShortCallerEncoder,     // Short caller (file and line)
	}

	logLevel := zap.NewAtomicLevelAt(zap.InfoLevel)

	return zap.Config{
		Level:            logLevel,
		Development:      true,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

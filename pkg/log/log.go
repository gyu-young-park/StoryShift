package log

import (
	"fmt"

	"github.com/gyu-young-park/VelogStoryShift/internal/config"
)

var (
	core logger = nil
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

func GetLogger(c config.LogConfigModel) logger {
	if core == nil {
		core = newLogger(c)
	}
	return core
}

func newLogger(c config.LogConfigModel) logger {
	if c.Library == "zap" {
		l, err := newZapLoggerImpl(setZapConfig(c))
		if err != nil {
			fmt.Printf("failed to create logger: %v", c.Library)
			return newFmtLoggerImpl(c)
		}

		return l
	}

	return newFmtLoggerImpl(c)
}

func SetLoggerLevel(l loggerLevel) error {
	if core == nil {
		return fmt.Errorf("logger is now nil")
	}
	return core.SetLevel(l)
}

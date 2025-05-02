package log

import (
	"fmt"

	"github.com/gyu-young-park/StoryShift/internal/config"
)

type fmtLoggerLevel int

func (f fmtLoggerLevel) isWritable(target fmtLoggerLevel) bool {
	return target >= f
}

const (
	ERR fmtLoggerLevel = iota
	WARN
	DEBUG
	INFO
)

func convertLoggerLevelToFmtLoggerLevel(l loggerLevel) fmtLoggerLevel {
	mapper := map[loggerLevel]fmtLoggerLevel{
		LoggerLevelEnum.ERROR: ERR,
		LoggerLevelEnum.WARN:  WARN,
		LoggerLevelEnum.DEBUG: DEBUG,
		LoggerLevelEnum.INFO:  INFO,
	}

	fl, ok := mapper[l]
	if !ok {
		fmt.Println("failed to write log level")
		return INFO
	}

	return fl
}

func newFmtLoggerImpl(c config.LogConfigModel) *fmtLoggerImpl {
	return &fmtLoggerImpl{
		level: convertLoggerLevelToFmtLoggerLevel(loggerLevel(c.Level)),
	}
}

type fmtLoggerImpl struct {
	level fmtLoggerLevel
}

func (f *fmtLoggerImpl) write(level fmtLoggerLevel, msg string, variables ...any) {
	if f.level.isWritable(level) {
		data := []any{msg}
		data = append(data, variables...)
		fmt.Println(data...)
	}
}

func (f *fmtLoggerImpl) writef(level fmtLoggerLevel, msg string, variables ...any) {
	if f.level.isWritable(level) {
		fmt.Printf(msg, variables...)
	}
}

func (f *fmtLoggerImpl) SetLevel(level loggerLevel) error {
	f.level = convertLoggerLevelToFmtLoggerLevel(level)
	return nil
}

func (f *fmtLoggerImpl) Info(msg string, variables ...any) {
	f.write(INFO, msg, variables...)
}

func (f *fmtLoggerImpl) Debug(msg string, variables ...any) {
	f.write(DEBUG, msg, variables...)
}

func (f *fmtLoggerImpl) Warn(msg string, variables ...any) {
	f.write(WARN, msg, variables...)
}

func (f *fmtLoggerImpl) Error(msg string, variables ...any) {
	f.write(ERR, msg, variables...)
}

func (f *fmtLoggerImpl) Infof(msg string, variables ...any) {
	f.writef(INFO, msg, variables...)
}

func (f *fmtLoggerImpl) Debugf(msg string, variables ...any) {
	f.writef(DEBUG, msg, variables...)
}

func (f *fmtLoggerImpl) Warnf(msg string, variables ...any) {
	f.writef(WARN, msg, variables...)
}

func (f *fmtLoggerImpl) Errorf(msg string, variables ...any) {
	f.writef(ERR, msg, variables...)
}

func (f *fmtLoggerImpl) Close() error {
	return nil
}

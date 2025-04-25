package log

import (
	"fmt"
)

func newFmtLoggerImpl() *fmtLoggerImpl {
	return &fmtLoggerImpl{}
}

type fmtLoggerImpl struct {
	level loggerLevel
}

func (f *fmtLoggerImpl) SetLevel(level loggerLevel) error {
	f.level = level
	return nil
}

func (f *fmtLoggerImpl) Info(msg string, variables ...any) {
	fmt.Println(msg, variables)
}

func (f *fmtLoggerImpl) Debug(msg string, variables ...any) {
	fmt.Print(msg, variables)
}

func (f *fmtLoggerImpl) Warn(msg string, variables ...any) {
	fmt.Print(msg, variables)
}

func (f *fmtLoggerImpl) Error(msg string, variables ...any) {
	fmt.Print(msg, variables)
}

func (f *fmtLoggerImpl) Infof(msg string, variables ...any) {
	fmt.Printf(msg, variables)
}

func (f *fmtLoggerImpl) Debugf(msg string, variables ...any) {
	fmt.Printf(msg, variables)
}

func (f *fmtLoggerImpl) Warnf(msg string, variables ...any) {
	fmt.Printf(msg, variables)
}

func (f *fmtLoggerImpl) Errorf(msg string, variables ...any) {
	fmt.Printf(msg, variables)
}

func (f *fmtLoggerImpl) Close() error {
	return nil
}

package log

import (
	"fmt"

	"github.com/gyu-young-park/StoryShift/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLoggerImpl struct {
	logger *zap.Logger
	level  zap.AtomicLevel
}

func newZapLoggerImpl(config zap.Config) (*zapLoggerImpl, error) {
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &zapLoggerImpl{
		logger: logger,
		level:  config.Level,
	}, nil
}

func convertAnyListToZapFieldList(variables ...any) ([]zap.Field, bool) {
	fields := []zap.Field{}
	for _, variable := range variables {
		field, ok := variable.(zap.Field)
		if !ok {
			return fields, false
		}
		fields = append(fields, field)
	}

	return fields, true
}

func setZapConfig(c config.LogConfigModel) zap.Config {
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

	zapLogLevel := convertLogLevelToZapLogLevel(loggerLevel(c.Level))
	logLevel := zap.NewAtomicLevelAt(zapLogLevel)

	return zap.Config{
		Level:            logLevel,
		Development:      true,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func convertLogLevelToZapLogLevel(level loggerLevel) zapcore.Level {
	logLevelMapper := map[loggerLevel]zapcore.Level{
		LoggerLevelEnum.INFO:  zapcore.InfoLevel,
		LoggerLevelEnum.DEBUG: zapcore.DebugLevel,
		LoggerLevelEnum.ERROR: zapcore.ErrorLevel,
		LoggerLevelEnum.WARN:  zapcore.WarnLevel,
	}

	zapLoggerLevel, ok := logLevelMapper[level]
	if !ok {
		panic(fmt.Sprintf("failed to change log level to: %s", level))
	}

	return zapLoggerLevel
}

func (z *zapLoggerImpl) SetLevel(level loggerLevel) error {
	zapLogLevel := convertLogLevelToZapLogLevel(level)
	z.level.SetLevel(zapLogLevel)
	return nil
}

func (z *zapLoggerImpl) Info(msg string, variables ...any) {
	fields, ok := convertAnyListToZapFieldList(variables...)
	if !ok {
		z.logger.Sugar().Info(msg)
		return
	}
	z.logger.Info(msg, fields...)
}

func (z *zapLoggerImpl) Debug(msg string, variables ...any) {
	fields, ok := convertAnyListToZapFieldList(variables...)
	if !ok {
		z.logger.Sugar().Debug(msg)
		return
	}
	z.logger.Debug(msg, fields...)
}

func (z *zapLoggerImpl) Warn(msg string, variables ...any) {
	fields, ok := convertAnyListToZapFieldList(variables...)
	if !ok {
		z.logger.Sugar().Warn(msg)
		return
	}
	z.logger.Warn(msg, fields...)
}

func (z *zapLoggerImpl) Error(msg string, variables ...any) {
	fields, ok := convertAnyListToZapFieldList(variables...)
	if !ok {
		z.logger.Sugar().Error(msg)
		return
	}
	z.logger.Error(msg, fields...)
}

func (z *zapLoggerImpl) Infof(msg string, variables ...any) {
	z.logger.Sugar().Infof(msg, variables...)
}

func (z *zapLoggerImpl) Debugf(msg string, variables ...any) {
	z.logger.Sugar().Debugf(msg, variables...)
}

func (z *zapLoggerImpl) Warnf(msg string, variables ...any) {
	z.logger.Sugar().Warnf(msg, variables...)
}

func (z *zapLoggerImpl) Errorf(msg string, variables ...any) {
	z.logger.Sugar().Errorf(msg, variables...)
}

func (z *zapLoggerImpl) Close() error {
	defer z.logger.Sync()
	return nil
}

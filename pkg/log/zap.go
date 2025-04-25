package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLoggerImpl struct {
	logger *zap.Logger
	level  zap.AtomicLevel
}

func newZapLoggerImpl(config zap.Config) (*ZapLoggerImpl, error) {
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &ZapLoggerImpl{
		logger: logger,
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

func (z *ZapLoggerImpl) SetLevel(level loggerLevel) error {
	logLevelMapper := map[loggerLevel]zapcore.Level{
		LoggerLevelEnum.INFO:  zapcore.InfoLevel,
		LoggerLevelEnum.DEBUG: zapcore.DebugLevel,
		LoggerLevelEnum.ERROR: zapcore.ErrorLevel,
		LoggerLevelEnum.WARN:  zapcore.WarnLevel,
	}

	zapLoggerLevel, ok := logLevelMapper[level]
	if !ok {
		return fmt.Errorf("failed to change log level to: %s", level)
	}

	z.level.SetLevel(zapLoggerLevel)
	return nil
}

func (z *ZapLoggerImpl) Info(msg string, variables ...any) {
	fields, ok := convertAnyListToZapFieldList(variables...)
	if !ok {
		z.logger.Sugar().Info(msg)
		return
	}
	z.logger.Info(msg, fields...)
}

func (z *ZapLoggerImpl) Debug(msg string, variables ...any) {
	fields, ok := convertAnyListToZapFieldList(variables...)
	if !ok {
		z.logger.Sugar().Debug(msg)
		return
	}
	z.logger.Debug(msg, fields...)
}

func (z *ZapLoggerImpl) Warn(msg string, variables ...any) {
	fields, ok := convertAnyListToZapFieldList(variables...)
	if !ok {
		z.logger.Sugar().Warn(msg)
		return
	}
	z.logger.Warn(msg, fields...)
}

func (z *ZapLoggerImpl) Error(msg string, variables ...any) {
	fields, ok := convertAnyListToZapFieldList(variables...)
	if !ok {
		z.logger.Sugar().Error(msg)
		return
	}
	z.logger.Error(msg, fields...)
}

func (z *ZapLoggerImpl) Infof(msg string, variables ...any) {
	z.logger.Sugar().Infof(msg, variables...)
}

func (z *ZapLoggerImpl) Debugf(msg string, variables ...any) {
	z.logger.Sugar().Debugf(msg, variables...)
}

func (z *ZapLoggerImpl) Warnf(msg string, variables ...any) {
	z.logger.Sugar().Warnf(msg, variables...)
}

func (z *ZapLoggerImpl) Errorf(msg string, variables ...any) {
	z.logger.Sugar().Errorf(msg, variables...)
}

func (z *ZapLoggerImpl) Close() error {
	defer z.logger.Sync()
	return nil
}

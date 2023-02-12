package logger

import (
	"go.uber.org/zap"
)

func New(fields ...zap.Field) (*Logger, error) {
	return newZapLogger(fields...)
}

func GetLogger() *zap.Logger {
	return log.logger
}

func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

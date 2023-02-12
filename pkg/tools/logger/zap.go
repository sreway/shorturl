package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Logger struct {
		logger *zap.Logger
	}
	options struct {
		IsProduction bool
		Level        zapcore.Level
	}
)

var (
	log *Logger
	opt *options
)

func init() {
	opt = &options{IsProduction: true}
	if strings.ToLower(strings.TrimSpace(os.Getenv("LOGGER_IS_PRODUCTION"))) == "false" {
		opt.IsProduction = false
	}

	switch strings.ToUpper(strings.TrimSpace(os.Getenv("LOGGER_LEVEL"))) {
	case "ERR", "ERROR":
		opt.Level = zapcore.ErrorLevel
	case "WARN", "WARNING":
		opt.Level = zapcore.WarnLevel
	case "INFO":
		opt.Level = zapcore.InfoLevel
	case "DEBUG":
		opt.Level = zapcore.DebugLevel
	case "FATAL":
		opt.Level = zapcore.FatalLevel
	default:
		opt.Level = zapcore.InfoLevel
	}

	if log == nil {
		newLogger, err := newZapLogger()
		if err != nil {
			panic(err)
		}

		log = newLogger
	}
}

func newZapLogger(fields ...zap.Field) (*Logger, error) {
	var config zap.Config

	if opt.IsProduction {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.Level = zap.NewAtomicLevelAt(opt.Level)
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	newLogger, err := config.Build(zap.AddCallerSkip(2))
	if err != nil {
		return nil, err
	}
	newLogger.Info("Set LOG_LEVEL", zap.Stringer("level", opt.Level))

	log = &Logger{logger: newLogger.With(fields...)}
	return log, nil
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(fmt.Sprintf("%v", msg), fields...)
}

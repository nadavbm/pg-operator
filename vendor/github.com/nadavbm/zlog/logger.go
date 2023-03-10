package zlog

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
	fields []zap.Field
}

func New(fields ...zapcore.Field) *Logger {
	encfg := zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "level",
		TimeKey:     "time",
		EncodeTime:  zapcore.ISO8601TimeEncoder,
		EncodeLevel: zapcore.CapitalLevelEncoder,
	}

	options := zapcore.NewCore(zapcore.NewJSONEncoder(encfg), os.Stdout, zap.DebugLevel)

	l := zap.New(options)

	logger := &Logger{
		logger: l,
		fields: fields,
	}
	return logger
}

func (l *Logger) Debug(text string, fields ...zapcore.Field) {
	l.logger.Debug(text, fields...)
}

func (l *Logger) Info(text string, fields ...zapcore.Field) {
	l.logger.Info(text, fields...)
}

func (l *Logger) Error(text string, fields ...zapcore.Field) {
	l.logger.Error(text, fields...)
}

func (l *Logger) Fatal(text string, fields ...zapcore.Field) {
	l.logger.Fatal(text, fields...)
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		logger: l.logger.With(fields...),
		fields: append(l.fields, fields...),
	}
}

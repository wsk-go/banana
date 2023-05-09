package logger

import (
	"github.com/JackWSK/banana/logger/field"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

var defaultLogger = New(Config{
	Level:       zapcore.DebugLevel,
	Writer:      []io.Writer{os.Stdout},
	LevelWriter: nil,
})

func Configure(config Config) {
	defaultLogger = New(config)
}

func Debug(msg string, fields ...field.Field) {
	defaultLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...field.Field) {
	defaultLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...field.Field) {
	defaultLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...field.Field) {
	defaultLogger.Error(msg, fields...)
}

func DPanic(msg string, fields ...field.Field) {
	defaultLogger.DPanic(msg, fields...)
}

func Panic(msg string, fields ...field.Field) {
	defaultLogger.Panic(msg, fields...)
}

package logger

import (
	"github.com/JackWSK/banana/logger/field"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
)

type Config struct {
	Level  Level
	Writer io.Writer

	LevelWriter map[Level]io.Writer
}

type Logger struct {
	defaultLogger  *zap.Logger
	loggerForLevel map[Level]*zap.Logger
}

// NewLogger defaultLogger 默认输出的logger
func NewLogger(config Config) *Logger {

	logger := &Logger{defaultLogger: NewZapLogger(config.Level, config.Writer),
		loggerForLevel: make(map[Level]*zap.Logger),
	}

	if len(config.LevelWriter) > 0 {
		for k, v := range config.LevelWriter {
			logger.ConfigureLoggerForLevel(k, v)
		}
	}

	return logger
}

// ConfigureLoggerForLevel 配置level对应的logger
// 如果没有找到，则使用defaultLogger
func (th *Logger) ConfigureLoggerForLevel(level Level, writer io.Writer) *Logger {
	th.loggerForLevel[level] = NewZapLogger(level, writer)
	return th
}

func (th *Logger) Debug(msg string, fields ...field.Field) {
	th.determineLogger(DebugLevel).Debug(msg, fields...)
}

func (th *Logger) Info(msg string, fields ...field.Field) {
	th.determineLogger(InfoLevel).Info(msg, fields...)
}

func (th *Logger) Warn(msg string, fields ...field.Field) {
	th.determineLogger(WarnLevel).Warn(msg, fields...)
}

func (th *Logger) Error(msg string, fields ...field.Field) {
	th.determineLogger(ErrorLevel).Error(msg, fields...)
}

func (th *Logger) DPanic(msg string, fields ...field.Field) {
	th.determineLogger(DPanicLevel).DPanic(msg, fields...)
}

func (th *Logger) Panic(msg string, fields ...field.Field) {
	th.determineLogger(PanicLevel).Panic(msg, fields...)
}

func (th *Logger) Enabled(level Level) bool {
	l := th.determineLogger(level)
	ce := l.Check(level, "")
	return ce != nil
}

func (th *Logger) determineLogger(level Level) *zap.Logger {
	if logger, ok := th.loggerForLevel[level]; ok {
		return logger
	}
	return th.defaultLogger
}

func NewZapLogger(level Level, writer io.Writer) *zap.Logger {
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)

	//公用编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "Level",
		NameKey:        "zap",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(writer)),
		atomicLevel, // 日志级别
	)

	return zap.New(core)
}

func NewFileWriter(filename string) io.Writer {
	return &lumberjack.Logger{
		Filename:   filename, // 日志文件路径
		MaxSize:    256,      // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,       // 日志文件最多保存多少个备份
		MaxAge:     7,        // 文件最多保存多少天
		Compress:   false,    // 是否压缩
	}
}

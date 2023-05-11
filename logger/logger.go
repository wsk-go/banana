package logger

import (
	"github.com/wsk-go/banana/logger/field"
	"github.com/wsk-go/banana/utils/stream"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
)

type Config struct {
	Level  Level
	Writer []io.Writer

	LevelWriter map[Level][]io.Writer
}

type Logger struct {
	defaultLogger  *zap.Logger
	loggerForLevel map[Level]*zap.Logger
}

// New logger
func New(config Config) *Logger {

	logger := &Logger{defaultLogger: newZapLogger(config.Level, config.Writer),
		loggerForLevel: make(map[Level]*zap.Logger),
	}

	if len(config.LevelWriter) > 0 {
		for k, v := range config.LevelWriter {
			logger.ConfigureLoggerForLevel(k, v)
		}
	}

	return logger
}

// ConfigureLoggerForLevel config level for each logger
func (th *Logger) ConfigureLoggerForLevel(level Level, writer []io.Writer) *Logger {
	th.loggerForLevel[level] = newZapLogger(level, writer)
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

func newZapLogger(level Level, writers []io.Writer) *zap.Logger {
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

	syncers := stream.Map(stream.Of(writers), func(in io.Writer) zapcore.WriteSyncer {
		return zapcore.AddSync(in)
	})
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(syncers...),
		atomicLevel, // 日志级别
	)

	return zap.New(core)
}

//type FilterWriter = lumberjack.Logger
//
//func NewFileWriter(filename string) io.Writer {
//	return &lumberjack.Logger{
//		Filename:   filename, // 日志文件路径
//		MaxSize:    256,      // 每个日志文件保存的最大尺寸 单位：M
//		MaxBackups: 30,       // 日志文件最多保存多少个备份
//		MaxAge:     7,        // 文件最多保存多少天
//		Compress:   false,    // 是否压缩
//	}
//}

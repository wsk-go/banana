package logger

import (
	"github.com/JackWSK/banana/errors"
	"github.com/JackWSK/banana/logger/field"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"testing"
)

func NewFileWriter(filename string) io.Writer {
	return &lumberjack.Logger{
		Filename:   filename, // 日志文件路径
		MaxSize:    256,      // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,       // 日志文件最多保存多少个备份
		MaxAge:     7,        // 文件最多保存多少天
		Compress:   false,    // 是否压缩
	}
}
func TestDefault(t *testing.T) {
	Configure(Config{
		Level: DebugLevel,
		Writer: []io.Writer{
			os.Stdout,
			NewFileWriter("log.log"),
		},
		LevelWriter: nil,
	})
	Info("hello", field.String("hello", "world"), field.Error(errors.New("xxx")))

}

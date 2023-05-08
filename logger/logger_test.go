package logger

import (
	"github.com/JackWSK/banana/errors"
	"github.com/JackWSK/banana/logger/field"
	"testing"
)

func TestDefault(t *testing.T) {
	Info("hello", field.String("hello", "world"), field.Error(errors.New("xxx")))

}

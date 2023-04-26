package errors

import (
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	err := NewLogicError("xxx", 12)
	fmt.Printf("%+v", err)
}

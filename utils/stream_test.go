package utils

import (
	"fmt"
	"testing"
)

func TestStream(t *testing.T) {
	aaa := []string{"123", "456"}
	aa := Map(Of(aaa).
		Filter(func(s string) bool {
			return s == "123"
		}),
		func(in string) int {
			return 1
		},
	)
	fmt.Println(aa)
}

package pointer

import "github.com/wsk-go/banana/utils/assert"

func To[T any](s T) *T {
	return &s
}

func Value[T any](t *T) T {
	if t == nil {
		var zero T
		return zero
	}
	return *t
}

// ZeroTo return nil if value is zero
func ZeroTo[T any](t T) *T {
	if assert.IsZero(t) {
		return nil
	}
	return &t
}

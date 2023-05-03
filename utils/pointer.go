package utils

import "reflect"

func ToPtr[T any](s T) *T {
	return &s
}

func IsZero(i any) bool {
	value := reflect.ValueOf(i)
	return value.IsZero()
}

func ToValue[T any](t *T) T {
	if t == nil {
		var zero T
		return zero
	}
	return *t
}

package utils

func ToPtr[T any](s T) *T {
	return &s
}

func ToValue[T any](t *T) T {
	if t == nil {
		var zero T
		return zero
	}
	return *t
}

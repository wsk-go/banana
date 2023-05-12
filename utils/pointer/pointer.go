package pointer

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

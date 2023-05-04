package utils

func Map[IN any, OUT any](in []IN, mapper func(IN) OUT) []OUT {
	out := make([]OUT, 0, len(in))

	for i, v := range in {
		out[i] = mapper(v)
	}

	return out
}

func Filter[IN any](in []IN, filter func(IN) bool) []IN {
	out := make([]IN, 0)

	for _, v := range in {
		if filter(v) {
			out = append(out, v)
		}
	}

	return out
}

func ToMap[IN any, KEY comparable](in []IN, keyMapper func(IN) KEY) map[KEY]IN {
	return ToMapWithValue(in, keyMapper, func(in IN) IN {
		return in
	})
}

func ToMapWithValue[IN any, KEY comparable, VALUE any](in []IN, keyMapper func(IN) KEY, valueMapper func(IN) VALUE) map[KEY]VALUE {
	m := make(map[KEY]VALUE)

	for _, v := range in {
		key := keyMapper(v)
		value := valueMapper(v)
		m[key] = value
	}

	return m
}

func Group[IN any, KEY comparable](in []IN, keyMapper func(IN) KEY) map[KEY][]IN {
	return GroupWithValue(in, keyMapper, func(in IN) IN {
		return in
	})
}

func GroupWithValue[IN any, KEY comparable, VALUE any](in []IN, keyMapper func(IN) KEY, valueMapper func(IN) VALUE) map[KEY][]VALUE {
	m := make(map[KEY][]VALUE)

	for _, v := range in {
		key := keyMapper(v)
		m[key] = append(m[key], valueMapper(v))
	}

	return m
}
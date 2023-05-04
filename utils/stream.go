package utils

type sink[T any] struct {
	next   *sink[T]
	accept func(e T, next *sink[T])
}

func (th *sink[T]) Accept(e T) {
	th.accept(e, th.next)
}

type SliceStream[T any] struct {
	//elements []T
	head func()
	last *sink[T]
}

func (th *SliceStream[T]) Filter(filter func(T) bool) *SliceStream[T] {
	th.addSink(func(e T, next *sink[T]) {
		if filter(e) {
			next.Accept(e)
		}
	})
	return th
}

func (th *SliceStream[T]) Map(mapper func(T) T) *SliceStream[T] {
	th.addSink(func(e T, next *sink[T]) {
		next.Accept(mapper(e))
	})
	return th
}

func (th *SliceStream[T]) accept(accept func(T)) {
	th.addSink(func(e T, next *sink[T]) {
		accept(e)
	})

	th.head()
}

func (th *SliceStream[T]) addSink(accept func(e T, next *sink[T])) {
	s := &sink[T]{
		accept: accept,
	}
	th.last.next = s
	th.last = s
}

func Stream[T any](elements []T) *SliceStream[T] {
	sin := &sink[T]{
		accept: func(e T, next *sink[T]) {
			if next != nil {
				next.Accept(e)
			}
		},
	}
	return &SliceStream[T]{
		head: func() {
			for _, element := range elements {
				sin.Accept(element)
			}
		},
		last: sin,
	}
}

//type Collector[IN any, OUT any] struct {
//	stages []func(in IN) OUT
//}
//
//func (th *Collector[IN, OUT]) Map() *Collector[IN, OUT] {
//	th.stages = append(th.stages)
//}

func (th *SliceStream[T]) ToList() []T {
	var out []T

	th.accept(func(in T) {
		out = append(out, in)
	})

	return out
}

func MapStream[IN any, OUT any](s *SliceStream[IN], mapper func(IN) OUT) *SliceStream[OUT] {

	last := &sink[OUT]{
		accept: func(e OUT, next *sink[OUT]) {
			if next != nil {
				next.Accept(e)
			}
		},
	}

	return &SliceStream[OUT]{
		head: func() {
			s.accept(func(in IN) {
				last.Accept(mapper(in))
			})
		},
		last: last,
	}
}

func Map[IN any, OUT any](s *SliceStream[IN], mapper func(IN) OUT) []OUT {
	var out []OUT
	s.accept(func(in IN) {
		out = append(out, mapper(in))
	})

	return out
}

func ToMap[IN any, KEY comparable](s *SliceStream[IN], keyMapper func(IN) KEY) map[KEY]IN {
	return ToMapV(s, keyMapper, func(in IN) IN {
		return in
	})
}

func ToMapV[IN any, KEY comparable, VALUE any](s *SliceStream[IN], keyMapper func(IN) KEY, valueMapper func(IN) VALUE) map[KEY]VALUE {
	m := make(map[KEY]VALUE)
	s.accept(func(in IN) {
		key := keyMapper(in)
		value := valueMapper(in)
		m[key] = value
	})

	return m
}

//func ToMap[IN any, KEY comparable](in []IN, keyMapper func(IN) KEY) map[KEY]IN {
//	return ToMapV(in, keyMapper, func(in IN) IN {
//		return in
//	})
//}

//func ToMapV[IN any, KEY comparable, VALUE any](in []IN, keyMapper func(IN) KEY, valueMapper func(IN) VALUE) map[KEY]VALUE {
//	m := make(map[KEY]VALUE)
//
//	for _, v := range in {
//		key := keyMapper(v)
//		value := valueMapper(v)
//		m[key] = value
//	}
//
//	return m
//}

func Group[IN any, KEY comparable](s *SliceStream[IN], keyMapper func(IN) KEY) map[KEY][]IN {
	return GroupV(s, keyMapper, func(in IN) IN {
		return in
	})
}

func GroupV[IN any, KEY comparable, VALUE any](s *SliceStream[IN], keyMapper func(IN) KEY, valueMapper func(IN) VALUE) map[KEY][]VALUE {
	m := make(map[KEY][]VALUE)

	s.accept(func(in IN) {
		key := keyMapper(in)
		m[key] = append(m[key], valueMapper(in))
	})
	return m
}
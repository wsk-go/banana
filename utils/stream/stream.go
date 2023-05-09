package stream

import "github.com/JackWSK/banana/types/set"

type pipeline[T any] struct {
	next   *pipeline[T]
	accept func(e T, next *pipeline[T])
}

func (th *pipeline[T]) Accept(e T) {
	th.accept(e, th.next)
}

type Stream[T any] struct {
	//elements []T
	head func()
	last *pipeline[T]
}

func (th *Stream[T]) Filter(filter func(T) bool) *Stream[T] {
	th.addPipeline(func(e T, next *pipeline[T]) {
		if filter(e) {
			next.Accept(e)
		}
	})
	return th
}

func (th *Stream[T]) Distinct() *Stream[T] {
	s := set.New[any]()
	th.addPipeline(func(e T, next *pipeline[T]) {
		if !s.Contain(e) {
			s.Add(e)
			next.Accept(e)
		}
	})
	return th
}

func (th *Stream[T]) DistinctByKey(keyMapper func(T) any) *Stream[T] {
	s := set.New[any]()
	th.addPipeline(func(e T, next *pipeline[T]) {
		k := keyMapper(e)
		if !s.Contain(k) {
			s.Add(k)
			next.Accept(e)
		}
	})
	return th
}

func (th *Stream[T]) Map(mapper func(T) T) *Stream[T] {
	th.addPipeline(func(e T, next *pipeline[T]) {
		next.Accept(mapper(e))
	})
	return th
}

func (th *Stream[T]) accept(accept func(T)) {
	th.addPipeline(func(e T, next *pipeline[T]) {
		accept(e)
	})

	th.head()
}

func (th *Stream[T]) addPipeline(accept func(e T, next *pipeline[T])) {
	s := &pipeline[T]{
		accept: accept,
	}
	th.last.next = s
	th.last = s
}

func Of[T any](elements []T) *Stream[T] {
	p := &pipeline[T]{
		accept: func(e T, next *pipeline[T]) {
			if next != nil {
				next.Accept(e)
			}
		},
	}
	return &Stream[T]{
		head: func() {
			for _, element := range elements {
				p.Accept(element)
			}
		},
		last: p,
	}
}

func (th *Stream[T]) ToList() []T {
	var out []T

	th.accept(func(in T) {
		out = append(out, in)
	})

	return out
}

func MapStream[IN any, OUT any](s *Stream[IN], mapper func(IN) OUT) *Stream[OUT] {

	last := &pipeline[OUT]{
		accept: func(e OUT, next *pipeline[OUT]) {
			if next != nil {
				next.Accept(e)
			}
		},
	}

	return &Stream[OUT]{
		head: func() {
			s.accept(func(in IN) {
				last.Accept(mapper(in))
			})
		},
		last: last,
	}
}

func Map[IN any, OUT any](s *Stream[IN], mapper func(IN) OUT) []OUT {
	var out []OUT
	s.accept(func(in IN) {
		out = append(out, mapper(in))
	})

	return out
}

func ToMap[IN any, KEY comparable](s *Stream[IN], keyMapper func(IN) KEY) map[KEY]IN {
	return ToMapV(s, keyMapper, func(in IN) IN {
		return in
	})
}

func ToMapV[IN any, KEY comparable, VALUE any](s *Stream[IN], keyMapper func(IN) KEY, valueMapper func(IN) VALUE) map[KEY]VALUE {
	m := make(map[KEY]VALUE)
	s.accept(func(in IN) {
		key := keyMapper(in)
		value := valueMapper(in)
		m[key] = value
	})

	return m
}

func Group[IN any, KEY comparable](s *Stream[IN], keyMapper func(IN) KEY) map[KEY][]IN {
	return GroupV(s, keyMapper, func(in IN) IN {
		return in
	})
}

func GroupV[IN any, KEY comparable, VALUE any](s *Stream[IN], keyMapper func(IN) KEY, valueMapper func(IN) VALUE) map[KEY][]VALUE {
	m := make(map[KEY][]VALUE)

	s.accept(func(in IN) {
		key := keyMapper(in)
		m[key] = append(m[key], valueMapper(in))
	})
	return m
}

package set

import "sync"

type Set[V comparable] interface {
	Add(v ...V)

	Remove(v V)

	Contain(v V) bool

	Iter(f func(V))

	Size() int
}

type UnsafeSet[V comparable] struct {
	m map[V]struct{}
}

func New[V comparable]() Set[V] {
	m := make(map[V]struct{})
	return &UnsafeSet[V]{
		m: m,
	}
}

func (th *UnsafeSet[V]) Add(v ...V) {
	for _, vv := range v {
		if !th.contain(vv) {
			th.m[vv] = struct{}{}
		}
	}
}

func (th *UnsafeSet[V]) Remove(v V) {
	if th.contain(v) {
		delete(th.m, v)
	}
}

func (th *UnsafeSet[V]) Contain(v V) bool {
	return th.contain(v)
}

func (th *UnsafeSet[V]) Iter(f func(V)) {
	for k, _ := range th.m {
		f(k)
	}
}

func (th *UnsafeSet[V]) Size() int {
	return len(th.m)
}

func (th *UnsafeSet[V]) contain(v V) bool {
	_, ok := th.m[v]
	return ok
}

// SafeSet Thread Safe
type SafeSet[V comparable] struct {
	m map[V]struct{}
	l sync.RWMutex
}

func NewSafe[V comparable]() Set[V] {
	m := make(map[V]struct{})
	return &SafeSet[V]{
		m: m,
	}
}

func (th *SafeSet[V]) Add(v ...V) {
	th.l.Lock()
	defer th.l.Unlock()
	for _, vv := range v {
		if !th.contain(vv) {
			th.m[vv] = struct{}{}
		}
	}
}

func (th *SafeSet[V]) Remove(v V) {
	th.l.Lock()
	defer th.l.Unlock()
	if th.contain(v) {
		delete(th.m, v)
	}
}

func (th *SafeSet[V]) Contain(v V) bool {
	th.l.RLock()
	defer th.l.RUnlock()
	return th.contain(v)
}

func (th *SafeSet[V]) Iter(f func(V)) {
	th.l.RLock()
	defer th.l.RUnlock()
	for k, _ := range th.m {
		f(k)
	}
}

func (th *SafeSet[V]) Size() int {
	return len(th.m)
}

func (th *SafeSet[V]) contain(v V) bool {
	_, ok := th.m[v]
	return ok
}

package types

import "sync"

type Set[V comparable] struct {
	m map[V]struct{}
	l sync.RWMutex
}

func New[V comparable]() *Set[V] {
	m := make(map[V]struct{})
	return &Set[V]{
		m: m,
	}
}

func (th *Set[V]) Add(v ...V) {
	th.l.Lock()
	defer th.l.Unlock()
	for _, vv := range v {
		if !th.contain(vv) {
			th.m[vv] = struct{}{}
		}
	}
}

func (th *Set[V]) Remove(v V) {
	th.l.Lock()
	defer th.l.Unlock()
	if th.contain(v) {
		delete(th.m, v)
	}
}

func (th *Set[V]) Contain(v V) bool {
	th.l.RLock()
	defer th.l.RUnlock()
	return th.contain(v)
}

func (th *Set[V]) Iter(f func(V)) {
	th.l.RLock()
	defer th.l.RUnlock()
	for k, _ := range th.m {
		f(k)
	}
}

func (th *Set[V]) contain(v V) bool {
	_, ok := th.m[v]
	return ok
}

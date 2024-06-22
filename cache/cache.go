package cache

import "sync"

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{
		m: map[string]T{},
		L: sync.Mutex{},
	}
}

type Cache[T any] struct {
	m map[string]T
	L sync.Mutex
}

func (s *Cache[T]) Load(key string) (T, bool) {
	s.L.Lock()
	defer s.L.Unlock()
	v, ok := s.m[key]
	return v, ok
}

func (s *Cache[T]) Store(key string, val T) {
	s.L.Lock()
	defer s.L.Unlock()
	s.m[key] = val
}

func (s *Cache[T]) Clear() {
	s.L.Lock()
	defer s.L.Unlock()
	s.m = make(map[string]T)
}

func (s *Cache[T]) Size() int {
	return len(s.m)
}

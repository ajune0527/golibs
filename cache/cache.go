package cache

import (
	"sync"
	"time"
)

type Option[T any] func(*Cache[T])

func WithTicker[T any](duration time.Duration) Option[T] {
	return func(c *Cache[T]) {
		if c.ticker != nil {
			c.ticker.Stop()
		}
		c.ticker = time.NewTicker(duration)
		go func() {
			for {
				select {
				case <-c.ticker.C:
					c.Clear() // 定时清空缓存
				case <-c.stopCh:
					return // 停止 ticker
				}
			}
		}()
	}
}

func NewCache[T any](opts ...Option[T]) *Cache[T] {
	cache := &Cache[T]{
		m:      sync.Map{},
		stopCh: make(chan struct{}),
	}
	for _, opt := range opts {
		opt(cache)
	}
	return cache
}

func (s *Cache[T]) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
		close(s.stopCh)
	}
}

type Cache[T any] struct {
	m      sync.Map
	ticker *time.Ticker
	stopCh chan struct{}
}

func (s *Cache[T]) Store(key string, val T) {
	s.m.Store(key, val)
}

func (s *Cache[T]) Load(key string) (T, bool) {
	value, ok := s.m.Load(key)
	if !ok {
		var zero T
		return zero, false
	}
	return value.(T), true
}

func (s *Cache[T]) Clear() {
	s.m.Range(func(key, value interface{}) bool {
		s.m.Delete(key)
		return true
	})
}

func (s *Cache[T]) Size() int {
	var size int
	s.m.Range(func(key, value interface{}) bool {
		size++
		return true
	})

	return size
}

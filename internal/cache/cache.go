package cache

import (
	"sync"
	"time"
)

type CacheItem[T any] struct {
	Value      T
	Expiration time.Time
}

func (item CacheItem[T]) IsExpired() bool {
	if item.Expiration.IsZero() {
		return false
	}
	return time.Now().After(item.Expiration)
}

type Cache[T any] struct {
	items map[string]CacheItem[T]
	mu    sync.RWMutex
	ttl   time.Duration
}

func NewCache[T any](ttl time.Duration) *Cache[T] {
	c := &Cache[T]{
		items: make(map[string]CacheItem[T]),
		ttl:   ttl,
	}

	go c.startCleanup()

	return c
}

func (c *Cache[T]) Set(key string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiration time.Time
	if c.ttl > 0 {
		expiration = time.Now().Add(c.ttl)
	}

	c.items[key] = CacheItem[T]{
		Value:      value,
		Expiration: expiration,
	}
}

func (c *Cache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		var zero T
		return zero, false
	}

	if item.IsExpired() {
		var zero T
		return zero, false
	}

	return item.Value, true
}

func (c *Cache[T]) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *Cache[T]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]CacheItem[T])
}

func (c *Cache[T]) DeleteByPrefix(prefix string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key := range c.items {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			delete(c.items, key)
		}
	}
}

func (c *Cache[T]) startCleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanup()
	}
}

func (c *Cache[T]) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, item := range c.items {
		if item.IsExpired() {
			delete(c.items, key)
		}
	}
}

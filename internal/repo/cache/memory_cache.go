package cache

import (
	"sync"
	"time"
)

type memItem struct {
	val string
	exp time.Time
}

type MemoryCache struct {
	mu sync.RWMutex
	m  map[string]memItem
}

func NewMemoryCache() *MemoryCache {
	c := &MemoryCache{m: map[string]memItem{}}
	go c.gcLoop()
	return c
}

func (c *MemoryCache) gcLoop() {
	t := time.NewTicker(30 * time.Second)
	defer t.Stop()

	for range t.C {
		now := time.Now()
		c.mu.Lock()
		for k, val := range c.m {
			if !val.exp.IsZero() && now.After(val.exp) {
				delete(c.m, k)
			}
		}
		c.mu.Unlock()
	}
}

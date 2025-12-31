package cache

import (
	"context"
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

func (c *MemoryCache) Get(ctx context.Context, key string) (string, bool, error) {
	c.mu.RLock()
	it, ok := c.m[key]
	c.mu.RUnlock()

	if !ok {
		return "", false, nil
	}
	if !it.exp.IsZero() && time.Now().After(it.exp) {
		_ = c.Del(ctx, key)
		return "", false, nil
	}
	return it.val, true, nil
}

func (c *MemoryCache) Set(ctx context.Context, key string, val string, ttl time.Duration) error {
	exp := time.Time{}
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}
	c.mu.Lock()
	c.m[key] = memItem{val: val, exp: exp}
	c.mu.Unlock()
	return nil
}

func (c *MemoryCache) Del(ctx context.Context, key string) error {
	c.mu.Lock()
	delete(c.m, key)
	c.mu.Unlock()
	return nil
}

func (c *MemoryCache) Incr(ctx context.Context, key string, by int64, ttl time.Duration) (int64, error) {
	// simple int64 stored as string
	c.mu.Lock()
	defer c.mu.Unlock()

	it, ok := c.m[key]
	if ok && !it.exp.IsZero() && time.Now().After(it.exp) {
		ok = false
	}
	var cur int64 = 0
	if ok {
		// parse
		var n int64
		for i := 0; i < len(it.val); i++ {
			ch := it.val[i]
			if ch < '0' || ch > '9' {
				n = 0
				break
			}
			n = n*10 + int64(ch-'0')
		}
		cur = n
	}
	cur += by
	exp := time.Time{}
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}
	c.m[key] = memItem{val: itoa(cur), exp: exp}
	return cur, nil
}

func (c *MemoryCache) Publish(ctx context.Context, channel, msg string) error {
	// no-op in memory
	return nil
}

func itoa(n int64) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 20)
	for n > 0 {
		buf = append(buf, byte('0'+(n%10)))
		n /= 10
	}
	// reverse
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}

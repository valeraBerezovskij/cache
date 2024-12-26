package cache

import (
	"context"
	"sync"
	"time"
)

type Cache struct {
	cache map[string]interface{}
	mu    *sync.Mutex
}

func New() *Cache {
	return &Cache{cache: make(map[string]interface{}), mu: &sync.Mutex{}}
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	deadline := time.Now().Add(ttl)
	ctx, _ := context.WithDeadline(context.Background(), deadline)

	c.cache[key] = value

	go c.living(ctx, key)
}

func (c *Cache) living(ctx context.Context, key string) {
	<-ctx.Done()
	if ctx.Err() == context.DeadlineExceeded {
		c.mu.Lock()
		c.Delete(key)
		c.mu.Unlock()
	}
}

func (c *Cache) Get(key string) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cache[key]
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
}
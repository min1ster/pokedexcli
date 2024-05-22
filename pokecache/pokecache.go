package pokecache 

import (
	"time"
	"sync"
)

type Cache struct {
	entries map[string]cacheEntry
	mu sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val: val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	entries := c.entries
	entry, ok := entries[key]
	if ok {
		return entry.val, ok
	}
	return []byte{}, ok
}

func (c *Cache) readLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		<-ticker.C
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.entries {
			if now.Sub(entry.createdAt) > interval {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		entries: make(map[string]cacheEntry),
	}
	go cache.readLoop(interval)
	return cache
}
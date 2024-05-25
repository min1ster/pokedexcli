package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Entries map[string]CacheEntry
	mu      sync.Mutex
}

type CacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Entries[key] = CacheEntry{
		CreatedAt: time.Now(),
		Val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	entries := c.Entries
	entry, ok := entries[key]
	if ok {
		return entry.Val, ok
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
		for key, entry := range c.Entries {
			if now.Sub(entry.CreatedAt) > interval {
				delete(c.Entries, key)
			}
		}
		c.mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		Entries: make(map[string]CacheEntry),
	}
	go cache.readLoop(interval)
	return cache
}

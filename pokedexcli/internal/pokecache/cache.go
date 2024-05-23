package pokecache

import (
	"log"
	"sync"
	"time"
)

type Cache struct {
	cache    map[string]cacheEntry
	interval time.Duration
	ticker   *time.Ticker
	lock     sync.RWMutex
}

type cacheEntry struct {
	data      []byte
	createdAt time.Time
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cache:    make(map[string]cacheEntry),
		interval: interval,
		ticker:   time.NewTicker(interval),
	}

	// Start the reap loop in a separate goroutine
	go func() {
		defer c.ticker.Stop()
		for range c.ticker.C {
			c.reap()
		}
	}()

	return c
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	entry, ok := c.cache[key]
	if !ok {
		log.Printf("%s: CACHE MISS", key)
		return nil, false
	}

	log.Printf("%s: CACHE HIT", key)
	return entry.data, true
}

func (c *Cache) Add(key string, data []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.cache[key] = cacheEntry{
		data:      data,
		createdAt: time.Now(),
	}

	log.Printf("Cache: Added key %s", key)
}

func (c *Cache) reap() {
	c.lock.Lock()
	defer c.lock.Unlock()

	log.Println("Cache: Running reap loop")

	cutoff := time.Now().Add(-c.interval)
	for key, entry := range c.cache {
		if entry.createdAt.Before(cutoff) {
			delete(c.cache, key)
			log.Printf("Cache: Entry %s reaped", key)
		}
	}
}

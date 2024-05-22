package pokecache

import (
	"log"
	"sync"
	"time"
)

type PokeCache interface {
	NewCache(interval time.Duration) *Cache
	Get(key string) ([]byte, bool)
	Add(key string, data []byte)
}

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
	log.Default().Println("Cache: Creating New Cache")
	c := &Cache{
		cache:    make(map[string]cacheEntry),
		interval: interval,
		ticker:   time.NewTicker(interval),
		lock:     sync.RWMutex{},
	}
	go c.reapLoop()
	return c
}

func (c *Cache) Get(key string) ([]byte, bool) {
	log.Default().Printf("Cache: get key:%v", key)
	c.lock.RLock()
	defer c.lock.RUnlock()
	entry, ok := c.cache[key]
	if !ok {
		log.Default().Println("\t: CACHE MISS!")
		return nil, false
	}
	log.Default().Println("\t: CACHE HIT!")
	return entry.data, true
}

func (c *Cache) Add(key string, data []byte) {
	c.lock.Lock()
	log.Default().Printf("Cache: Adding\n\tkey:%v\n", key)
	defer c.lock.Unlock()
	c.cache[key] = cacheEntry{
		data:      data,
		createdAt: time.Now(),
	}
}

func (c *Cache) reapLoop() {
	for range c.ticker.C {
		c.lock.Lock()
		log.Default().Println("Cache: running reap loop now!")

		cutoff := time.Now().Add(-c.interval)
		for k, v := range c.cache {
			if v.createdAt.Before(cutoff) {
				delete(c.cache, k)
			}
		}
		c.lock.Unlock()
	}
	c.ticker.Stop()
	log.Println("reapLoop stopped")
}

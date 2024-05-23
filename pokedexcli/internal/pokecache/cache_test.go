package pokecache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}

func TestCache(t *testing.T) {
	cacheDuration := 100 * time.Millisecond
	c := NewCache(cacheDuration)

	t.Run("Add and Get cache entry", func(t *testing.T) {
		key := "testKey"
		expectedData := []byte("testData")
		c.Add(key, expectedData)

		data, ok := c.Get(key)
		if !ok {
			t.Fatalf("expected key %s to be present in cache", key)
		}
		if string(data) != string(expectedData) {
			t.Fatalf("expected data %s, got %s", expectedData, data)
		}
	})

	t.Run("Cache miss", func(t *testing.T) {
		key := "missingKey"
		_, ok := c.Get(key)
		if ok {
			t.Fatalf("expected key %s to be absent in cache", key)
		}
	})

	t.Run("Reap old entries", func(t *testing.T) {
		key1 := "oldKey"
		key2 := "newKey"
		data := []byte("data")

		c.Add(key1, data)
		time.Sleep(cacheDuration + 10*time.Millisecond)
		c.Add(key2, data)

		time.Sleep(cacheDuration + 10*time.Millisecond)

		_, ok := c.Get(key1)
		if ok {
			t.Fatalf("expected key %s to be reaped from cache", key1)
		}

		stillInCacheData, stillInCache := c.Get(key2)
		if !stillInCache {
			t.Fatalf("expected key %s to be present in cache", key2)
		}
		if string(stillInCacheData) != string(data) {
			t.Fatalf("expected data %s, got %s", data, stillInCacheData)
		}
	})

	t.Run("Concurrent access to cache", func(t *testing.T) {
		var wg sync.WaitGroup
		key := "concurrentKey"
		data := []byte("concurrentData")
		c.Add(key, data)

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				c.Get(key)
			}()
		}

		wg.Wait()

		retrievedData, ok := c.Get(key)
		if !ok {
			t.Fatalf("expected key %s to be present in cache", key)
		}
		if string(retrievedData) != string(data) {
			t.Fatalf("expected data %s, got %s", data, retrievedData)
		}
	})
}

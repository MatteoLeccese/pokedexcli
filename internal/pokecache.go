package pokecache

import (
	"time"
	"sync"
)

// CacheEntry struct to hold the cache entry data
type cacheEntry struct {
	createdAt 	time.Time
	val 				[]byte
}

// Cache struct to hold the cache entries
type Cache struct {
	interval 		time.Duration
	mutex				sync.Mutex
	entries  		map[string]cacheEntry
}

// Adding a method to add an entry to the cache
func (c Cache) Add(key string, val []byte) {
	// Lock the mutex to prevent concurrent access
	c.mutex.Lock();

	// Add the entry to the cache
	c[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}

	// Unlock the mutex
	c.mutex.Unlock();
}

// Addig a method to get the value from the cache
func (c Cache) Get(key string) ([]byte, bool) {
	if entry, ok := c[key]; ok {
		return entry.val, true;
	}

	return nil, false;
}

// Adding a method to remove expired entries
func (c Cache) reapLoop() {
	// Creating a ticker to check for expired entries
	ticker := time.NewTicker(c.interval);

	// Defer the ticker stop
	defer ticker.Stop()

	done := make(chan bool);

	// Spawning a goroutine to check for expired entries
	go func() {
		time.Sleep(c.interval);
		done <- true;
	}();

	for {
		select {
			case <-done:
				return;
			case <-ticker.C:
				// Lock the mutex to prevent concurrent access
				c.mutex.Lock();

				// Iterate over the cache entries
				for key, entry := range c.entries {
					// Check if the entry is expired
					if time.Since(entry.createdAt) > c.interval {
						// Remove the expired entry
						delete(c.entries, key);
					}
				}

				// Unlock the mutex
				c.mutex.Unlock();
		}
	}
}

// Function to create a new cache
func NewCache(interval time.Duration) *Cache {
	// Create a new cache with the specified interval
	cache := &Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
		// mutex is initialized to its zero value, which is an unlocked mutex
	}

	// Start the cache reaper loop
	go cache.reapLoop();

	// Return the new cache
	return cache;
}

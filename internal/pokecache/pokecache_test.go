// filepath: internal/pokecache/pokecache_test.go
package pokecache

import (
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	cache := NewCache(time.Millisecond * 100)
	key := "testKey"
	value := []byte("testValue")

	// Test case 1: Key exists in cache
	cache.Add(key, value)
	retrievedValue, ok := cache.Get(key)

	if !ok {
		t.Errorf("Expected key to be found in cache")
	}

	if string(retrievedValue) != string(value) {
		t.Errorf("Expected value %s, but got %s", value, retrievedValue)
	}

	// Test case 2: Key does not exist in cache
	nonExistentKey := "nonExistentKey"
	_, ok = cache.Get(nonExistentKey)

	if ok {
		t.Errorf("Expected key to not be found in cache")
	}

	// Test case 3: Key exists but is expired
	time.Sleep(time.Millisecond * 200) // Wait for the entry to expire
	_, ok = cache.Get(key)

	if ok {
		t.Errorf("Expected key to not be found in cache after expiration")
	}
}

func TestAdd(t *testing.T) {
	cache := NewCache(time.Millisecond * 100)
	key := "testKey"
	value := []byte("testValue")

	// Test case 1: Add a new entry to the cache
	cache.Add(key, value)
	retrievedValue, ok := cache.Get(key)

	if !ok {
		t.Errorf("Expected key to be found in cache")
	}

	if string(retrievedValue) != string(value) {
		t.Errorf("Expected value %s, but got %s", value, retrievedValue)
	}

	// Test case 2: Add an existing entry to the cache with a new value
	newValue := []byte("newValue")
	cache.Add(key, newValue)
	retrievedValue, ok = cache.Get(key)

	if !ok {
		t.Errorf("Expected key to be found in cache")
	}

	if string(retrievedValue) != string(newValue) {
		t.Errorf("Expected value %s, but got %s", newValue, retrievedValue)
	}
}

func TestNewCache(t *testing.T) {
	interval := time.Millisecond * 100
	cache := NewCache(interval)

	// Test case 1: Check if the cache is initialized with the correct interval
	if cache.interval != interval {
		t.Errorf("Expected interval %v, but got %v", interval, cache.interval)
	}

	// Test case 2: Check if the cache entries map is initialized
	if cache.entries == nil {
		t.Errorf("Expected cache entries map to be initialized")
	}
}

func TestReapLoop(t *testing.T) {
	interval := time.Millisecond * 100
	cache := NewCache(interval)
	key := "testKey"
	value := []byte("testValue")

	cache.Add(key, value)

	// Wait for the entry to expire
	time.Sleep(time.Millisecond * 200)

	// Check if the entry is removed from the cache
	_, ok := cache.Get(key)

	if ok {
		t.Errorf("Expected key to not be found in cache after expiration")
	}
}
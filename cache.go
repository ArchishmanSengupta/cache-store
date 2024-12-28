/*
	Purpose: Implement a thread-safe in-memory cache with optional expiration times.

	This CacheStore struct manages cached items using a sync.Map for concurrency safety.
	It supports operations like setting values with expiration, retrieving values,
	iterating over cache items, and removing specific items.

	Background:
	- CacheStore is initialized with a cleaningInterval, which determines how often expired items are removed.
	- A cleaning goroutine periodically checks and removes expired items based on their expiration time.

	Operations:
	- Get retrieves a cached value for a given key, ensuring it's not expired before returning.
	- Set sets a key-value pair in the cache with an optional expiration duration.
	- Iterate allows iterating over non-expired items in the cache using a provided function.
	- RemoveKey removes a specific key-value pair from the cache.
	- CloseCacheStore stops the cleaning process and clears the cache.

	Thread Safety:
	- Utilizes sync.Map to ensure thread safety for concurrent access to cache operations.
	- Cleaning operations and cache access are synchronized to prevent data races.

	Error Handling:
	- Errors are returned for nil keys or values in SetValue and GetValue.
	- Error returned if the cleaning interval is non-positive during initialization.

*/

package cache

import (
	"context"
	"errors"
	"sync"
	"time"
)

// CacheStore represents a thread-safe cache store with optional expiration times.
type CacheStore struct {
	items            sync.Map // Map for storing cached items (key-value pairs)
	ctx              context.Context
	cancel           context.CancelFunc
	cleaningInterval time.Duration // Interval for cleaning expired items
}

// CachedItem represents an item stored in the cache with associated data and expiration time.
type CachedItem struct {
	data    interface{} // Data stored in the cache
	expires int64       // Expiration time in nanoseconds since Unix epoch
}

// NewCacheStore creates a new CacheStore instance with a specified cleaning interval.
func NewCacheStore(cleaningInterval time.Duration) (*CacheStore, error) {
	if cleaningInterval <= 0 {
		return nil, errors.New("cleaning interval must be positive")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cacheStore := &CacheStore{
		ctx:              ctx,
		cancel:           cancel,
		cleaningInterval: cleaningInterval,
	}

	go cacheStore.cleanupExpiredItems()

	return cacheStore, nil
}

// cleanupExpiredItems periodically removes expired items from the cache.
func (cacheStore *CacheStore) cleanupExpiredItems() {
	ticker := time.NewTicker(cacheStore.cleaningInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			currentTime := time.Now().UnixNano()
			cacheStore.items.Range(func(key, value interface{}) bool {
				item := value.(CachedItem)
				if item.expires > 0 && currentTime > item.expires {
					cacheStore.items.Delete(key)
				}
				return true
			})
		case <-cacheStore.ctx.Done():
			return
		}
	}
}

// GetValue retrieves a value from the cache for a given key.
func (cacheStore *CacheStore) Get(key interface{}) (interface{}, bool, error) {
	if key == nil {
		return nil, false, errors.New("key cannot be nil")
	}

	obj, exists := cacheStore.items.Load(key)
	if !exists {
		return nil, false, nil
	}

	item := obj.(CachedItem)
	if item.expires > 0 && time.Now().UnixNano() > item.expires {
		cacheStore.items.Delete(key)
		return nil, false, nil
	}

	return item.data, true, nil
}

// SetValue sets a value in the cache for a given key with an optional expiration duration.
func (cacheStore *CacheStore) Set(key interface{}, value interface{}, duration time.Duration) error {
	if key == nil {
		return errors.New("key cannot be nil")
	}
	if value == nil {
		return errors.New("value cannot be nil")
	}

	var expires int64
	if duration > 0 {
		expires = time.Now().Add(duration).UnixNano()
	}

	cacheStore.items.Store(key, CachedItem{
		data:    value,
		expires: expires,
	})

	return nil
}

// Iterate iterates over all non-expired items in the cache and applies the given function.
func (cacheStore *CacheStore) Iterate(f func(key, value interface{}) bool) error {
	if f == nil {
		return errors.New("function cannot be nil")
	}

	currentTime := time.Now().UnixNano()
	fn := func(key, value interface{}) bool {
		item := value.(CachedItem)
		if item.expires > 0 && currentTime > item.expires {
			cacheStore.items.Delete(key)
			return true
		}
		return f(key, item.data)
	}

	cacheStore.items.Range(fn)
	return nil
}

// RemoveKey removes a key-value pair from the cache.
func (cacheStore *CacheStore) RemoveKey(key interface{}) error {
	if key == nil {
		return errors.New("key cannot be nil")
	}

	cacheStore.items.Delete(key)
	return nil
}

// CloseCacheStore stops the cache cleaning process and clears all stored items.
func (cacheStore *CacheStore) CloseCacheStore() {
	cacheStore.cancel()
	cacheStore.items = sync.Map{}
}

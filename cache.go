/*
	The purpose of this code is to implement a thread-safe in-memory cache with support
	for optional expiration times for stored items.

	This cache can store any type of data and automatically remove expired items at specified intervals.
*/

package cache

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

// CacheStore represents a cache store that holds arbitrary data with expiration time.
type CacheStore struct {
	items            sync.Map // for thread-safe
	ctx              context.Context
	cancel           context.CancelFunc
	cleaningInterval time.Duration
}

// CachedItem represents an item in the cache store with arbitrary data and expiration time.
type CachedItem struct {
	data    interface{}
	expires int64
}

// NewCacheStore creates a new cache store asynchronously cleans expired entries after the given time passes
func NewCacheStore(cleaningInterval time.Duration) *CacheStore {
	// Validate that the cleaningInterval is positive; logs a fatal error and stops the program if not.
	if cleaningInterval <= 0 {
		log.Fatal("cleaning interval must be positive")
	}

	// Create a cancellable context using context.WithCancel.
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize a CacheStore with the context, cancel function, and cleaning interval.
	cacheStore := &CacheStore{
		ctx:              ctx,
		cancel:           cancel,
		cleaningInterval: cleaningInterval,
	}

	// Start a goroutine that runs the cleanupExpiredItems method to periodically remove expired items.
	go cacheStore.cleanupExpiredItems()

	// Return the initialized CacheStore.
	return cacheStore
}

// Go routine to periodically clean up expired items
func (cacheStore *CacheStore) cleanupExpiredItems() {
	// create a time.Ticker that ticks at intervals specified by 'cleaning Intervals'
	ticker := time.NewTicker(cacheStore.cleaningInterval)

	// Stop the ticker after the function exits
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// on each tick it gets the current time in nanoseconds
			currentTime := time.Now().UnixNano()

			// iterate over all the items in the cache, checking each item's expiration time
			cacheStore.items.Range(func(key, value interface{}) bool {
				item := value.(CachedItem)

				// if an item is expired it deletes the item from the cache
				if item.expires > 0 && currentTime > item.expires {
					cacheStore.items.Delete(key)
				}
				return true
			})
		}
	}
}

// Retrieves value for a given key
func (cacheStore *CacheStore) GetValue(key interface{}) (interface{}, bool, error) {
	if key == nil {
		return nil, false, errors.New("key cannot be nil")
	}

	obj, exists := cacheStore.items.Load(key)

	if !exists {
		return nil, false, nil
	}

	// type asserts the loaded item
	item := obj.(CachedItem)

	// check if the item is expired, if yes delete it
	if item.expires > 0 && time.Now().UnixNano() > item.expires {
		cacheStore.items.Delete(key)
		return nil, false, nil
	}

	// if not expired, return the item's data
	return item.data, true, nil
}

// set the value for a given key with optional expiration duration
func (cacheStore *CacheStore) SetValue(key interface{}, value interface{}, duration time.Duration) error {
	if key == nil {
		return errors.New("key cannot be nil")
	}

	if value == nil {
		return errors.New("value cannot be nil")
	}

	var expires int64

	// calculate expiration timestamp
	if duration > 0 {
		expires = time.Now().Add(duration).UnixNano()
	}

	// store the kv pair in the cache
	cacheStore.items.Store(key, CachedItem{
		data:    value,
		expires: expires,
	})

	return nil
}

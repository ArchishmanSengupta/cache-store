/*
	The purpose of this code is to implement a thread-safe in-memory cache with support
	for optional expiration times for stored items.

	This cache can store any type of data and automatically remove expired items at specified intervals.
*/

package cache

import (
	"context"
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
	// Validates that the cleaningInterval is positive; logs a fatal error and stops the program if not.
	if cleaningInterval <= 0 {
		log.Fatal("cleaning interval must be positive")
	}

	// Creates a cancellable context using context.WithCancel.
	ctx, cancel := context.WithCancel(context.Background())

	// Initializes a CacheStore with the context, cancel function, and cleaning interval.
	cacheStore := &CacheStore{
		ctx:              ctx,
		cancel:           cancel,
		cleaningInterval: cleaningInterval,
	}

	// Starts a goroutine that runs the cleanupExpiredItems method to periodically remove expired items.
	// go cacheStore.cleanupExpiredItems();

	// Returns the initialized CacheStore.
	return cacheStore
}

// Go routine to periodically clean up expired items
func (cacheStore *CacheStore) cleanupExpiredItems() {
	// creates a time.Ticker that ticks at intervals specified by 'cleaning Intervals'
	ticker := time.NewTicker(cacheStore.cleaningInterval)

	// Stop the ticker after the function exits
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// on each tick it gets the current time in nanoseconds
			currentTime := time.Now().UnixNano()

			// iterates over all the items in the cache, checking each item's expiration time
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

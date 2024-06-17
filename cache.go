package cache

import (
	"context"
	"log"
	"sync"
	"time"
)

// CacheStore represents a cache store that holds arbitrary data with expiration time.
type CacheStore struct {
	items            sync.Map
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
	if cleaningInterval <= 0 {
		log.Fatal("cleaning interval must be positive")
	}

	ctx, cancel := context.WithCancel(context.Background())

	cacheStore := &CacheStore{
		ctx:              ctx,
		cancel:           cancel,
		cleaningInterval: cleaningInterval,
	}

	//go cacheStore.cleanupExpiredItems();

	return cacheStore
}

package cache

import (
	"context"
	"sync"
)

// CacheStore represents a cache store that holds arbitrary data with expiration time.
type CacheStore struct {
	items  sync.Map
	ctx    context.Context
	cancel context.CancelFunc
}

package cache

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkCache_SetValue(b *testing.B) {
	cache, _ := NewCacheStore(time.Second * 1)
	defer cache.CloseCacheStore()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		duration := time.Second * 5 // Example expiration duration

		b.Run(fmt.Sprintf("Set-%d", i), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				cache.SetValue(key, value, duration)
			}
		})
	}
}

func BenchmarkCache_GetValue(b *testing.B) {
	cache, _ := NewCacheStore(time.Second * 1)
	defer cache.CloseCacheStore()

	key := "key"
	value := "value"
	duration := time.Second * 5

	cache.SetValue(key, value, duration)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.Run(fmt.Sprintf("Get-%d", i), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				_, _, _ = cache.GetValue(key)
			}
		})
	}
}

func BenchmarkCache_Iterate(b *testing.B) {
	cache, _ := NewCacheStore(time.Second * 1)
	defer cache.CloseCacheStore()

	// Populate the cache with some items for iteration
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		cache.SetValue(key, value, 0)
	}

	b.ResetTimer()

	// Run the benchmark for Iterate
	b.Run("Iterate", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := cache.Iterate(func(key, value interface{}) bool {
				// Do nothing in the callback function for benchmarking purposes
				return true
			})
			if err != nil {
				b.Fatalf("Error during iteration: %v", err)
			}
		}
	})
}

package cache

import (
	"testing"
	"time"
)

func TestCacheStore(t *testing.T) {
	// Test NewCacheStore function
	t.Run("TestNewCacheStore", func(t *testing.T) {
		// Valid interval
		_, err := NewCacheStore(time.Second)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Invalid interval
		_, err = NewCacheStore(0)
		if err == nil || err.Error() != "cleaning interval must be positive" {
			t.Errorf("Expected error 'cleaning interval must be positive', got %v", err)
		}
	})

	// Initialize a cache store for subsequent tests
	cacheStore, _ := NewCacheStore(time.Second)

	// Test SetValue and GetValue
	t.Run("TestSetValueGetValue", func(t *testing.T) {
		key := "test_key"
		value := "test_value"
		expiration := time.Millisecond * 500

		err := cacheStore.SetValue(key, value, expiration)
		if err != nil {
			t.Errorf("Expected no error on SetValue, got %v", err)
		}

		// Retrieve the value
		v, exists, err := cacheStore.GetValue(key)
		if err != nil {
			t.Errorf("Expected no error on GetValue, got %v", err)
		}
		if !exists {
			t.Error("Expected key to exist, but it doesn't")
		}
		if v != value {
			t.Errorf("Expected value %v, got %v", value, v)
		}

		// Wait for expiration
		time.Sleep(time.Millisecond * 600)

		// Check if the value expired
		_, exists, err = cacheStore.GetValue(key)
		if exists {
			t.Error("Expected key to be expired, but it still exists")
		}
	})

	// Test Iterate
	t.Run("TestIterate", func(t *testing.T) {
		// Add a couple of items
		cacheStore.SetValue("key1", "value1", 0)
		cacheStore.SetValue("key2", "value2", 0)

		// Iterate and ensure all items are iterated over
		err := cacheStore.Iterate(func(key, value interface{}) bool {
			t.Logf("Iterating over key: %v, value: %v", key, value)
			return true
		})
		if err != nil {
			t.Errorf("Expected no error on Iterate, got %v", err)
		}
	})

	// Test RemoveKey
	t.Run("TestRemoveKey", func(t *testing.T) {
		// Add a key
		cacheStore.SetValue("key_to_remove", "value", 0)

		// Remove existing key
		err := cacheStore.RemoveKey("key_to_remove")
		if err != nil {
			t.Errorf("Expected no error on RemoveKey, got %v", err)
		}

		// Try to remove non-existent key
		err = cacheStore.RemoveKey("non_existent_key")
		if err != nil {
			t.Errorf("Expected no error on RemoveKey for non-existent key, got %v", err)
		}

		// Try to remove with nil key
		err = cacheStore.RemoveKey(nil)
		if err == nil || err.Error() != "key cannot be nil" {
			t.Errorf("Expected error 'key cannot be nil', got %v", err)
		}
	})

	// Test CloseCacheStore
	t.Run("TestCloseCacheStore", func(t *testing.T) {
		cacheStore.CloseCacheStore()

		// Ensure all items are cleared
		err := cacheStore.Iterate(func(key, value interface{}) bool {
			t.Error("Expected cache to be empty, but found items")
			return true // Stop iteration immediately upon finding an item
		})
		if err != nil {
			t.Errorf("Expected no error on Iterate, got %v", err)
		}
	})
}

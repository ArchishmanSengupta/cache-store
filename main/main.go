package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/archishmansengupta/cache"
)

func main() {
	// Create a new cache store with a cleaning interval of 1 minute
	cacheStore, err := cache.NewCacheStore(time.Minute)
	if err != nil {
		fmt.Println("Error creating cache store:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter key: ")
	key, _ := reader.ReadString('\n')
	key = strings.TrimSpace(key)

	fmt.Print("Enter value: ")
	value, _ := reader.ReadString('\n')
	value = strings.TrimSpace(value)

	fmt.Print("Enter expiration time (seconds): ")
	expiration, _ := reader.ReadString('\n')
	expiration = strings.TrimSpace(expiration)

	duration, err := time.ParseDuration(expiration + "s")
	if err != nil {
		fmt.Println("Error parsing expiration time:", err)
		return
	}

	// Set a value in the cache with the key and an expiration time
	err = cacheStore.SetValue(key, value, duration)
	if err != nil {
		fmt.Println("Error setting cache value:", err)
		return
	}

	// Get the value from the cache
	cachedValue, found, err := cacheStore.GetValue(key)
	if err != nil {
		fmt.Println("Error getting cache value:", err)
		return
	}

	if found {
		fmt.Println("Cache value:", cachedValue)
	} else {
		fmt.Println("Key not found in cache")
	}

	// Wait for the expiration time plus 1 second to let the cache entry expire
	time.Sleep(duration + time.Second)

	// Try to get the value from the cache again
	cachedValue, found, err = cacheStore.GetValue(key)
	if err != nil {
		fmt.Println("Error getting cache value:", err)
		return
	}

	if found {
		fmt.Println("Cache value:", cachedValue)
	} else {
		fmt.Println("Key not found in cache")
	}

	// Close the cache store
	cacheStore.CloseCacheStore()
}

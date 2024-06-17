# Cache Store
Cache Store is a simple, efficient, and thread-safe in-memory cache library for Go. It allows you to store key-value pairs with optional expiration times, making it an ideal choice for caching frequently accessed data to improve application performance.

### Use Case
Cache Store can be used in various scenarios where caching is beneficial, such as:

- Caching API responses to reduce latency and minimize the number of requests to external services.
- Storing frequently accessed data, like session data or user preferences, to improve application performance.
- Implementing a least recently used (LRU) cache for efficient memory management.

Here's an example of using Cache Store:

```Enter key: a3f5d76b-908c-4b7b-9b0c-8a6ee2533e6a
Enter value: {"name": "John Doe", "age": 30, "email": "johndoe@example.com"}
Enter expiration time (seconds): 30
Cache value: {"name": "John Doe", "age": 30, "email": "johndoe@example.com"}
Key not found in cache
```

### Features:
1. Optional expiration caching
2. Thread-safe design
3. Automatic expired cleanup
4. Comprehensive cache ops
5. Proper error handling
   

### Benchmarking
To run benchmarks, use the following command:

`go test -bench=. -benchmem -cpu 12`

Example benchmark results on a macOS (Darwin) system with an ARM64 architecture:
```
goos: darwin
goarch: arm64
pkg: github.com/archishmansengupta/cache
BenchmarkCache_SetValue/Set-0-12                 8359509               135.5 ns/op            72 B/op          4 allocs/op
BenchmarkCache_GetValue/Get-0-12                31501023                38.18 ns/op            0 B/op          0 allocs/op
BenchmarkCache_Iterate/Iterate-12                 108660             10925 ns/op              16 B/op          1 allocs/op
PASS
ok      github.com/archishmansengupta/cache     5.959s
```
These benchmarks demonstrate the performance of Cache Store's core functions, such as setting, getting, and iterating over cache items.

### Contribution
Contribute to the project, report issues, or request features by opening a pull request or an issue.

Happy caching!

## cache random data with set expiration time

### Use Case:
```
Enter key: a3f5d76b-908c-4b7b-9b0c-8a6ee2533e6a
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

`go test -bench=. -benchmem -cpu 12`

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


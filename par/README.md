Additional synchronization primitives:
* generic thread-safe map;
* safer waitgroup;
* singleflight (duplicate call suppression).

**Thread-safe map**
```go
Len() int
Get(key K) (V, bool)
Set(key K, value V)
SetIf(key K, cond func(value V, exists bool) bool, valfunc func(prev V) V) (value V, ok bool)
Delete(key K)
DeleteIf(key K, cond func(value V) bool) bool
Clear()
ForEach(fun func(key K, value V) bool) bool
```

**Safer waitgroup**
```go
Go(fun func())
Wait()
```

**Singleflight**
```go
Do(key K, fun func() V) V
```

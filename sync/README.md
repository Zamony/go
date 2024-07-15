Additional synchronization primitives:
* generic goroutine-safe map;
* safer waitgroup;
* singleflight (duplicate call suppression);
* context-aware mutex.

**Thread-safe map**
```
Len() int
Get(key K) (V, bool)
Set(key K, value V)
SetIf(key K, newValue V, cond Condition) (actual V, ok bool)
Delete(key K)
DeleteIf(key K, cond Condition) bool
Clear()
All(yield func(key K, value V) bool) bool
```

**Safer waitgroup**
```
Go(fun func())
Wait()
```

**Singleflight**
```
Group.Do(key K, fun func() V) V
```

**Mutex**
```
TryLock(ctx context.Context) error
TryRLock(ctx context.Context) error
Lock()
Unlock()
RLock()
RUnlock()
```

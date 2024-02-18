Additional synchronization primitives:
* generic goroutine-safe map;
* safer waitgroup;
* singleflight (duplicate call suppression).

**Thread-safe map**
```
New() Map

Map.Len() int
Map.Get(key K) (V, bool)
Map.Set(key K, value V)
Map.SetIf(key K, cond func(value V, exists bool) bool, valfunc func(prev V) V) (value V, ok bool)
Map.Delete(key K)
Map.DeleteIf(key K, cond func(value V) bool) bool
Map.Clear()
Map.All(yield func(key K, value V) bool) bool
```

**Safer waitgroup**
```
Go(fun func())
Wait()
```

**Singleflight**
```
New() Group

Group.Do(key K, fun func() V) V
```


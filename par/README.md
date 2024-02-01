Additional synchronization primitives:
* generic thread-safe map;
* safer waitgroup;
* singleflight (duplicate call suppression);
* generic memory pool;
* genetic atomic value.

**Thread-safe map**
```
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
```
Go(fun func())
Wait()
```

**Singleflight**
```
Do(key K, fun func() V) V
```

**Memory pool**
```
Init(factory func() T)
Get() T
Put(x T)
```

**Atomic value**
```
CompareAndSwap(old, new T) (swapped bool)
Load() (val T, ok bool)
Store(val T)
Swap(new T) (old T, ok bool)
```

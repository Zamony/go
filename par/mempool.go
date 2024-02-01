package par

import (
	"sync"
)

// A MemoryPool is a generic wrapper around a sync.Pool.
//
// A MemoryPool must not be copied after first use.
type MemoryPool[T any] struct {
	pool sync.Pool
}

// Init initializes a pool with a generator function.
func (p *MemoryPool[T]) Init(factory func() T) {
	p.pool = sync.Pool{New: func() any {
		return factory()
	}}
}

// Get is a generic wrapper around sync.Pool's Get method.
func (p *MemoryPool[T]) Get() T {
	if v := p.pool.Get(); v != nil {
		return v.(T)
	}
	var t T
	return t
}

// Put is a generic wrapper around sync.Pool's Put method.
func (p *MemoryPool[T]) Put(x T) {
	p.pool.Put(x)
}

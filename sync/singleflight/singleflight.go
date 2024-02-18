package singleflight

import (
	"sync/atomic"

	"github.com/Zamony/go/sync/muxmap"
)

type result[T any] struct {
	Value   T
	Done    chan struct{}
	Waiters int64
}

// Group suppresses duplicate function calls.
//
// Group must not be copied after first use.
type Group[K comparable, V any] struct {
	results muxmap.Map[K, *result[V]]
}

// New creates a new singleflight Group.
func New[K comparable, V any]() Group[K, V] {
	return Group[K, V]{muxmap.New[K, *result[V]]()}
}

// Do executes and returns the results of the given function, making
// sure that only one execution is in-flight for a given key at a
// time. If a duplicate comes in, the duplicate caller waits for the
// original to complete and receives the same results.
func (g *Group[K, V]) Do(key K, fun func() V) V {
	res, isPrimary := g.results.SetIf(key, func(_ *result[V], exists bool) bool {
		return !exists
	}, func(*result[V]) *result[V] {
		return &result[V]{Waiters: 1, Done: make(chan struct{})}
	})
	if !isPrimary {
		atomic.AddInt64(&res.Waiters, 1)
	} else {
		res.Value = fun()
		close(res.Done)
	}

	<-res.Done
	atomic.AddInt64(&res.Waiters, -1)
	g.results.DeleteIf(key, func(value *result[V]) bool {
		return atomic.LoadInt64(&value.Waiters) == 0
	})

	return res.Value
}

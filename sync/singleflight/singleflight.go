// Package singleflight provides duplicate function calls suppression.
package singleflight

import (
	"context"
	"sync/atomic"

	"github.com/Zamony/go/sync/hatmap"
)

type flight[T any] struct {
	Result       T
	Done         chan struct{}
	WaitersCount int64
}

// Group suppresses duplicate function calls.
//
// Group must not be copied after first use.
type Group[K comparable, V any] struct {
	flights hatmap.Map[K, *flight[V]]
}

// Do executes and returns the results of the given function, making
// sure that only one execution is in-flight for a given key at a
// time. If a duplicate comes in, the duplicate caller waits for the
// original to complete and receives the same results.
func (g *Group[K, V]) Do(ctx context.Context, key K, fun func() V) (result V, err error) {
	newFlight := &flight[V]{WaitersCount: 1, Done: make(chan struct{})}
	currFlight, isSet := g.flights.SetIf(key, newFlight, notExists)
	if !isSet {
		atomic.AddInt64(&currFlight.WaitersCount, 1)
	} else {
		currFlight.Result = fun()
		close(currFlight.Done)
	}

	select {
	case <-currFlight.Done:
		result = currFlight.Result
	case <-ctx.Done():
		err = ctx.Err()
	}

	g.flights.DeleteIf(key, func(v *flight[V]) bool {
		return atomic.LoadInt64(&v.WaitersCount) == 0
	})

	return result, err
}

func notExists[V any](currValue *flight[V]) bool {
	return currValue == nil
}

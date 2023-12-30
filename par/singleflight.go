package par

import (
	"sync/atomic"
)

type flightResult[T any] struct {
	Result  T
	Done    chan struct{}
	Waiters int64
}

// Singleflight suppresses duplicate function calls.
type Singleflight[K comparable, V any] struct {
	flights Map[K, *flightResult[V]]
}

// Do executes and returns the results of the given function, making
// sure that only one execution is in-flight for a given key at a
// time. If a duplicate comes in, the duplicate caller waits for the
// original to complete and receives the same results.
func (f *Singleflight[K, V]) Do(key K, fun func() V) V {
	flight, isPrimary := f.flights.SetIf(key, func(_ *flightResult[V], exists bool) bool {
		return !exists
	}, func(*flightResult[V]) *flightResult[V] {
		return &flightResult[V]{Waiters: 1, Done: make(chan struct{})}
	})
	if !isPrimary {
		atomic.AddInt64(&flight.Waiters, 1)
	} else {
		flight.Result = fun()
		close(flight.Done)
	}

	<-flight.Done
	atomic.AddInt64(&flight.Waiters, -1)
	f.flights.DeleteIf(key, func(value *flightResult[V]) bool {
		return atomic.LoadInt64(&value.Waiters) == 0
	})

	return flight.Result
}

package piter

import (
	"iter"
)

// New creates parallel iterator.
// Unused parallel iterator is a memory leak.
func New[T any](seq iter.Seq[T]) iter.Seq[T] {
	ch := make(chan T)
	done := make(chan struct{})
	go func() {
		defer close(ch)
		for elem := range seq {
			select {
			case ch <- elem:
			case <-done:
				return
			}
		}
	}()
	return func(yield func(T) bool) {
		for elem := range ch {
			if !yield(elem) {
				close(done)
				return
			}
		}
	}
}

type pair[T, E any] struct {
	Key   T
	Value E
}

func flatten[T, E any](seq iter.Seq2[T, E]) iter.Seq[pair[T, E]] {
	return func(yield func(pair[T, E]) bool) {
		for key, val := range seq {
			if !yield(pair[T, E]{Key: key, Value: val}) {
				return
			}
		}
	}
}

// New2 creates parallel iterator for two values.
// Unused parallel iterator is a memory leak.
func New2[T, E any](seq iter.Seq2[T, E]) iter.Seq2[T, E] {
	return func(yield func(T, E) bool) {
		for pair := range New(flatten(seq)) {
			if !yield(pair.Key, pair.Value) {
				return
			}
		}
	}
}

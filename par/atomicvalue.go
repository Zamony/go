package par

import "sync/atomic"

// A Value provides an atomic load and store of a consistently typed value.
//
// A Value must not be copied after first use.
type AtomicValue[T any] struct {
	val atomic.Value
}

// CompareAndSwap executes the compare-and-swap operation for the Value.
func (v *AtomicValue[T]) CompareAndSwap(old, new T) (swapped bool) {
	return v.val.CompareAndSwap(old, new)
}

// Load returns the value set by the most recent Store.
func (v *AtomicValue[T]) Load() (val T, ok bool) {
	if r := v.val.Load(); r != nil {
		return r.(T), true
	}
	return val, false
}

// Store sets the value of the Value v to val.
func (v *AtomicValue[T]) Store(val T) {
	v.val.Store(val)
}

// Swap stores new into Value and returns the previous value.
func (v *AtomicValue[T]) Swap(new T) (old T, ok bool) {
	if r := v.val.Swap(new); r != nil {
		return r.(T), true
	}
	return old, false
}

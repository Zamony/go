// Package hatmap provides generic goroutine-safe map.
// It is implemented as an ordinary map protected by mutex ("mutex hat" idiom).
package hatmap

import (
	"iter"
	"sync"
)

// Condition to operate on a current value by the given key.
// May be called multiple times.
type Condition[V any] func(currValue V) bool

// Map is a goroutine-safe map.
//
// Map must not be copied after first use.
type Map[K comparable, V any] struct {
	data map[K]V
	mu   sync.RWMutex // protects data
}

// New creates a new map with size hint.
func New[K comparable, V any](size int) Map[K, V] {
	return Map[K, V]{data: make(map[K]V, size)}
}

// Len returns the number of elements in the map.
func (m *Map[K, V]) Len() int {
	m.mu.RLock()
	n := len(m.data)
	m.mu.RUnlock()
	return n
}

// Set sets the value by the given key.
func (m *Map[K, V]) Set(key K, value V) {
	m.mu.Lock()
	if m.data == nil {
		m.data = make(map[K]V)
	}
	m.data[key] = value
	m.mu.Unlock()
}

func (m *Map[K, V]) canSet(key K, cond Condition[V]) (V, bool) {
	m.mu.RLock()
	value := m.data[key]
	m.mu.RUnlock()
	return value, cond(value)
}

// SetIf conditionally sets the value by the given key.
// Returns final value and condition result.
func (m *Map[K, V]) SetIf(key K, newValue V, cond Condition[V]) (actual V, ok bool) {
	value, ok := m.canSet(key, cond)
	if !ok {
		return value, false
	}
	return m.setIfSlow(key, newValue, cond)
}

func (m *Map[K, V]) setIfSlow(key K, newValue V, cond Condition[V]) (actual V, ok bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	value, ok := m.data[key]
	if !cond(value) {
		return value, false
	}
	if m.data == nil {
		m.data = make(map[K]V)
	}
	m.data[key] = newValue
	return newValue, true
}

// Get gets value y the given key.
func (m *Map[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	value, ok := m.data[key]
	m.mu.RUnlock()
	return value, ok
}

// Delete deletes the value by the given key.
// If the key doesn't exist does nothing.
func (m *Map[K, V]) Delete(key K) {
	m.mu.Lock()
	delete(m.data, key)
	m.mu.Unlock()
}

func (m *Map[K, V]) canDelete(key K, cond Condition[V]) bool {
	m.mu.RLock()
	value, ok := m.data[key]
	m.mu.RUnlock()
	return ok && cond(value)
}

// DeleteIf conditionally deletes the value by the given key.
// Returns true if the value was deleted.
func (m *Map[K, V]) DeleteIf(key K, cond Condition[V]) bool {
	if !m.canDelete(key, cond) {
		return false
	}

	return m.deleteIfSlow(key, cond)
}

func (m *Map[K, V]) deleteIfSlow(key K, cond Condition[V]) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	value, ok := m.data[key]
	if !ok || !cond(value) {
		return false
	}

	delete(m.data, key)
	return true
}

// Clear clears the map.
func (m *Map[K, V]) Clear() {
	m.mu.Lock()
	clear(m.data)
	m.mu.Unlock()
}

// All iterates over a map.
// Don't modify the map while iterating over it.
func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.mu.RLock()
		defer m.mu.RUnlock()

		for k, v := range m.data {
			if !yield(k, v) {
				return
			}
		}
	}
}

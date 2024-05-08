package hatmap

import "sync"

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

func (m *Map[K, V]) canSet(key K, cond func(value V, exists bool) bool) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, ok := m.data[key]
	condOk := cond(value, ok)
	return value, condOk
}

// SetIf conditionally sets the value by the given key.
// Condition function must be pure.
// Returns final value and condition result.
func (m *Map[K, V]) SetIf(key K, cond func(value V, exists bool) bool, valfunc func(prev V) V) (value V, ok bool) {
	value, ok = m.canSet(key, cond)
	if !ok {
		return value, false
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	value, ok = m.data[key]
	if !cond(value, ok) {
		return value, false
	}

	value = valfunc(value)
	if m.data == nil {
		m.data = make(map[K]V)
	}
	m.data[key] = value
	return value, true
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

func (m *Map[K, V]) canDelete(key K, cond func(value V) bool) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, ok := m.data[key]
	return ok && cond(value)
}

// DeleteIf conditionally deletes the value by the given key.
// Condition function must be pure.
// Returns true if the value was deleted.
func (m *Map[K, V]) DeleteIf(key K, cond func(value V) bool) bool {
	if !m.canDelete(key, cond) {
		return false
	}

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
func (m *Map[K, V]) All(yield func(K, V) bool) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for k, v := range m.data {
		if !yield(k, v) {
			return false
		}
	}

	return true
}

package par

import "sync"

// Map is a thread-safe map. N is a size hint for the map.
type Map[K comparable, V any] struct {
	data map[K]V
	mu   sync.RWMutex // protects data
	N    int
}

// Len returns number of elements in the map.
func (m *Map[K, V]) Len() int {
	m.mu.RLock()
	n := len(m.data)
	m.mu.RUnlock()
	return n
}

func (m *Map[K, V]) initOnce() {
	switch {
	case m.data != nil:
	case m.N > 0:
		m.data = make(map[K]V, m.N)
	default:
		m.data = make(map[K]V)
	}
}

// Set sets the value by the given key.
func (m *Map[K, V]) Set(key K, value V) {
	m.mu.Lock()
	m.initOnce()
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

	m.initOnce()
	value = valfunc(value)
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
	m.data = nil
	m.mu.Unlock()
}

// ForEach iterates over map and calls provided function for each key and value.
// Iteration is aborted after provided function returns false.
// Returns false if iteration was aborted.
func (m *Map[K, V]) ForEach(fun func(key K, value V) bool) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for k, v := range m.data {
		if !fun(k, v) {
			return false
		}
	}

	return true
}

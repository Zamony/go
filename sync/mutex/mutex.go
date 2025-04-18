package mutex

import (
	"context"
	"sync"
)

// Mutex is a custom mutual exclusion lock that uses a channel for synchronization.
type Mutex struct {
	ch chan struct{}
}

// New creates and returns a new Mutex instance.
func New() Mutex {
	return Mutex{ch: make(chan struct{}, 1)}
}

// TryLock attempts to acquire the lock.
// It returns an error if the context is canceled.
func (m Mutex) TryLock(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err() // don't lock if context is canceled
	default:
		select {
		case <-ctx.Done():
			return ctx.Err()
		case m.ch <- struct{}{}:
			return nil
		}
	}
}

// Lock acquires the lock, blocking until it is available.
func (m Mutex) Lock() {
	m.ch <- struct{}{}
}

// Unlock releases the lock.
func (m Mutex) Unlock() {
	<-m.ch
}

// A RWMutex is a reader/writer mutual exclusion lock.
// The lock can be held by an arbitrary number of readers or a single writer.
// The zero value for a RWMutex is an unlocked mutex.
//
// A RWMutex must not be copied after first use.
type RWMutex struct {
	mu sync.RWMutex
}

// TryLock waits to lock RWMutex for writing.
// A context error is returned on canceled context.
func (m *RWMutex) TryLock(ctx context.Context) error {
	if m.mu.TryLock() {
		return nil
	}

	return wait(ctx, &m.mu)
}

// TryRLock waits to lock RWMutex for reading.
// A context error is returned on canceled context.
func (m *RWMutex) TryRLock(ctx context.Context) error {
	if m.mu.TryRLock() {
		return nil
	}
	return wait(ctx, m.mu.RLocker())
}

// Lock locks RWMutex for writing.
func (m *RWMutex) Lock() {
	m.mu.Lock()
}

// Unlock unlocks RWMutex for writing.
func (m *RWMutex) Unlock() {
	m.mu.Unlock()
}

// RLock locks RWMutex for reading.
func (m *RWMutex) RLock() {
	m.mu.RLock()
}

// RUnlock undoes a single RLock call.
func (m *RWMutex) RUnlock() {
	m.mu.RUnlock()
}

func wait(ctx context.Context, mu sync.Locker) error {
	done := make(chan struct{})
	go func() {
		mu.Lock()
		select {
		case <-ctx.Done():
			mu.Unlock()
		case done <- struct{}{}:
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

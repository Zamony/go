package rwmutex

import (
	"context"
	"sync"
)

type Mutex struct {
	mu sync.RWMutex
}

func (m *Mutex) TryLock(ctx context.Context) error {
	if m.mu.TryLock() {
		return nil
	}

	return wait(ctx, &m.mu)
}

func (m *Mutex) TryRLock(ctx context.Context) error {
	if m.mu.TryRLock() {
		return nil
	}
	return wait(ctx, (&m.mu).RLocker())
}

func (m *Mutex) Lock() {
	m.mu.Lock()
}

func (m *Mutex) Unlock() {
	m.mu.Unlock()
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

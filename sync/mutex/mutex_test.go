package mutex_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Zamony/go/sync/mutex"
)

const timeout = 10 * time.Millisecond

func TestRWMutexContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(cancel)

	var mu mutex.RWMutex
	mu.Lock()
	defer mu.Unlock()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := mu.TryLock(ctx); err != ctx.Err() {
			t.Error("lock", err)
		}
	}()
	go func() {
		defer wg.Done()
		if err := mu.TryRLock(ctx); err != ctx.Err() {
			t.Error("rlock", err)
		}
	}()

	wg.Wait()
}

func TestRWMutexLock(t *testing.T) {
	ctx := context.Background()
	var mu mutex.RWMutex
	if err := mu.TryLock(ctx); err != nil {
		t.Error("lock", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := mu.TryRLock(ctx); err != nil {
			t.Error("rlock", err)
		}
		mu.RUnlock()
	}()

	time.Sleep(timeout)
	mu.Unlock()
	wg.Wait()
}

func TestMutex(t *testing.T) {
	m := mutex.New()

	// Test TryLock with a context that is not canceled
	ctx := context.Background()
	err := m.TryLock(ctx)
	if err != nil {
		t.Fatalf("TryLock failed: %v", err)
	}
	m.Unlock()

	// Test Lock
	m.Lock()
	m.Unlock()

	// Test TryLock with a canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err = m.TryLock(ctx)
	if err != context.Canceled {
		t.Fatalf("TryLock with canceled context should return context.Canceled, got: %v", err)
	}

	// Test that the mutex is actually locked
	m.Lock()
	done := make(chan bool)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()
		err := m.TryLock(ctx)
		if err != context.DeadlineExceeded {
			t.Errorf("Expected TryLock to fail with context.DeadlineExceeded, got: %v", err)
		}
		done <- true
	}()
	<-done
	m.Unlock()
}

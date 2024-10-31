package mutex_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Zamony/go/sync/mutex"
)

const timeout = 10 * time.Millisecond

func TestMutexContext(t *testing.T) {
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

func TestMutexLock(t *testing.T) {
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

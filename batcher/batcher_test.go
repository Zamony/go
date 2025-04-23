package batcher_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Zamony/go/batcher"
)

func TestBatchByCount(t *testing.T) {
	var batches [][]int
	var waiter sync.WaitGroup
	waiter.Add(1)
	flushFunc := func(_ context.Context, items []int) {
		defer waiter.Done()
		batches = append(batches, items)
	}

	b := batcher.New(3, time.Hour, flushFunc)
	for i := range 5 { // Add items to trigger count-based flush
		if err := b.Add(context.Background(), i); err != nil {
			t.Fatalf("Add failed: %v", err)
		}
	}

	waiter.Wait()
	if len(batches) != 1 {
		t.Errorf("Expected 1 batch, got %v", batches)
	} else if len(batches[0]) != 3 {
		t.Errorf("Expected 1 batch with 3 items, got %v", batches)
	}
}

func TestBatchByTimeout(t *testing.T) {
	var batches [][]string
	var waiter sync.WaitGroup
	waiter.Add(1)
	flushFunc := func(_ context.Context, items []string) {
		defer waiter.Done()
		batches = append(batches, items)
	}

	b := batcher.New(100, 50*time.Millisecond, flushFunc)
	if err := b.Add(context.Background(), "test"); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	waiter.Wait()
	if len(batches) != 1 || len(batches[0]) != 1 || batches[0][0] != "test" {
		t.Errorf("Expected 1 batch with ['test'], got %v", batches)
	}
}

func TestAddWithContextCancellation(t *testing.T) {
	flushFunc := func(_ context.Context, items []int) {
		t.Fatal("flush should not be called")
	}

	b := batcher.New(10, time.Minute, flushFunc)
	defer b.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := b.Add(ctx, 1)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	if err != context.Canceled {
		t.Fatalf("Expected context.Canceled error, got %v", err)
	}
}

func TestClose(t *testing.T) {
	var waiter sync.WaitGroup
	waiter.Add(1)
	var flushed bool
	flushFunc := func(_ context.Context, items []int) {
		defer waiter.Done()
		flushed = true
	}

	b := batcher.New(100, time.Hour, flushFunc)
	if err := b.Add(context.Background(), 1); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	b.Close()
	waiter.Wait()

	if !flushed {
		t.Error("Close didn't trigger flush")
	}
}

package batcher

import (
	"context"
	"sync"
	"time"
)

// Batcher represents a batching mechanism that collects items and flushes them
// either when a maximum count is reached or after a timeout period.
// The generic type T represents the type of items to be batched.
type Batcher[T any] struct {
	ctx       context.Context
	cancel    context.CancelFunc
	items     []T
	itemsCh   chan []T
	maxCount  int
	ticker    *time.Ticker
	flushFunc FlushFunc[T]
	waitGroup sync.WaitGroup
}

// FlushFunc defines the signature for the function that will be called
// when the batch is flushed. It receives the context and the batched items.
type FlushFunc[T any] func(context.Context, []T)

// New creates and starts a new Batcher instance.
// maxCount specifies the maximum number of items to collect before flushing.
// timeout specifies the maximum duration to wait before flushing.
// fun is the function that will be called when the batch is flushed.
// Returns a pointer to the newly created Batcher.
func New[T any](maxCount int, timeout time.Duration, fun FlushFunc[T]) *Batcher[T] {
	ctx, cancel := context.WithCancel(context.Background())
	b := &Batcher[T]{
		ctx:       ctx,
		cancel:    cancel,
		itemsCh:   make(chan []T),
		maxCount:  maxCount,
		ticker:    time.NewTicker(timeout),
		flushFunc: fun,
	}

	b.waitGroup.Add(1)
	go func() {
		defer b.waitGroup.Done()
		b.run()
	}()

	return b
}

// Add adds one or more items to the batch. The operation respects the provided context.
// Returns an error if the context is canceled before the items can be added.
func (b *Batcher[T]) Add(ctx context.Context, items ...T) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case b.itemsCh <- items:
		return nil
	}
}

// run is the main processing loop that handles adding items, timeouts, and shutdown.
func (b *Batcher[T]) run() {
	for {
		select {
		case <-b.ctx.Done():
			b.flush()
			return
		case <-b.ticker.C:
			b.flush()
		case items := <-b.itemsCh:
			b.items = append(b.items, items...)
			if len(b.items) >= b.maxCount {
				b.flush()
			}
		}
	}
}

// flush calls the flush function with the current batch of items and resets the collection.
func (b *Batcher[T]) flush() {
	if len(b.items) > 0 {
		b.flushFunc(b.ctx, b.items)
		b.items = nil
	}
}

// Close stops the batch processor, flushes any remaining items, and cleans up resources.
// It blocks until all pending operations are complete.
func (b *Batcher[T]) Close() {
	b.cancel()
	b.waitGroup.Wait()
	b.ticker.Stop()
}

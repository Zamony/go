package singleflight_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Zamony/go/sync/singleflight"
)

func TestSingleFlight(t *testing.T) {
	var ncalls int64
	var wg sync.WaitGroup
	defer wg.Wait()
	single := singleflight.Group[string, int64]{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := single.Do(context.Background(), "a", func() int64 {
				time.Sleep(10 * time.Millisecond)
				return atomic.AddInt64(&ncalls, 1)
			})
			if err != nil {
				t.Errorf("No error is expected, got: %+v", err)
				return
			}
			if want := int64(1); got != want {
				t.Errorf("Not equal (-want, +got):\n- %+v\n+ %+v\n", want, got)
			}
		}()
	}
}

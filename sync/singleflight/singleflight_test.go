package singleflight_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Zamony/go/sync/singleflight"
)

func TestSingleFlight(t *testing.T) {
	t.Parallel()

	var ncalls int64
	var wg sync.WaitGroup
	defer wg.Wait()
	single := singleflight.New[string, int64]()
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got := single.Do("a", func() int64 {
				time.Sleep(10 * time.Millisecond)
				return atomic.AddInt64(&ncalls, 1)
			})
			if want := int64(1); got != want {
				t.Errorf("Not equal (-want, +got):\n- %+v\n+ %+v\n", want, got)
			}
		}()
	}
}

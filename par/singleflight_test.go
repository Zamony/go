package par_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/Zamony/go/par"
)

func TestSingleFlight(t *testing.T) {
	t.Parallel()

	var ncalls int64
	var wg par.WaitGroup
	single := par.Singleflight[string, int64]{}
	for i := 0; i < 5; i++ {
		wg.Go(func() {
			got := single.Do("a", func() int64 {
				time.Sleep(10 * time.Millisecond)
				return atomic.AddInt64(&ncalls, 1)
			})
			equal(t, got, int64(1))
		})
	}

	wg.Wait()
}

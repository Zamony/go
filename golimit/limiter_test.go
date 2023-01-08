package golimit

import (
	"testing"
	"time"
)

func TestLimiter_Limit(t *testing.T) {
	const (
		niter  = 3
		limit  = 10
		total  = 100 * limit
		period = 20 * time.Millisecond
	)
	allowed := 0
	lim := New(limit, period)
	for i := 0; i < niter; i++ {
		for j := 0; j < total; j++ {
			if !lim.Limit(1) {
				allowed++
			}
		}

		time.Sleep(period)
	}

	want := niter * limit
	if allowed != want {
		t.Fatalf("Got %d allowed, want %d", allowed, want)
	}
}

func TestLimiter_Up(t *testing.T) {
	const (
		limit = 3
	)

	allowed := 0
	lim := New(limit, time.Second)
	if !lim.Limit(limit) {
		allowed += limit
	}
	if !lim.Limit(limit) {
		allowed += limit
	}

	lim.Up(limit)
	if !lim.Limit(limit) {
		allowed += limit
	}

	want := limit + limit
	if allowed != want {
		t.Fatalf("Got %d allowed, want %d", allowed, want)
	}
}

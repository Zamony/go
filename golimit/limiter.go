package golimit

import (
	"sync"
	"time"
)

// Limiter is a goroutine-safe rate-limiter,
// which implements token-bucket algorithm.
type Limiter struct {
	mu     sync.Mutex
	limit  float64
	curr   float64
	period float64
	last   float64
}

// New creates a new limiter with specified limit and period.
func New(limit float64, period time.Duration) *Limiter {
	return &Limiter{
		mu:     sync.Mutex{},
		limit:  limit,
		curr:   limit,
		period: float64(period.Nanoseconds()),
		last:   0,
	}
}

// Limit returns true if an action was rejected.
// It accepts positive weight of an action as an argument.
func (l *Limiter) Limit(n float64) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := float64(time.Now().UnixNano())
	l.curr += (now - l.last) * l.limit / l.period
	l.last = now
	if l.curr > l.limit {
		l.curr = l.limit
	}

	if l.curr < n {
		return true
	}

	l.curr -= n
	return false
}

// Up increases the current possible weight by n.
func (l *Limiter) Up(n float64) {
	l.mu.Lock()
	l.curr += n
	if l.curr > l.limit {
		l.curr = l.limit
	}
	l.mu.Unlock()
}

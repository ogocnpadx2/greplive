// Package throttle provides a line-rate throttler that drops lines when
// the incoming rate exceeds a configured maximum lines-per-second.
package throttle

import (
	"sync"
	"time"
)

// Throttle tracks a sliding count of lines and reports whether a new line
// should be allowed through.
type Throttle struct {
	mu       sync.Mutex
	maxLines int
	window   time.Duration
	times    []time.Time
}

// New creates a Throttle that allows at most maxLines per window duration.
// If maxLines is zero the throttle is disabled and every line is allowed.
func New(maxLines int, window time.Duration) *Throttle {
	if window <= 0 {
		window = time.Second
	}
	return &Throttle{
		maxLines: maxLines,
		window:   window,
		times:    make([]time.Time, 0, maxLines+1),
	}
}

// Allow returns true if the line should be forwarded, false if it should be
// dropped. It is safe for concurrent use.
func (t *Throttle) Allow() bool {
	if t.maxLines == 0 {
		return true
	}
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-t.window)

	// evict stale entries
	valid := t.times[:0]
	for _, ts := range t.times {
		if ts.After(cutoff) {
			valid = append(valid, ts)
		}
	}
	t.times = valid

	if len(t.times) >= t.maxLines {
		return false
	}
	t.times = append(t.times, now)
	return true
}

// Enabled reports whether throttling is active.
func (t *Throttle) Enabled() bool {
	return t.maxLines > 0
}

// Reset clears the internal counter.
func (t *Throttle) Reset() {
	t.mu.Lock()
	t.times = t.times[:0]
	t.mu.Unlock()
}

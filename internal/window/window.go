// Package window provides a sliding time-window counter for tracking
// line throughput over a configurable rolling duration.
package window

import (
	"sync"
	"time"
)

// entry records a single event timestamp.
type entry struct {
	at time.Time
}

// Window is a thread-safe sliding time-window counter.
type Window struct {
	mu       sync.Mutex
	bucket   []entry
	duration time.Duration
	now      func() time.Time // injectable for testing
}

// New creates a Window that tracks events within the given duration.
// A zero or negative duration disables eviction (all events are kept).
func New(d time.Duration) *Window {
	return &Window{
		duration: d,
		now:      time.Now,
	}
}

// Add records a new event at the current time and evicts stale entries.
func (w *Window) Add() {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := w.now()
	w.bucket = append(w.bucket, entry{at: now})
	w.evict(now)
}

// Count returns the number of events within the current window.
func (w *Window) Count() int {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.evict(w.now())
	return len(w.bucket)
}

// Rate returns the average events-per-second over the window duration.
// Returns 0 if the window duration is zero or no events are recorded.
func (w *Window) Rate() float64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.duration <= 0 {
		return 0
	}
	w.evict(w.now())
	if len(w.bucket) == 0 {
		return 0
	}
	return float64(len(w.bucket)) / w.duration.Seconds()
}

// Reset clears all recorded events.
func (w *Window) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.bucket = w.bucket[:0]
}

// evict removes entries older than the window duration.
// Caller must hold w.mu.
func (w *Window) evict(now time.Time) {
	if w.duration <= 0 {
		return
	}
	cutoff := now.Add(-w.duration)
	i := 0
	for i < len(w.bucket) && w.bucket[i].at.Before(cutoff) {
		i++
	}
	w.bucket = w.bucket[i:]
}

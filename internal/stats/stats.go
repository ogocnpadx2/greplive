// Package stats tracks runtime metrics for greplive sessions.
package stats

import (
	"sync"
	"sync/atomic"
	"time"
)

// Counters holds atomic counters for a streaming session.
type Counters struct {
	LinesRead    atomic.Int64
	LinesMatched atomic.Int64
	LinesDropped atomic.Int64
	StartTime    time.Time
	mu           sync.RWMutex
	severityCounts map[string]int64
}

// New creates a new Counters instance with the start time set to now.
func New() *Counters {
	return &Counters{
		StartTime:      time.Now(),
		severityCounts: make(map[string]int64),
	}
}

// IncrRead increments the lines-read counter.
func (c *Counters) IncrRead() { c.LinesRead.Add(1) }

// IncrMatched increments the lines-matched counter.
func (c *Counters) IncrMatched() { c.LinesMatched.Add(1) }

// IncrDropped increments the lines-dropped counter.
func (c *Counters) IncrDropped() { c.LinesDropped.Add(1) }

// IncrSeverity increments the counter for a named severity level.
func (c *Counters) IncrSeverity(level string) {
	c.mu.Lock()
	c.severityCounts[level]++
	c.mu.Unlock()
}

// SeverityCounts returns a copy of the per-severity counters.
func (c *Counters) SeverityCounts() map[string]int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[string]int64, len(c.severityCounts))
	for k, v := range c.severityCounts {
		out[k] = v
	}
	return out
}

// Elapsed returns the duration since the session started.
func (c *Counters) Elapsed() time.Duration {
	return time.Since(c.StartTime)
}

// MatchRate returns the fraction of read lines that matched, in the range
// [0.0, 1.0]. Returns 0 if no lines have been read yet.
func (c *Counters) MatchRate() float64 {
	read := c.LinesRead.Load()
	if read == 0 {
		return 0
	}
	return float64(c.LinesMatched.Load()) / float64(read)
}

// Snapshot returns an immutable point-in-time view of the counters.
type Snapshot struct {
	LinesRead      int64
	LinesMatched   int64
	LinesDropped   int64
	Elapsed        time.Duration
	SeverityCounts map[string]int64
	MatchRate      float64
}

// Snapshot captures the current state of the counters.
func (c *Counters) Snapshot() Snapshot {
	return Snapshot{
		LinesRead:      c.LinesRead.Load(),
		LinesMatched:   c.LinesMatched.Load(),
		LinesDropped:   c.LinesDropped.Load(),
		Elapsed:        c.Elapsed(),
		SeverityCounts: c.SeverityCounts(),
		MatchRate:      c.MatchRate(),
	}
}

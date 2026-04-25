// Package clamp limits the number of consecutive matching lines emitted
// within a sliding time window. Once the cap is reached, further matches
// are silently dropped until the window advances and capacity is restored.
package clamp

import (
	"regexp"
	"sync"
	"time"
)

// Clamp holds state for the consecutive-match limiter.
type Clamp struct {
	re      *regexp.Regexp
	max     int
	window  time.Duration
	clock   func() time.Time
	mu      sync.Mutex
	buckets []time.Time
	enabled bool
}

// New creates a Clamp that allows at most max matching lines per window
// duration. A zero or negative max disables clamping (all lines pass).
// pattern must be a valid regular expression; an error is returned otherwise.
func New(pattern string, max int, window time.Duration) (*Clamp, error) {
	return newWithClock(pattern, max, window, time.Now)
}

func newWithClock(pattern string, max int, window time.Duration, clock func() time.Time) (*Clamp, error) {
	if pattern == "" || max <= 0 {
		return &Clamp{enabled: false, clock: clock}, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	if window <= 0 {
		window = time.Minute
	}
	return &Clamp{
		re:      re,
		max:     max,
		window:  window,
		clock:   clock,
		buckets: make([]time.Time, 0, max),
		enabled: true,
	}, nil
}

// Enabled reports whether the clamp is active.
func (c *Clamp) Enabled() bool { return c.enabled }

// Allow returns true if line should be emitted. Non-matching lines always
// pass. Matching lines pass only while fewer than max matches have occurred
// within the current window.
func (c *Clamp) Allow(line string) bool {
	if !c.enabled {
		return true
	}
	if !c.re.MatchString(line) {
		return true
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.clock()
	cutoff := now.Add(-c.window)
	// evict stale buckets
	valid := c.buckets[:0]
	for _, t := range c.buckets {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	c.buckets = valid
	if len(c.buckets) >= c.max {
		return false
	}
	c.buckets = append(c.buckets, now)
	return true
}

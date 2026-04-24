// Package skip provides a transformer that suppresses the first N lines
// of a stream, useful for skipping headers or boilerplate output.
package skip

import "sync/atomic"

// Skipper drops the first N lines and passes all subsequent lines through.
type Skipper struct {
	max     int64
	counted atomic.Int64
	enabled bool
}

// New creates a Skipper that will drop the first n lines.
// If n <= 0, the Skipper is disabled and all lines pass through.
func New(n int) (*Skipper, error) {
	if n < 0 {
		n = 0
	}
	return &Skipper{
		max:     int64(n),
		enabled: n > 0,
	}, nil
}

// Enabled reports whether the Skipper will drop any lines.
func (s *Skipper) Enabled() bool {
	return s.enabled
}

// Allow returns true if the line should be passed downstream.
// The first N lines return false; all subsequent lines return true.
func (s *Skipper) Allow(line string) bool {
	if !s.enabled {
		return true
	}
	current := s.counted.Add(1)
	return current > s.max
}

// Reset resets the internal counter so that the next N lines are dropped again.
func (s *Skipper) Reset() {
	s.counted.Store(0)
}

// Remaining returns how many lines will still be dropped before passing through.
func (s *Skipper) Remaining() int64 {
	if !s.enabled {
		return 0
	}
	v := s.max - s.counted.Load()
	if v < 0 {
		return 0
	}
	return v
}

// Package head limits output to the first N lines, then signals done.
package head

// Limiter emits at most Max lines, then stops.
type Limiter struct {
	max     int
	count   int
	enabled bool
}

// New creates a Limiter that allows at most max lines through.
// If max is zero or negative, the limiter is disabled (all lines pass).
func New(max int) *Limiter {
	return &Limiter{
		max:     max,
		enabled: max > 0,
	}
}

// Enabled reports whether the limiter is active.
func (l *Limiter) Enabled() bool { return l.enabled }

// Allow returns true if the line should be emitted and false once the
// limit has been reached. The second return value is true when the
// limit has just been hit, signalling that the caller should stop.
func (l *Limiter) Allow() (emit bool, done bool) {
	if !l.enabled {
		return true, false
	}
	if l.count >= l.max {
		return false, true
	}
	l.count++
	return true, l.count == l.max
}

// Reset resets the internal counter so the limiter can be reused.
func (l *Limiter) Reset() { l.count = 0 }

// Remaining returns how many lines can still be emitted.
// Returns -1 when the limiter is disabled.
func (l *Limiter) Remaining() int {
	if !l.enabled {
		return -1
	}
	r := l.max - l.count
	if r < 0 {
		return 0
	}
	return r
}

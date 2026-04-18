// Package repeat suppresses lines that repeat more than a configured
// number of times within a sliding window, emitting a summary instead.
package repeat

import (
	"fmt"
	"sync"
	"time"
)

// Repeater tracks consecutive identical lines and suppresses excess repetitions.
type Repeater struct {
	mu       sync.Mutex
	last     string
	count    int
	max      int
	window   time.Duration
	firstAt  time.Time
	clock    func() time.Time
}

// New returns a Repeater that allows up to max occurrences of the same line
// within window before suppressing further copies.
// A max of 0 disables suppression.
func New(max int, window time.Duration) *Repeater {
	return &Repeater{
		max:    max,
		window: window,
		clock:  time.Now,
	}
}

// WithClock replaces the time source (for testing).
func (r *Repeater) WithClock(fn func() time.Time) *Repeater {
	r.clock = fn
	return r
}

// Enabled reports whether suppression is active.
func (r *Repeater) Enabled() bool {
	return r.max > 0
}

// Push evaluates line and returns (output, emit).
// output is the string to print; emit is false when the line should be dropped.
// When the window expires the counter resets automatically.
func (r *Repeater) Push(line string) (string, bool) {
	if !r.Enabled() {
		return line, true
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.clock()

	// Reset if line changed or window expired.
	if line != r.last || now.Sub(r.firstAt) > r.window {
		r.last = line
		r.count = 1
		r.firstAt = now
		return line, true
	}

	r.count++
	if r.count <= r.max {
		return line, true
	}
	if r.count == r.max+1 {
		return fmt.Sprintf("[repeated message suppressed after %d occurrences]", r.max), true
	}
	return "", false
}

// Flush returns a pending summary if the last line was suppressed, then resets.
func (r *Repeater) Flush() (string, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.count > r.max {
		s := fmt.Sprintf("[%d total occurrences of last repeated message]", r.count)
		r.count = 0
		r.last = ""
		return s, true
	}
	return "", false
}

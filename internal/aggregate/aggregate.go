// Package aggregate groups consecutive matching lines into counted summaries.
package aggregate

import (
	"fmt"
	"regexp"
	"sync"
)

// Aggregator buffers lines matching a pattern and emits summaries.
type Aggregator struct {
	mu      sync.Mutex
	re      *regexp.Regexp
	current string
	count   int
}

// New creates an Aggregator for the given compiled pattern.
func New(re *regexp.Regexp) *Aggregator {
	return &Aggregator{re: re}
}

// Push accepts a line. If it matches the pattern it is buffered.
// Returns (summary, true) when a buffered run ends, otherwise ("", false).
func (a *Aggregator) Push(line string) (string, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	matches := a.re != nil && a.re.MatchString(line)

	if matches {
		if a.current == "" {
			a.current = line
		}
		a.count++
		return "", false
	}

	// Non-matching line — flush any pending group first.
	if a.count > 0 {
		summary := a.flush()
		return summary, true
	}
	return "", false
}

// Flush emits any pending buffered group regardless of what follows.
func (a *Aggregator) Flush() (string, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.count == 0 {
		return "", false
	}
	return a.flush(), true
}

// flush must be called with a.mu held.
func (a *Aggregator) flush() string {
	s := fmt.Sprintf("%s [x%d]", a.current, a.count)
	a.current = ""
	a.count = 0
	return s
}

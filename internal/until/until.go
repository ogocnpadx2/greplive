// Package until provides a filter that passes lines through until a
// pattern matches, then stops emitting further lines.
package until

import (
	"fmt"
	"regexp"
	"sync"
)

// Until emits lines until a trigger pattern is matched. Once triggered,
// all subsequent lines are dropped. The trigger line itself is not emitted.
type Until struct {
	re      *regexp.Regexp
	triggered bool
	mu      sync.Mutex
	enabled bool
}

// New creates a new Until filter. If pattern is empty, the filter is
// disabled and all lines are passed through unchanged.
func New(pattern string) (*Until, error) {
	if pattern == "" {
		return &Until{enabled: false}, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("until: invalid pattern %q: %w", pattern, err)
	}
	return &Until{re: re, enabled: true}, nil
}

// Enabled reports whether the filter is active.
func (u *Until) Enabled() bool {
	return u.enabled
}

// Allow returns true if the line should be emitted. Once the trigger
// pattern has been matched, Allow returns false for all subsequent lines.
// The matching line itself is not emitted.
func (u *Until) Allow(line string) bool {
	if !u.enabled {
		return true
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.triggered {
		return false
	}
	if u.re.MatchString(line) {
		u.triggered = true
		return false
	}
	return true
}

// Reset clears the triggered state, allowing lines to flow again.
func (u *Until) Reset() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.triggered = false
}

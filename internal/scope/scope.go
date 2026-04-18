// Package scope provides line-range filtering: only lines within a
// specified start/end window (by line number or regex) are passed through.
package scope

import (
	"fmt"
	"regexp"
)

// Scope filters lines to a contiguous range defined by start and end
// patterns. Once the start pattern matches, lines are emitted until the
// end pattern matches (inclusive).
type Scope struct {
	start   *regexp.Regexp
	end     *regexp.Regexp
	active  bool
	enabled bool
}

// New creates a Scope from start and end regex patterns. If both patterns
// are empty the Scope is disabled and every line passes through.
func New(startPattern, endPattern string) (*Scope, error) {
	if startPattern == "" && endPattern == "" {
		return &Scope{enabled: false}, nil
	}
	var s, e *regexp.Regexp
	var err error
	if startPattern != "" {
		if s, err = regexp.Compile(startPattern); err != nil {
			return nil, fmt.Errorf("scope: invalid start pattern: %w", err)
		}
	}
	if endPattern != "" {
		if e, err = regexp.Compile(endPattern); err != nil {
			return nil, fmt.Errorf("scope: invalid end pattern: %w", err)
		}
	}
	return &Scope{start: s, end: e, enabled: true}, nil
}

// Enabled reports whether the scope filter is active.
func (sc *Scope) Enabled() bool { return sc.enabled }

// Allow returns true if the line should be emitted.
func (sc *Scope) Allow(line string) bool {
	if !sc.enabled {
		return true
	}
	if !sc.active && sc.start != nil && sc.start.MatchString(line) {
		sc.active = true
	}
	if !sc.active {
		return false
	}
	emit := true
	if sc.end != nil && sc.end.MatchString(line) {
		sc.active = false
	}
	return emit
}

// Reset returns the Scope to its initial (inactive) state.
func (sc *Scope) Reset() { sc.active = false }

// Package after provides a transformer that emits only lines
// that appear after a matching trigger line, with an optional
// maximum count of lines to emit.
package after

import "regexp"

// After emits lines that follow a line matching a trigger pattern.
type After struct {
	re      *regexp.Regexp
	max     int
	count   int
	active  bool
	enabled bool
}

// New creates a new After filter. pattern is the trigger regex;
// max is the maximum number of lines to emit after the match
// (0 means unlimited). Returns an error if the pattern is invalid.
func New(pattern string, max int) (*After, error) {
	if pattern == "" {
		return &After{enabled: false}, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	if max < 0 {
		max = 0
	}
	return &After{
		re:      re,
		max:     max,
		enabled: true,
	}, nil
}

// Enabled reports whether the After filter is active.
func (a *After) Enabled() bool { return a.enabled }

// Allow returns true if the line should be emitted.
// The trigger line itself is not emitted; only subsequent lines are.
func (a *After) Allow(line string) bool {
	if !a.enabled {
		return true
	}
	if a.re.MatchString(line) {
		a.active = true
		a.count = 0
		return false
	}
	if !a.active {
		return false
	}
	if a.max > 0 && a.count >= a.max {
		a.active = false
		a.count = 0
		return false
	}
	a.count++
	return true
}

// Reset clears the active state, stopping emission until the next trigger.
func (a *After) Reset() {
	a.active = false
	a.count = 0
}

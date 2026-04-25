// Package since provides a filter that discards log lines whose embedded
// timestamp falls before a configurable cut-off time.
package since

import (
	"time"
)

// Filter drops lines whose parsed timestamp is before the cut-off.
// Lines that contain no recognisable timestamp are always allowed through.
type Filter struct {
	enabled bool
	cutoff  time.Time
	layouts []string
}

// defaultLayouts are tried in order when parsing a timestamp from a line.
var defaultLayouts = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05",
	"02/Jan/2006:15:04:05 -0700",
}

// New returns a Filter that allows only lines at or after cutoff.
// If cutoff is zero the filter is disabled and all lines are allowed.
func New(cutoff time.Time, layouts ...string) *Filter {
	f := &Filter{
		enabled: !cutoff.IsZero(),
		cutoff:  cutoff,
		layouts: defaultLayouts,
	}
	if len(layouts) > 0 {
		f.layouts = layouts
	}
	return f
}

// Enabled reports whether the filter is active.
func (f *Filter) Enabled() bool { return f.enabled }

// Allow returns true when the line should be emitted.
// A line is allowed when:
//   - the filter is disabled, OR
//   - no timestamp can be parsed from the line, OR
//   - the parsed timestamp is at or after the cut-off.
func (f *Filter) Allow(line string) bool {
	if !f.enabled {
		return true
	}
	t, ok := f.parse(line)
	if !ok {
		return true
	}
	return !t.Before(f.cutoff)
}

// parse attempts to extract a timestamp from line by trying each layout
// against every possible sub-string starting position.
func (f *Filter) parse(line string) (time.Time, bool) {
	for _, layout := range f.layouts {
		width := len(layout)
		if width > len(line) {
			continue
		}
		for i := 0; i <= len(line)-width; i++ {
			t, err := time.Parse(layout, line[i:i+width])
			if err == nil {
				return t, true
			}
		}
	}
	return time.Time{}, false
}

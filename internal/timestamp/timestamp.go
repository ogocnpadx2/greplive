// Package timestamp provides line timestamping utilities that prepend
// a formatted timestamp to each log line as it passes through the pipeline.
package timestamp

import (
	"strings"
	"time"
)

// Stamper prepends a timestamp to each line.
type Stamper struct {
	format  string
	enabled bool
	clock   func() time.Time
}

// Option configures a Stamper.
type Option func(*Stamper)

// WithClock overrides the clock used by the Stamper (useful for testing).
func WithClock(fn func() time.Time) Option {
	return func(s *Stamper) { s.clock = fn }
}

// New creates a Stamper with the given Go time format string.
// If format is empty the Stamper is disabled and lines pass through unchanged.
func New(format string, opts ...Option) *Stamper {
	s := &Stamper{
		format:  format,
		enabled: format != "",
		clock:   time.Now,
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Enabled reports whether timestamping is active.
func (s *Stamper) Enabled() bool { return s.enabled }

// Apply prepends the current timestamp to line, or returns line unchanged
// when the Stamper is disabled.
func (s *Stamper) Apply(line string) string {
	if !s.enabled {
		return line
	}
	var b strings.Builder
	b.WriteString(s.clock().Format(s.format))
	b.WriteByte(' ')
	b.WriteString(line)
	return b.String()
}

// ApplyAll applies each Stamper in turn to line.
func ApplyAll(line string, stampers []*Stamper) string {
	for _, st := range stampers {
		line = st.Apply(line)
	}
	return line
}

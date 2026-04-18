// Package linenum prepends a formatted line number to each log line.
package linenum

import "fmt"

// Numberer tracks a running line count and optionally prepends it to lines.
type Numberer struct {
	enabled bool
	pad     int
	count   int64
}

// New creates a Numberer. If pad > 0 the number is zero-padded to that width.
func New(enabled bool, pad int) *Numberer {
	if pad <= 0 {
		pad = 0
	}
	return &Numberer{enabled: enabled, pad: pad}
}

// Enabled reports whether line numbering is active.
func (n *Numberer) Enabled() bool { return n.enabled }

// Apply increments the internal counter and, when enabled, prepends the
// current line number to line. The original string is returned unchanged
// when disabled.
func (n *Numberer) Apply(line string) string {
	n.count++
	if !n.enabled {
		return line
	}
	var prefix string
	if n.pad > 0 {
		prefix = fmt.Sprintf("%0*d ", n.pad, n.count)
	} else {
		prefix = fmt.Sprintf("%d ", n.count)
	}
	return prefix + line
}

// Reset zeroes the internal counter.
func (n *Numberer) Reset() { n.count = 0 }

// Count returns the number of lines processed so far.
func (n *Numberer) Count() int64 { return n.count }

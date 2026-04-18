// Package indent provides line indentation for structured log output.
package indent

import "strings"

// Indenter prepends a fixed string to each line.
type Indenter struct {
	prefix  string
	enabled bool
}

// New creates an Indenter with the given prefix string.
// If prefix is empty the Indenter is disabled and lines pass through unchanged.
func New(prefix string) *Indenter {
	return &Indenter{
		prefix:  prefix,
		enabled: prefix != "",
	}
}

// Enabled reports whether indentation is active.
func (i *Indenter) Enabled() bool { return i.enabled }

// Apply prepends the prefix to line when enabled.
func (i *Indenter) Apply(line string) string {
	if !i.enabled {
		return line
	}
	return i.prefix + line
}

// ApplyAll applies each Indenter in order to line.
func ApplyAll(line string, indenters []*Indenter) string {
	for _, ind := range indenters {
		line = ind.Apply(line)
	}
	return line
}

// Repeat returns a prefix string made of n repetitions of unit (e.g. "  ").
func Repeat(unit string, n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(unit, n)
}

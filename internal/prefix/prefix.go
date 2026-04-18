// Package prefix prepends a static or dynamic string to each log line.
package prefix

import "strings"

// Prefixer prepends a string to every line it processes.
type Prefixer struct {
	prefix string
	enabled bool
}

// New returns a Prefixer that prepends p to every line.
// If p is empty the Prefixer is a no-op.
func New(p string) *Prefixer {
	return &Prefixer{
		prefix:  p,
		enabled: p != "",
	}
}

// Enabled reports whether the prefixer will modify lines.
func (p *Prefixer) Enabled() bool { return p.enabled }

// Apply prepends the configured prefix to line.
// If the Prefixer is disabled the original line is returned unchanged.
func (p *Prefixer) Apply(line string) string {
	if !p.enabled {
		return line
	}
	var b strings.Builder
	b.WriteString(p.prefix)
	b.WriteString(line)
	return b.String()
}

// ApplyAll applies each Prefixer in ps to line in order.
func ApplyAll(line string, ps []*Prefixer) string {
	for _, p := range ps {
		line = p.Apply(line)
	}
	return line
}

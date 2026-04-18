// Package labelmap attaches static key=value labels to every log line.
package labelmap

import "strings"

// Labeler prepends a fixed set of labels to each line.
type Labeler struct {
	prefix string
	enabled bool
}

// New creates a Labeler from a map of label key/value pairs.
// Labels are rendered as "key=value" pairs separated by spaces and
// prepended to the line inside square brackets, e.g. "[env=prod svc=api] line".
func New(labels map[string]string) *Labeler {
	if len(labels) == 0 {
		return &Labeler{}
	}
	parts := make([]string, 0, len(labels))
	for k, v := range labels {
		parts = append(parts, k+"="+v)
	}
	return &Labeler{
		prefix:  "[" + strings.Join(parts, " ") + "] ",
		enabled: true,
	}
}

// Enabled reports whether any labels are configured.
func (l *Labeler) Enabled() bool { return l.enabled }

// Apply prepends the label prefix to line.
func (l *Labeler) Apply(line string) string {
	if !l.enabled {
		return line
	}
	return l.prefix + line
}

// ApplyAll applies each Labeler in sequence to line.
func ApplyAll(line string, labelers []*Labeler) string {
	for _, lb := range labelers {
		line = lb.Apply(line)
	}
	return line
}

// Package timefmt provides timestamp formatting helpers for log output.
package timefmt

import (
	"time"
)

// Format defines a named timestamp layout.
type Format struct {
	name   string
	layout string
}

// Name returns the format name.
func (f Format) Name() string { return f.name }

// Layout returns the Go time layout string.
func (f Format) Layout() string { return f.layout }

// Format formats t using the layout.
func (f Format) Format(t time.Time) string {
	return t.Format(f.layout)
}

// Predefined formats.
var (
	RFC3339   = Format{name: "rfc3339", layout: time.RFC3339}
	RFC3339Ms = Format{name: "rfc3339ms", layout: "2006-01-02T15:04:05.000Z07:00"}
	Short     = Format{name: "short", layout: "15:04:05"}
	ShortMs   = Format{name: "shortms", layout: "15:04:05.000"}
	Date      = Format{name: "date", layout: "2006-01-02 15:04:05"}
)

var known = []Format{RFC3339, RFC3339Ms, Short, ShortMs, Date}

// Parse returns the Format matching name (case-insensitive).
// Returns RFC3339 and false if not found.
func Parse(name string) (Format, bool) {
	for _, f := range known {
		if f.name == name {
			return f, true
		}
	}
	return RFC3339, false
}

// Names returns all known format names.
func Names() []string {
	out := make([]string, len(known))
	for i, f := range known {
		out[i] = f.name
	}
	return out
}

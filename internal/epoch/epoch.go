// Package epoch converts Unix epoch timestamps embedded in log lines
// into human-readable time strings.
package epoch

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Converter finds Unix epoch values in a line and replaces them with
// formatted time strings.
type Converter struct {
	re     *regexp.Regexp
	fmt    string
	utc    bool
	enabled bool
}

// New returns a Converter that replaces epoch timestamps matching the
// default pattern (a 10- or 13-digit integer) with times formatted
// according to layout.  If layout is empty the converter is disabled.
func New(layout string, utc bool) *Converter {
	if layout == "" {
		return &Converter{}
	}
	return &Converter{
		re:      regexp.MustCompile(`\b(\d{10}|\d{13})\b`),
		fmt:     layout,
		utc:     utc,
		enabled: true,
	}
}

// Enabled reports whether the converter will modify lines.
func (c *Converter) Enabled() bool { return c.enabled }

// Apply replaces epoch timestamps in line with formatted time strings.
func (c *Converter) Apply(line string) string {
	if !c.enabled {
		return line
	}
	return c.re.ReplaceAllStringFunc(line, func(s string) string {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return s
		}
		var t time.Time
		if len(s) == 13 {
			t = time.UnixMilli(n)
		} else {
			t = time.Unix(n, 0)
		}
		if c.utc {
			t = t.UTC()
		}
		return fmt.Sprintf("%s", t.Format(c.fmt))
	})
}

// ApplyAll applies each Converter in cs to line in order.
func ApplyAll(line string, cs []*Converter) string {
	for _, c := range cs {
		line = c.Apply(line)
	}
	return line
}

// ParseLayout maps a short name to a Go time layout string.
// Unknown names are returned unchanged so callers may pass raw layouts.
func ParseLayout(name string) string {
	switch strings.ToLower(name) {
	case "rfc3339":
		return time.RFC3339
	case "rfc3339ms":
		return "2006-01-02T15:04:05.000Z07:00"
	case "short":
		return "2006-01-02 15:04:05"
	case "date":
		return "2006-01-02"
	default:
		return name
	}
}

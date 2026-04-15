// Package highlight provides regex-based term highlighting for log lines.
package highlight

import (
	"regexp"
	"strings"
)

// Highlighter applies color highlighting to matched terms within a line.
type Highlighter struct {
	pattern *regexp.Regexp
	color   string
}

// ANSI color codes for highlighting.
const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Red    = "\033[31m"
)

// New creates a Highlighter for the given regex pattern and ANSI color code.
// Returns an error if the pattern is invalid.
func New(pattern, color string) (*Highlighter, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &Highlighter{pattern: re, color: color}, nil
}

// Apply wraps all matches of the pattern in the line with the configured color.
func (h *Highlighter) Apply(line string) string {
	if h == nil || h.pattern == nil {
		return line
	}
	return h.pattern.ReplaceAllStringFunc(line, func(match string) string {
		return h.color + Bold + match + Reset
	})
}

// ApplyAll applies multiple highlighters to a line in order.
func ApplyAll(line string, highlighters []*Highlighter) string {
	result := line
	for _, h := range highlighters {
		result = h.Apply(result)
	}
	return result
}

// StripANSI removes ANSI escape sequences from a string.
func StripANSI(s string) string {
	ansi := regexp.MustCompile(`\033\[[0-9;]*m`)
	return ansi.ReplaceAllString(s, "")
}

// ColorByName maps a color name string to its ANSI code.
func ColorByName(name string) string {
	switch strings.ToLower(name) {
	case "yellow":
		return Yellow
	case "cyan":
		return Cyan
	case "red":
		return Red
	default:
		return Yellow
	}
}

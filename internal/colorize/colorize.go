// Package colorize provides per-pattern colorisation of log lines using
// ANSI escape codes. Each Colorizer wraps a compiled regular expression and
// applies a fixed color to every substring that matches.
package colorize

import (
	"fmt"
	"regexp"
)

// ANSI color codes supported by ColorByName.
var namedColors = map[string]string{
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",
	"white":   "\033[37m",
	"bold":    "\033[1m",
}

const reset = "\033[0m"

// Colorizer colorises substrings of a log line that match a regex.
type Colorizer struct {
	re    *regexp.Regexp
	color string
}

// New compiles pattern and looks up colorName in the built-in palette.
// An empty pattern returns a disabled (no-op) Colorizer without an error.
func New(pattern, colorName string) (*Colorizer, error) {
	if pattern == "" {
		return &Colorizer{}, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("colorize: invalid pattern %q: %w", pattern, err)
	}
	code, ok := namedColors[colorName]
	if !ok {
		code = namedColors["white"]
	}
	return &Colorizer{re: re, color: code}, nil
}

// Enabled reports whether the Colorizer has an active pattern.
func (c *Colorizer) Enabled() bool { return c.re != nil }

// Apply wraps every match in line with the configured ANSI color code.
// If the Colorizer is disabled the original line is returned unchanged.
func (c *Colorizer) Apply(line string) string {
	if !c.Enabled() {
		return line
	}
	return c.re.ReplaceAllStringFunc(line, func(match string) string {
		return c.color + match + reset
	})
}

// ApplyAll runs each Colorizer in cs over line in order.
func ApplyAll(line string, cs []*Colorizer) string {
	for _, c := range cs {
		line = c.Apply(line)
	}
	return line
}

// ColorByName returns the ANSI escape sequence for a named color, or an empty
// string when the name is not recognised.
func ColorByName(name string) string { return namedColors[name] }

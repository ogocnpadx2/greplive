// Package fold provides case-folding for line matching.
// When enabled every comparison is performed in a case-insensitive manner
// by lower-casing both the line and the pattern before matching.
package fold

import (
	"regexp"
	"strings"
)

// Folder wraps a compiled regexp and applies case-insensitive folding.
type Folder struct {
	enabled bool
	re      *regexp.Regexp
}

// New returns a Folder for the given pattern.
// If caseInsensitive is false the Folder is a no-op pass-through.
// An error is returned when the pattern fails to compile.
func New(pattern string, caseInsensitive bool) (*Folder, error) {
	if pattern == "" || !caseInsensitive {
		return &Folder{enabled: false}, nil
	}
	re, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		return nil, err
	}
	return &Folder{enabled: true, re: re}, nil
}

// Enabled reports whether case-folding is active.
func (f *Folder) Enabled() bool { return f.enabled }

// Match reports whether line contains a case-insensitive match for the pattern.
// When the Folder is disabled it always returns true so the caller can use it
// as an unconditional pass-through.
func (f *Folder) Match(line string) bool {
	if !f.enabled {
		return true
	}
	return f.re.MatchString(line)
}

// Fold returns the lower-cased form of s.
// Useful when callers need to normalise a string independently.
func Fold(s string) string { return strings.ToLower(s) }

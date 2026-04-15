// Package truncate provides line-length limiting for streamed log output.
// Lines exceeding the configured maximum are truncated and annotated with
// a configurable suffix so the reader is aware content was cut.
package truncate

import (
	"strings"
	"unicode/utf8"
)

const (
	// DefaultMaxRunes is the default maximum rune count per line.
	DefaultMaxRunes = 512
	// DefaultSuffix is appended to lines that have been truncated.
	DefaultSuffix = " …"
)

// Truncator holds the configuration for line truncation.
type Truncator struct {
	maxRunes int
	suffix   string
}

// New returns a Truncator that cuts lines at maxRunes runes and appends
// suffix. If maxRunes is zero or negative no truncation is applied.
// If suffix is empty DefaultSuffix is used.
func New(maxRunes int, suffix string) *Truncator {
	if suffix == "" {
		suffix = DefaultSuffix
	}
	return &Truncator{
		maxRunes: maxRunes,
		suffix:   suffix,
	}
}

// Apply returns the (possibly truncated) form of line.
// ANSI escape sequences are counted as zero-width so that colour codes do
// not consume the rune budget.
func (t *Truncator) Apply(line string) string {
	if t.maxRunes <= 0 {
		return line
	}

	count := 0
	for i, r := range line {
		// Skip ANSI CSI sequences: ESC '[' … final-byte (0x40-0x7E).
		if r == '\x1b' && i+1 < len(line) && line[i+1] == '[' {
			continue
		}
		_ = r
		count++
		if count > t.maxRunes {
			// Find the byte offset of the (maxRunes+1)-th rune.
			byteIdx := runeOffset(line, t.maxRunes)
			return line[:byteIdx] + t.suffix
		}
	}
	return line
}

// runeOffset returns the byte index of the n-th rune in s.
func runeOffset(s string, n int) int {
	for i := range s {
		if n == 0 {
			return i
		}
		_, size := utf8.DecodeRuneInString(s[i:])
		_ = size
		n--
	}
	return len(s)
}

// ApplyAll applies t.Apply to every element of lines in place and returns
// the slice for convenience.
func (t *Truncator) ApplyAll(lines []string) []string {
	for i, l := range lines {
		lines[i] = t.Apply(l)
	}
	return lines
}

// Enabled reports whether truncation is active.
func (t *Truncator) Enabled() bool { return t.maxRunes > 0 }

// MaxRunes returns the configured rune limit.
func (t *Truncator) MaxRunes() int { return t.maxRunes }

// Suffix returns the configured truncation suffix.
func (t *Truncator) Suffix() string { return t.suffix }

// StripSuffix removes the default suffix from s if present, returning the
// cleaned string and whether a suffix was found.
func StripSuffix(s string) (string, bool) {
	if strings.HasSuffix(s, DefaultSuffix) {
		return strings.TrimSuffix(s, DefaultSuffix), true
	}
	return s, false
}

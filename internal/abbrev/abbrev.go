// Package abbrev provides line abbreviation by collapsing long repeated
// whitespace runs and optionally trimming a line to a maximum rune width.
package abbrev

import (
	"strings"
	"unicode/utf8"
)

// Abbreviator shortens lines according to configured rules.
type Abbreviator struct {
	maxRunes      int
	suffix        string
	collapseSpace bool
	enabled       bool
}

// New returns an Abbreviator. maxRunes <= 0 disables length truncation.
// collapseSpace replaces runs of whitespace with a single space.
func New(maxRunes int, suffix string, collapseSpace bool) *Abbreviator {
	if suffix == "" {
		suffix = "…"
	}
	enabled := maxRunes > 0 || collapseSpace
	return &Abbreviator{
		maxRunes:      maxRunes,
		suffix:        suffix,
		collapseSpace: collapseSpace,
		enabled:       enabled,
	}
}

// Enabled reports whether any abbreviation step is active.
func (a *Abbreviator) Enabled() bool { return a.enabled }

// Apply abbreviates line according to the configured rules.
func (a *Abbreviator) Apply(line string) string {
	if !a.enabled {
		return line
	}
	if a.collapseSpace {
		line = collapseWhitespace(line)
	}
	if a.maxRunes > 0 {
		line = truncate(line, a.maxRunes, a.suffix)
	}
	return line
}

// ApplyAll applies each Abbreviator in order to line.
func ApplyAll(line string, abbrevs []*Abbreviator) string {
	for _, ab := range abbrevs {
		line = ab.Apply(line)
	}
	return line
}

// collapseWhitespace replaces every run of whitespace characters with a
// single ASCII space and trims leading/trailing whitespace.
func collapseWhitespace(s string) string {
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}

// truncate shortens s to at most maxRunes runes, appending suffix when
// the string is actually shortened.
func truncate(s string, maxRunes int, suffix string) string {
	if utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	suffixRunes := utf8.RuneCountInString(suffix)
	keep := maxRunes - suffixRunes
	if keep <= 0 {
		return suffix
	}
	var sb strings.Builder
	count := 0
	for _, r := range s {
		if count >= keep {
			break
		}
		sb.WriteRune(r)
		count++
	}
	sb.WriteString(suffix)
	return sb.String()
}

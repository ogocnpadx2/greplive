// Package mask provides field-level masking of structured log values.
// It replaces matched key=value or JSON field values with a fixed placeholder.
package mask

import (
	"fmt"
	"regexp"
	"strings"
)

const defaultPlaceholder = "***"

// Masker replaces the value of a named field in a log line.
type Masker struct {
	re          *regexp.Regexp
	placeholder string
	replacement string
}

// New creates a Masker that redacts values for the given field name.
// fieldPattern is a regex matching the field name (e.g. "password|secret").
// placeholder overrides the default replacement text; pass "" for default.
func New(fieldPattern, placeholder string) (*Masker, error) {
	if fieldPattern == "" {
		return nil, fmt.Errorf("mask: field pattern must not be empty")
	}
	if placeholder == "" {
		placeholder = defaultPlaceholder
	}
	// Matches key=value and key="value" forms as well as JSON "key":"value".
	pattern := fmt.Sprintf(
		`(?i)(%s)(\s*[:=]\s*)("?)([^"\s,}]+)("?)`,
		fieldPattern,
	)
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("mask: invalid field pattern: %w", err)
	}
	return &Masker{
		re:          re,
		placeholder: placeholder,
		replacement: placeholder,
	}, nil
}

// Apply returns line with matched field values replaced by the placeholder.
func (m *Masker) Apply(line string) string {
	return m.re.ReplaceAllStringFunc(line, func(match string) string {
		subs := m.re.FindStringSubmatch(match)
		if len(subs) < 6 {
			return match
		}
		// subs: [full, key, sep, openQuote, value, closeQuote]
		return subs[1] + subs[2] + subs[3] + m.replacement + subs[5]
	})
}

// ApplyAll runs every Masker in order over line.
func ApplyAll(maskers []*Masker, line string) string {
	for _, m := range maskers {
		line = m.Apply(line)
	}
	return line
}

// Fields returns a comma-separated summary of active field patterns.
func (m *Masker) Fields() string {
	return strings.TrimPrefix(
		strings.TrimSuffix(m.re.String(), `)(\s*[:=]\s*)("?)([^"\s,}]+)("?)`),
		`(?i)(`)
}

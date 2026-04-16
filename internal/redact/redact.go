// Package redact provides regex-based redaction of sensitive data in log lines.
package redact

import (
	"fmt"
	"regexp"
)

// Redactor replaces matches of a compiled pattern with a fixed replacement string.
type Redactor struct {
	re          *regexp.Regexp
	replacement string
}

// New compiles pattern and returns a Redactor that substitutes matches with
// replacement. An error is returned if the pattern is invalid.
func New(pattern, replacement string) (*Redactor, error) {
	if pattern == "" {
		return nil, fmt.Errorf("redact: pattern must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("redact: invalid pattern %q: %w", pattern, err)
	}
	if replacement == "" {
		replacement = "[REDACTED]"
	}
	return &Redactor{re: re, replacement: replacement}, nil
}

// Apply returns a copy of line with all pattern matches replaced.
func (r *Redactor) Apply(line string) string {
	return r.re.ReplaceAllString(line, r.replacement)
}

// ApplyAll runs each Redactor in order against line and returns the result.
func ApplyAll(redactors []*Redactor, line string) string {
	for _, r := range redactors {
		line = r.Apply(line)
	}
	return line
}

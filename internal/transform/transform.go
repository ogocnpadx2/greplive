// Package transform provides line transformation functions such as
// trimming whitespace, replacing substrings, and applying regex-based
// substitutions before lines are written to output.
package transform

import (
	"regexp"
	"strings"
)

// Transformer applies a single transformation to a log line.
type Transformer interface {
	Apply(line string) string
}

// TrimTransformer removes leading and trailing whitespace from a line.
type TrimTransformer struct{}

// Apply trims whitespace from both ends of line.
func (t TrimTransformer) Apply(line string) string {
	return strings.TrimSpace(line)
}

// ReplaceTransformer performs a literal string replacement.
type ReplaceTransformer struct {
	Old string
	New string
}

// Apply replaces all occurrences of Old with New.
func (r ReplaceTransformer) Apply(line string) string {
	if r.Old == "" {
		return line
	}
	return strings.ReplaceAll(line, r.Old, r.New)
}

// RegexTransformer replaces regex matches with a substitution string.
type RegexTransformer struct {
	pattern *regexp.Regexp
	repl    string
}

// NewRegex compiles pattern and returns a RegexTransformer.
// Returns an error if the pattern is invalid.
func NewRegex(pattern, repl string) (*RegexTransformer, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &RegexTransformer{pattern: re, repl: repl}, nil
}

// Apply replaces all regex matches in line with the substitution string.
func (r *RegexTransformer) Apply(line string) string {
	return r.pattern.ReplaceAllString(line, r.repl)
}

// Chain applies a slice of Transformers in order, passing the output of
// each as the input to the next.
func Chain(line string, transformers []Transformer) string {
	for _, t := range transformers {
		line = t.Apply(line)
	}
	return line
}

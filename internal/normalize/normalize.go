// Package normalize provides line normalization utilities that strip or
// replace common noise patterns before further processing.
package normalize

import (
	"strings"
	"unicode"
)

// Normalizer applies a chain of normalization steps to a log line.
type Normalizer struct {
	steps []func(string) string
}

// Option is a functional option for configuring a Normalizer.
type Option func(*Normalizer)

// WithTrimSpace trims leading and trailing whitespace from each line.
func WithTrimSpace() Option {
	return func(n *Normalizer) {
		n.steps = append(n.steps, strings.TrimSpace)
	}
}

// WithCollapseSpaces replaces runs of whitespace with a single space.
func WithCollapseSpaces() Option {
	return func(n *Normalizer) {
		n.steps = append(n.steps, func(s string) string {
			var b strings.Builder
			prevSpace := false
			for _, r := range s {
				if unicode.IsSpace(r) {
					if !prevSpace {
						b.WriteRune(' ')
					}
					prevSpace = true
				} else {
					b.WriteRune(r)
					prevSpace = false
				}
			}
			return b.String()
		})
	}
}

// WithLowercase converts all characters to lowercase.
func WithLowercase() Option {
	return func(n *Normalizer) {
		n.steps = append(n.steps, strings.ToLower)
	}
}

// WithReplace replaces all occurrences of old with new.
func WithReplace(old, new string) Option {
	return func(n *Normalizer) {
		n.steps = append(n.steps, func(s string) string {
			return strings.ReplaceAll(s, old, new)
		})
	}
}

// New creates a Normalizer with the given options.
func New(opts ...Option) *Normalizer {
	n := &Normalizer{}
	for _, o := range opts {
		o(n)
	}
	return n
}

// Apply runs all normalization steps on line and returns the result.
func (n *Normalizer) Apply(line string) string {
	for _, step := range n.steps {
		line = step(line)
	}
	return line
}

// Enabled reports whether any normalization steps are configured.
func (n *Normalizer) Enabled() bool {
	return len(n.steps) > 0
}

// Package bracket provides a transformer that wraps matched substrings
// or entire lines in configurable left/right bracket strings.
package bracket

import (
	"regexp"
	"strings"
)

// Bracket wraps matched text in a left and right bracket string.
type Bracket struct {
	re      *regexp.Regexp
	left    string
	right   string
	enabled bool
}

// New creates a Bracket transformer. If pattern is empty the transformer is
// disabled and Apply returns lines unchanged. left and right default to "["
// and "]" when empty.
func New(pattern, left, right string) (*Bracket, error) {
	if pattern == "" {
		return &Bracket{}, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	if left == "" {
		left = "["
	}
	if right == "" {
		right = "]"
	}
	return &Bracket{re: re, left: left, right: right, enabled: true}, nil
}

// Enabled reports whether the transformer is active.
func (b *Bracket) Enabled() bool { return b.enabled }

// Apply wraps every non-overlapping match in the line with the bracket
// strings. If the transformer is disabled the original line is returned.
func (b *Bracket) Apply(line string) string {
	if !b.enabled {
		return line
	}
	var sb strings.Builder
	last := 0
	for _, loc := range b.re.FindAllStringIndex(line, -1) {
		sb.WriteString(line[last:loc[0]])
		sb.WriteString(b.left)
		sb.WriteString(line[loc[0]:loc[1]])
		sb.WriteString(b.right)
		last = loc[1]
	}
	sb.WriteString(line[last:])
	return sb.String()
}

// ApplyAll applies each Bracket transformer in order to line.
func ApplyAll(brackets []*Bracket, line string) string {
	for _, b := range brackets {
		line = b.Apply(line)
	}
	return line
}

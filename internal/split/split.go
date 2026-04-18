// Package split provides a transformer that splits a line into fields
// by a configurable delimiter and rejoins a subset of them.
package split

import "strings"

// Splitter splits a line by Delimiter and keeps only the fields at the
// given Indices (0-based). If Indices is empty all fields are kept.
// Fields are rejoined with Join (defaults to a single space).
type Splitter struct {
	delimiter string
	indices   []int
	join      string
	enabled   bool
}

// New creates a Splitter. delimiter must be non-empty.
// indices may be nil/empty to keep all fields.
// join is the string used when reassembling; defaults to " ".
func New(delimiter string, indices []int, join string) (*Splitter, error) {
	if delimiter == "" {
		return &Splitter{}, nil
	}
	if join == "" {
		join = " "
	}
	copy := make([]int, len(indices))
	copy = append(copy[:0], indices...)
	return &Splitter{
		delimiter: delimiter,
		indices:   copy,
		join:      join,
		enabled:   true,
	}, nil
}

// Enabled reports whether the splitter will modify lines.
func (s *Splitter) Enabled() bool { return s.enabled }

// Apply splits line by the delimiter, selects the configured fields and
// rejoins them. If the line cannot be split (no delimiter found) it is
// returned unchanged.
func (s *Splitter) Apply(line string) string {
	if !s.enabled {
		return line
	}
	parts := strings.Split(line, s.delimiter)
	if len(s.indices) == 0 {
		return strings.Join(parts, s.join)
	}
	out := make([]string, 0, len(s.indices))
	for _, idx := range s.indices {
		if idx >= 0 && idx < len(parts) {
			out = append(out, parts[idx])
		}
	}
	if len(out) == 0 {
		return line
	}
	return strings.Join(out, s.join)
}

// ApplyAll applies each Splitter in order to line.
func ApplyAll(splitters []*Splitter, line string) string {
	for _, sp := range splitters {
		line = sp.Apply(line)
	}
	return line
}

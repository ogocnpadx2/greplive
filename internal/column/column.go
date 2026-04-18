// Package column provides line column extraction and formatting.
// It splits a log line by a delimiter and emits only the selected columns.
package column

import (
	"strings"
)

// Extractor selects specific columns from a delimited log line.
type Extractor struct {
	delimiter string
	columns   []int
	enabled   bool
}

// New creates a new Extractor. delimiter is the field separator and columns
// is the zero-based list of column indices to keep. If columns is empty the
// Extractor is disabled and lines pass through unchanged.
func New(delimiter string, columns []int) *Extractor {
	if delimiter == "" {
		delimiter = " "
	}
	return &Extractor{
		delimiter: delimiter,
		columns:   columns,
		enabled:   len(columns) > 0,
	}
}

// Enabled reports whether column extraction is active.
func (e *Extractor) Enabled() bool { return e.enabled }

// Apply extracts the configured columns from line, joining them with the same
// delimiter. If the extractor is disabled, or a requested index is out of
// range, the original line is returned unchanged.
func (e *Extractor) Apply(line string) string {
	if !e.enabled {
		return line
	}
	parts := strings.Split(line, e.delimiter)
	out := make([]string, 0, len(e.columns))
	for _, idx := range e.columns {
		if idx < 0 || idx >= len(parts) {
			return line
		}
		out = append(out, parts[idx])
	}
	return strings.Join(out, e.delimiter)
}

// ApplyAll applies each extractor in order to line.
func ApplyAll(extractors []*Extractor, line string) string {
	for _, ex := range extractors {
		line = ex.Apply(line)
	}
	return line
}

// Count returns the number of fields in line when split by the extractor's
// delimiter. This is useful for inspecting line structure before selecting
// column indices.
func (e *Extractor) Count(line string) int {
	return len(strings.Split(line, e.delimiter))
}

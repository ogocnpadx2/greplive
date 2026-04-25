// Package fieldmap renames or aliases fields in key=value log lines.
package fieldmap

import (
	"strings"
)

// Mapper renames fields in key=value formatted log lines.
type Mapper struct {
	mappings map[string]string
	enabled  bool
}

// New creates a Mapper from a map of old→new field name pairs.
// If mappings is empty the mapper is disabled and Apply is a no-op.
func New(mappings map[string]string) *Mapper {
	clean := make(map[string]string, len(mappings))
	for k, v := range mappings {
		if k != "" && v != "" {
			clean[k] = v
		}
	}
	return &Mapper{
		mappings: clean,
		enabled:  len(clean) > 0,
	}
}

// Enabled reports whether any mappings are configured.
func (m *Mapper) Enabled() bool { return m.enabled }

// Apply rewrites any key=value pairs whose key appears in the mapping.
// Tokens that are not in key=value form are left untouched.
func (m *Mapper) Apply(line string) string {
	if !m.enabled {
		return line
	}
	tokens := strings.Fields(line)
	for i, tok := range tokens {
		eq := strings.IndexByte(tok, '=')
		if eq <= 0 {
			continue
		}
		key := tok[:eq]
		if newKey, ok := m.mappings[key]; ok {
			tokens[i] = newKey + tok[eq:]
		}
	}
	return strings.Join(tokens, " ")
}

// ApplyAll applies every mapper in ms to line in order.
func ApplyAll(ms []*Mapper, line string) string {
	for _, m := range ms {
		line = m.Apply(line)
	}
	return line
}

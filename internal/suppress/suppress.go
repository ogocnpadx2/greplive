// Package suppress drops log lines that match a given regex pattern,
// acting as the inverse of a grep/filter: lines that match are silently
// discarded while non-matching lines pass through unchanged.
package suppress

import "regexp"

// Suppressor drops lines whose text matches a compiled regular expression.
type Suppressor struct {
	re      *regexp.Regexp
	enabled bool
}

// New compiles pattern into a Suppressor. An empty pattern returns a
// disabled Suppressor that never drops any line. A non-empty but invalid
// pattern returns an error.
func New(pattern string) (*Suppressor, error) {
	if pattern == "" {
		return &Suppressor{}, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &Suppressor{re: re, enabled: true}, nil
}

// Enabled reports whether the suppressor has an active pattern.
func (s *Suppressor) Enabled() bool { return s.enabled }

// Drop returns true when the line should be discarded, i.e. the pattern
// matches. When the suppressor is disabled it always returns false.
func (s *Suppressor) Drop(line string) bool {
	if !s.enabled {
		return false
	}
	return s.re.MatchString(line)
}

// ApplyAll returns a copy of lines with any entry matched by at least one
// Suppressor removed. A nil or empty slice of suppressors returns lines
// unchanged.
func ApplyAll(suppressors []*Suppressor, lines []string) []string {
	if len(suppressors) == 0 {
		return lines
	}
	out := lines[:0:len(lines)]
	for _, l := range lines {
		dropped := false
		for _, s := range suppressors {
			if s.Drop(l) {
				dropped = true
				break
			}
		}
		if !dropped {
			out = append(out, l)
		}
	}
	return out
}

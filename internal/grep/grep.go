// Package grep provides multi-pattern line matching with named groups.
package grep

import "regexp"

// Matcher holds one or more compiled patterns.
type Matcher struct {
	patterns []*regexp.Regexp
	any      bool // true = OR logic, false = AND logic
}

// New compiles the supplied patterns. If any is true a line matches when at
// least one pattern matches; otherwise all patterns must match.
func New(patterns []string, any bool) (*Matcher, error) {
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		if p == "" {
			continue
		}
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, re)
	}
	return &Matcher{patterns: compiled, any: any}, nil
}

// Enabled reports whether the matcher has at least one pattern.
func (m *Matcher) Enabled() bool { return len(m.patterns) > 0 }

// Match reports whether line satisfies the matcher's logic.
// If the matcher has no patterns it always returns true.
func (m *Matcher) Match(line string) bool {
	if !m.Enabled() {
		return true
	}
	if m.any {
		for _, re := range m.patterns {
			if re.MatchString(line) {
				return true
			}
		}
		return false
	}
	for _, re := range m.patterns {
		if !re.MatchString(line) {
			return false
		}
	}
	return true
}

// Groups returns named sub-match groups for the first pattern that matches.
func (m *Matcher) Groups(line string) map[string]string {
	for _, re := range m.patterns {
		match := re.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		result := make(map[string]string)
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}
		return result
	}
	return nil
}

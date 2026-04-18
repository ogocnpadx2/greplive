// Package squash merges consecutive matching lines into a single output line.
package squash

import "regexp"

// Squasher merges consecutive lines that match a pattern into one.
type Squasher struct {
	re      *regexp.Regexp
	pending []string
	enabled bool
}

// New creates a Squasher for the given regex pattern.
// If pattern is empty, the Squasher is disabled and lines pass through unchanged.
func New(pattern string) (*Squasher, error) {
	if pattern == "" {
		return &Squasher{}, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &Squasher{re: re, enabled: true}, nil
}

// Enabled reports whether squashing is active.
func (s *Squasher) Enabled() bool { return s.enabled }

// Push accepts a line. If the line matches, it is buffered and ("" , false) is
// returned. When a non-matching line arrives after buffered matches, the merged
// line is flushed first, then the current line is returned.
// Returns (line, emit).
func (s *Squasher) Push(line string) (string, bool) {
	if !s.enabled {
		return line, true
	}
	if s.re.MatchString(line) {
		s.pending = append(s.pending, line)
		return "", false
	}
	if len(s.pending) > 0 {
		merged := s.flush()
		s.pending = s.pending[:0]
		// We must emit merged first; caller should call Flush after Push.
		_ = merged
		// Store the non-matching line to emit after flush.
		s.pending = append(s.pending[:0:0], line) // reuse trick: store non-match
		return merged, true
	}
	return line, true
}

// Flush returns any remaining buffered lines merged into one and clears the buffer.
func (s *Squasher) Flush() (string, bool) {
	if len(s.pending) == 0 {
		return "", false
	}
	merged := s.flush()
	s.pending = s.pending[:0]
	return merged, true
}

func (s *Squasher) flush() string {
	if len(s.pending) == 1 {
		return s.pending[0]
	}
	out := s.pending[0]
	for _, l := range s.pending[1:] {
		out += " | " + l
	}
	return out
}

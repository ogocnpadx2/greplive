// Package strip removes ANSI escape sequences and other non-printable
// control characters from log lines before further processing.
package strip

import (
	"regexp"
)

var (
	// ansiEscape matches ANSI CSI escape sequences (colours, cursor movement, etc.).
	ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	// ansiOSC matches OSC (Operating System Command) sequences.
	ansiOSC = regexp.MustCompile(`\x1b\][^\x07]*\x07`)
	// controlChars matches raw control characters except newline/tab.
	controlChars = regexp.MustCompile(`[\x00-\x08\x0b-\x0c\x0e-\x1f\x7f]`)
)

// Stripper removes escape sequences and control characters from lines.
type Stripper struct {
	ansi    bool
	control bool
	enabled bool
}

// New returns a Stripper configured by opts.
// If neither ANSI nor control stripping is requested the Stripper is disabled
// and Apply becomes a no-op.
func New(ansi, control bool) *Stripper {
	return &Stripper{
		ansi:    ansi,
		control: control,
		enabled: ansi || control,
	}
}

// Enabled reports whether the Stripper will modify lines.
func (s *Stripper) Enabled() bool { return s.enabled }

// Apply removes configured sequences from line and returns the result.
func (s *Stripper) Apply(line string) string {
	if !s.enabled {
		return line
	}
	if s.ansi {
		line = ansiOSC.ReplaceAllString(line, "")
		line = ansiEscape.ReplaceAllString(line, "")
	}
	if s.control {
		line = controlChars.ReplaceAllString(line, "")
	}
	return line
}

// ApplyAll runs all strippers in order against line.
func ApplyAll(line string, strippers []*Stripper) string {
	for _, s := range strippers {
		line = s.Apply(line)
	}
	return line
}

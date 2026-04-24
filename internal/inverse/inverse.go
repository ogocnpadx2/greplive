// Package inverse provides a line filter that inverts match logic,
// passing through lines that do NOT match a given regular expression.
// This mirrors grep's -v / --invert-match behaviour.
package inverse

import "regexp"

// Inverter drops lines that match the compiled pattern and passes
// through lines that do not match. When disabled (no pattern) every
// line is allowed.
type Inverter struct {
	re      *regexp.Regexp
	enabled bool
}

// New compiles pattern and returns an Inverter. An empty pattern
// returns a disabled (pass-through) Inverter. A non-empty but invalid
// pattern returns an error.
func New(pattern string) (*Inverter, error) {
	if pattern == "" {
		return &Inverter{}, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &Inverter{re: re, enabled: true}, nil
}

// Enabled reports whether the inverter has an active pattern.
func (inv *Inverter) Enabled() bool {
	return inv.enabled
}

// Allow returns true when the line should be passed downstream.
// A disabled Inverter always returns true.
// An enabled Inverter returns true only when the line does NOT match.
func (inv *Inverter) Allow(line string) bool {
	if !inv.enabled {
		return true
	}
	return !inv.re.MatchString(line)
}

// ApplyAll returns only the lines from input that are allowed by
// every Inverter in the slice. If inverters is empty all lines pass.
func ApplyAll(inverters []*Inverter, lines []string) []string {
	if len(inverters) == 0 {
		return lines
	}
	out := lines[:0:len(lines)]
	for _, l := range lines {
		allow := true
		for _, inv := range inverters {
			if !inv.Allow(l) {
				allow = false
				break
			}
		}
		if allow {
			out = append(out, l)
		}
	}
	return out
}

// Package fieldextract provides utilities for extracting named fields
// from structured log lines (e.g. key=value or JSON key:value pairs).
package fieldextract

import (
	"regexp"
	"strings"
)

// Extractor extracts named fields from a log line.
type Extractor struct {
	pattern *regexp.Regexp
	fields  []string
}

// New creates an Extractor from a regex pattern that must contain at least one
// named capture group (e.g. `(?P<level>\w+)`).
func New(pattern string) (*Extractor, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	var fields []string
	for _, name := range re.SubexpNames() {
		if name != "" {
			fields = append(fields, name)
		}
	}
	return &Extractor{pattern: re, fields: fields}, nil
}

// Extract returns a map of named capture group values found in line.
// Fields not present in the match are omitted from the result.
func (e *Extractor) Extract(line string) map[string]string {
	match := e.pattern.FindStringSubmatch(line)
	if match == nil {
		return nil
	}
	result := make(map[string]string, len(e.fields))
	for i, name := range e.pattern.SubexpNames() {
		if name != "" && i < len(match) && match[i] != "" {
			result[name] = match[i]
		}
	}
	return result
}

// Fields returns the named capture groups defined in the pattern.
func (e *Extractor) Fields() []string {
	out := make([]string, len(e.fields))
	copy(out, e.fields)
	return out
}

// ExtractKV parses simple key=value pairs from a line and returns them as a map.
func ExtractKV(line string) map[string]string {
	result := make(map[string]string)
	parts := strings.Fields(line)
	for _, p := range parts {
		idx := strings.IndexByte(p, '=')
		if idx <= 0 || idx == len(p)-1 {
			continue
		}
		key := p[:idx]
		val := strings.Trim(p[idx+1:], `"`)
		result[key] = val
	}
	return result
}

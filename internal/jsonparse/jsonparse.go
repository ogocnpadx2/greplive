// Package jsonparse provides utilities for detecting and pretty-printing
// JSON log lines, with optional field extraction and colorisation.
package jsonparse

import (
	"bytes"
	"encoding/json"
)

// Formatter pretty-prints JSON lines and optionally extracts a message field.
type Formatter struct {
	messageKey string
	pretty     bool
}

// Option configures a Formatter.
type Option func(*Formatter)

// WithMessageKey sets the JSON key used to extract the primary message.
func WithMessageKey(key string) Option {
	return func(f *Formatter) { f.messageKey = key }
}

// WithPretty enables indented JSON output.
func WithPretty(enabled bool) Option {
	return func(f *Formatter) { f.pretty = enabled }
}

// New returns a Formatter with the supplied options applied.
func New(opts ...Option) *Formatter {
	f := &Formatter{messageKey: "message"}
	for _, o := range opts {
		o(f)
	}
	return f
}

// IsJSON reports whether line is a valid JSON object.
func IsJSON(line string) bool {
	trimmed := []byte(line)
	trimmed = bytes.TrimSpace(trimmed)
	if len(trimmed) == 0 || trimmed[0] != '{' {
		return false
	}
	return json.Valid(trimmed)
}

// Format returns the formatted representation of line.
// If the line is not valid JSON it is returned unchanged.
func (f *Formatter) Format(line string) string {
	if !IsJSON(line) {
		return line
	}
	if f.pretty {
		var buf bytes.Buffer
		if err := json.Indent(&buf, []byte(line), "", "  "); err == nil {
			return buf.String()
		}
	}
	return line
}

// ExtractMessage returns the value of the configured message key from a JSON
// line, or the original line if the key is absent or the line is not JSON.
func (f *Formatter) ExtractMessage(line string) string {
	if !IsJSON(line) {
		return line
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return line
	}
	if v, ok := m[f.messageKey]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return line
}

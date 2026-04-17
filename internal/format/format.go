// Package format provides log line formatting utilities for greplive,
// supporting JSON and plain-text output modes.
package format

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Mode represents an output format mode.
type Mode string

const (
	ModePlain Mode = "plain"
	ModeJSON  Mode = "json"
)

// Entry holds the structured data for a single log line.
type Entry struct {
	Timestamp time.Time
	Level     string
	Message   string
	Fields    map[string]string
}

// Formatter formats a log Entry into a string.
type Formatter interface {
	Format(e Entry) string
}

// New returns a Formatter for the given mode.
// Defaults to plain if the mode is unrecognised.
func New(mode Mode) Formatter {
	switch mode {
	case ModeJSON:
		return &jsonFormatter{}
	default:
		return &plainFormatter{}
	}
}

// plainFormatter emits lines as "TIMESTAMP [LEVEL] MESSAGE key=value ...".
type plainFormatter struct{}

func (p *plainFormatter) Format(e Entry) string {
	var sb strings.Builder
	if !e.Timestamp.IsZero() {
		sb.WriteString(e.Timestamp.Format(time.RFC3339))
		sb.WriteByte(' ')
	}
	if e.Level != "" {
		fmt.Fprintf(&sb, "[%s] ", strings.ToUpper(e.Level))
	}
	sb.WriteString(e.Message)
	for k, v := range e.Fields {
		fmt.Fprintf(&sb, " %s=%s", k, v)
	}
	return sb.String()
}

// jsonFormatter emits lines as a JSON object.
type jsonFormatter struct{}

func (j *jsonFormatter) Format(e Entry) string {
	m := make(map[string]interface{}, len(e.Fields)+3)
	if !e.Timestamp.IsZero() {
		m["ts"] = e.Timestamp.Format(time.RFC3339)
	}
	if e.Level != "" {
		m["level"] = strings.ToLower(e.Level)
	}
	m["msg"] = e.Message
	for k, v := range e.Fields {
		m[k] = v
	}
	b, err := json.Marshal(m)
	if err != nil {
		return e.Message
	}
	return string(b)
}

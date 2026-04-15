package severity_test

import (
	"testing"

	"github.com/greplive/greplive/internal/severity"
)

func TestDetect(t *testing.T) {
	tests := []struct {
		line     string
		expected severity.Level
	}{
		{"2024/01/15 ERROR: connection refused", severity.Error},
		{"[WARN] disk usage above 90%", severity.Warn},
		{"INFO server started on :8080", severity.Info},
		{"debug mode enabled", severity.Debug},
		{"FATAL: out of memory", severity.Fatal},
		{"panic: runtime error", severity.Fatal},
		{"just a regular log line", severity.Unknown},
		{"failed to connect to database", severity.Error},
		{"warning: deprecated API used", severity.Warn},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			got := severity.Detect(tt.line)
			if got != tt.expected {
				t.Errorf("Detect(%q) = %v, want %v", tt.line, got, tt.expected)
			}
		})
	}
}

func TestColorize(t *testing.T) {
	line := "ERROR: something went wrong"
	colored := severity.Colorize(line, severity.Error)
	if colored == line {
		t.Error("expected colorized output to differ from plain line")
	}
	if len(colored) <= len(line) {
		t.Error("expected colorized output to be longer than plain line")
	}
}

func TestColorizeUnknown(t *testing.T) {
	line := "plain log line"
	colored := severity.Colorize(line, severity.Unknown)
	if colored != line {
		t.Errorf("expected Unknown level to return line unchanged, got %q", colored)
	}
}

func TestLevelString(t *testing.T) {
	tests := []struct {
		level    severity.Level
		expected string
	}{
		{severity.Debug, "DEBUG"},
		{severity.Info, "INFO"},
		{severity.Warn, "WARN"},
		{severity.Error, "ERROR"},
		{severity.Fatal, "FATAL"},
		{severity.Unknown, "UNKNOWN"},
	}
	for _, tt := range tests {
		if got := tt.level.String(); got != tt.expected {
			t.Errorf("Level(%d).String() = %q, want %q", tt.level, got, tt.expected)
		}
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected severity.Level
	}{
		{"debug", severity.Debug},
		{"INFO", severity.Info},
		{"Warning", severity.Warn},
		{"ERR", severity.Error},
		{"PANIC", severity.Fatal},
		{"garbage", severity.Unknown},
	}
	for _, tt := range tests {
		got := severity.ParseLevel(tt.input)
		if got != tt.expected {
			t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

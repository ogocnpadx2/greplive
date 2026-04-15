package filter_test

import (
	"testing"

	"github.com/user/greplive/internal/filter"
	"github.com/user/greplive/internal/severity"
)

func TestNewConfig_ValidPattern(t *testing.T) {
	cfg, err := filter.NewConfig("error", severity.LevelInfo, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
}

func TestNewConfig_InvalidPattern(t *testing.T) {
	_, err := filter.NewConfig("[invalid", severity.LevelInfo, false)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestNewConfig_EmptyPattern(t *testing.T) {
	cfg, err := filter.NewConfig("", severity.LevelDebug, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Pattern != nil {
		t.Error("expected nil pattern for empty string")
	}
}

func TestMatch_PatternAndLevel(t *testing.T) {
	cfg, _ := filter.NewConfig("timeout", severity.LevelWarn, false)

	tests := []struct {
		line    string
		wantMatch bool
	}{
		{"WARN: connection timeout occurred", true},
		{"ERROR: timeout while reading", true},
		{"INFO: timeout info message", false}, // INFO < WARN
		{"ERROR: something else", false},      // no pattern match
		{"DEBUG: timeout debug", false},       // DEBUG < WARN
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			got := cfg.Match(tt.line)
			if got != tt.wantMatch {
				t.Errorf("Match(%q) = %v, want %v", tt.line, got, tt.wantMatch)
			}
		})
	}
}

func TestMatch_NoPattern(t *testing.T) {
	cfg, _ := filter.NewConfig("", severity.LevelError, false)

	if cfg.Match("ERROR: something bad") != true {
		t.Error("expected ERROR line to match with no pattern")
	}
	if cfg.Match("INFO: something fine") != false {
		t.Error("expected INFO line to not match when minLevel is ERROR")
	}
}

func TestHighlight(t *testing.T) {
	cfg, _ := filter.NewConfig("error", severity.LevelDebug, false)
	line := "an error occurred"
	result := cfg.Highlight(line)
	if result == line {
		t.Error("expected highlighted output to differ from input")
	}
	if len(result) <= len(line) {
		t.Error("expected highlighted output to be longer due to escape codes")
	}
}

func TestHighlight_NoPattern(t *testing.T) {
	cfg, _ := filter.NewConfig("", severity.LevelDebug, false)
	line := "plain log line"
	if cfg.Highlight(line) != line {
		t.Error("expected unchanged line when no pattern is set")
	}
}

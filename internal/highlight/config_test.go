package highlight_test

import (
	"testing"

	"greplive/internal/highlight"
)

func TestConfig_Build_Valid(t *testing.T) {
	cfg := &highlight.Config{
		Terms: []highlight.TermConfig{
			{Pattern: `error`, Color: "red"},
			{Pattern: `warn`, Color: "yellow"},
		},
	}
	highlighters, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(highlighters) != 2 {
		t.Errorf("expected 2 highlighters, got %d", len(highlighters))
	}
}

func TestConfig_Build_InvalidPattern(t *testing.T) {
	cfg := &highlight.Config{
		Terms: []highlight.TermConfig{
			{Pattern: `[bad`, Color: "red"},
		},
	}
	_, err := cfg.Build()
	if err == nil {
		t.Fatal("expected error for invalid pattern, got nil")
	}
}

func TestConfig_Build_Empty(t *testing.T) {
	cfg := &highlight.Config{}
	highlighters, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(highlighters) != 0 {
		t.Errorf("expected 0 highlighters, got %d", len(highlighters))
	}
}

func TestDefaultConfig_Build(t *testing.T) {
	cfg := highlight.DefaultConfig()
	highlighters, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error building default config: %v", err)
	}
	if len(highlighters) == 0 {
		t.Error("expected at least one highlighter from default config")
	}
	line := "ERROR: something went wrong, warn level"
	result := highlight.ApplyAll(line, highlighters)
	stripped := highlight.StripANSI(result)
	if stripped != line {
		t.Errorf("stripped mismatch: got %q, want %q", stripped, line)
	}
	if result == line {
		t.Error("expected ANSI codes in result, but line was unchanged")
	}
}

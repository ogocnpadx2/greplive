package grep_test

import (
	"testing"

	"github.com/user/greplive/internal/grep"
)

func TestDefaultConfig_Build(t *testing.T) {
	cfg := grep.DefaultConfig()
	m, err := cfg.Build()
	if err != nil {
		t.Fatal(err)
	}
	if m.Enabled() {
		t.Fatal("default config should produce disabled matcher")
	}
}

func TestConfig_Build_ValidPatterns(t *testing.T) {
	cfg := grep.Config{Patterns: []string{"error", "warn"}, Any: true}
	m, err := cfg.Build()
	if err != nil {
		t.Fatal(err)
	}
	if !m.Enabled() {
		t.Fatal("expected enabled matcher")
	}
	if !m.Match("warn: disk full") {
		t.Error("expected match on warn")
	}
}

func TestConfig_Build_InvalidPattern(t *testing.T) {
	cfg := grep.Config{Patterns: []string{"[bad"}}
	_, err := cfg.Build()
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestConfig_Build_AndLogic(t *testing.T) {
	cfg := grep.Config{Patterns: []string{"foo", "bar"}, Any: false}
	m, err := cfg.Build()
	if err != nil {
		t.Fatal(err)
	}
	if m.Match("only foo here") {
		t.Error("AND logic: should not match with only one pattern")
	}
	if !m.Match("foo and bar together") {
		t.Error("AND logic: should match with both patterns")
	}
}

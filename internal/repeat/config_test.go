package repeat

import (
	"testing"
	"time"
)

func TestDefaultConfig_Build(t *testing.T) {
	r, err := DefaultConfig().Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Enabled() {
		t.Fatal("default config should be disabled")
	}
}

func TestConfig_Build_Valid(t *testing.T) {
	c := Config{Max: 5, Window: 30 * time.Second}
	r, err := c.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Enabled() {
		t.Fatal("expected enabled")
	}
}

func TestConfig_Build_NegativeMax(t *testing.T) {
	c := Config{Max: -1, Window: time.Minute}
	_, err := c.Build()
	if err == nil {
		t.Fatal("expected error for negative max")
	}
}

func TestConfig_Build_ZeroWindow_DefaultsToMinute(t *testing.T) {
	c := Config{Max: 2, Window: 0}
	r, err := c.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.window != time.Minute {
		t.Fatalf("expected window=1m, got %v", r.window)
	}
}

func TestConfig_Build_CustomWindow(t *testing.T) {
	c := Config{Max: 3, Window: 10 * time.Second}
	r, err := c.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.window != 10*time.Second {
		t.Fatalf("expected 10s window, got %v", r.window)
	}
}

package multiline

import "testing"

func TestDefaultConfig_Build(t *testing.T) {
	cfg := DefaultConfig()
	c, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil collector")
	}
}

func TestConfig_Build_InvalidStart(t *testing.T) {
	cfg := Config{StartPattern: "[invalid", MaxLines: 0}
	_, err := cfg.Build()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestConfig_Build_CustomMaxLines(t *testing.T) {
	cfg := Config{StartPattern: "^START", MaxLines: 5}
	c, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.max != 5 {
		t.Fatalf("expected max=5, got %d", c.max)
	}
}

func TestConfig_Build_ZeroMaxLines(t *testing.T) {
	cfg := Config{StartPattern: "^START", MaxLines: 0}
	c, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.max != 0 {
		t.Fatalf("expected max=0, got %d", c.max)
	}
}

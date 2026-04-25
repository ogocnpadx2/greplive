package batch

import (
	"testing"
	"time"
)

func TestDefaultConfig_MaxSize(t *testing.T) {
	c := DefaultConfig()
	if c.MaxSize != 100 {
		t.Fatalf("expected MaxSize 100, got %d", c.MaxSize)
	}
}

func TestDefaultConfig_Interval(t *testing.T) {
	c := DefaultConfig()
	if c.Interval != time.Second {
		t.Fatalf("expected 1s interval, got %v", c.Interval)
	}
}

func TestConfig_Build_Valid(t *testing.T) {
	c := Config{MaxSize: 50, Interval: 500 * time.Millisecond}
	b, err := c.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer b.Stop()
	if b.maxSize != 50 {
		t.Fatalf("expected maxSize 50, got %d", b.maxSize)
	}
}

func TestConfig_Build_NegativeMaxSize_ReturnsError(t *testing.T) {
	c := Config{MaxSize: -1, Interval: time.Second}
	_, err := c.Build()
	if err == nil {
		t.Fatal("expected error for negative MaxSize")
	}
}

func TestConfig_Build_ZeroMaxSize_Allowed(t *testing.T) {
	c := Config{MaxSize: 0, Interval: time.Second}
	b, err := c.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer b.Stop()
	if b.maxSize != 0 {
		t.Fatalf("expected maxSize 0, got %d", b.maxSize)
	}
}

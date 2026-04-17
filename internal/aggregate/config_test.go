package aggregate

import "testing"

func TestConfig_Build_ValidPattern(t *testing.T) {
	cfg := Config{Pattern: `ERROR`}
	a, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil aggregator")
	}
}

func TestConfig_Build_InvalidPattern(t *testing.T) {
	cfg := Config{Pattern: `[invalid`}
	_, err := cfg.Build()
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestConfig_Build_EmptyPattern_ReturnsNoop(t *testing.T) {
	cfg := Config{}
	a, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// nil pattern — Push should never aggregate.
	s, ok := a.Push("ERROR whatever")
	if ok || s != "" {
		t.Fatalf("noop aggregator should not aggregate: %q %v", s, ok)
	}
}

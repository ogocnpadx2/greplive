package scope

import "testing"

func TestDefaultConfig_Build_Disabled(t *testing.T) {
	sc, err := DefaultConfig().Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sc.Enabled() {
		t.Fatal("default config should produce disabled scope")
	}
}

func TestConfig_Build_Valid(t *testing.T) {
	cfg := Config{Start: "BEGIN", End: "END"}
	sc, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sc.Enabled() {
		t.Fatal("expected enabled scope")
	}
}

func TestConfig_Build_InvalidPattern(t *testing.T) {
	cfg := Config{Start: "[", End: ""}
	_, err := cfg.Build()
	if err == nil {
		t.Fatal("expected error for invalid start pattern")
	}
}

func TestConfig_Build_StartOnly(t *testing.T) {
	cfg := Config{Start: "MARK"}
	sc, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sc.Enabled() {
		t.Fatal("expected enabled scope with start only")
	}
	if sc.Allow("before") {
		t.Fatal("line before start should be blocked")
	}
	if !sc.Allow("MARK") {
		t.Fatal("start line should be allowed")
	}
}

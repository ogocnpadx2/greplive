package timefmt_test

import (
	"testing"
	"time"

	"greplive/internal/timefmt"
)

func TestDefaultConfig_Build(t *testing.T) {
	cfg := timefmt.DefaultConfig()
	f, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.FormatName() != "rfc3339" {
		t.Fatalf("expected rfc3339, got %q", f.FormatName())
	}
}

func TestConfig_Build_CustomFormat(t *testing.T) {
	cfg := timefmt.Config{Name: "short"}
	f, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.FormatName() != "short" {
		t.Fatalf("expected short, got %q", f.FormatName())
	}
}

func TestConfig_Build_UnknownFormat(t *testing.T) {
	cfg := timefmt.Config{Name: "nope"}
	_, err := cfg.Build()
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestConfig_Build_EmptyName_DefaultsToRFC3339(t *testing.T) {
	cfg := timefmt.Config{}
	f, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.FormatName() != "rfc3339" {
		t.Fatalf("expected rfc3339, got %q", f.FormatName())
	}
}

func TestConfig_Build_UTC(t *testing.T) {
	cfg := timefmt.Config{Name: "rfc3339", UTC: true}
	f, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// format a non-UTC time and ensure output ends with Z
	loc, _ := time.LoadLocation("America/New_York")
	t1 := time.Date(2024, 1, 1, 10, 0, 0, 0, loc)
	out := f.FormatTime(t1)
	if out[len(out)-1] != 'Z' {
		t.Fatalf("expected UTC suffix Z, got %q", out)
	}
}

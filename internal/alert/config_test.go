package alert

import (
	"testing"
	"time"
)

func TestConfig_Build_EmptyPattern_Disabled(t *testing.T) {
	cfg := Config{Pattern: "", Threshold: 5, Window: time.Second}
	a, err := cfg.Build(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Enabled() {
		t.Fatal("expected disabled alert for empty pattern")
	}
}

func TestConfig_Build_ZeroWindow_DefaultsToMinute(t *testing.T) {
	cfg := Config{Pattern: `ERR`, Threshold: 5, Window: 0}
	a, err := cfg.Build(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.window != time.Minute {
		t.Errorf("expected window to default to 1m, got %s", a.window)
	}
}

func TestConfig_Build_InvalidPattern(t *testing.T) {
	cfg := Config{Pattern: `[bad`, Threshold: 1, Window: time.Second}
	_, err := cfg.Build(nil)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

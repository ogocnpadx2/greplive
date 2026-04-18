package alert

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(`[invalid`, 1, time.Second, nil)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_EmptyPattern_Noop(t *testing.T) {
	a, err := New("", 5, time.Second, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Enabled() {
		t.Fatal("expected alert to be disabled for empty pattern")
	}
}

func TestEnabled_WithPattern(t *testing.T) {
	a, _ := New(`ERROR`, 3, time.Second, nil)
	if !a.Enabled() {
		t.Fatal("expected alert to be enabled")
	}
}

func TestCheck_TriggersAtThreshold(t *testing.T) {
	var buf bytes.Buffer
	a, _ := New(`ERROR`, 3, time.Minute, &buf)
	a.Check("INFO ok")
	a.Check("ERROR boom")
	a.Check("ERROR boom")
	if buf.Len() != 0 {
		t.Fatalf("alert fired too early: %s", buf.String())
	}
	a.Check("ERROR boom")
	if buf.Len() == 0 {
		t.Fatal("expected alert to fire at threshold")
	}
	if !strings.Contains(buf.String(), "[ALERT]") {
		t.Errorf("unexpected alert output: %s", buf.String())
	}
}

func TestCheck_NoAlert_BelowThreshold(t *testing.T) {
	var buf bytes.Buffer
	a, _ := New(`WARN`, 10, time.Minute, &buf)
	for i := 0; i < 9; i++ {
		a.Check("WARN something")
	}
	if buf.Len() != 0 {
		t.Fatalf("unexpected alert: %s", buf.String())
	}
}

func TestCheck_WindowReset(t *testing.T) {
	var buf bytes.Buffer
	a, _ := New(`ERROR`, 2, 10*time.Millisecond, &buf)
	a.Check("ERROR one")
	time.Sleep(20 * time.Millisecond)
	a.Check("ERROR two") // window resets here
	if buf.Len() != 0 {
		t.Fatalf("alert should not fire after window reset: %s", buf.String())
	}
}

func TestConfig_Build_Valid(t *testing.T) {
	cfg := Config{Pattern: `CRIT`, Threshold: 1, Window: time.Second}
	a, err := cfg.Build(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !a.Enabled() {
		t.Fatal("expected alert enabled")
	}
}

func TestConfig_Build_NegativeThreshold(t *testing.T) {
	cfg := Config{Pattern: `X`, Threshold: -1, Window: time.Second}
	_, err := cfg.Build(nil)
	if err == nil {
		t.Fatal("expected error for negative threshold")
	}
}

func TestDefaultConfig_Build(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Threshold != 10 {
		t.Errorf("expected threshold 10, got %d", cfg.Threshold)
	}
	if cfg.Window != time.Minute {
		t.Errorf("expected window 1m, got %s", cfg.Window)
	}
}

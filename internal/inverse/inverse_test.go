package inverse

import (
	"testing"
)

func TestNew_EmptyPattern_Disabled(t *testing.T) {
	inv, err := New("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inv.Enabled() {
		t.Fatal("expected inverter to be disabled for empty pattern")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := New("[invalid")
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_ValidPattern_Enabled(t *testing.T) {
	inv, err := New(`ERROR`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !inv.Enabled() {
		t.Fatal("expected inverter to be enabled")
	}
}

func TestAllow_Disabled_AlwaysTrue(t *testing.T) {
	inv, _ := New("")
	for _, line := range []string{"", "ERROR", "some line"} {
		if !inv.Allow(line) {
			t.Errorf("disabled inverter should allow %q", line)
		}
	}
}

func TestAllow_MatchingLine_ReturnsFalse(t *testing.T) {
	inv, _ := New(`ERROR`)
	if inv.Allow("2024/01/01 ERROR something broke") {
		t.Fatal("expected matching line to be dropped")
	}
}

func TestAllow_NonMatchingLine_ReturnsTrue(t *testing.T) {
	inv, _ := New(`ERROR`)
	if !inv.Allow("2024/01/01 INFO all good") {
		t.Fatal("expected non-matching line to be allowed")
	}
}

func TestApplyAll_FiltersMatchingLines(t *testing.T) {
	inv1, _ := New(`ERROR`)
	inv2, _ := New(`SKIP`)
	input := []string{
		"INFO hello",
		"ERROR bad thing",
		"SKIP this line",
		"DEBUG ok",
	}
	got := ApplyAll([]*Inverter{inv1, inv2}, input)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(got), got)
	}
	if got[0] != "INFO hello" || got[1] != "DEBUG ok" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestApplyAll_NoInverters_PassesAll(t *testing.T) {
	input := []string{"a", "b", "c"}
	got := ApplyAll(nil, input)
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
}

func TestApplyAll_DisabledInverter_PassesAll(t *testing.T) {
	inv, _ := New("")
	input := []string{"ERROR line", "INFO line"}
	got := ApplyAll([]*Inverter{inv}, input)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}

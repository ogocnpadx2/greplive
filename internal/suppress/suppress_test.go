package suppress

import (
	"testing"
)

func TestNew_EmptyPattern_Disabled(t *testing.T) {
	s, err := New("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Enabled() {
		t.Fatal("expected suppressor to be disabled for empty pattern")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := New("[invalid")
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_ValidPattern_Enabled(t *testing.T) {
	s, err := New(`ERROR`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.Enabled() {
		t.Fatal("expected suppressor to be enabled")
	}
}

func TestDrop_Disabled_AlwaysFalse(t *testing.T) {
	s, _ := New("")
	if s.Drop("ERROR: something went wrong") {
		t.Fatal("disabled suppressor should never drop")
	}
}

func TestDrop_MatchingLine_ReturnsTrue(t *testing.T) {
	s, _ := New(`(?i)debug`)
	if !s.Drop("DEBUG: verbose output") {
		t.Fatal("expected matching line to be dropped")
	}
}

func TestDrop_NonMatchingLine_ReturnsFalse(t *testing.T) {
	s, _ := New(`(?i)debug`)
	if s.Drop("INFO: service started") {
		t.Fatal("expected non-matching line to pass through")
	}
}

func TestApplyAll_NoSuppressors_Unchanged(t *testing.T) {
	lines := []string{"a", "b", "c"}
	got := ApplyAll(nil, lines)
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
}

func TestApplyAll_DropsMatchingLines(t *testing.T) {
	s1, _ := New(`DEBUG`)
	s2, _ := New(`TRACE`)
	lines := []string{
		"INFO: ok",
		"DEBUG: noisy",
		"TRACE: very noisy",
		"WARN: watch out",
	}
	got := ApplyAll([]*Suppressor{s1, s2}, lines)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(got), got)
	}
	if got[0] != "INFO: ok" || got[1] != "WARN: watch out" {
		t.Fatalf("unexpected output: %v", got)
	}
}

func TestApplyAll_DisabledSuppressor_PassesAll(t *testing.T) {
	s, _ := New("")
	lines := []string{"DEBUG: a", "INFO: b"}
	got := ApplyAll([]*Suppressor{s}, lines)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}

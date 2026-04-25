package after

import (
	"testing"
)

func TestNew_EmptyPattern_Disabled(t *testing.T) {
	a, err := New("", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Enabled() {
		t.Error("expected disabled for empty pattern")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := New("[invalid", 0)
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestNew_ValidPattern_Enabled(t *testing.T) {
	a, err := New("START", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !a.Enabled() {
		t.Error("expected enabled for valid pattern")
	}
}

func TestAllow_Disabled_AlwaysTrue(t *testing.T) {
	a, _ := New("", 0)
	for _, line := range []string{"foo", "bar", "START"} {
		if !a.Allow(line) {
			t.Errorf("disabled: expected Allow(%q) = true", line)
		}
	}
}

func TestAllow_TriggerLineNotEmitted(t *testing.T) {
	a, _ := New("START", 0)
	if a.Allow("this is the START line") {
		t.Error("trigger line should not be emitted")
	}
}

func TestAllow_LinesBeforeTrigger_Dropped(t *testing.T) {
	a, _ := New("START", 0)
	if a.Allow("before trigger") {
		t.Error("lines before trigger should be dropped")
	}
}

func TestAllow_LinesAfterTrigger_Emitted(t *testing.T) {
	a, _ := New("START", 0)
	a.Allow("START")
	for i, line := range []string{"line1", "line2", "line3"} {
		if !a.Allow(line) {
			t.Errorf("line %d after trigger should be emitted: %q", i, line)
		}
	}
}

func TestAllow_MaxLines_StopsAfterLimit(t *testing.T) {
	a, _ := New("START", 2)
	a.Allow("START")
	if !a.Allow("line1") {
		t.Error("first line after trigger should be emitted")
	}
	if !a.Allow("line2") {
		t.Error("second line after trigger should be emitted")
	}
	if a.Allow("line3") {
		t.Error("third line should be dropped (exceeds max)")
	}
}

func TestAllow_RetriggerResetsCount(t *testing.T) {
	a, _ := New("START", 1)
	a.Allow("START")
	a.Allow("line1") // consumes the 1 allowed line
	a.Allow("START") // re-trigger
	if !a.Allow("line2") {
		t.Error("line after re-trigger should be emitted")
	}
}

func TestReset_StopsEmission(t *testing.T) {
	a, _ := New("START", 0)
	a.Allow("START")
	a.Reset()
	if a.Allow("after reset") {
		t.Error("line after Reset should be dropped")
	}
}

package colorize

import (
	"strings"
	"testing"
)

func TestNew_EmptyPattern_Disabled(t *testing.T) {
	c, err := New("", "red")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Enabled() {
		t.Fatal("expected disabled colorizer for empty pattern")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := New("[invalid", "red")
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_ValidPattern_Enabled(t *testing.T) {
	c, err := New(`\d+`, "blue")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !c.Enabled() {
		t.Fatal("expected enabled colorizer")
	}
}

func TestApply_Disabled_Unchanged(t *testing.T) {
	c, _ := New("", "red")
	line := "no change expected"
	if got := c.Apply(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestApply_MatchesColored(t *testing.T) {
	c, err := New(`ERROR`, "red")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := c.Apply("ERROR: something went wrong")
	if !strings.Contains(result, namedColors["red"]) {
		t.Error("expected red ANSI code in output")
	}
	if !strings.Contains(result, reset) {
		t.Error("expected reset ANSI code in output")
	}
	if !strings.Contains(result, "ERROR") {
		t.Error("expected original match text preserved")
	}
}

func TestApply_NoMatch_Unchanged(t *testing.T) {
	c, _ := New(`ERROR`, "red")
	line := "everything is fine"
	if got := c.Apply(line); got != line {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestApply_UnknownColor_FallsBackToWhite(t *testing.T) {
	c, err := New(`WARN`, "ultraviolet")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := c.Apply("WARN: check this")
	if !strings.Contains(result, namedColors["white"]) {
		t.Error("expected white fallback ANSI code in output")
	}
}

func TestApplyAll_MultipleColorizers(t *testing.T) {
	c1, _ := New(`ERROR`, "red")
	c2, _ := New(`\d+`, "yellow")
	line := "ERROR code 42"
	result := ApplyAll(line, []*Colorizer{c1, c2})
	if !strings.Contains(result, namedColors["red"]) {
		t.Error("expected red code for ERROR")
	}
	if !strings.Contains(result, namedColors["yellow"]) {
		t.Error("expected yellow code for digits")
	}
}

func TestColorByName_KnownName(t *testing.T) {
	if got := ColorByName("cyan"); got != namedColors["cyan"] {
		t.Fatalf("expected cyan code, got %q", got)
	}
}

func TestColorByName_UnknownName_Empty(t *testing.T) {
	if got := ColorByName("rainbow"); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

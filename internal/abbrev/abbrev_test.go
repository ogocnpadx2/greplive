package abbrev

import (
	"strings"
	"testing"
)

func TestNew_Disabled_WhenNoRulesActive(t *testing.T) {
	a := New(0, "", false)
	if a.Enabled() {
		t.Fatal("expected disabled")
	}
}

func TestNew_Enabled_WhenMaxRunesSet(t *testing.T) {
	a := New(10, "", false)
	if !a.Enabled() {
		t.Fatal("expected enabled")
	}
}

func TestNew_Enabled_WhenCollapseSpace(t *testing.T) {
	a := New(0, "", true)
	if !a.Enabled() {
		t.Fatal("expected enabled")
	}
}

func TestNew_DefaultSuffix(t *testing.T) {
	a := New(5, "", false)
	if a.suffix != "…" {
		t.Fatalf("expected default suffix '…', got %q", a.suffix)
	}
}

func TestApply_Disabled_Unchanged(t *testing.T) {
	a := New(0, "", false)
	line := "hello   world"
	if got := a.Apply(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestApply_ShortLine_Unchanged(t *testing.T) {
	a := New(20, "…", false)
	line := "short"
	if got := a.Apply(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestApply_LongLine_Truncated(t *testing.T) {
	a := New(10, "…", false)
	line := "hello world this is a long line"
	got := a.Apply(line)
	if len([]rune(got)) > 10 {
		t.Fatalf("expected at most 10 runes, got %d: %q", len([]rune(got)), got)
	}
	if !strings.HasSuffix(got, "…") {
		t.Fatalf("expected suffix '…', got %q", got)
	}
}

func TestApply_CollapseSpace(t *testing.T) {
	a := New(0, "", true)
	line := "  foo   bar   baz  "
	got := a.Apply(line)
	if got != "foo bar baz" {
		t.Fatalf("expected 'foo bar baz', got %q", got)
	}
}

func TestApply_CollapseAndTruncate(t *testing.T) {
	a := New(7, "…", true)
	line := "  hello   world  "
	got := a.Apply(line) // collapse -> "hello world", truncate to 7 runes
	if len([]rune(got)) > 7 {
		t.Fatalf("expected at most 7 runes, got %d: %q", len([]rune(got)), got)
	}
}

func TestApplyAll_MultipleAbbreviators(t *testing.T) {
	abbrevs := []*Abbreviator{
		New(0, "", true),
		New(8, "…", false),
	}
	line := "  one   two   three  "
	got := ApplyAll(line, abbrevs)
	if len([]rune(got)) > 8 {
		t.Fatalf("expected at most 8 runes, got %d: %q", len([]rune(got)), got)
	}
}

func TestApply_CustomSuffix(t *testing.T) {
	a := New(6, "[…]", false)
	line := "hello world"
	got := a.Apply(line)
	if !strings.HasSuffix(got, "[…]") {
		t.Fatalf("expected custom suffix, got %q", got)
	}
}

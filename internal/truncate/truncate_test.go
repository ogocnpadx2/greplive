package truncate_test

import (
	"strings"
	"testing"

	"greplive/internal/truncate"
)

func TestNew_DefaultSuffix(t *testing.T) {
	tr := truncate.New(80, "")
	if tr.Suffix() != truncate.DefaultSuffix {
		t.Errorf("expected default suffix %q, got %q", truncate.DefaultSuffix, tr.Suffix())
	}
}

func TestNew_CustomSuffix(t *testing.T) {
	tr := truncate.New(40, "[cut]")
	if tr.Suffix() != "[cut]" {
		t.Errorf("unexpected suffix: %q", tr.Suffix())
	}
}

func TestEnabled(t *testing.T) {
	if truncate.New(0, "").Enabled() {
		t.Error("expected Enabled()=false for maxRunes=0")
	}
	if !truncate.New(10, "").Enabled() {
		t.Error("expected Enabled()=true for maxRunes=10")
	}
}

func TestApply_ShortLine_Unchanged(t *testing.T) {
	tr := truncate.New(20, "…")
	line := "hello world"
	if got := tr.Apply(line); got != line {
		t.Errorf("expected %q unchanged, got %q", line, got)
	}
}

func TestApply_LongLine_Truncated(t *testing.T) {
	tr := truncate.New(10, "[…]")
	line := strings.Repeat("a", 30)
	got := tr.Apply(line)
	if !strings.HasSuffix(got, "[…]") {
		t.Errorf("expected suffix '[…]' in %q", got)
	}
	// Content before suffix should be exactly 10 runes.
	content := strings.TrimSuffix(got, "[…]")
	if len([]rune(content)) != 10 {
		t.Errorf("expected 10 runes before suffix, got %d", len([]rune(content)))
	}
}

func TestApply_ZeroMaxRunes_NoTruncation(t *testing.T) {
	tr := truncate.New(0, "")
	line := strings.Repeat("x", 1000)
	if got := tr.Apply(line); got != line {
		t.Error("expected line unchanged when maxRunes=0")
	}
}

func TestApply_ExactlyAtLimit_Unchanged(t *testing.T) {
	tr := truncate.New(5, "[…]")
	line := "abcde" // exactly 5 runes
	if got := tr.Apply(line); got != line {
		t.Errorf("expected %q unchanged at exact limit, got %q", line, got)
	}
}

func TestApply_Unicode(t *testing.T) {
	tr := truncate.New(4, "…")
	line := "日本語テスト" // 6 runes
	got := tr.Apply(line)
	if !strings.HasSuffix(got, "…") {
		t.Errorf("expected truncation suffix, got %q", got)
	}
	content := strings.TrimSuffix(got, "…")
	if len([]rune(content)) != 4 {
		t.Errorf("expected 4 runes, got %d", len([]rune(content)))
	}
}

func TestApplyAll(t *testing.T) {
	tr := truncate.New(5, "…")
	lines := []string{"short", strings.Repeat("z", 20), "hi"}
	result := tr.ApplyAll(lines)
	if result[0] != "short" {
		t.Errorf("first line should be unchanged, got %q", result[0])
	}
	if !strings.HasSuffix(result[1], "…") {
		t.Errorf("second line should be truncated, got %q", result[1])
	}
	if result[2] != "hi" {
		t.Errorf("third line should be unchanged, got %q", result[2])
	}
}

func TestStripSuffix(t *testing.T) {
	clean, found := truncate.StripSuffix("some log line …")
	if !found {
		t.Error("expected suffix to be found")
	}
	if strings.HasSuffix(clean, truncate.DefaultSuffix) {
		t.Error("suffix should have been stripped")
	}

	_, found = truncate.StripSuffix("no suffix here")
	if found {
		t.Error("expected no suffix found")
	}
}

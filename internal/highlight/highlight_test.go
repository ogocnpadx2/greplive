package highlight_test

import (
	"testing"

	"greplive/internal/highlight"
)

func TestNew_ValidPattern(t *testing.T) {
	h, err := highlight.New(`error`, highlight.Yellow)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil Highlighter")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := highlight.New(`[invalid`, highlight.Yellow)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestApply_MatchesAndColors(t *testing.T) {
	h, _ := highlight.New(`error`, highlight.Red)
	line := "this is an error message"
	result := h.Apply(line)
	stripped := highlight.StripANSI(result)
	if stripped != line {
		t.Errorf("stripped result mismatch: got %q, want %q", stripped, line)
	}
	if result == line {
		t.Error("expected ANSI codes to be added, but line was unchanged")
	}
}

func TestApply_NoMatch(t *testing.T) {
	h, _ := highlight.New(`fatal`, highlight.Red)
	line := "everything is fine"
	result := h.Apply(line)
	if result != line {
		t.Errorf("expected unchanged line, got %q", result)
	}
}

func TestApplyAll_MultipleHighlighters(t *testing.T) {
	h1, _ := highlight.New(`error`, highlight.Red)
	h2, _ := highlight.New(`warn`, highlight.Yellow)
	line := "warn: an error occurred"
	result := highlight.ApplyAll(line, []*highlight.Highlighter{h1, h2})
	stripped := highlight.StripANSI(result)
	if stripped != line {
		t.Errorf("stripped mismatch: got %q, want %q", stripped, line)
	}
}

func TestStripANSI(t *testing.T) {
	input := "\033[33m\033[1mwarn\033[0m: message"
	expected := "warn: message"
	if got := highlight.StripANSI(input); got != expected {
		t.Errorf("StripANSI got %q, want %q", got, expected)
	}
}

func TestColorByName(t *testing.T) {
	cases := []struct {
		name     string
		expected string
	}{
		{"yellow", highlight.Yellow},
		{"cyan", highlight.Cyan},
		{"red", highlight.Red},
		{"unknown", highlight.Yellow},
	}
	for _, tc := range cases {
		if got := highlight.ColorByName(tc.name); got != tc.expected {
			t.Errorf("ColorByName(%q) = %q, want %q", tc.name, got, tc.expected)
		}
	}
}

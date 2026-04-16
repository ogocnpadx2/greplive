package redact

import (
	"testing"
)

func TestNew_ValidPattern(t *testing.T) {
	r, err := New(`\d{4}-\d{4}-\d{4}-\d{4}`, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil Redactor")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(`[invalid`, "")
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New("", "")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_DefaultReplacement(t *testing.T) {
	r, _ := New(`secret`, "")
	out := r.Apply("my secret value")
	if out != "my [REDACTED] value" {
		t.Errorf("got %q", out)
	}
}

func TestApply_ReplacesMatches(t *testing.T) {
	r, _ := New(`\d+`, "<NUM>")
	out := r.Apply("error on line 42 at col 7")
	expected := "error on line <NUM> at col <NUM>"
	if out != expected {
		t.Errorf("got %q, want %q", out, expected)
	}
}

func TestApply_NoMatch(t *testing.T) {
	r, _ := New(`secret`, "[REDACTED]")
	line := "nothing sensitive here"
	out := r.Apply(line)
	if out != line {
		t.Errorf("expected unchanged line, got %q", out)
	}
}

func TestApplyAll_MultipleRedactors(t *testing.T) {
	r1, _ := New(`password=\S+`, "password=[REDACTED]")
	r2, _ := New(`token=\S+`, "token=[REDACTED]")
	line := "auth password=hunter2 token=abc123"
	out := ApplyAll([]*Redactor{r1, r2}, line)
	expected := "auth password=[REDACTED] token=[REDACTED]"
	if out != expected {
		t.Errorf("got %q, want %q", out, expected)
	}
}

func TestApplyAll_EmptySlice(t *testing.T) {
	line := "unchanged line"
	out := ApplyAll(nil, line)
	if out != line {
		t.Errorf("expected unchanged, got %q", out)
	}
}

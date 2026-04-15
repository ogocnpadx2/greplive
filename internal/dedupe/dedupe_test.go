package dedupe

import (
	"strings"
	"testing"
)

func TestCheck_FirstLineAlwaysAllowed(t *testing.T) {
	d := New(false)
	flush, allow := d.Check("hello")
	if !allow {
		t.Fatal("expected first line to be allowed")
	}
	if flush != "" {
		t.Fatalf("unexpected flush: %q", flush)
	}
}

func TestCheck_DuplicateSuppressed(t *testing.T) {
	d := New(false)
	d.Check("hello")
	_, allow := d.Check("hello")
	if allow {
		t.Fatal("expected duplicate to be suppressed")
	}
}

func TestCheck_NewLineAfterDuplicate(t *testing.T) {
	d := New(false)
	d.Check("hello")
	d.Check("hello")
	_, allow := d.Check("world")
	if !allow {
		t.Fatal("expected new line to be allowed")
	}
}

func TestCheck_SummaryEmittedOnChange(t *testing.T) {
	d := New(true)
	d.Check("line")
	d.Check("line") // dup 1
	d.Check("line") // dup 2
	flush, allow := d.Check("other")
	if !allow {
		t.Fatal("expected new line to be allowed")
	}
	if !strings.Contains(flush, "3") {
		t.Fatalf("expected summary to mention count 3, got %q", flush)
	}
}

func TestCheck_NoSummaryWhenDisabled(t *testing.T) {
	d := New(false)
	d.Check("line")
	d.Check("line")
	flush, _ := d.Check("other")
	if flush != "" {
		t.Fatalf("expected no summary, got %q", flush)
	}
}

func TestCheck_NoSummaryForSingleOccurrence(t *testing.T) {
	d := New(true)
	d.Check("a")
	flush, _ := d.Check("b")
	if flush != "" {
		t.Fatalf("expected no summary for single occurrence, got %q", flush)
	}
}

func TestFlush_ReturnsSummaryAndResets(t *testing.T) {
	d := New(true)
	d.Check("x")
	d.Check("x")
	d.Check("x")
	s := d.Flush()
	if !strings.Contains(s, "3") {
		t.Fatalf("expected flush summary with count 3, got %q", s)
	}
	// second flush should be empty
	if s2 := d.Flush(); s2 != "" {
		t.Fatalf("expected empty second flush, got %q", s2)
	}
}

func TestReset_ClearsState(t *testing.T) {
	d := New(true)
	d.Check("line")
	d.Check("line")
	d.Reset()
	// After reset the same line should be treated as new.
	_, allow := d.Check("line")
	if !allow {
		t.Fatal("expected line to be allowed after reset")
	}
}

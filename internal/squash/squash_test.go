package squash

import (
	"testing"
)

func TestNew_EmptyPattern_Disabled(t *testing.T) {
	s, err := New("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Enabled() {
		t.Fatal("expected disabled")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New("[invalid")
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_ValidPattern_Enabled(t *testing.T) {
	s, err := New(`^DEBUG`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.Enabled() {
		t.Fatal("expected enabled")
	}
}

func TestPush_Disabled_PassThrough(t *testing.T) {
	s, _ := New("")
	line, ok := s.Push("hello")
	if !ok || line != "hello" {
		t.Fatalf("expected passthrough, got %q %v", line, ok)
	}
}

func TestPush_MatchingLines_Buffered(t *testing.T) {
	s, _ := New(`^DEBUG`)
	_, ok := s.Push("DEBUG a")
	if ok {
		t.Fatal("expected buffered (no emit)")
	}
	_, ok = s.Push("DEBUG b")
	if ok {
		t.Fatal("expected buffered (no emit)")
	}
}

func TestPush_NonMatchFlushesBuffer(t *testing.T) {
	s, _ := New(`^DEBUG`)
	s.Push("DEBUG a")
	s.Push("DEBUG b")
	merged, ok := s.Push("INFO done")
	if !ok {
		t.Fatal("expected emit on non-match")
	}
	if merged != "DEBUG a | DEBUG b" {
		t.Fatalf("unexpected merged: %q", merged)
	}
}

func TestFlush_ReturnsPending(t *testing.T) {
	s, _ := New(`^DEBUG`)
	s.Push("DEBUG x")
	out, ok := s.Flush()
	if !ok {
		t.Fatal("expected flush to emit")
	}
	if out != "DEBUG x" {
		t.Fatalf("unexpected: %q", out)
	}
}

func TestFlush_EmptyBuffer_ReturnsFalse(t *testing.T) {
	s, _ := New(`^DEBUG`)
	_, ok := s.Flush()
	if ok {
		t.Fatal("expected false on empty flush")
	}
}

func TestPush_SingleMatch_MergesAsIs(t *testing.T) {
	s, _ := New(`^WARN`)
	s.Push("WARN only")
	out, ok := s.Flush()
	if !ok || out != "WARN only" {
		t.Fatalf("unexpected: %q %v", out, ok)
	}
}

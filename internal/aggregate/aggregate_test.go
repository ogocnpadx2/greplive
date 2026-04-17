package aggregate

import (
	"regexp"
	"testing"
)

func re(pat string) *regexp.Regexp { return regexp.MustCompile(pat) }

func TestPush_NonMatchingLine_NoSummary(t *testing.T) {
	a := New(re(`ERROR`))
	s, ok := a.Push("INFO hello")
	if ok || s != "" {
		t.Fatalf("expected no summary, got %q %v", s, ok)
	}
}

func TestPush_MatchingLines_Buffered(t *testing.T) {
	a := New(re(`ERROR`))
	for i := 0; i < 3; i++ {
		s, ok := a.Push("ERROR boom")
		if ok || s != "" {
			t.Fatalf("unexpected summary on match: %q %v", s, ok)
		}
	}
}

func TestPush_RunEnds_EmitsSummary(t *testing.T) {
	a := New(re(`ERROR`))
	a.Push("ERROR boom")
	a.Push("ERROR boom")
	s, ok := a.Push("INFO ok")
	if !ok {
		t.Fatal("expected summary")
	}
	if s != "ERROR boom [x2]" {
		t.Fatalf("unexpected summary: %q", s)
	}
}

func TestFlush_ReturnsPending(t *testing.T) {
	a := New(re(`WARN`))
	a.Push("WARN low disk")
	a.Push("WARN low disk")
	a.Push("WARN low disk")
	s, ok := a.Flush()
	if !ok {
		t.Fatal("expected flush summary")
	}
	if s != "WARN low disk [x3]" {
		t.Fatalf("got %q", s)
	}
}

func TestFlush_EmptyBuffer_ReturnsFalse(t *testing.T) {
	a := New(re(`WARN`))
	_, ok := a.Flush()
	if ok {
		t.Fatal("expected false on empty flush")
	}
}

func TestPush_NilPattern_NeverMatches(t *testing.T) {
	a := New(nil)
	s, ok := a.Push("ERROR anything")
	if ok || s != "" {
		t.Fatalf("nil pattern should never match: %q %v", s, ok)
	}
}

package strip

import (
	"testing"
)

func TestNew_Disabled_WhenBothFalse(t *testing.T) {
	s := New(false, false)
	if s.Enabled() {
		t.Fatal("expected disabled when both flags false")
	}
}

func TestNew_Enabled_WhenANSI(t *testing.T) {
	s := New(true, false)
	if !s.Enabled() {
		t.Fatal("expected enabled when ansi=true")
	}
}

func TestNew_Enabled_WhenControl(t *testing.T) {
	s := New(false, true)
	if !s.Enabled() {
		t.Fatal("expected enabled when control=true")
	}
}

func TestApply_Disabled_Unchanged(t *testing.T) {
	s := New(false, false)
	input := "\x1b[31mhello\x1b[0m"
	if got := s.Apply(input); got != input {
		t.Fatalf("expected %q, got %q", input, got)
	}
}

func TestApply_StripANSI_Colour(t *testing.T) {
	s := New(true, false)
	got := s.Apply("\x1b[31mERROR\x1b[0m: something failed")
	want := "ERROR: something failed"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestApply_StripANSI_OSC(t *testing.T) {
	s := New(true, false)
	got := s.Apply("\x1b]0;title\x07line")
	want := "line"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestApply_StripControl(t *testing.T) {
	s := New(false, true)
	got := s.Apply("hel\x01lo\x7f")
	want := "hello"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestApply_StripBoth(t *testing.T) {
	s := New(true, true)
	got := s.Apply("\x1b[32m\x01ok\x1b[0m")
	want := "ok"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestApply_PlainLine_Unchanged(t *testing.T) {
	s := New(true, true)
	want := "plain log line"
	if got := s.Apply(want); got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestApplyAll_MultipleStrippers(t *testing.T) {
	strippers := []*Stripper{
		New(true, false),
		New(false, true),
	}
	got := ApplyAll("\x1b[31m\x01data\x1b[0m", strippers)
	want := "data"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

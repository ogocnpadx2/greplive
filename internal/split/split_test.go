package split

import (
	"testing"
)

func TestNew_EmptyDelimiter_Disabled(t *testing.T) {
	sp, err := New("", nil, " ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sp.Enabled() {
		t.Fatal("expected disabled")
	}
}

func TestNew_NonEmpty_Enabled(t *testing.T) {
	sp, err := New("|", nil, " ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sp.Enabled() {
		t.Fatal("expected enabled")
	}
}

func TestApply_Disabled_Unchanged(t *testing.T) {
	sp, _ := New("", nil, " ")
	got := sp.Apply("hello|world")
	if got != "hello|world" {
		t.Fatalf("expected unchanged, got %q", got)
	}
}

func TestApply_AllFields(t *testing.T) {
	sp, _ := New("|", nil, "-")
	got := sp.Apply("a|b|c")
	if got != "a-b-c" {
		t.Fatalf("expected %q got %q", "a-b-c", got)
	}
}

func TestApply_SelectFields(t *testing.T) {
	sp, _ := New(",", []int{0, 2}, " ")
	got := sp.Apply("foo,bar,baz")
	if got != "foo baz" {
		t.Fatalf("expected %q got %q", "foo baz", got)
	}
}

func TestApply_OutOfRangeIndex_Skipped(t *testing.T) {
	sp, _ := New(",", []int{0, 99}, " ")
	got := sp.Apply("x,y")
	if got != "x" {
		t.Fatalf("expected %q got %q", "x", got)
	}
}

func TestApply_NoDelimiterInLine_Unchanged(t *testing.T) {
	sp, _ := New("|", []int{1}, " ")
	got := sp.Apply("nodeli miter")
	// only one field produced, index 1 out of range → fallback
	if got != "nodeli miter" {
		t.Fatalf("expected unchanged, got %q", got)
	}
}

func TestApply_DefaultJoin(t *testing.T) {
	sp, _ := New(":", nil, "")
	got := sp.Apply("a:b:c")
	if got != "a b c" {
		t.Fatalf("expected %q got %q", "a b c", got)
	}
}

func TestApplyAll_MultipleSlitters(t *testing.T) {
	sp1, _ := New("|", []int{0, 1}, "-")
	sp2, _ := New("-", []int{1}, " ")
	got := ApplyAll([]*Splitter{sp1, sp2}, "hello|world|extra")
	if got != "world" {
		t.Fatalf("expected %q got %q", "world", got)
	}
}

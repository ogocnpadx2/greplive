package column

import (
	"testing"
)

func TestNew_EmptyColumns_Disabled(t *testing.T) {
	ex := New(" ", nil)
	if ex.Enabled() {
		t.Fatal("expected disabled when no columns provided")
	}
}

func TestNew_WithColumns_Enabled(t *testing.T) {
	ex := New(" ", []int{0, 2})
	if !ex.Enabled() {
		t.Fatal("expected enabled")
	}
}

func TestApply_Disabled_Unchanged(t *testing.T) {
	ex := New(" ", nil)
	const line = "hello world foo"
	if got := ex.Apply(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestApply_SelectsColumns(t *testing.T) {
	ex := New(" ", []int{0, 2})
	got := ex.Apply("hello world foo")
	const want = "hello foo"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestApply_SingleColumn(t *testing.T) {
	ex := New(":", []int{1})
	got := ex.Apply("level:info:msg")
	if got != "info" {
		t.Fatalf("expected %q, got %q", "info", got)
	}
}

func TestApply_OutOfRange_ReturnsOriginal(t *testing.T) {
	ex := New(" ", []int{0, 10})
	const line = "only three words here"
	// index 10 does not exist — original returned
	if got := ex.Apply("a b"); got != "a b" {
		t.Fatalf("expected original line, got %q", got)
	}
	_ = line
}

func TestApply_DefaultDelimiter(t *testing.T) {
	ex := New("", []int{1})
	if got := ex.Apply("foo bar baz"); got != "bar" {
		t.Fatalf("expected %q, got %q", "bar", got)
	}
}

func TestApplyAll_MultipleExtractors(t *testing.T) {
	// first extractor splits on space and keeps col 1 onward via a single-col
	e1 := New(" ", []int{0})
	e2 := New("-", []int{0})
	got := ApplyAll([]*Extractor{e1, e2}, "hello-world there")
	if got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

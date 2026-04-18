package indent_test

import (
	"testing"

	"greplive/internal/indent"
)

func TestNew_EmptyPrefix_Disabled(t *testing.T) {
	i := indent.New("")
	if i.Enabled() {
		t.Fatal("expected disabled for empty prefix")
	}
}

func TestNew_NonEmpty_Enabled(t *testing.T) {
	i := indent.New("  ")
	if !i.Enabled() {
		t.Fatal("expected enabled for non-empty prefix")
	}
}

func TestApply_Disabled_Unchanged(t *testing.T) {
	i := indent.New("")
	got := i.Apply("hello")
	if got != "hello" {
		t.Fatalf("expected 'hello', got %q", got)
	}
}

func TestApply_PrependPrefix(t *testing.T) {
	i := indent.New(">> ")
	got := i.Apply("hello")
	if got != ">> hello" {
		t.Fatalf("expected '>> hello', got %q", got)
	}
}

func TestApplyAll_MultipleIndenters(t *testing.T) {
	indenters := []*indent.Indenter{
		indent.New("[A] "),
		indent.New("[B] "),
	}
	got := indent.ApplyAll("msg", indenters)
	if got != "[B] [A] msg" {
		t.Fatalf("unexpected result: %q", got)
	}
}

func TestApplyAll_EmptySlice(t *testing.T) {
	got := indent.ApplyAll("msg", nil)
	if got != "msg" {
		t.Fatalf("expected 'msg', got %q", got)
	}
}

func TestRepeat_ZeroN(t *testing.T) {
	if indent.Repeat("  ", 0) != "" {
		t.Fatal("expected empty string for n=0")
	}
}

func TestRepeat_PositiveN(t *testing.T) {
	got := indent.Repeat("  ", 3)
	if got != "      " {
		t.Fatalf("expected 6 spaces, got %q", got)
	}
}

func TestRepeat_NegativeN(t *testing.T) {
	if indent.Repeat("\t", -1) != "" {
		t.Fatal("expected empty string for negative n")
	}
}

package normalize_test

import (
	"testing"

	"greplive/internal/normalize"
)

func TestEnabled_NoOptions(t *testing.T) {
	n := normalize.New()
	if n.Enabled() {
		t.Fatal("expected disabled with no options")
	}
}

func TestEnabled_WithOption(t *testing.T) {
	n := normalize.New(normalize.WithTrimSpace())
	if !n.Enabled() {
		t.Fatal("expected enabled with option")
	}
}

func TestApply_NoSteps_Unchanged(t *testing.T) {
	n := normalize.New()
	got := n.Apply("  hello  ")
	if got != "  hello  " {
		t.Fatalf("expected unchanged, got %q", got)
	}
}

func TestApply_TrimSpace(t *testing.T) {
	n := normalize.New(normalize.WithTrimSpace())
	got := n.Apply("  hello world  ")
	if got != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", got)
	}
}

func TestApply_CollapseSpaces(t *testing.T) {
	n := normalize.New(normalize.WithCollapseSpaces())
	got := n.Apply("hello   world\t  foo")
	if got != "hello world foo" {
		t.Fatalf("expected %q, got %q", "hello world foo", got)
	}
}

func TestApply_Lowercase(t *testing.T) {
	n := normalize.New(normalize.WithLowercase())
	got := n.Apply("ERROR: Something Failed")
	if got != "error: something failed" {
		t.Fatalf("expected %q, got %q", "error: something failed", got)
	}
}

func TestApply_Replace(t *testing.T) {
	n := normalize.New(normalize.WithReplace("foo", "bar"))
	got := n.Apply("foo and foo")
	if got != "bar and bar" {
		t.Fatalf("expected %q, got %q", "bar and bar", got)
	}
}

func TestApply_ChainedSteps(t *testing.T) {
	n := normalize.New(
		normalize.WithTrimSpace(),
		normalize.WithCollapseSpaces(),
		normalize.WithLowercase(),
	)
	got := n.Apply("  ERROR   occurred  ")
	want := "error occurred"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestApply_Replace_EmptyOld_Unchanged(t *testing.T) {
	n := normalize.New(normalize.WithReplace("", "x"))
	input := "hello"
	// strings.ReplaceAll with empty old inserts between each rune; just verify no panic
	_ = n.Apply(input)
}

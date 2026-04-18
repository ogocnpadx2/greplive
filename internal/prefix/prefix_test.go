package prefix_test

import (
	"testing"

	"greplive/internal/prefix"
)

func TestNew_EmptyPrefix_Disabled(t *testing.T) {
	p := prefix.New("")
	if p.Enabled() {
		t.Fatal("expected disabled for empty prefix")
	}
}

func TestNew_NonEmpty_Enabled(t *testing.T) {
	p := prefix.New("[app] ")
	if !p.Enabled() {
		t.Fatal("expected enabled for non-empty prefix")
	}
}

func TestApply_PrependsPrefixToLine(t *testing.T) {
	p := prefix.New("[app] ")
	got := p.Apply("hello world")
	want := "[app] hello world"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestApply_EmptyPrefix_Unchanged(t *testing.T) {
	p := prefix.New("")
	line := "no change"
	if got := p.Apply(line); got != line {
		t.Fatalf("got %q, want %q", got, line)
	}
}

func TestApplyAll_MultiplePrefix(t *testing.T) {
	ps := []*prefix.Prefixer{
		prefix.New("A:"),
		prefix.New("B:"),
	}
	got := prefix.ApplyAll("msg", ps)
	want := "B:A:msg"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestApplyAll_EmptySlice(t *testing.T) {
	line := "unchanged"
	got := prefix.ApplyAll(line, nil)
	if got != line {
		t.Fatalf("got %q, want %q", got, line)
	}
}

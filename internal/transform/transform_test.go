package transform_test

import (
	"testing"

	"greplive/internal/transform"
)

func TestTrimTransformer(t *testing.T) {
	tr := transform.TrimTransformer{}
	cases := []struct {
		input, want string
	}{
		{"  hello  ", "hello"},
		{"\t log line \n", "log line"},
		{"no spaces", "no spaces"},
		{"", ""},
	}
	for _, c := range cases {
		got := tr.Apply(c.input)
		if got != c.want {
			t.Errorf("Trim(%q) = %q; want %q", c.input, got, c.want)
		}
	}
}

func TestReplaceTransformer(t *testing.T) {
	tr := transform.ReplaceTransformer{Old: "foo", New: "bar"}
	got := tr.Apply("foo and foo")
	if got != "bar and bar" {
		t.Errorf("got %q; want %q", got, "bar and bar")
	}
}

func TestReplaceTransformer_EmptyOld(t *testing.T) {
	tr := transform.ReplaceTransformer{Old: "", New: "x"}
	input := "unchanged"
	if got := tr.Apply(input); got != input {
		t.Errorf("expected unchanged, got %q", got)
	}
}

func TestNewRegex_InvalidPattern(t *testing.T) {
	_, err := transform.NewRegex("[invalid", "")
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestRegexTransformer_Apply(t *testing.T) {
	tr, err := transform.NewRegex(`\d+`, "NUM")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := tr.Apply("error 404 on line 12")
	want := "error NUM on line NUM"
	if got != want {
		t.Errorf("got %q; want %q", got, want)
	}
}

func TestChain_AppliesInOrder(t *testing.T) {
	re, _ := transform.NewRegex(`\d+`, "NUM")
	transformers := []transform.Transformer{
		transform.TrimTransformer{},
		transform.ReplaceTransformer{Old: "error", New: "ERR"},
		re,
	}
	got := transform.Chain("  error 404  ", transformers)
	want := "ERR NUM"
	if got != want {
		t.Errorf("got %q; want %q", got, want)
	}
}

func TestChain_EmptyTransformers(t *testing.T) {
	input := "unchanged line"
	got := transform.Chain(input, nil)
	if got != input {
		t.Errorf("expected %q, got %q", input, got)
	}
}

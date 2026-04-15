package transform_test

import (
	"testing"

	"greplive/internal/transform"
)

func TestBuild_Trim(t *testing.T) {
	steps := []transform.Step{{Kind: "trim"}}
	transformers, err := transform.Build(steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(transformers) != 1 {
		t.Fatalf("expected 1 transformer, got %d", len(transformers))
	}
	got := transformers[0].Apply("  hello  ")
	if got != "hello" {
		t.Errorf("got %q; want %q", got, "hello")
	}
}

func TestBuild_Replace(t *testing.T) {
	steps := []transform.Step{{Kind: "replace", Old: "foo", New: "baz"}}
	transformers, err := transform.Build(steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := transformers[0].Apply("foo bar foo")
	if got != "baz bar baz" {
		t.Errorf("got %q", got)
	}
}

func TestBuild_Regex(t *testing.T) {
	steps := []transform.Step{{Kind: "regex", Pattern: `\bERROR\b`, New: "[ERROR]"}}
	transformers, err := transform.Build(steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := transformers[0].Apply("ERROR: disk full")
	want := "[ERROR]: disk full"
	if got != want {
		t.Errorf("got %q; want %q", got, want)
	}
}

func TestBuild_InvalidRegex(t *testing.T) {
	steps := []transform.Step{{Kind: "regex", Pattern: "[bad", New: ""}}
	_, err := transform.Build(steps)
	if err == nil {
		t.Fatal("expected error for invalid regex pattern")
	}
}

func TestBuild_UnknownKind(t *testing.T) {
	steps := []transform.Step{{Kind: "unknown"}}
	_, err := transform.Build(steps)
	if err == nil {
		t.Fatal("expected error for unknown kind")
	}
}

func TestBuild_MultipleSteps(t *testing.T) {
	steps := []transform.Step{
		{Kind: "trim"},
		{Kind: "replace", Old: "WARN", New: "WARNING"},
	}
	transformers, err := transform.Build(steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := transform.Chain("  WARN: low memory  ", transformers)
	want := "WARNING: low memory"
	if got != want {
		t.Errorf("got %q; want %q", got, want)
	}
}

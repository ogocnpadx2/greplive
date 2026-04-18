package labelmap_test

import (
	"strings"
	"testing"

	"greplive/internal/labelmap"
)

func TestConfig_Build_Empty(t *testing.T) {
	c := labelmap.Config{}
	l, err := c.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Enabled() {
		t.Fatal("expected disabled labeler for empty config")
	}
}

func TestConfig_Build_ValidLabels(t *testing.T) {
	c := labelmap.Config{Labels: []string{"env=prod", "svc=api"}}
	l, err := c.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !l.Enabled() {
		t.Fatal("expected enabled labeler")
	}
	got := l.Apply("msg")
	if !strings.Contains(got, "env=prod") || !strings.Contains(got, "svc=api") {
		t.Fatalf("labels missing in %q", got)
	}
}

func TestConfig_Build_MissingEquals(t *testing.T) {
	c := labelmap.Config{Labels: []string{"badlabel"}}
	_, err := c.Build()
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestConfig_Build_EmptyKey(t *testing.T) {
	c := labelmap.Config{Labels: []string{"=value"}}
	_, err := c.Build()
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestConfig_Build_ValueWithEquals(t *testing.T) {
	c := labelmap.Config{Labels: []string{"url=http://x?a=1"}}
	l, err := c.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := l.Apply("line")
	if !strings.Contains(got, "url=http://x?a=1") {
		t.Fatalf("expected full value in %q", got)
	}
}

package labelmap_test

import (
	"strings"
	"testing"

	"greplive/internal/labelmap"
)

func TestNew_EmptyLabels_Disabled(t *testing.T) {
	l := labelmap.New(nil)
	if l.Enabled() {
		t.Fatal("expected disabled for nil labels")
	}
}

func TestNew_WithLabels_Enabled(t *testing.T) {
	l := labelmap.New(map[string]string{"env": "prod"})
	if !l.Enabled() {
		t.Fatal("expected enabled")
	}
}

func TestApply_Disabled_Unchanged(t *testing.T) {
	l := labelmap.New(map[string]string{})
	got := l.Apply("hello")
	if got != "hello" {
		t.Fatalf("expected unchanged, got %q", got)
	}
}

func TestApply_PrependsSingleLabel(t *testing.T) {
	l := labelmap.New(map[string]string{"env": "prod"})
	got := l.Apply("msg")
	if !strings.HasPrefix(got, "[") {
		t.Fatalf("expected bracket prefix, got %q", got)
	}
	if !strings.Contains(got, "env=prod") {
		t.Fatalf("expected env=prod in %q", got)
	}
	if !strings.HasSuffix(got, "] msg") {
		t.Fatalf("expected line appended, got %q", got)
	}
}

func TestApply_MultipleLabels(t *testing.T) {
	l := labelmap.New(map[string]string{"env": "prod", "svc": "api"})
	got := l.Apply("line")
	if !strings.Contains(got, "env=prod") || !strings.Contains(got, "svc=api") {
		t.Fatalf("missing labels in %q", got)
	}
}

func TestApplyAll_MultipleLabelers(t *testing.T) {
	labelers := []*labelmap.Labeler{
		labelmap.New(map[string]string{"env": "prod"}),
		labelmap.New(map[string]string{"svc": "api"}),
	}
	got := labelmap.ApplyAll("msg", labelers)
	if !strings.Contains(got, "env=prod") || !strings.Contains(got, "svc=api") {
		t.Fatalf("expected both labels in %q", got)
	}
}

func TestApplyAll_NoLabelers_Unchanged(t *testing.T) {
	got := labelmap.ApplyAll("msg", nil)
	if got != "msg" {
		t.Fatalf("expected unchanged, got %q", got)
	}
}

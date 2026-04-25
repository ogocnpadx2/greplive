package fieldmap_test

import (
	"testing"

	"greplive/internal/fieldmap"
)

func TestNew_EmptyMappings_Disabled(t *testing.T) {
	m := fieldmap.New(map[string]string{})
	if m.Enabled() {
		t.Fatal("expected disabled for empty mappings")
	}
}

func TestNew_WithMappings_Enabled(t *testing.T) {
	m := fieldmap.New(map[string]string{"lvl": "level"})
	if !m.Enabled() {
		t.Fatal("expected enabled")
	}
}

func TestNew_SkipsEmptyKeys(t *testing.T) {
	m := fieldmap.New(map[string]string{"": "level", "lvl": ""})
	if m.Enabled() {
		t.Fatal("expected disabled — all entries had empty key or value")
	}
}

func TestApply_Disabled_Unchanged(t *testing.T) {
	m := fieldmap.New(map[string]string{})
	line := "lvl=info msg=hello"
	if got := m.Apply(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestApply_RenamesSingleField(t *testing.T) {
	m := fieldmap.New(map[string]string{"lvl": "level"})
	got := m.Apply("lvl=info msg=hello")
	want := "level=info msg=hello"
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestApply_RenamesMultipleFields(t *testing.T) {
	m := fieldmap.New(map[string]string{"lvl": "level", "msg": "message"})
	got := m.Apply("lvl=warn msg=oops ts=now")
	want := "level=warn message=oops ts=now"
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestApply_NonKVToken_Untouched(t *testing.T) {
	m := fieldmap.New(map[string]string{"lvl": "level"})
	got := m.Apply("INFO lvl=info msg=hello")
	want := "INFO level=info msg=hello"
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestApply_ValueContainsEquals(t *testing.T) {
	m := fieldmap.New(map[string]string{"url": "endpoint"})
	got := m.Apply("url=http://x.com?a=1 status=200")
	want := "endpoint=http://x.com?a=1 status=200"
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestApplyAll_MultipleMappers(t *testing.T) {
	ms := []*fieldmap.Mapper{
		fieldmap.New(map[string]string{"lvl": "level"}),
		fieldmap.New(map[string]string{"msg": "message"}),
	}
	got := fieldmap.ApplyAll(ms, "lvl=error msg=boom")
	want := "level=error message=boom"
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

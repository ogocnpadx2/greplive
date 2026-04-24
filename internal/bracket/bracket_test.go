package bracket

import (
	"testing"
)

func TestNew_EmptyPattern_Disabled(t *testing.T) {
	b, err := New("", "[", "]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Enabled() {
		t.Fatal("expected disabled for empty pattern")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := New("[invalid", "[", "]")
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_ValidPattern_Enabled(t *testing.T) {
	b, err := New(`\d+`, "[", "]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !b.Enabled() {
		t.Fatal("expected enabled")
	}
}

func TestNew_DefaultBrackets(t *testing.T) {
	b, err := New(`\d+`, "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := b.Apply("port 8080 ready")
	want := "port [8080] ready"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestApply_Disabled_Unchanged(t *testing.T) {
	b, _ := New("", "(", ")")
	line := "hello world"
	if got := b.Apply(line); got != line {
		t.Fatalf("got %q, want %q", got, line)
	}
}

func TestApply_WrapsMatches(t *testing.T) {
	b, err := New(`ERROR|WARN`, "<<", ">>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := b.Apply("ERROR: something WARN happened")
	want := "<<ERROR>>: something <<WARN>> happened"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestApply_NoMatch_Unchanged(t *testing.T) {
	b, err := New(`FATAL`, "[", "]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := "INFO: all good"
	if got := b.Apply(line); got != line {
		t.Fatalf("got %q, want %q", got, line)
	}
}

func TestApplyAll_MultipleTransformers(t *testing.T) {
	b1, _ := New(`\d+`, "[", "]")
	b2, _ := New(`ERROR`, "{", "}")
	line := "ERROR code 42"
	got := ApplyAll([]*Bracket{b1, b2}, line)
	want := "{ERROR} code [42]"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestApplyAll_EmptySlice_Unchanged(t *testing.T) {
	line := "no change expected"
	if got := ApplyAll(nil, line); got != line {
		t.Fatalf("got %q, want %q", got, line)
	}
}

package linenum

import (
	"testing"
)

func TestNew_Disabled(t *testing.T) {
	n := New(false, 0)
	if n.Enabled() {
		t.Fatal("expected disabled")
	}
}

func TestNew_Enabled(t *testing.T) {
	n := New(true, 0)
	if !n.Enabled() {
		t.Fatal("expected enabled")
	}
}

func TestApply_Disabled_Unchanged(t *testing.T) {
	n := New(false, 0)
	got := n.Apply("hello")
	if got != "hello" {
		t.Fatalf("want %q got %q", "hello", got)
	}
	if n.Count() != 1 {
		t.Fatalf("count should still increment: got %d", n.Count())
	}
}

func TestApply_Enabled_PrependNumber(t *testing.T) {
	n := New(true, 0)
	got := n.Apply("hello")
	if got != "1 hello" {
		t.Fatalf("want %q got %q", "1 hello", got)
	}
	got = n.Apply("world")
	if got != "2 world" {
		t.Fatalf("want %q got %q", "2 world", got)
	}
}

func TestApply_Padded(t *testing.T) {
	n := New(true, 4)
	got := n.Apply("line")
	if got != "0001 line" {
		t.Fatalf("want %q got %q", "0001 line", got)
	}
}

func TestReset_ZerosCount(t *testing.T) {
	n := New(true, 0)
	n.Apply("a")
	n.Apply("b")
	n.Reset()
	if n.Count() != 0 {
		t.Fatalf("expected 0 after reset, got %d", n.Count())
	}
	got := n.Apply("c")
	if got != "1 c" {
		t.Fatalf("want %q got %q", "1 c", got)
	}
}

func TestApply_NegativePad_TreatedAsZero(t *testing.T) {
	n := New(true, -5)
	got := n.Apply("x")
	if got != "1 x" {
		t.Fatalf("want %q got %q", "1 x", got)
	}
}

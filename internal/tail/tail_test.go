package tail

import (
	"testing"
)

func TestNew_ClampsZero(t *testing.T) {
	b := New(0)
	if b.cap != 1 {
		t.Fatalf("expected cap 1, got %d", b.cap)
	}
}

func TestNew_PositiveCapacity(t *testing.T) {
	b := New(5)
	if b.cap != 5 {
		t.Fatalf("expected cap 5, got %d", b.cap)
	}
}

func TestPush_BelowCapacity(t *testing.T) {
	b := New(4)
	b.Push("a")
	b.Push("b")
	got := b.Lines()
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("unexpected lines: %v", got)
	}
}

func TestPush_AtCapacity(t *testing.T) {
	b := New(3)
	b.Push("a")
	b.Push("b")
	b.Push("c")
	got := b.Lines()
	if len(got) != 3 || got[0] != "a" || got[2] != "c" {
		t.Fatalf("unexpected lines: %v", got)
	}
}

func TestPush_Overflow_OldestEvicted(t *testing.T) {
	b := New(3)
	b.Push("a")
	b.Push("b")
	b.Push("c")
	b.Push("d")
	got := b.Lines()
	if got[0] != "b" || got[1] != "c" || got[2] != "d" {
		t.Fatalf("expected [b c d], got %v", got)
	}
}

func TestLen_BeforeAndAfterFull(t *testing.T) {
	b := New(3)
	if b.Len() != 0 {
		t.Fatal("expected 0")
	}
	b.Push("x")
	if b.Len() != 1 {
		t.Fatal("expected 1")
	}
	b.Push("y")
	b.Push("z")
	b.Push("w")
	if b.Len() != 3 {
		t.Fatalf("expected 3, got %d", b.Len())
	}
}

func TestReset_ClearsBuffer(t *testing.T) {
	b := New(3)
	b.Push("a")
	b.Push("b")
	b.Reset()
	if b.Len() != 0 {
		t.Fatal("expected 0 after reset")
	}
	if got := b.Lines(); len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

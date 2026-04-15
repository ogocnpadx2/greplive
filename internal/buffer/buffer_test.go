package buffer_test

import (
	"testing"

	"github.com/user/greplive/internal/buffer"
)

func TestNew_ZeroSize(t *testing.T) {
	b := buffer.New(0)
	if b.Cap() != 0 {
		t.Fatalf("expected cap 0, got %d", b.Cap())
	}
	b.Push("ignored")
	if got := b.Snapshot(); len(got) != 0 {
		t.Fatalf("expected empty snapshot, got %v", got)
	}
}

func TestNew_NegativeSize(t *testing.T) {
	b := buffer.New(-5)
	if b.Cap() != 0 {
		t.Fatalf("expected cap 0 for negative size, got %d", b.Cap())
	}
}

func TestPush_BelowCapacity(t *testing.T) {
	b := buffer.New(4)
	b.Push("a")
	b.Push("b")
	snap := b.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap))
	}
	if snap[0] != "a" || snap[1] != "b" {
		t.Errorf("unexpected order: %v", snap)
	}
}

func TestPush_AtCapacity(t *testing.T) {
	b := buffer.New(3)
	for _, l := range []string{"x", "y", "z"} {
		b.Push(l)
	}
	snap := b.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected 3, got %d", len(snap))
	}
	if snap[0] != "x" || snap[2] != "z" {
		t.Errorf("unexpected snapshot: %v", snap)
	}
}

func TestPush_Overflow_OldestEvicted(t *testing.T) {
	b := buffer.New(3)
	for _, l := range []string{"a", "b", "c", "d"} {
		b.Push(l)
	}
	snap := b.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected 3, got %d", len(snap))
	}
	if snap[0] != "b" || snap[1] != "c" || snap[2] != "d" {
		t.Errorf("expected [b c d], got %v", snap)
	}
}

func TestLen(t *testing.T) {
	b := buffer.New(5)
	if b.Len() != 0 {
		t.Fatalf("expected 0, got %d", b.Len())
	}
	b.Push("one")
	b.Push("two")
	if b.Len() != 2 {
		t.Fatalf("expected 2, got %d", b.Len())
	}
}

func TestReset_ClearsBuffer(t *testing.T) {
	b := buffer.New(4)
	b.Push("line1")
	b.Push("line2")
	b.Reset()
	if b.Len() != 0 {
		t.Fatalf("expected 0 after reset, got %d", b.Len())
	}
	if snap := b.Snapshot(); len(snap) != 0 {
		t.Errorf("expected empty snapshot after reset, got %v", snap)
	}
}

func TestSnapshot_IsACopy(t *testing.T) {
	b := buffer.New(3)
	b.Push("original")
	snap := b.Snapshot()
	snap[0] = "mutated"
	snap2 := b.Snapshot()
	if snap2[0] != "original" {
		t.Errorf("snapshot mutation affected ring buffer")
	}
}

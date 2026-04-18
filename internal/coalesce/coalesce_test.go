package coalesce

import (
	"testing"
	"time"
)

func collect(c *Coalescer, lines []string) []string {
	for _, l := range lines {
		c.Push(l)
	}
	c.Flush()
	return nil // results captured via closure
}

func TestPush_UniqueLines_EmittedAsIs(t *testing.T) {
	var got []string
	c := New(50*time.Millisecond, func(s string) { got = append(got, s) })
	c.Push("alpha")
	c.Push("beta")
	c.Flush()
	if len(got) != 2 || got[0] != "alpha" || got[1] != "beta" {
		t.Fatalf("unexpected output: %v", got)
	}
}

func TestPush_RepeatedLine_Annotated(t *testing.T) {
	var got []string
	c := New(50*time.Millisecond, func(s string) { got = append(got, s) })
	c.Push("hello")
	c.Push("hello")
	c.Push("hello")
	c.Flush()
	if len(got) != 1 {
		t.Fatalf("expected 1 line, got %d: %v", len(got), got)
	}
	expected := "hello  [x3]"
	if got[0] != expected {
		t.Fatalf("got %q, want %q", got[0], expected)
	}
}

func TestPush_SingleOccurrence_NoAnnotation(t *testing.T) {
	var got []string
	c := New(50*time.Millisecond, func(s string) { got = append(got, s) })
	c.Push("once")
	c.Flush()
	if len(got) != 1 || got[0] != "once" {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestPush_BurstThenDifferent_FlushesFirst(t *testing.T) {
	var got []string
	c := New(50*time.Millisecond, func(s string) { got = append(got, s) })
	c.Push("a")
	c.Push("a")
	c.Push("b")
	c.Flush()
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %v", got)
	}
	if got[0] != "a  [x2]" {
		t.Fatalf("got[0] = %q", got[0])
	}
	if got[1] != "b" {
		t.Fatalf("got[1] = %q", got[1])
	}
}

func TestFlush_Empty_NoEmit(t *testing.T) {
	called := false
	c := New(50*time.Millisecond, func(s string) { called = true })
	c.Flush()
	if called {
		t.Fatal("emit should not be called on empty flush")
	}
}

func TestNew_ZeroWindow_UsesDefault(t *testing.T) {
	c := New(0, func(string) {})
	if c.window != 200*time.Millisecond {
		t.Fatalf("expected default window, got %v", c.window)
	}
}

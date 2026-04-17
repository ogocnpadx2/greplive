package contextline

import (
	"testing"
)

func lines(c *Collector, pairs ...interface{}) []string {
	var out []string
	for i := 0; i < len(pairs); i += 2 {
		l := pairs[i].(string)
		m := pairs[i+1].(bool)
		out = append(out, c.Feed(l, m)...)
	}
	return out
}

func TestNew_NegativeClampedToZero(t *testing.T) {
	c := New(-1, -2)
	if c.before != 0 || c.after != 0 {
		t.Fatal("expected before and after to be clamped to 0")
	}
}

func TestFeed_MatchOnly(t *testing.T) {
	c := New(0, 0)
	out := c.Feed("hit", true)
	if len(out) != 1 || out[0] != "hit" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestFeed_NoMatch_NoContext(t *testing.T) {
	c := New(0, 0)
	out := c.Feed("miss", false)
	if len(out) != 0 {
		t.Fatalf("expected no output, got %v", out)
	}
}

func TestFeed_PreContext(t *testing.T) {
	c := New(2, 0)
	out := lines(c, "a", false, "b", false, "c", false, "hit", true)
	// expect b, c, hit
	if len(out) != 3 || out[0] != "b" || out[1] != "c" || out[2] != "hit" {
		t.Fatalf("unexpected pre-context output: %v", out)
	}
}

func TestFeed_PostContext(t *testing.T) {
	c := New(0, 2)
	out := lines(c, "hit", true, "p1", false, "p2", false, "p3", false)
	// expect hit, p1, p2
	if len(out) != 3 || out[0] != "hit" || out[1] != "p1" || out[2] != "p2" {
		t.Fatalf("unexpected post-context output: %v", out)
	}
}

func TestFeed_PreAndPost(t *testing.T) {
	c := New(1, 1)
	out := lines(c, "pre", false, "hit", true, "post", false, "gone", false)
	if len(out) != 3 {
		t.Fatalf("expected 3 lines, got %v", out)
	}
	if out[0] != "pre" || out[1] != "hit" || out[2] != "post" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestReset_ClearsPending(t *testing.T) {
	c := New(1, 3)
	c.Feed("hit", true)
	c.Reset()
	out := c.Feed("after", false)
	if len(out) != 0 {
		t.Fatalf("expected no output after reset, got %v", out)
	}
}

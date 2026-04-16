package multiline

import (
	"testing"
)

func TestNew_InvalidStartPattern(t *testing.T) {
	_, err := New("[", "", 0)
	if err == nil {
		t.Fatal("expected error for invalid start pattern")
	}
}

func TestNew_InvalidContinuePattern(t *testing.T) {
	_, err := New("^ERROR", "[", 0)
	if err == nil {
		t.Fatal("expected error for invalid continue pattern")
	}
}

func TestPush_SingleLine(t *testing.T) {
	c, _ := New("^ERROR", "", 0)
	event, ok := c.Push("INFO hello")
	if !ok || event != "INFO hello" {
		t.Fatalf("expected standalone line, got %q %v", event, ok)
	}
}

func TestPush_MultiLine_StackTrace(t *testing.T) {
	c, _ := New("^ERROR", "", 0)
	// Start buffering
	_, ok := c.Push("ERROR something failed")
	if ok {
		t.Fatal("expected buffered, not flushed")
	}
	_, ok = c.Push("\tat foo.go:10")
	if ok {
		t.Fatal("expected buffered")
	}
	// New start flushes previous
	event, ok := c.Push("ERROR second error")
	if !ok {
		t.Fatal("expected flush on new start")
	}
	if event != "ERROR something failed\n\tat foo.go:10" {
		t.Fatalf("unexpected event: %q", event)
	}
}

func TestFlush_ReturnsPending(t *testing.T) {
	c, _ := New("^ERROR", "", 0)
	c.Push("ERROR boom")
	c.Push("\tcaused by: x")
	event := c.Flush()
	if event != "ERROR boom\n\tcaused by: x" {
		t.Fatalf("unexpected flush: %q", event)
	}
	if second := c.Flush(); second != "" {
		t.Fatalf("expected empty second flush, got %q", second)
	}
}

func TestPush_MaxLines_Respected(t *testing.T) {
	c, _ := New("^ERROR", "", 2)
	c.Push("ERROR start")
	c.Push("\tline1")
	// Third continuation should NOT be buffered (max=2 already hit)
	event, ok := c.Push("\tline2")
	if !ok {
		t.Fatal("expected line2 to be emitted standalone when max exceeded")
	}
	_ = event
}

func TestPush_CustomContinuePattern(t *testing.T) {
	c, _ := New("^\\[", "^  ", 0)
	c.Push("[2024] start")
	c.Push("  detail")
	event := c.Flush()
	if event != "[2024] start\n  detail" {
		t.Fatalf("got %q", event)
	}
}

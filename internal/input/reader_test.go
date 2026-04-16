package input_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"greplive/internal/input"
)

func TestLines_Basic(t *testing.T) {
	src := strings.NewReader("line one\nline two\nline three\n")
	r := input.New(src)

	ctx := context.Background()
	lines, errs := r.Lines(ctx)

	var got []string
	for line := range lines {
		got = append(got, line)
	}

	if err := <-errs; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"line one", "line two", "line three"}
	if len(got) != len(expected) {
		t.Fatalf("expected %d lines, got %d", len(expected), len(got))
	}
	for i, e := range expected {
		if got[i] != e {
			t.Errorf("line %d: expected %q, got %q", i, e, got[i])
		}
	}
}

func TestLines_SkipsEmptyLines(t *testing.T) {
	src := strings.NewReader("alpha\n\nbeta\n\n\ngamma")
	r := input.New(src)

	lines, _ := r.Lines(context.Background())

	var got []string
	for line := range lines {
		got = append(got, line)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 non-empty lines, got %d: %v", len(got), got)
	}
}

func TestLines_ContextCancellation(t *testing.T) {
	// Write many lines and cancel early to verify the goroutine stops.
	var sb strings.Builder
	for i := 0; i < 1000; i++ {
		sb.WriteString("log line\n")
	}

	r := input.New(strings.NewReader(sb.String()))
	ctx, cancel := context.WithCancel(context.Background())

	lines, errs := r.Lines(ctx)

	// Read one line then cancel.
	<-lines
	cancel()

	// Drain remaining lines with a timeout to avoid hanging.
	done := make(chan struct{})
	go func() {
		for range lines {
		}
		<-errs
		close(done)
	}()

	select {
	case <-done:
		// success
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Lines to stop after context cancellation")
	}
}

func TestLines_EmptyReader(t *testing.T) {
	r := input.New(strings.NewReader(""))
	lines, errs := r.Lines(context.Background())

	var got []string
	for line := range lines {
		got = append(got, line)
	}

	if err := <-errs; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected no lines, got %v", got)
	}
}

func TestLines_NoTrailingNewline(t *testing.T) {
	// Ensure the last line is emitted even when not terminated by a newline.
	src := strings.NewReader("first\nsecond\nthird")
	r := input.New(src)

	lines, errs := r.Lines(context.Background())

	var got []string
	for line := range lines {
		got = append(got, line)
	}

	if err := <-errs; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"first", "second", "third"}
	if len(got) != len(expected) {
		t.Fatalf("expected %d lines, got %d: %v", len(expected), len(got), got)
	}
	for i, e := range expected {
		if got[i] != e {
			t.Errorf("line %d: expected %q, got %q", i, e, got[i])
		}
	}
}

package splice_test

import (
	"context"
	"sort"
	"testing"
	"time"

	"greplive/internal/splice"
)

func collect(t *testing.T, ch <-chan string, timeout time.Duration) []string {
	t.Helper()
	var lines []string
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		select {
		case line, ok := <-ch:
			if !ok {
				return lines
			}
			lines = append(lines, line)
		case <-timer.C:
			t.Fatal("collect timed out")
			return lines
		}
	}
}

func TestSplice_MergesTwoChannels(t *testing.T) {
	ctx := context.Background()
	a := make(chan string, 2)
	b := make(chan string, 2)
	a <- "alpha"
	a <- "beta"
	close(a)
	b <- "gamma"
	close(b)

	s := splice.New(ctx, a, b)
	got := collect(t, s.Out(), 2*time.Second)
	sort.Strings(got)

	want := []string{"alpha", "beta", "gamma"}
	if len(got) != len(want) {
		t.Fatalf("want %d lines, got %d: %v", len(want), len(got), got)
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("index %d: want %q, got %q", i, w, got[i])
		}
	}
}

func TestSplice_EmptyInputs_ClosesOutput(t *testing.T) {
	ctx := context.Background()
	a := make(chan string)
	close(a)
	s := splice.New(ctx, a)
	_, ok := <-s.Out()
	if ok {
		t.Fatal("expected output channel to be closed")
	}
}

func TestSplice_ContextCancellation_StopsEarly(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	a := make(chan string) // never sends
	s := splice.New(ctx, a)
	cancel()

	select {
	case <-s.Out():
		// drained or closed – acceptable
	case <-time.After(time.Second):
		t.Fatal("output channel not closed after context cancel")
	}
}

func TestSplice_SingleSource_PassesThrough(t *testing.T) {
	ctx := context.Background()
	a := make(chan string, 3)
	a <- "one"
	a <- "two"
	a <- "three"
	close(a)

	s := splice.New(ctx, a)
	got := collect(t, s.Out(), 2*time.Second)
	if len(got) != 3 {
		t.Fatalf("want 3 lines, got %d", len(got))
	}
}

func TestSplice_NoSources_ClosesImmediately(t *testing.T) {
	ctx := context.Background()
	s := splice.New(ctx)
	select {
	case _, ok := <-s.Out():
		if ok {
			t.Fatal("expected closed channel")
		}
	case <-time.After(time.Second):
		t.Fatal("channel not closed with no sources")
	}
}

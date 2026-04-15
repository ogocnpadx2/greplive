package input

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeLine(t *testing.T, f *os.File, line string) {
	t.Helper()
	_, err := f.WriteString(line + "\n")
	if err != nil {
		t.Fatalf("writeLine: %v", err)
	}
}

func TestTailFile_ReadsExistingLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	writeLine(t, f, "line one")
	writeLine(t, f, "line two")
	f.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch, err := TailFile(ctx, path, TailOptions{PollInterval: 50 * time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}

	var got []string
	for line := range ch {
		got = append(got, line)
		if len(got) == 2 {
			cancel()
		}
	}

	if len(got) < 2 {
		t.Fatalf("expected at least 2 lines, got %d", len(got))
	}
	if got[0] != "line one" {
		t.Errorf("got[0] = %q, want %q", got[0], "line one")
	}
	if got[1] != "line two" {
		t.Errorf("got[1] = %q, want %q", got[1], "line two")
	}
}

func TestTailFile_NonExistentFile(t *testing.T) {
	ctx := context.Background()
	_, err := TailFile(ctx, "/nonexistent/path/to/file.log", TailOptions{})
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
}

func TestTailFile_ContextCancellation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cancel.log")

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	ctx, cancel := context.WithCancel(context.Background())
	ch, err := TailFile(ctx, path, TailOptions{PollInterval: 20 * time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}

	cancel()

	select {
	case <-ch:
		// channel closed as expected
	case <-time.After(time.Second):
		t.Fatal("channel was not closed after context cancellation")
	}
}

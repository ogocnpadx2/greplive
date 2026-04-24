package reopen_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/greplive/internal/reopen"
)

func writeTmp(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "test.log")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestNew_OpensFile(t *testing.T) {
	p := writeTmp(t, "hello\n")
	r, err := reopen.New(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer r.Close()
}

func TestNew_NonExistent_ReturnsError(t *testing.T) {
	_, err := reopen.New("/no/such/file.log")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRead_ReturnsContent(t *testing.T) {
	p := writeTmp(t, "data")
	r, _ := reopen.New(p)
	defer r.Close()

	buf := make([]byte, 4)
	n, err := r.Read(buf)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if string(buf[:n]) != "data" {
		t.Fatalf("got %q, want %q", buf[:n], "data")
	}
}

func TestReopen_ReadsFromBeginning(t *testing.T) {
	p := writeTmp(t, "abcd")
	r, _ := reopen.New(p)
	defer r.Close()

	buf := make([]byte, 4)
	r.Read(buf) // consume

	if err := r.Reopen(); err != nil {
		t.Fatalf("reopen error: %v", err)
	}
	n, _ := r.Read(buf)
	if string(buf[:n]) != "abcd" {
		t.Fatalf("got %q after reopen", buf[:n])
	}
}

func TestClose_SecondCloseNoError(t *testing.T) {
	p := writeTmp(t, "x")
	r, _ := reopen.New(p)
	r.Close()
	if err := r.Close(); err != nil {
		t.Fatalf("second close: %v", err)
	}
}

func TestWatchAndReopen_ReopensOnSignal(t *testing.T) {
	p := writeTmp(t, "first")
	r, _ := reopen.New(p)
	defer r.Close()

	ch := make(chan struct{}, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.WatchAndReopen(ctx, ch)

	// overwrite with new content
	os.WriteFile(p, []byte("second"), 0o644)
	ch <- struct{}{}
	time.Sleep(20 * time.Millisecond)

	buf := make([]byte, 6)
	n, _ := r.Read(buf)
	if string(buf[:n]) != "second" {
		t.Fatalf("got %q, want %q", buf[:n], "second")
	}
}

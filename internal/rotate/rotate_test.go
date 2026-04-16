package rotate_test

import (
	"os"
	"testing"
	"time"

	"greplive/internal/rotate"
)

func TestNew_DefaultInterval(t *testing.T) {
	w := rotate.New("file.log", 0)
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}

func TestNotify_ChannelNotNil(t *testing.T) {
	w := rotate.New("file.log", time.Second)
	if w.Notify() == nil {
		t.Fatal("Notify channel should not be nil")
	}
}

func TestStop_DoesNotPanic(t *testing.T) {
	w := rotate.New("file.log", 50*time.Millisecond)
	w.Start()
	time.Sleep(20 * time.Millisecond)
	w.Stop()
}

func TestDetectsRotation(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "log")
	if err != nil {
		t.Fatal(err)
	}
	path := f.Name()
	_, _ = f.WriteString("initial content\n")
	f.Close()

	w := rotate.New(path, 30*time.Millisecond)
	w.Start()
	defer w.Stop()

	// Simulate rotation: truncate file
	time.Sleep(20 * time.Millisecond)
	_ = os.WriteFile(path, []byte("new content\n"), 0o644)

	select {
	case <-w.Notify():
		// success
	case <-time.After(500 * time.Millisecond):
		t.Fatal("expected rotation notification")
	}
}

func TestNoNotify_WhenFileUnchanged(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "log")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString("same\n")
	f.Close()

	w := rotate.New(f.Name(), 30*time.Millisecond)
	w.Start()
	defer w.Stop()

	select {
	case <-w.Notify():
		t.Fatal("unexpected rotation notification")
	case <-time.After(150 * time.Millisecond):
		// success
	}
}

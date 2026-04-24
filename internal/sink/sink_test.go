package sink

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNew_NilWriter_DefaultsToStdout(t *testing.T) {
	s := New(nil)
	if s.w != os.Stdout {
		t.Fatal("expected os.Stdout as default writer")
	}
}

func TestWrite_AppendsNewline(t *testing.T) {
	var buf bytes.Buffer
	s := New(&buf)
	if err := s.Write("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got != "hello\n" {
		t.Fatalf("want %q, got %q", "hello\n", got)
	}
}

func TestWriteAll_MultipleLines(t *testing.T) {
	var buf bytes.Buffer
	s := New(&buf)
	lines := []string{"alpha", "beta", "gamma"}
	if err := s.WriteAll(lines); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	for _, l := range lines {
		if !strings.Contains(got, l) {
			t.Errorf("output missing line %q", l)
		}
	}
}

func TestNewFile_CreatesAndWrites(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.log")

	s, close, err := NewFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer close()

	if err := s.Write("file line"); err != nil {
		t.Fatalf("write error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if !strings.Contains(string(data), "file line") {
		t.Errorf("file content missing expected line, got %q", string(data))
	}
}

func TestNewFile_InvalidPath_ReturnsError(t *testing.T) {
	_, _, err := NewFile("/no/such/dir/out.log")
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestMulti_WritesToAllSinks(t *testing.T) {
	var b1, b2 bytes.Buffer
	s1 := New(&b1)
	s2 := New(&b2)
	m := Multi(s1, s2)

	if err := m.Write("broadcast"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, b := range []*bytes.Buffer{&b1, &b2} {
		if !strings.Contains(b.String(), "broadcast") {
			t.Errorf("sink %d missing output", i+1)
		}
	}
}

func TestWrite_ConcurrentSafe(t *testing.T) {
	var buf bytes.Buffer
	s := New(&buf)
	done := make(chan struct{})
	for i := 0; i < 50; i++ {
		go func() {
			_ = s.Write("concurrent")
			done <- struct{}{}
		}()
	}
	for i := 0; i < 50; i++ {
		<-done
	}
}

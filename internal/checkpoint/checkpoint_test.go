package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/greplive/internal/checkpoint"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestNew_NoFile_ReturnsEmpty(t *testing.T) {
	c, err := checkpoint.New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := c.Get("/var/log/app.log"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestSet_PersistsOffset(t *testing.T) {
	p := tempPath(t)
	c, _ := checkpoint.New(p)
	if err := c.Set("/var/log/app.log", 1024); err != nil {
		t.Fatalf("Set error: %v", err)
	}
	if got := c.Get("/var/log/app.log"); got != 1024 {
		t.Fatalf("expected 1024, got %d", got)
	}
}

func TestNew_LoadsExistingFile(t *testing.T) {
	p := tempPath(t)
	c1, _ := checkpoint.New(p)
	_ = c1.Set("/var/log/app.log", 4096)

	c2, err := checkpoint.New(p)
	if err != nil {
		t.Fatalf("reload error: %v", err)
	}
	if got := c2.Get("/var/log/app.log"); got != 4096 {
		t.Fatalf("expected 4096, got %d", got)
	}
}

func TestSet_MultipleFiles(t *testing.T) {
	p := tempPath(t)
	c, _ := checkpoint.New(p)
	_ = c.Set("a.log", 10)
	_ = c.Set("b.log", 20)
	if got := c.Get("a.log"); got != 10 {
		t.Fatalf("a.log: expected 10, got %d", got)
	}
	if got := c.Get("b.log"); got != 20 {
		t.Fatalf("b.log: expected 20, got %d", got)
	}
}

func TestNew_CorruptFile_ReturnsError(t *testing.T) {
	p := tempPath(t)
	_ = os.WriteFile(p, []byte("not json"), 0o644)
	_, err := checkpoint.New(p)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}

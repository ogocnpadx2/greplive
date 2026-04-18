package offset_test

import (
	"testing"

	"github.com/user/greplive/internal/offset"
)

func TestNew_StartsAtInitial(t *testing.T) {
	tr := offset.New(100)
	if got := tr.Get(); got != 100 {
		t.Fatalf("expected 100, got %d", got)
	}
}

func TestNew_ZeroInitial(t *testing.T) {
	tr := offset.New(0)
	if got := tr.Get(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestAdd_IncrementsOffset(t *testing.T) {
	tr := offset.New(0)
	got := tr.Add(42)
	if got != 42 {
		t.Fatalf("expected 42, got %d", got)
	}
	if tr.Get() != 42 {
		t.Fatalf("Get mismatch")
	}
}

func TestAdd_Cumulative(t *testing.T) {
	tr := offset.New(10)
	tr.Add(5)
	tr.Add(5)
	if got := tr.Get(); got != 20 {
		t.Fatalf("expected 20, got %d", got)
	}
}

func TestReset_ZerosOffset(t *testing.T) {
	tr := offset.New(0)
	tr.Add(999)
	tr.Reset()
	if got := tr.Get(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestSnapshot_IndependentOfFutureAdds(t *testing.T) {
	tr := offset.New(0)
	tr.Add(10)
	snap := tr.Snapshot()
	tr.Add(50)
	if snap != 10 {
		t.Fatalf("snapshot should be 10, got %d", snap)
	}
	if tr.Get() != 60 {
		t.Fatalf("current should be 60, got %d", tr.Get())
	}
}

package tee_test

import (
	"testing"

	"greplive/internal/tee"
)

func TestNew_FansOutToAllConsumers(t *testing.T) {
	src := make(chan string, 4)
	tr := tee.New(src, 3, 8)
	outs := tr.Outputs()
	if len(outs) != 3 {
		t.Fatalf("expected 3 outputs, got %d", len(outs))
	}
	lines := []string{"alpha", "beta", "gamma"}
	for _, l := range lines {
		src <- l
	}
	close(src)
	tr.Wait()
	for i, ch := range outs {
		for _, want := range lines {
			got, ok := <-ch
			if !ok {
				t.Fatalf("output %d closed early", i)
			}
			if got != want {
				t.Errorf("output %d: got %q, want %q", i, got, want)
			}
		}
		if _, ok := <-ch; ok {
			t.Errorf("output %d: expected channel closed", i)
		}
	}
}

func TestNew_ZeroN_DefaultsToOne(t *testing.T) {
	src := make(chan string, 2)
	tr := tee.New(src, 0, 4)
	if len(tr.Outputs()) != 1 {
		t.Fatalf("expected 1 output for n=0")
	}
	close(src)
	tr.Wait()
}

func TestNew_EmptySource_ClosesOutputs(t *testing.T) {
	src := make(chan string)
	tr := tee.New(src, 2, 4)
	close(src)
	tr.Wait()
	for i, ch := range tr.Outputs() {
		if _, ok := <-ch; ok {
			t.Errorf("output %d should be closed", i)
		}
	}
}

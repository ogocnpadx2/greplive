package pause_test

import (
	"sync"
	"testing"
	"time"

	"greplive/internal/pause"
)

func TestNew_NotPaused(t *testing.T) {
	p := pause.New()
	if p.Paused() {
		t.Fatal("expected not paused on creation")
	}
}

func TestPause_SetsPaused(t *testing.T) {
	p := pause.New()
	p.Pause()
	if !p.Paused() {
		t.Fatal("expected paused after Pause()")
	}
}

func TestResume_ClearsPaused(t *testing.T) {
	p := pause.New()
	p.Pause()
	p.Resume()
	if p.Paused() {
		t.Fatal("expected not paused after Resume()")
	}
}

func TestToggle_FlipsState(t *testing.T) {
	p := pause.New()
	p.Toggle()
	if !p.Paused() {
		t.Fatal("expected paused after first Toggle()")
	}
	p.Toggle()
	if p.Paused() {
		t.Fatal("expected not paused after second Toggle()")
	}
}

func TestWait_BlocksWhilePaused(t *testing.T) {
	p := pause.New()
	p.Pause()

	var wg sync.WaitGroup
	wg.Add(1)
	proceed := make(chan struct{})

	go func() {
		close(proceed)
		p.Wait()
		wg.Done()
	}()

	<-proceed
	time.Sleep(20 * time.Millisecond)
	p.Resume()

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()

	select {
	case <-done:
		// ok
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Wait did not unblock after Resume")
	}
}

func TestWait_DoesNotBlockWhenRunning(t *testing.T) {
	p := pause.New()
	done := make(chan struct{})
	go func() {
		p.Wait()
		close(done)
	}()
	select {
	case <-done:
		// ok
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Wait blocked unexpectedly when not paused")
	}
}

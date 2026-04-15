package stats

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Reporter periodically prints a summary of session counters to a writer.
type Reporter struct {
	counters  *Counters
	interval  time.Duration
	out       io.Writer
	stopCh    chan struct{}
}

// NewReporter creates a Reporter that writes to w every interval.
// If w is nil it defaults to os.Stderr.
func NewReporter(c *Counters, interval time.Duration, w io.Writer) *Reporter {
	if w == nil {
		w = os.Stderr
	}
	return &Reporter{
		counters: c,
		interval: interval,
		out:      w,
		stopCh:   make(chan struct{}),
	}
}

// Start begins periodic reporting in a background goroutine.
func (r *Reporter) Start() {
	go func() {
		ticker := time.NewTicker(r.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				r.Print()
			case <-r.stopCh:
				return
			}
		}
	}()
}

// Stop halts periodic reporting.
func (r *Reporter) Stop() {
	close(r.stopCh)
}

// Print writes a single stats summary line to the configured writer.
func (r *Reporter) Print() {
	s := r.counters.Snapshot()
	fmt.Fprintf(
		r.out,
		"[stats] elapsed=%-8s read=%-6d matched=%-6d dropped=%-6d\n",
		s.Elapsed.Round(time.Millisecond),
		s.LinesRead,
		s.LinesMatched,
		s.LinesDropped,
	)
}

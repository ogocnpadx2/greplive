// Package alert emits a notification when a line matches a pattern
// and the match count crosses a configured threshold within a window.
package alert

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sync/atomic"
	"time"
)

// Alert watches a stream of lines and writes a message to Writer when
// the number of pattern matches within Window exceeds Threshold.
type Alert struct {
	re        *regexp.Regexp
	threshold int
	window    time.Duration
	writer    io.Writer
	buckets   []int64
	tick      time.Duration
	stop      chan struct{}
	count     atomic.Int64
	last      time.Time
}

// New creates an Alert. threshold=0 disables alerting.
func New(pattern string, threshold int, window time.Duration, w io.Writer) (*Alert, error) {
	if pattern == "" {
		return &Alert{}, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("alert: invalid pattern: %w", err)
	}
	if w == nil {
		w = os.Stderr
	}
	a := &Alert{
		re:        re,
		threshold: threshold,
		window:    window,
		writer:    w,
		stop:      make(chan struct{}),
		last:      time.Now(),
	}
	return a, nil
}

// Check tests line against the pattern and fires an alert if the threshold
// has been crossed within the configured window.
func (a *Alert) Check(line string) {
	if a.re == nil || a.threshold == 0 {
		return
	}
	if !a.re.MatchString(line) {
		return
	}
	n := a.count.Add(1)
	now := time.Now()
	if now.Sub(a.last) > a.window {
		a.count.Store(1)
		a.last = now
		return
	}
	if int(n) == a.threshold {
		fmt.Fprintf(a.writer, "[ALERT] pattern %q matched %d times within %s\n",
			a.re.String(), n, a.window)
	}
}

// Enabled reports whether alerting is active.
func (a *Alert) Enabled() bool {
	return a.re != nil && a.threshold > 0
}

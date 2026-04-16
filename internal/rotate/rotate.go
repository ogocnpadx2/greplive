// Package rotate detects log file rotation and signals the caller to reopen.
package rotate

import (
	"os"
	"time"
)

// Watcher polls a file path and detects rotation by comparing inode / size.
type Watcher struct {
	path     string
	interval time.Duration
	notify   chan struct{}
	stop     chan struct{}
}

// New creates a Watcher for path that polls every interval.
func New(path string, interval time.Duration) *Watcher {
	if interval <= 0 {
		interval = 2 * time.Second
	}
	return &Watcher{
		path:     path,
		interval: interval,
		notify:   make(chan struct{}, 1),
		stop:     make(chan struct{}),
	}
}

// Notify returns a channel that receives a value when rotation is detected.
func (w *Watcher) Notify() <-chan struct{} { return w.notify }

// Start begins polling in a background goroutine.
func (w *Watcher) Start() {
	go w.poll()
}

// Stop halts the background goroutine.
func (w *Watcher) Stop() {
	close(w.stop)
}

func (w *Watcher) poll() {
	ref, _ := stat(w.path)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-w.stop:
			return
		case <-ticker.C:
			cur, err := stat(w.path)
			if err != nil {
				continue
			}
			if cur != ref {
				ref = cur
				select {
				case w.notify <- struct{}{}:
				default:
				}
			}
		}
	}
}

type fileID struct {
	dev, ino uint64
	size     int64
}

func stat(path string) (fileID, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return fileID{}, err
	}
	return fileID{size: fi.Size(), ino: inode(fi)}, nil
}

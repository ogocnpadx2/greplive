// Package checkpoint tracks the last-read byte offset for a file so that
// greplive can resume tailing after a restart without replaying old lines.
package checkpoint

import (
	"encoding/json"
	"os"
	"sync"
)

// Checkpoint persists and retrieves a byte offset for a named file.
type Checkpoint struct {
	mu      sync.Mutex
	path    string
	offsets map[string]int64
}

type store struct {
	Offsets map[string]int64 `json:"offsets"`
}

// New loads an existing checkpoint file from path, or returns an empty
// Checkpoint if the file does not yet exist.
func New(path string) (*Checkpoint, error) {
	c := &Checkpoint{path: path, offsets: make(map[string]int64)}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return c, nil
	}
	if err != nil {
		return nil, err
	}
	var s store
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	if s.Offsets != nil {
		c.offsets = s.Offsets
	}
	return c, nil
}

// Get returns the saved offset for file, or 0 if none exists.
func (c *Checkpoint) Get(file string) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.offsets[file]
}

// Set updates the in-memory offset for file and flushes to disk.
func (c *Checkpoint) Set(file string, offset int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.offsets[file] = offset
	return c.flush()
}

// flush writes the current offsets to disk. Caller must hold mu.
func (c *Checkpoint) flush() error {
	data, err := json.Marshal(store{Offsets: c.offsets})
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0o644)
}

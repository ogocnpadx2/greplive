package output

import (
	"bufio"
	"io"

	"greplive/internal/filter"
)

// Pipeline reads lines from src, applies cfg to filter them, and writes
// matching lines via w. It blocks until src reaches EOF or returns an error.
// The first non-EOF read error is returned to the caller.
func Pipeline(src io.Reader, cfg *filter.Config, w *Writer) error {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()
		if cfg.Match(line) {
			w.WriteLine(line)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

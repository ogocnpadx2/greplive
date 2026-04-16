// Package checkpoint provides persistent byte-offset tracking for tailed files.
//
// When greplive tails a file it can optionally record the last successfully
// processed byte offset to a small JSON file on disk.  On the next invocation
// the offset is restored so that only new content is processed, avoiding
// duplicate output across restarts.
//
// Usage:
//
//	cp, err := checkpoint.New("/tmp/greplive.checkpoint")
//	if err != nil { ... }
//	offset := cp.Get("/var/log/app.log")   // 0 on first run
//	// ... seek file to offset, tail from there ...
//	_ = cp.Set("/var/log/app.log", newOffset)
package checkpoint

// Package offset provides a thread-safe byte-offset tracker for log stream
// consumers.
//
// It is intended to be used alongside the checkpoint package: after processing
// each line the caller adds the line's byte length to the Tracker, then
// periodically persists the current offset via checkpoint.Set so that the
// stream can be resumed from the correct position after a restart.
//
// Example:
//
//	tr := offset.New(savedOffset)
//	for line := range lines {
//		process(line)
//		tr.Add(len(line) + 1) // +1 for newline
//	}
//	checkpoint.Set(path, tr.Get())
package offset

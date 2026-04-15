// Package buffer implements a fixed-capacity ring buffer for recent log lines.
//
// The Ring type stores up to N lines in insertion order, automatically
// discarding the oldest entry when the buffer is full. It is safe for
// concurrent use and is intended to support context-window features such as
// printing the last N lines before a regex match.
//
// Basic usage:
//
//	b := buffer.New(5)
//	b.Push("line one")
//	b.Push("line two")
//	for _, l := range b.Snapshot() {
//		fmt.Println(l)
//	}
package buffer

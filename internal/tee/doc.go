// Package tee provides a fan-out multiplexer for log line channels.
//
// It reads from a single source channel and duplicates every line to
// N independent output channels, enabling multiple downstream consumers
// (e.g. a filter pipeline and a stats reporter) to process the same
// stream without coordination.
//
// Usage:
//
//	src := make(chan string, 64)
//	t := tee.New(src, 2, 64)
//	outs := t.Outputs() // outs[0] and outs[1] each receive every line
package tee

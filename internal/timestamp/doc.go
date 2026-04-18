// Package timestamp provides a lightweight Stamper type that prepends a
// formatted wall-clock timestamp to log lines as they flow through the
// greplive pipeline.
//
// Usage:
//
//	s := timestamp.New(time.RFC3339)
//	outLine := s.Apply(inLine)
//
// When an empty format string is supplied the Stamper is disabled and
// Apply returns the line unchanged, making it safe to wire unconditionally
// into a processing chain.
package timestamp

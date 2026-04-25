// Package strip provides a Stripper that removes ANSI escape sequences and
// raw control characters from log lines.
//
// # Usage
//
// Build a Stripper from a Config:
//
//	s := strip.DefaultConfig().Build()
//	clean := s.Apply(rawLine)
//
// Or construct one directly:
//
//	s := strip.New(true, false)  // strip ANSI, keep control chars
//
// # Pipeline integration
//
// ApplyAll runs a slice of Strippers in order, making it easy to compose
// stripping with other line transforms inside an output.Pipeline.
package strip

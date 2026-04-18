// Package squash provides a Squasher that merges consecutive log lines
// matching a regex pattern into a single pipe-delimited output line.
//
// This is useful for collapsing repetitive debug or trace lines that share
// a common prefix (e.g. "DEBUG ...") into one readable summary line.
//
// Usage:
//
//	s, err := squash.New(`^DEBUG`)
//	if err != nil { ... }
//	for _, line := range lines {
//		if out, ok := s.Push(line); ok {
//			fmt.Println(out)
//		}
//	}
//	if out, ok := s.Flush(); ok {
//		fmt.Println(out)
//	}
package squash

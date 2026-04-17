// Package contextline provides grep-style context line capture around matches.
//
// A Collector buffers up to N lines before a match (pre-context) and emits up
// to M lines after a match (post-context), mirroring the behaviour of
// grep -B N -A M.
//
// Usage:
//
//	c := contextline.New(2, 3) // 2 before, 3 after
//	for _, line := range logLines {
//		matched := filter.Match(line)
//		for _, out := range c.Feed(line, matched) {
//			fmt.Println(out)
//		}
//	}
package contextline

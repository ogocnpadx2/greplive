// Package inverse implements invert-match filtering for log lines.
//
// An Inverter compiles a regular expression and passes through only
// the lines that do NOT match — the logical complement of a normal
// grep filter. This is equivalent to `grep -v` on the command line.
//
// Usage:
//
//	inv, err := inverse.New(`DEBUG`)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if inv.Allow(line) {
//		// line does not contain DEBUG — forward it
//	}
//
// Multiple Inverters can be composed with ApplyAll, which drops any
// line matched by at least one pattern in the slice.
package inverse

// Package throttle limits the number of log lines emitted per second.
//
// A Throttle counts lines within a rolling time window and drops lines
// once the configured maximum has been reached for that window.  When
// the maximum is zero the throttle is disabled and every line is
// allowed through.
//
// Usage:
//
//	t := throttle.New(100, time.Second)
//	for _, line := range lines {
//		if t.Allow(line) {
//			fmt.Println(line)
//		}
//	}
package throttle

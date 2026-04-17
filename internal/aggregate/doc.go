// Package aggregate provides run-length aggregation of consecutive log lines
// that match a given regular expression.
//
// When many identical (or pattern-matching) lines appear in a burst, the
// Aggregator suppresses individual lines and instead emits a single summary
// line of the form:
//
//	<first matched line> [xN]
//
// where N is the number of consecutive matching lines seen.
//
// Usage:
//
//	cfg := aggregate.Config{Pattern: `ERROR`}
//	a, err := cfg.Build()
//	// ...
//	for _, line := range lines {
//		if summary, ok := a.Push(line); ok {
//			fmt.Println(summary)
//		}
//	}
//	if summary, ok := a.Flush(); ok {
//		fmt.Println(summary)
//	}
package aggregate

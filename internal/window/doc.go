// Package window implements a sliding time-window counter used by greplive
// to measure line throughput in real-time.
//
// A Window records event timestamps and automatically evicts entries that fall
// outside the configured rolling duration. It is safe for concurrent use.
//
// Typical usage:
//
//	w := window.New(5 * time.Second)
//
//	// Record a matched line
//	w.Add()
//
//	// Query current throughput
//	fmt.Printf("lines/sec: %.2f\n", w.Rate())
package window

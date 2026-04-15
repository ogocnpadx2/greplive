// Package dedupe implements consecutive-line deduplication for greplive.
//
// When log streams produce bursts of identical lines (e.g. a tight retry loop
// or a repeating health-check message) the terminal quickly becomes unusable.
// Deduper solves this by suppressing repeated lines and, optionally, emitting
// a human-readable summary such as:
//
//	[repeated 42 times]
//
// Usage:
//
//	d := dedupe.New(true)   // true = emit summary messages
//	for _, line := range lines {
//	    flush, allow := d.Check(line)
//	    if flush != "" {
//	        fmt.Println(flush)
//	    }
//	    if allow {
//	        fmt.Println(line)
//	    }
//	}
//	if s := d.Flush(); s != "" {
//	    fmt.Println(s) // trailing run
//	}
package dedupe

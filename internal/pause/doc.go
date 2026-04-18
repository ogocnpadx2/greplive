// Package pause provides a concurrency-safe toggle that can suspend
// processing in a pipeline stage without closing channels or losing
// data from the upstream source.
//
// Typical usage:
//
//	p := pause.New()
//
//	// In the processing goroutine, call Wait before handling each line:
//	for line := range lines {
//		p.Wait() // blocks if paused
//		process(line)
//	}
//
//	// From a key-press handler:
//	p.Toggle()
package pause

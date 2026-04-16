// Package multiline collapses multi-line log entries — such as Java stack
// traces or Go panic output — into a single logical line for downstream
// filtering and display.
//
// Usage:
//
//	cfg := multiline.DefaultConfig()
//	collector, err := cfg.Build()
//	if err != nil { ... }
//
//	for _, raw := range lines {
//		if event, ok := collector.Push(raw); ok {
//			// process event
//		}
//	}
//	if tail := collector.Flush(); tail != "" {
//		// process final event
//	}
package multiline

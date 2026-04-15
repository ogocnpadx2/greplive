// Package stats provides lightweight runtime counters and a periodic reporter
// for greplive streaming sessions.
//
// Usage:
//
//	// Create a counter set at the start of a session.
//	 c := stats.New()
//
//	// Increment counters as lines flow through the pipeline.
//	 c.IncrRead()
//	 c.IncrMatched()
//	 c.IncrSeverity("ERROR")
//
//	// Optionally attach a periodic reporter that prints to stderr.
//	 r := stats.NewReporter(c, 5*time.Second, nil)
//	 r.Start()
//	 defer r.Stop()
//
//	// Take an immutable snapshot at any point.
//	 snap := c.Snapshot()
//	 fmt.Println(snap.LinesRead, snap.LinesMatched)
package stats

// Package ratelimit implements a token-bucket rate limiter for greplive.
//
// It is used to cap the number of log lines processed per second, which
// prevents downstream consumers (writers, matchers) from being overwhelmed
// when tailing high-volume log sources.
//
// Usage:
//
//	limiter := ratelimit.New(500) // allow 500 lines/sec
//	defer limiter.Stop()
//
//	for line := range lines {
//		if !limiter.Wait(ctx) {
//			break // context cancelled
//		}
//		process(line)
//	}
//
// A rate of zero or less disables limiting entirely, letting all lines
// pass through without delay.
package ratelimit

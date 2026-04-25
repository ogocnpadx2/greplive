// Package epoch replaces raw Unix epoch timestamps (seconds or
// milliseconds) found anywhere in a log line with human-readable time
// strings.
//
// Both 10-digit (second) and 13-digit (millisecond) epochs are
// recognised automatically.  The output format and timezone (UTC vs
// local) are configurable.
//
// Usage:
//
//	c := epoch.New(time.RFC3339, true /* UTC */)
//	fmt.Println(c.Apply("ts=1700000000 msg=started"))
//	// ts=2023-11-14T22:13:20Z msg=started
package epoch

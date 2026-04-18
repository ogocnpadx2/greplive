// Package alert provides threshold-based alerting for greplive.
//
// An Alert watches each log line against a compiled regex pattern and
// counts how many times it matches within a rolling time window. When
// the count reaches the configured threshold an alert message is written
// to the configured writer (defaulting to os.Stderr).
//
// Example usage:
//
//	a, err := alert.New(`ERROR`, 5, 30*time.Second, os.Stderr)
//	if err != nil { /* handle */ }
//	for _, line := range lines {
//		a.Check(line)
//	}
package alert

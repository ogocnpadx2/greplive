// Package output provides a Writer that formats and prints filtered log lines
// to a terminal. Each line is colorized according to its detected severity
// level (using the severity package) and can optionally be prefixed with a
// wall-clock timestamp and a severity label.
//
// Basic usage:
//
//	w := output.New(os.Stdout,
//		output.WithTimestamp(true),
//		output.WithLevel(true),
//	)
//	w.WriteLine(line)
//
// The Writer is safe for sequential use from a single goroutine. Callers that
// write from multiple goroutines must synchronize externally.
package output

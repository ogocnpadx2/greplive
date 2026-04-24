// Package sink provides output destinations for processed log lines.
//
// A Sink wraps an io.Writer and exposes a simple Write(line string) API that
// appends a newline after each entry. Sinks are safe for concurrent use.
//
// # Basic usage
//
//	s := sink.New(os.Stdout)
//	s.Write("processed log line")
//
// # File sink
//
//	s, close, err := sink.NewFile("/var/log/greplive.log")
//	if err != nil { ... }
//	defer close()
//
// # Fan-out
//
// Multi fans a single Write call out to several sinks simultaneously,
// which is useful for mirroring output to both a file and stdout:
//
//	m := sink.Multi(sink.New(os.Stdout), fileSink)
//	m.Write(line)
package sink

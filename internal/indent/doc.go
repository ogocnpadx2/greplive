// Package indent prepends a configurable prefix string to each log line,
// allowing structured or hierarchical display of log output in the terminal.
//
// Usage:
//
//	i := indent.New("  ")   // two-space indent
//	formatted := i.Apply(line)
//
// Multiple indenters can be chained via ApplyAll. Use Repeat to build
// a prefix from a repeated unit such as spaces or tabs.
package indent

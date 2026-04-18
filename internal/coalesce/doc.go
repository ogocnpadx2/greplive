// Package coalesce provides burst-coalescing for repeated log lines.
//
// When the same line appears multiple times in quick succession, the
// Coalescer suppresses duplicates and emits a single annotated line
// of the form:
//
//	<original line>  [xN]
//
// where N is the repeat count. Lines that appear only once are passed
// through unchanged.
//
// A configurable quiet window controls how long the Coalescer waits
// after the last duplicate before emitting. Calling Flush forces
// immediate emission of any buffered line, which is useful at shutdown.
package coalesce

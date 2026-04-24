// Package abbrev shortens log lines by collapsing runs of whitespace and
// optionally truncating lines that exceed a maximum rune width.
//
// # Usage
//
//	abbreviator := abbrev.New(120, "…", true)
//	if abbreviator.Enabled() {
//		line = abbreviator.Apply(line)
//	}
//
// Multiple abbreviators can be chained with [ApplyAll].
package abbrev

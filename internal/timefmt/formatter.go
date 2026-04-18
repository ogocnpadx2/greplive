package timefmt

import "time"

// Formatter applies a Format to a time.Time, optionally converting to UTC.
type Formatter struct {
	fmt Format
	utc bool
}

// New returns a Formatter using the named format.
// Falls back to RFC3339 for unknown names.
func New(name string, utc bool) Formatter {
	f, _ := Parse(name)
	return Formatter{fmt: f, utc: utc}
}

// FormatTime formats t according to the configured layout.
func (f Formatter) FormatTime(t time.Time) string {
	if f.utc {
		t = t.UTC()
	}
	return f.fmt.Format(t)
}

// Now formats the current time.
func (f Formatter) Now() string {
	return f.FormatTime(time.Now())
}

// FormatName returns the underlying format name.
func (f Formatter) FormatName() string {
	return f.fmt.Name()
}

// IsUTC reports whether the Formatter converts times to UTC before formatting.
func (f Formatter) IsUTC() bool {
	return f.utc
}

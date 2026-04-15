package severity

import (
	"regexp"
	"strings"
)

// Level represents a log severity level.
type Level int

const (
	Unknown Level = iota
	Debug
	Info
	Warn
	Error
	Fatal
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorGray   = "\033[90m"
	ColorCyan   = "\033[36m"
	ColorYellow = "\033[33m"
	ColorRed    = "\033[31m"
	ColorMagenta = "\033[35m"
)

var levelPatterns = []struct {
	level   Level
	pattern *regexp.Regexp
}{
	{Fatal, regexp.MustCompile(`(?i)\b(fatal|panic|critical)\b`)},
	{Error, regexp.MustCompile(`(?i)\b(error|err|exception|fail(ed)?)\b`)},
	{Warn, regexp.MustCompile(`(?i)\b(warn(ing)?|caution)\b`)},
	{Info, regexp.MustCompile(`(?i)\b(info(rmation)?|notice)\b`)},
	{Debug, regexp.MustCompile(`(?i)\b(debug|trace|verbose)\b`)},
}

// Detect returns the severity level detected in the given log line.
func Detect(line string) Level {
	for _, lp := range levelPatterns {
		if lp.pattern.MatchString(line) {
			return lp.level
		}
	}
	return Unknown
}

// Colorize wraps the line with the appropriate ANSI color for its severity.
func Colorize(line string, level Level) string {
	color := colorFor(level)
	if color == "" {
		return line
	}
	return color + line + ColorReset
}

// String returns a human-readable name for the level.
func (l Level) String() string {
	switch l {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	case Fatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// ParseLevel converts a string to a Level, case-insensitively.
func ParseLevel(s string) Level {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return Debug
	case "INFO":
		return Info
	case "WARN", "WARNING":
		return Warn
	case "ERROR", "ERR":
		return Error
	case "FATAL", "PANIC", "CRITICAL":
		return Fatal
	default:
		return Unknown
	}
}

func colorFor(level Level) string {
	switch level {
	case Debug:
		return ColorGray
	case Info:
		return ColorCyan
	case Warn:
		return ColorYellow
	case Error:
		return ColorRed
	case Fatal:
		return ColorMagenta
	default:
		return ""
	}
}

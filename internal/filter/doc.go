// Package filter provides log line filtering based on regex patterns
// and minimum severity levels.
//
// A Config is created with a regex pattern string and a minimum severity.Level.
// Lines are matched against both criteria: the regex pattern (if provided) must
// match, and the detected severity of the line must be >= the configured minimum.
//
// Example usage:
//
//	cfg, err := filter.NewConfig("timeout", severity.LevelWarn, false)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	scanner := bufio.NewScanner(reader)
//	for scanner.Scan() {
//		line := scanner.Text()
//		if cfg.Match(line) {
//			fmt.Println(cfg.Highlight(line))
//		}
//	}
package filter

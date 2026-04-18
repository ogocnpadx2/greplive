package labelmap

import "fmt"

// Config holds raw label definitions for building a Labeler.
type Config struct {
	// Labels is a slice of "key=value" strings.
	Labels []string
}

// Build parses the Config and returns a ready-to-use Labeler.
// Each entry in Labels must be in "key=value" format.
func (c Config) Build() (*Labeler, error) {
	if len(c.Labels) == 0 {
		return New(nil), nil
	}
	m := make(map[string]string, len(c.Labels))
	for _, entry := range c.Labels {
		for i := 0; i < len(entry); i++ {
			if entry[i] == '=' {
				k, v := entry[:i], entry[i+1:]
				if k == "" {
					return nil, fmt.Errorf("labelmap: empty key in %q", entry)
				}
				m[k] = v
				goto next
			}
		}
		return nil, fmt.Errorf("labelmap: missing '=' in label %q", entry)
	next:
	}
	return New(m), nil
}

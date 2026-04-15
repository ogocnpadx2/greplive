package transform

import "fmt"

// Step describes a single transformation step in a pipeline config.
type Step struct {
	// Kind is one of "trim", "replace", or "regex".
	Kind string

	// Old is the string to replace (used by "replace" kind).
	Old string

	// New is the replacement string (used by "replace" and "regex" kinds).
	New string

	// Pattern is the regex pattern (used by "regex" kind).
	Pattern string
}

// Build converts a slice of Steps into a slice of Transformers.
// Returns an error if any step is misconfigured.
func Build(steps []Step) ([]Transformer, error) {
	out := make([]Transformer, 0, len(steps))
	for i, s := range steps {
		switch s.Kind {
		case "trim":
			out = append(out, TrimTransformer{})
		case "replace":
			out = append(out, ReplaceTransformer{Old: s.Old, New: s.New})
		case "regex":
			rt, err := NewRegex(s.Pattern, s.New)
			if err != nil {
				return nil, fmt.Errorf("transform step %d: %w", i, err)
			}
			out = append(out, rt)
		default:
			return nil, fmt.Errorf("transform step %d: unknown kind %q", i, s.Kind)
		}
	}
	return out, nil
}

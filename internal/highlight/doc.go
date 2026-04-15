// Package highlight provides regex-based term highlighting for log output.
//
// It allows one or more patterns to be matched within a log line and wrapped
// with ANSI color escape codes for terminal display. Multiple Highlighter
// instances can be composed and applied in sequence via ApplyAll.
//
// Example usage:
//
//	h, err := highlight.New(`error`, highlight.Red)
//	if err != nil {
//		log.Fatal(err)
//	}
//	colored := h.Apply("an error occurred")
//	fmt.Println(colored)
package highlight

// Package bracket provides a line transformer that wraps regex-matched
// substrings in configurable left and right bracket strings.
//
// It is useful for visually calling out tokens such as identifiers, error
// codes, or numeric values without changing the surrounding text.
//
// Example usage:
//
//	b, err := bracket.New(`\d+`, "[", "]")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(b.Apply("connected on port 8080"))
//	// Output: connected on port [8080]
package bracket

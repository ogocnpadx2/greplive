// Package suppress provides line suppression based on regular-expression
// matching. It is the logical complement of the highlight and grep packages:
// where those packages select lines to keep, suppress discards lines whose
// content matches one or more patterns.
//
// Basic usage:
//
//	s, err := suppress.New(`(?i)debug`)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, line := range inputLines {
//		if !s.Drop(line) {
//			fmt.Println(line)
//		}
//	}
//
// Multiple suppressors can be applied in one pass with ApplyAll.
package suppress

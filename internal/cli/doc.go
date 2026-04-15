// Package cli provides the command-line interface for greplive.
//
// It is responsible for:
//   - Parsing command-line flags via ParseFlags.
//   - Wiring together the input, filter, highlight, output and stats packages.
//   - Exposing a single Run function that the main package calls.
//
// Typical usage:
//
//	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
//	defer cancel()
//	if err := cli.Run(ctx, os.Args[1:]); err != nil {
//		log.Fatal(err)
//	}
package cli

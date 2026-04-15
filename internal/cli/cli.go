// Package cli wires together all greplive components and exposes a Run
// function that is the single entry-point for the binary.
package cli

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/user/greplive/internal/filter"
	"github.com/user/greplive/internal/highlight"
	"github.com/user/greplive/internal/input"
	"github.com/user/greplive/internal/output"
	"github.com/user/greplive/internal/severity"
	"github.com/user/greplive/internal/stats"
)

// Run parses args, builds the pipeline and blocks until ctx is cancelled or
// all input has been consumed.
func Run(ctx context.Context, args []string) error {
	cfg, err := ParseFlags(args)
	if err != nil {
		return err
	}

	filterCfg, err := filter.NewConfig(cfg.Pattern, cfg.MinLevel)
	if err != nil {
		return fmt.Errorf("invalid filter: %w", err)
	}

	highlighters, err := highlight.DefaultConfig().Build()
	if err != nil {
		return fmt.Errorf("highlight config: %w", err)
	}

	st := stats.New()

	var r io.ReadCloser
	if cfg.File != "" {
		lines, err := input.TailFile(ctx, cfg.File)
		if err != nil {
			return fmt.Errorf("tail %s: %w", cfg.File, err)
		}
		return runLines(ctx, lines, filterCfg, highlighters, st, cfg)
	}

	r = os.Stdin
	lines := input.New(r).Lines(ctx)
	return runLines(ctx, lines, filterCfg, highlighters, st, cfg)
}

func runLines(
	ctx context.Context,
	lines <-chan string,
	filterCfg *filter.Config,
	highlighters []highlight.Highlighter,
	st *stats.Stats,
	cfg *Config,
) error {
	writerOpts := []output.Option{}
	if cfg.ShowTimestamp {
		writerOpts = append(writerOpts, output.WithTimestamp())
	}
	if cfg.ShowLevel {
		writerOpts = append(writerOpts, output.WithLevel())
	}
	w := output.New(os.Stdout, writerOpts...)

	var reporter *stats.Reporter
	if cfg.ShowStats {
		reporter = stats.NewReporter(st, os.Stderr)
		reporter.Start(ctx, cfg.StatsInterval)
	}

	for {
		select {
		case <-ctx.Done():
			if reporter != nil {
				reporter.Stop()
				reporter.Print()
			}
			return nil
		case line, ok := <-lines:
			if !ok {
				if reporter != nil {
					reporter.Stop()
					reporter.Print()
				}
				return nil
			}
			st.IncrRead()
			lvl := severity.Detect(line)
			if !filterCfg.Match(line, lvl) {
				st.IncrDropped()
				continue
			}
			st.IncrMatched()
			st.IncrSeverity(lvl)
			coloured := highlight.ApplyAll(line, highlighters)
			w.WriteLine(coloured, lvl)
		}
	}
}

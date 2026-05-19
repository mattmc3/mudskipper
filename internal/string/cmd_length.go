package stringcmd

import (
	"flag"
	"fmt"
	"io"

	"rsc.io/getopt"
)

const usageLength = `Usage: string length [-h] [-q] [-V] [STRING ...]

  Print the length of each STRING.

Options:
  -V, --visible    Count visible width (strip ANSI escape sequences)
  -q, --quiet      Suppress output; exit 0 if any non-empty, 1 if all empty
  -h, --help       Show this help message
`

func runLength(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := getopt.NewFlagSet("length", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	quiet := fs.Bool("quiet", false, "")
	visible := fs.Bool("visible", false, "")
	fs.Aliases("h", "help", "q", "quiet", "V", "visible")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: length: %v\n", err)
		return 1
	}
	if *help {
		fmt.Fprint(stdout, usageLength)
		return 0
	}

	any := false
	next := newStringIter(fs.Args(), stdin)
	for s, ok := next(); ok; s, ok = next() {
		if *visible {
			for _, w := range visualWidthOfLines(s) {
				if w > 0 {
					any = true
				}
				if !*quiet {
					fmt.Fprintln(stdout, w)
				}
			}
		} else {
			n := len(s)
			if n > 0 {
				any = true
			}
			if !*quiet {
				fmt.Fprintln(stdout, n)
			}
		}
		if *quiet && any {
			return 0
		}
	}
	if any {
		return 0
	}
	return 1
}

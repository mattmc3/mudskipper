package stringcmd

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"rsc.io/getopt"
)

const usageUpper = `Usage: string upper [-h] [-q] [STRING ...]

  Convert STRING to uppercase.

Options:
  -q, --quiet    Suppress output; exit 0 if any string changed, 1 if none
  -h, --help     Show this help message
`

const usageLower = `Usage: string lower [-h] [-q] [STRING ...]

  Convert STRING to lowercase.

Options:
  -q, --quiet    Suppress output; exit 0 if any string changed, 1 if none
  -h, --help     Show this help message
`

func runUpper(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	return runTransform("upper", args, stdin, stdout, stderr, strings.ToUpper, usageUpper)
}

func runLower(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	return runTransform("lower", args, stdin, stdout, stderr, strings.ToLower, usageLower)
}

func runTransform(name string, args []string, stdin io.Reader, stdout, stderr io.Writer, transform func(string) string, usage string) int {
	fs := getopt.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	quiet := fs.Bool("quiet", false, "")
	fs.Aliases("h", "help", "q", "quiet")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: %s: %v\n", name, err)
		return 1
	}
	if *help {
		fmt.Fprint(stdout, usage)
		return 0
	}

	changed := false
	next := newStringIter(fs.Args(), stdin)
	for s, ok := next(); ok; s, ok = next() {
		result := transform(s)
		if !*quiet {
			fmt.Fprintln(stdout, result)
		}
		if result != s {
			changed = true
		}
		if *quiet && changed {
			return 0
		}
	}
	if changed {
		return 0
	}
	return 1
}

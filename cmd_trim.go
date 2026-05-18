package main

import (
	"flag"
	"fmt"
	"io"
	"strings"
	"unicode"

	"rsc.io/getopt"
)

func runTrim(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := getopt.NewFlagSet("trim", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	left := fs.Bool("left", false, "")
	right := fs.Bool("right", false, "")
	quiet := fs.Bool("quiet", false, "")
	chars := fs.String("chars", "", "")
	fs.Aliases("h", "help", "l", "left", "r", "right", "q", "quiet", "c", "chars")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: trim: %v\n", err)
		return 1
	}
	if *help {
		fmt.Fprint(stdout, `Usage: string trim [-h] [-l] [-r] [-q] [-c CHARS] [STRING ...]

  Remove leading and trailing whitespace from STRING.

Options:
  -l, --left          Trim leading whitespace only
  -r, --right         Trim trailing whitespace only
  -c, --chars CHARS   Trim CHARS instead of whitespace
  -q, --quiet         Suppress output; exit 0 if any string trimmed, 1 if none
  -h, --help          Show this help message
`)
		return 0
	}

	if !*left && !*right {
		*left = true
		*right = true
	}

	changed := false
	next := newStringIter(fs.Args(), stdin)
	for s, ok := next(); ok; s, ok = next() {
		var result string
		if *chars == "" {
			result = trimWhitespace(s, *left, *right)
		} else {
			result = trimChars(s, *left, *right, *chars)
		}
		if !*quiet {
			fmt.Fprintln(stdout, result)
		}
		if len(result) < len(s) {
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

func trimWhitespace(s string, left, right bool) string {
	switch {
	case left && right:
		return strings.TrimFunc(s, unicode.IsSpace)
	case left:
		return strings.TrimLeftFunc(s, unicode.IsSpace)
	case right:
		return strings.TrimRightFunc(s, unicode.IsSpace)
	default:
		return s
	}
}

func trimChars(s string, left, right bool, chars string) string {
	switch {
	case left && right:
		return strings.Trim(s, chars)
	case left:
		return strings.TrimLeft(s, chars)
	case right:
		return strings.TrimRight(s, chars)
	default:
		return s
	}
}

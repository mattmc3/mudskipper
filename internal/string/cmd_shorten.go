package stringcmd

import (
	"flag"
	"fmt"
	"io"
	"strconv"

	"rsc.io/getopt"
)

const usageShorten = `Usage: string shorten [-h] [-l] [-N] [-q] [(-c | --char) CHARS] [(-m | --max) INTEGER] [STRING ...]

  Shorten strings to a maximum width, appending an ellipsis if truncated.

Options:
  -l, --left          Shorten from the left (ellipsis at start)
  -N, --no-newline    Omit newline after last output
  -c, --char CHARS    Ellipsis string (default: …)
  -m, --max INT       Maximum width (default: no limit)
  -q, --quiet         Suppress output; exit 0 if any shortened, 1 if none
  -h, --help          Show this help message
`

func runShorten(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := getopt.NewFlagSet("shorten", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	quiet := fs.Bool("quiet", false, "")
	noNewline := fs.Bool("no-newline", false, "")
	left := fs.Bool("left", false, "")
	charStr := fs.String("char", "…", "")
	maxStr := fs.String("max", "", "")
	fs.Aliases("h", "help", "q", "quiet", "N", "no-newline", "l", "left", "c", "char", "m", "max")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: shorten: %v\n", err)
		return 1
	}
	if *help {
		fmt.Fprint(stdout, usageShorten)
		return 0
	}

	ellipsis := *charStr

	maxWidth := -1
	maxSet := false
	if *maxStr != "" {
		n, err := strconv.Atoi(*maxStr)
		if err != nil || n < 0 {
			fmt.Fprintf(stderr, "error: shorten: Invalid max value '%s'\n", *maxStr)
			return 1
		}
		if n > 0 {
			maxWidth = n
		}
		maxSet = true
	}

	strs := inputStrings(fs.Args(), stdin)

	if !maxSet {
		autoMax := -1
		for _, s := range strs {
			w := visualWidthOf(s)
			if w > 0 && (autoMax < 0 || w < autoMax) {
				autoMax = w
			}
		}
		if autoMax >= 0 {
			maxWidth = autoMax
		}
	}

	changed := false
	for i, s := range strs {
		result := shortenStr(s, maxWidth, ellipsis, *left)
		if result != s {
			changed = true
		}
		if !*quiet {
			isLast := i == len(strs)-1
			if *noNewline && isLast {
				fmt.Fprint(stdout, result)
			} else {
				fmt.Fprintln(stdout, result)
			}
		}
	}

	if *quiet {
		if changed {
			return 0
		}
		return 1
	}
	if len(strs) > 0 {
		return 0
	}
	return 1
}

func shortenStr(s string, maxWidth int, ellipsis string, left bool) string {
	if maxWidth < 0 || visualWidthOf(s) <= maxWidth {
		return s
	}
	ellipsisWidth := visualWidthOf(ellipsis)
	if ellipsisWidth > maxWidth {
		if left {
			return visualTakeRight(s, maxWidth)
		}
		return visualTakeLeft(s, maxWidth)
	}
	keep := maxWidth - ellipsisWidth
	if left {
		return ellipsis + visualTakeRight(s, keep)
	}
	return visualTakeLeft(s, keep) + ellipsis
}

package main

import (
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"

	"rsc.io/getopt"
)

func runRepeat(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := getopt.NewFlagSet("repeat", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	countStr := fs.String("count", "", "")
	maxStr := fs.String("max", "", "")
	noNewline := fs.Bool("no-newline", false, "")
	quiet := fs.Bool("quiet", false, "")
	fs.Aliases("h", "help", "n", "count", "m", "max", "N", "no-newline", "q", "quiet")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: repeat: %v\n", err)
		return 1
	}
	if *help {
		fmt.Fprintln(stdout, "Usage: string repeat [-h] -n COUNT [-m MAX] [-N] [-q] [STRING ...]")
		fmt.Fprintln(stdout, "")
		fmt.Fprintln(stdout, "  Repeat STRING COUNT times.")
		fmt.Fprintln(stdout, "")
		fmt.Fprintln(stdout, "Options:")
		fmt.Fprintln(stdout, "  -n, --count COUNT   Number of times to repeat (required)")
		fmt.Fprintln(stdout, "  -m, --max MAX       Maximum length of result")
		fmt.Fprintln(stdout, "  -N, --no-newline    Omit newline after last output")
		fmt.Fprintln(stdout, "  -q, --quiet         Suppress output; exit 0 if any output, 1 if none")
		fmt.Fprintln(stdout, "  -h, --help          Show this help message")
		return 0
	}

	count := 0
	if *countStr != "" {
		n, err := strconv.Atoi(*countStr)
		if err != nil || n < 0 {
			fmt.Fprintf(stderr, "error: repeat: Invalid count value '%s'\n", *countStr)
			return 1
		}
		count = n
	}

	maxLen := -1
	if *maxStr != "" {
		n, err := strconv.Atoi(*maxStr)
		if err != nil || n < 0 {
			fmt.Fprintf(stderr, "error: repeat: Invalid max value '%s'\n", *maxStr)
			return 1
		}
		maxLen = n
	}

	inputs := fs.Args()

	// positional count: if -n not given, try first positional arg as count
	if count == 0 && len(inputs) > 0 {
		n, err := strconv.Atoi(inputs[0])
		if err == nil {
			if n < 0 {
				fmt.Fprintf(stderr, "error: repeat: Invalid count value '%s'\n", inputs[0])
				return 1
			}
			count = n
			inputs = inputs[1:]
		} else if maxLen < 0 {
			fmt.Fprintf(stderr, "error: repeat: Invalid count value '%s'\n", inputs[0])
			return 1
		}
	}

	maxOnly := count == 0 && maxLen >= 0
	if count == 0 && !maxOnly {
		return 1
	}

	strs := inputStrings(inputs, stdin)
	if len(strs) == 0 {
		return 1
	}

	changed := false
	for i, s := range strs {
		effectiveCount := count
		if maxOnly {
			if len(s) > 0 {
				effectiveCount = maxLen/len(s) + 1
			} else {
				effectiveCount = 0
			}
		}
		repeated := strings.Repeat(s, effectiveCount)
		if maxLen >= 0 && len(repeated) > maxLen {
			repeated = repeated[:maxLen]
		}
		if len(repeated) > 0 {
			changed = true
		}
		if !*quiet && (len(repeated) > 0 || len(strs) > 1) {
			isLast := i == len(strs)-1
			if *noNewline && isLast {
				fmt.Fprint(stdout, repeated)
			} else {
				fmt.Fprintln(stdout, repeated)
			}
		}
	}

	if changed {
		return 0
	}
	return 1
}

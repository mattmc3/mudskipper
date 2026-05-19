package stringcmd

import (
	"flag"
	"fmt"
	"io"
	"math"
	"strconv"

	"rsc.io/getopt"
)

const usageSub = `Usage: string sub [-h] [(-s | --start) START] [(-e | --end) END] [(-l | --length) LENGTH] [-q] [STRING ...]

  Extract substrings from STRING.

Options:
  -s, --start INT      Start position (1-based; negative counts from end)
  -e, --end INT        End position, inclusive (1-based; negative counts from end)
  -l, --length INT     Length of substring (overrides --end)
  -q, --quiet          Suppress output; exit 0 if found result non-empty, 1 if all empty
  -h, --help           Show this help message
`

func runSub(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := getopt.NewFlagSet("sub", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	quiet := fs.Bool("quiet", false, "")
	startStr := fs.String("start", "", "")
	endStr := fs.String("end", "", "")
	lengthStr := fs.String("length", "", "")
	fs.Aliases("h", "help", "q", "quiet", "s", "start", "e", "end", "l", "length")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: sub: %v\n", err)
		return 1
	}
	if *help {
		fmt.Fprint(stdout, usageSub)
		return 0
	}

	var start, end, length *int

	if *startStr != "" {
		n, err := strconv.ParseInt(*startStr, 10, 64)
		if err != nil {
			fmt.Fprintf(stderr, "error: sub: Invalid start value '%s'\n", *startStr)
			return 1
		}
		v := clampInt64(n)
		if v == 0 {
			fmt.Fprintf(stderr, "error: sub: Invalid start value '0'\n")
			return 1
		}
		start = &v
	}

	if *endStr != "" {
		n, err := strconv.ParseInt(*endStr, 10, 64)
		if err != nil {
			fmt.Fprintf(stderr, "error: sub: Invalid end value '%s'\n", *endStr)
			return 1
		}
		v := clampInt64(n)
		if v == 0 {
			fmt.Fprintf(stderr, "error: sub: Invalid end value '0'\n")
			return 1
		}
		end = &v
	}

	if *lengthStr != "" {
		n, err := strconv.Atoi(*lengthStr)
		if err != nil {
			fmt.Fprintf(stderr, "error: sub: Invalid length value '%s'\n", *lengthStr)
			return 1
		}
		if n < 0 {
			fmt.Fprintf(stderr, "error: sub: Invalid length value '%d'\n", n)
			return 1
		}
		length = &n
	}

	if end != nil && length != nil {
		fmt.Fprintln(stderr, "error: sub: invalid option combination, --end and --length are mutually exclusive")
		return 1
	}

	found := false
	next := newStringIter(fs.Args(), stdin)
	for s, ok := next(); ok; s, ok = next() {
		result := substring(s, start, end, length)
		if len(result) > 0 {
			found = true
		}
		if !*quiet {
			fmt.Fprintln(stdout, result)
		}
		if *quiet && found {
			return 0
		}
	}

	if found {
		return 0
	}
	return 1
}

func clampInt64(n int64) int {
	return int(max(int64(math.MinInt32), min(int64(math.MaxInt32), n)))
}

func substring(s string, start, end, length *int) string {
	n := len(s)

	startIdx := 0
	if start != nil {
		startIdx = subToIndex(*start, n)
	}
	startIdx = max(0, min(n, startIdx))

	if length != nil {
		take := max(0, min(*length, n-startIdx))
		return s[startIdx : startIdx+take]
	}

	endIdx := n
	if end != nil {
		if *end >= 1 {
			endIdx = *end
		} else {
			endIdx = max(0, min(n, n+*end))
		}
	}
	endIdx = max(0, min(n, endIdx))

	if endIdx <= startIdx {
		return ""
	}
	return s[startIdx:endIdx]
}

func subToIndex(pos, length int) int {
	var idx int
	if pos >= 1 {
		idx = pos - 1
	} else {
		idx = length + pos
	}
	return max(0, min(length, idx))
}

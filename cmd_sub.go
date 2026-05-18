package main

import (
	"flag"
	"fmt"
	"io"
	"strconv"
	"math"

	"rsc.io/getopt"
)

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
		fmt.Fprintln(stdout, "Usage: string sub [-h] [(-s | --start) START] [(-e | --end) END] [(-l | --length) LENGTH] [-q] [STRING ...]")
		fmt.Fprintln(stdout, "")
		fmt.Fprintln(stdout, "  Extract substrings from STRING.")
		fmt.Fprintln(stdout, "")
		fmt.Fprintln(stdout, "Options:")
		fmt.Fprintln(stdout, "  -s, --start INT      Start position (1-based; negative counts from end)")
		fmt.Fprintln(stdout, "  -e, --end INT        End position, inclusive (1-based; negative counts from end)")
		fmt.Fprintln(stdout, "  -l, --length INT     Length of substring (overrides --end)")
		fmt.Fprintln(stdout, "  -q, --quiet          Suppress output; exit 0 if any result non-empty, 1 if all empty")
		fmt.Fprintln(stdout, "  -h, --help           Show this help message")
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

	any := false
	for _, s := range inputStrings(fs.Args(), stdin) {
		result := substring(s, start, end, length)
		if len(result) > 0 {
			any = true
		}
		if !*quiet {
			fmt.Fprintln(stdout, result)
		}
	}

	if any {
		return 0
	}
	return 1
}

func clampInt64(n int64) int {
	if n > math.MaxInt32 {
		return math.MaxInt32
	}
	if n < math.MinInt32 {
		return math.MinInt32
	}
	return int(n)
}

func substring(s string, start, end, length *int) string {
	n := len(s)

	startIdx := 0
	if start != nil {
		startIdx = subToIndex(*start, n)
	}
	startIdx = clampIdx(startIdx, 0, n)

	if length != nil {
		take := *length
		if take > n-startIdx {
			take = n - startIdx
		}
		if take < 0 {
			take = 0
		}
		return s[startIdx : startIdx+take]
	}

	endIdx := n
	if end != nil {
		if *end >= 1 {
			endIdx = *end
		} else {
			endIdx = clampIdx(n+*end, 0, n)
		}
	}
	endIdx = clampIdx(endIdx, 0, n)

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
	return clampIdx(idx, 0, length)
}

func clampIdx(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

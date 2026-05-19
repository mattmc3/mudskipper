package stringcmd

import (
	"flag"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"rsc.io/getopt"
)

const usageMatch = `Usage: string match [-h] [-a] [-e] [-i] [-g] [-n] [-r] [-q] [-v] [-m MAX] PATTERN [STRING ...]

  Match STRING against PATTERN (glob by default).

Options:
  -a, --all              Find all matches (not just the first)
  -e, --entire           Print the entire STRING for matches
  -i, --ignore-case      Ignore case when matching
  -g, --groups-only      Print only capture groups (requires -r)
  -n, --index            Print match position and length instead of value
  -r, --regex            Treat PATTERN as a regular expression
  -v, --invert           Print strings that do NOT match
  -m, --max-matches MAX  Maximum number of matches to output
  -q, --quiet            Suppress output; exit 0 if any match, 1 if none
  -h, --help             Show this help message
`

func runMatch(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := getopt.NewFlagSet("match", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	all := fs.Bool("all", false, "")
	entire := fs.Bool("entire", false, "")
	ignoreCase := fs.Bool("ignore-case", false, "")
	groupsOnly := fs.Bool("groups-only", false, "")
	useIndex := fs.Bool("index", false, "")
	useRegex := fs.Bool("regex", false, "")
	quiet := fs.Bool("quiet", false, "")
	invert := fs.Bool("invert", false, "")
	maxMatchesStr := fs.String("max-matches", "", "")
	fs.Aliases(
		"h", "help",
		"a", "all",
		"e", "entire",
		"i", "ignore-case",
		"g", "groups-only",
		"n", "index",
		"r", "regex",
		"q", "quiet",
		"v", "invert",
		"m", "max-matches",
	)

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: match: %v\n", err)
		return 1
	}
	if *help {
		fmt.Fprint(stdout, usageMatch)
		return 0
	}

	limit := -1
	if *maxMatchesStr != "" {
		n, err := strconv.Atoi(*maxMatchesStr)
		if err != nil || n <= 0 {
			fmt.Fprintf(stderr, "error: match: Invalid max matches value '%s'\n", *maxMatchesStr)
			return 1
		}
		limit = n
	}

	if *invert && *groupsOnly {
		fmt.Fprintln(stderr, "error: match: invalid option combination, --invert and --groups-only are mutually exclusive")
		return 1
	}
	if *entire && *useIndex {
		fmt.Fprintln(stderr, "error: match: invalid option combination, --entire and --index are mutually exclusive")
		return 1
	}
	if *entire && *groupsOnly {
		fmt.Fprintln(stderr, "error: match: invalid option combination, --entire and --groups-only are mutually exclusive")
		return 1
	}

	rest := fs.Args()
	if len(rest) == 0 {
		fmt.Fprintln(stderr, "error: match requires a pattern")
		return 1
	}

	re, err := buildMatchRegex(rest[0], *useRegex, *ignoreCase, *entire)
	if err != nil {
		fmt.Fprintf(stderr, "error: invalid pattern: %v\n", err)
		return 1
	}

	matched := false
	total := 0

	next := newStringIter(rest[1:], stdin)
	for s, ok := next(); ok; s, ok = next() {
		if limit >= 0 && total >= limit {
			break
		}

		if *invert {
			if !re.MatchString(s) {
				matched = true
				total++
				if !*quiet {
					if *useIndex {
						fmt.Fprintf(stdout, "1 %d\n", len(s))
					} else {
						fmt.Fprintln(stdout, s)
					}
				}
				if *quiet {
					return 0
				}
			}
			continue
		}

		if *all {
			for _, loc := range re.FindAllStringSubmatchIndex(s, -1) {
				if limit >= 0 && total >= limit {
					break
				}
				if writeMatchAt(stdout, s, loc, *groupsOnly, *useIndex, *entire, *quiet) {
					matched = true
					total++
				}
			}
		} else {
			if loc := re.FindStringSubmatchIndex(s); loc != nil && (limit < 0 || total < limit) {
				if writeMatchAt(stdout, s, loc, *groupsOnly, *useIndex, *entire, *quiet) {
					matched = true
					total++
				}
			}
		}
		if *quiet && matched {
			return 0
		}
	}

	if matched {
		return 0
	}
	return 1
}

func writeMatchAt(w io.Writer, s string, loc []int, groupsOnly, useIndex, entire, quiet bool) bool {
	if groupsOnly {
		found := false
		for g := 1; g < len(loc)/2; g++ {
			start, end := loc[g*2], loc[g*2+1]
			if start < 0 {
				continue
			}
			if !quiet {
				if useIndex {
					fmt.Fprintf(w, "%d %d\n", start+1, end-start)
				} else {
					fmt.Fprintln(w, s[start:end])
				}
			}
			found = true
		}
		return found
	}

	matchStart, matchEnd := loc[0], loc[1]
	if !quiet {
		if useIndex {
			fmt.Fprintf(w, "%d %d\n", matchStart+1, matchEnd-matchStart)
		} else if entire {
			fmt.Fprintln(w, s)
		} else {
			fmt.Fprintln(w, s[matchStart:matchEnd])
		}
		if !useIndex {
			for g := 1; g < len(loc)/2; g++ {
				start, end := loc[g*2], loc[g*2+1]
				if start >= 0 {
					fmt.Fprintln(w, s[start:end])
				}
			}
		}
	}
	return true
}

func buildMatchRegex(pattern string, useRegex, ignoreCase, entire bool) (*regexp.Regexp, error) {
	var pat string
	if useRegex {
		pat = pattern
	} else {
		pat = globToRegex(pattern, entire)
	}

	var flags strings.Builder
	flags.WriteString("(?s)")
	if ignoreCase {
		flags.WriteString("(?i)")
	}

	return regexp.Compile(flags.String() + pat)
}

func globToRegex(glob string, entire bool) string {
	var sb strings.Builder
	if !entire {
		sb.WriteString("^")
	}
	for i := 0; i < len(glob); {
		switch glob[i] {
		case '*':
			sb.WriteString(".*")
			i++
		case '?':
			sb.WriteByte('.')
			i++
		case '[':
			sb.WriteByte('[')
			i++
			if i < len(glob) && glob[i] == '!' {
				sb.WriteByte('^')
				i++
			}
			for i < len(glob) && glob[i] != ']' {
				sb.WriteByte(glob[i])
				i++
			}
			if i < len(glob) {
				sb.WriteByte(']')
				i++
			}
		default:
			sb.WriteString(regexp.QuoteMeta(string(glob[i])))
			i++
		}
	}
	if !entire {
		sb.WriteString("$")
	}
	return sb.String()
}

package main

import (
	"flag"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"rsc.io/getopt"
)

func runReplace(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := getopt.NewFlagSet("replace", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	all := fs.Bool("all", false, "")
	filter := fs.Bool("filter", false, "")
	ignoreCase := fs.Bool("ignore-case", false, "")
	useRegex := fs.Bool("regex", false, "")
	quiet := fs.Bool("quiet", false, "")
	maxStr := fs.String("max-matches", "", "")
	fs.Aliases("h", "help", "a", "all", "f", "filter", "i", "ignore-case", "r", "regex", "q", "quiet", "m", "max-matches")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: replace: %v\n", err)
		return 1
	}
	if *help {
		fmt.Fprintln(stdout, "Usage: string replace [-h] [-a] [-f] [-i] [-r] [-q] [(-m | --max-matches) MAX] PATTERN REPLACEMENT [STRING ...]")
		fmt.Fprintln(stdout, "")
		fmt.Fprintln(stdout, "  Replace PATTERN with REPLACEMENT in each STRING.")
		fmt.Fprintln(stdout, "")
		fmt.Fprintln(stdout, "Options:")
		fmt.Fprintln(stdout, "  -a, --all                Replace all occurrences (default: first only)")
		fmt.Fprintln(stdout, "  -f, --filter             Only print strings where a replacement was made")
		fmt.Fprintln(stdout, "  -i, --ignore-case        Case-insensitive matching")
		fmt.Fprintln(stdout, "  -r, --regex              Treat PATTERN as a regular expression")
		fmt.Fprintln(stdout, "  -m, --max-matches MAX    Maximum number of replacements per string")
		fmt.Fprintln(stdout, "  -q, --quiet              Suppress output; exit 0 if any replacement, 1 if none")
		fmt.Fprintln(stdout, "  -h, --help               Show this help message")
		return 0
	}

	maxMatches := 0
	if *maxStr != "" {
		n, err := strconv.Atoi(*maxStr)
		if err != nil || n < 0 {
			fmt.Fprintf(stderr, "error: replace: Invalid max matches value '%s'\n", *maxStr)
			return 1
		}
		maxMatches = n
	}

	rest := fs.Args()
	if len(rest) < 2 {
		fmt.Fprintln(stderr, "error: replace requires a pattern and replacement")
		return 1
	}

	pattern := rest[0]
	replacement := rest[1]
	limit := maxMatches
	if limit == 0 {
		if *all {
			limit = -1
		} else {
			limit = 1
		}
	}

	var re *regexp.Regexp
	if *useRegex {
		flags := "(?s)"
		if *ignoreCase {
			flags = "(?si)"
		}
		var err error
		re, err = regexp.Compile(flags + pattern)
		if err != nil {
			fmt.Fprintf(stderr, "error: invalid pattern: %v\n", err)
			return 1
		}
	}

	strs := inputStrings(rest[2:], stdin)
	changed := false

	for _, s := range strs {
		var result string
		if re != nil {
			result = replaceRegexLimited(s, re, replacement, limit)
		} else {
			result = replaceLiteral(s, pattern, replacement, limit, *ignoreCase)
		}

		didChange := result != s
		if didChange {
			changed = true
		}

		if !*filter || didChange {
			if !*quiet {
				fmt.Fprintln(stdout, result)
			}
		}
	}

	if changed {
		return 0
	}
	return 1
}

func replaceLiteral(s, pattern, replacement string, limit int, ignoreCase bool) string {
	if pattern == "" {
		return s
	}
	var sb strings.Builder
	count := 0
	pos := 0
	cmp := s
	pat := pattern
	if ignoreCase {
		cmp = strings.ToLower(s)
		pat = strings.ToLower(pattern)
	}
	for limit < 0 || count < limit {
		idx := strings.Index(cmp[pos:], pat)
		if idx < 0 {
			break
		}
		idx += pos
		sb.WriteString(s[pos:idx])
		sb.WriteString(replacement)
		pos = idx + len(pattern)
		count++
	}
	sb.WriteString(s[pos:])
	return sb.String()
}

func replaceRegexLimited(s string, re *regexp.Regexp, repl string, limit int) string {
	var sb strings.Builder
	pos := 0
	count := 0
	for limit < 0 || count < limit {
		sub := s[pos:]
		match := re.FindStringSubmatchIndex(sub)
		if match == nil {
			break
		}
		start, end := match[0], match[1]
		sb.WriteString(sub[:start])
		var expanded []byte
		expanded = re.ExpandString(expanded, repl, sub, match)
		sb.Write(expanded)
		pos += end
		if start == end {
			if end < len(sub) {
				sb.WriteByte(sub[end])
				pos++
			} else {
				break
			}
		}
		count++
	}
	sb.WriteString(s[pos:])
	return sb.String()
}

package stringcmd

import (
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"

	"rsc.io/getopt"
)

type splitOpts struct {
	help, quiet, noEmpty, right, allowEmpty bool
	fields, maxStr                          string
}

func parseSplitFlags(name string, args []string) (*splitOpts, *getopt.FlagSet, error) {
	var opts splitOpts
	fs := getopt.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.BoolVar(&opts.help, "help", false, "")
	fs.BoolVar(&opts.quiet, "quiet", false, "")
	fs.BoolVar(&opts.noEmpty, "no-empty", false, "")
	fs.BoolVar(&opts.right, "right", false, "")
	fs.BoolVar(&opts.allowEmpty, "allow-empty", false, "")
	fs.StringVar(&opts.fields, "fields", "", "")
	fs.StringVar(&opts.maxStr, "max", "", "")
	fs.Aliases("h", "help", "q", "quiet", "n", "no-empty", "r", "right", "a", "allow-empty", "f", "fields", "m", "max")
	return &opts, fs, fs.Parse(args)
}

func runSplit(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	opts, fs, err := parseSplitFlags("split", args)
	if err != nil {
		fmt.Fprintf(stderr, "error: split: %v\n", err)
		return 1
	}
	if opts.help {
		fmt.Fprint(stdout, `Usage: string split [-h] [-n] [-r] [-q] [(-f | --fields) FIELDS [-a]] [(-m | --max) MAX] SEP [STRING ...]

  Split each STRING by SEP.

Options:
  SEP                   Separator string
  -n, --no-empty        Suppress empty results
  -r, --right           Split from the right (useful with --max)
  -f, --fields FIELDS   Output only specified fields (e.g. 1,3-5)
  -a, --allow-empty     With --fields, skip missing fields instead of failing
  -m, --max MAX         Maximum number of splits per string
  -q, --quiet           Suppress output; exit 0 if any splits, 1 if none
  -h, --help            Show this help message
`)
		return 0
	}
	rest := fs.Args()
	if len(rest) == 0 {
		fmt.Fprintln(stderr, "error: split requires a separator")
		return 1
	}
	return splitCore(rest[0], false, rest[1:], stdin, stdout, stderr, opts)
}

func runSplit0(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	opts, fs, err := parseSplitFlags("split0", args)
	if err != nil {
		fmt.Fprintf(stderr, "error: split0: %v\n", err)
		return 1
	}
	if opts.help {
		fmt.Fprint(stdout, `Usage: string split0 [-h] [-n] [-r] [-q] [(-f | --fields) FIELDS [-a]] [(-m | --max) MAX] [STRING ...]

  Split each STRING by NUL (\0). Trailing NUL is ignored.

Options:
  -n, --no-empty        Suppress empty results
  -r, --right           Split from the right (useful with --max)
  -f, --fields FIELDS   Output only specified fields (e.g. 1,3-5)
  -a, --allow-empty     With --fields, skip missing fields instead of failing
  -m, --max MAX         Maximum number of splits per string
  -q, --quiet           Suppress output; exit 0 if any splits, 1 if none
  -h, --help            Show this help message
`)
		return 0
	}
	return splitCore("\x00", true, fs.Args(), stdin, stdout, stderr, opts)
}

func splitCore(sep string, nul0Mode bool, inputs []string, stdin io.Reader, stdout, stderr io.Writer, opts *splitOpts) int {
	maxSplits := 0
	if opts.maxStr != "" {
		n, err := strconv.Atoi(opts.maxStr)
		if err != nil || n < 0 {
			fmt.Fprintf(stderr, "error: split: Invalid max value '%s'\n", opts.maxStr)
			return 1
		}
		maxSplits = n
	}

	var fields []int
	if opts.fields != "" {
		var err error
		fields, err = parseFields(opts.fields)
		if err != nil {
			fmt.Fprintf(stderr, "error: %v\n", err)
			return 1
		}
	}

	if opts.allowEmpty && fields == nil {
		fmt.Fprintln(stderr, "error: split: --allow-empty is only valid with --fields")
		return 1
	}

	var strs []string
	if len(inputs) > 0 {
		if nul0Mode {
			for _, s := range inputs {
				strs = append(strs, strings.TrimRight(s, "\x00"))
			}
		} else {
			strs = inputs
		}
	} else {
		if nul0Mode {
			b, _ := io.ReadAll(stdin)
			content := string(b)
			if len(content) > 0 && content[len(content)-1] == '\x00' {
				content = content[:len(content)-1]
			}
			strs = []string{content}
		} else {
			strs = readLines(stdin)
		}
	}

	anySplit := false
	allFieldsFound := true

	for _, s := range strs {
		var parts []string
		if opts.right {
			parts = splitRight(s, sep, maxSplits)
		} else {
			parts = splitLeft(s, sep, maxSplits)
		}
		if len(parts) > 1 {
			anySplit = true
		}
		if opts.noEmpty {
			var filtered []string
			for _, p := range parts {
				if p != "" {
					filtered = append(filtered, p)
				}
			}
			parts = filtered
		}

		var selected []string
		if fields != nil {
			var ok bool
			selected, ok = selectFields(parts, fields, opts.allowEmpty)
			if !ok {
				allFieldsFound = false
			}
		} else {
			selected = parts
		}

		for _, p := range selected {
			if !opts.quiet {
				fmt.Fprintln(stdout, p)
			}
		}
	}

	if anySplit && allFieldsFound {
		return 0
	}
	return 1
}

func splitLeft(s, sep string, max int) []string {
	if sep == "" {
		if s == "" {
			return []string{""}
		}
		return strings.Split(s, "")
	}
	if max > 0 {
		return strings.SplitN(s, sep, max+1)
	}
	return strings.Split(s, sep)
}

func splitRight(s, sep string, max int) []string {
	if sep == "" {
		chars := strings.Split(s, "")
		if s == "" {
			chars = []string{""}
		}
		if max > 0 && len(chars) > max+1 {
			head := strings.Join(chars[:len(chars)-max], "")
			return append([]string{head}, chars[len(chars)-max:]...)
		}
		return chars
	}
	var parts []string
	count := 0
	remaining := s
	for {
		idx := strings.LastIndex(remaining, sep)
		if idx < 0 || (max > 0 && count >= max) {
			parts = append([]string{remaining}, parts...)
			break
		}
		parts = append([]string{remaining[idx+len(sep):]}, parts...)
		remaining = remaining[:idx]
		count++
	}
	return parts
}

func selectFields(parts []string, fields []int, allowEmpty bool) ([]string, bool) {
	var items []string
	allFound := true
	for _, f := range fields {
		if f >= 1 && f <= len(parts) {
			items = append(items, parts[f-1])
		} else if !allowEmpty {
			allFound = false
		}
	}
	return items, allFound
}

func parseFields(spec string) ([]int, error) {
	var fields []int
	for _, part := range strings.Split(spec, ",") {
		dashIdx := strings.Index(part[1:], "-")
		if dashIdx >= 0 {
			dashIdx++ // adjust for skipping first char
			fromStr := part[:dashIdx]
			toStr := part[dashIdx+1:]
			from, err1 := strconv.Atoi(fromStr)
			to, err2 := strconv.Atoi(toStr)
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("invalid field spec: %s", part)
			}
			if from <= 0 || to <= 0 {
				return nil, fmt.Errorf("split: Invalid range value for field '%s'", part)
			}
			step := 1
			if from > to {
				step = -1
			}
			for i := from; (step > 0 && i <= to) || (step < 0 && i >= to); i += step {
				fields = append(fields, i)
			}
		} else {
			f, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid field spec: %s", part)
			}
			if f <= 0 {
				return nil, fmt.Errorf("split: Invalid fields value '%s'", part)
			}
			fields = append(fields, f)
		}
	}
	return fields, nil
}

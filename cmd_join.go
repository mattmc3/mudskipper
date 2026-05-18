package main

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"rsc.io/getopt"
)

func runJoin(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs, help, quiet, noEmpty := joinFlags("join")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: join: %v\n", err)
		return 1
	}
	if *help {
		writeJoinHelp("join", stdout)
		return 0
	}
	rest := fs.Args()
	if len(rest) == 0 {
		fmt.Fprintln(stderr, "error: join requires a separator")
		return 1
	}
	return joinCore(rest[0], false, rest[1:], stdin, stdout, *quiet, *noEmpty)
}

func runJoin0(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs, help, quiet, noEmpty := joinFlags("join0")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: join0: %v\n", err)
		return 1
	}
	if *help {
		writeJoinHelp("join0", stdout)
		return 0
	}
	return joinCore("\x00", true, fs.Args(), stdin, stdout, *quiet, *noEmpty)
}

func joinFlags(name string) (*getopt.FlagSet, *bool, *bool, *bool) {
	fs := getopt.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	quiet := fs.Bool("quiet", false, "")
	noEmpty := fs.Bool("no-empty", false, "")
	fs.Aliases("h", "help", "q", "quiet", "n", "no-empty")
	return fs, help, quiet, noEmpty
}

func writeJoinHelp(name string, w io.Writer) {
	if name == "join0" {
		fmt.Fprintln(w, "Usage: string join0 [-h] [-q] [-n] [--] [STRING ...]")
	} else {
		fmt.Fprintln(w, "Usage: string join [-h] [-q] [-n] [--] SEP [STRING ...]")
	}
	fmt.Fprintln(w, "")
	if name == "join0" {
		fmt.Fprintln(w, "  Join strings with NUL (\\0) separator and a trailing NUL.")
	} else {
		fmt.Fprintln(w, "  Join strings with SEP separator.")
	}
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Options:")
	fmt.Fprintln(w, "  -n, --no-empty    Exclude empty strings")
	fmt.Fprintln(w, "  -q, --quiet       Suppress output; exit 0 if any strings joined, 1 if none")
	fmt.Fprintln(w, "  -h, --help        Show this help message")
}

func joinCore(sep string, appendNul bool, inputs []string, stdin io.Reader, stdout io.Writer, quiet, noEmpty bool) int {
	strs := inputStrings(inputs, stdin)
	if noEmpty {
		var filtered []string
		for _, s := range strs {
			if s != "" {
				filtered = append(filtered, s)
			}
		}
		strs = filtered
	}

	if len(strs) == 0 {
		return 1
	}

	if !quiet {
		fmt.Fprint(stdout, strings.Join(strs, sep))
		if appendNul {
			fmt.Fprint(stdout, "\x00")
		} else {
			fmt.Fprintln(stdout)
		}
	}

	if len(strs) >= 2 {
		return 0
	}
	return 1
}

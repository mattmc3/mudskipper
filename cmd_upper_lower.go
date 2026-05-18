package main

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"rsc.io/getopt"
)

func runUpper(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	return runTransform("upper", args, stdin, stdout, stderr, strings.ToUpper)
}

func runLower(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	return runTransform("lower", args, stdin, stdout, stderr, strings.ToLower)
}

func runTransform(name string, args []string, stdin io.Reader, stdout, stderr io.Writer, transform func(string) string) int {
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
		writeTransformHelp(name, stdout)
		return 0
	}

	changed := false
	for _, s := range inputStrings(fs.Args(), stdin) {
		result := transform(s)
		if !*quiet {
			fmt.Fprintln(stdout, result)
		}
		if result != s {
			changed = true
		}
	}
	if changed {
		return 0
	}
	return 1
}

func writeTransformHelp(name string, w io.Writer) {
	var verb string
	if name == "upper" {
		verb = "uppercase"
	} else {
		verb = "lowercase"
	}
	fmt.Fprintf(w, "Usage: string %s [-h] [-q] [STRING ...]\n", name)
	fmt.Fprintln(w)
	fmt.Fprintf(w, "  Convert STRING to %s.\n", verb)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Options:")
	fmt.Fprintln(w, "  -q, --quiet    Suppress output; exit 0 if any string changed, 1 if none")
	fmt.Fprintln(w, "  -h, --help     Show this help message")
}

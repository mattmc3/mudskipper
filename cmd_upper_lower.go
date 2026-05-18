package main

import (
	"fmt"
	"io"
	"strings"
)

func runUpper(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	return runTransform("upper", args, stdin, stdout, stderr, strings.ToUpper)
}

func runLower(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	return runTransform("lower", args, stdin, stdout, stderr, strings.ToLower)
}

func runTransform(name string, args []string, stdin io.Reader, stdout, stderr io.Writer, transform func(string) string) int {
	var quiet bool
	var inputs []string

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-h", "--help":
			writeTransformHelp(name, stdout)
			return 0
		case "-q", "--quiet":
			quiet = true
		case "--":
			inputs = append(inputs, args[i+1:]...)
			i = len(args)
		default:
			if strings.HasPrefix(args[i], "-") {
				fmt.Fprintf(stderr, "error: %s: unknown option '%s'\n", name, args[i])
				return 1
			}
			inputs = append(inputs, args[i])
		}
	}

	changed := false
	for _, s := range inputStrings(inputs, stdin) {
		result := transform(s)
		if !quiet {
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

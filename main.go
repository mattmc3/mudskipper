package main

import (
	"fmt"
	"io"
	"os"
)

const version = "0.0.1"

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) < 1 {
		fmt.Fprintln(stderr, "string: missing subcommand")
		return 1
	}

	cmd := args[0]
	switch cmd {
	case "--help", "-h", "help":
		writeHelp(stdout)
		return 0
	case "--version", "-v":
		fmt.Fprintf(stdout, "string %s\n", version)
		return 0
	}

	rest := args[1:]

	switch cmd {
	case "upper":
		return runUpper(rest, stdin, stdout, stderr)
	case "lower":
		return runLower(rest, stdin, stdout, stderr)
	case "length":
		return runLength(rest, stdin, stdout, stderr)
	case "trim":
		return runTrim(rest, stdin, stdout, stderr)
	case "match":
		return runMatch(rest, stdin, stdout, stderr)
	case "collect":
		return runCollect(rest, stdin, stdout, stderr)
	case "repeat":
		return runRepeat(rest, stdin, stdout, stderr)
	case "pad":
		return runPad(rest, stdin, stdout, stderr)
	case "sub":
		return runSub(rest, stdin, stdout, stderr)
	case "shorten":
		return runShorten(rest, stdin, stdout, stderr)
	case "replace":
		return runReplace(rest, stdin, stdout, stderr)
	default:
		fmt.Fprintf(stderr, "string %s: invalid subcommand\n", cmd)
		return 1
	}
}

func writeHelp(w io.Writer) {
	fmt.Fprintln(w, "Usage: string <command> [options] [STRING ...]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Commands:")
	fmt.Fprintln(w, "  upper    Convert strings to uppercase")
	fmt.Fprintln(w, "  lower    Convert strings to lowercase")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Use 'string <command> --help' for more information about a specific command.")
}

package stringcmd

import (
	"fmt"
	"io"
)

const version = "0.0.1"

// Run is the entry point for the string command.
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	return run(args, stdin, stdout, stderr)
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
	case "split":
		return runSplit(rest, stdin, stdout, stderr)
	case "split0":
		return runSplit0(rest, stdin, stdout, stderr)
	case "join":
		return runJoin(rest, stdin, stdout, stderr)
	case "join0":
		return runJoin0(rest, stdin, stdout, stderr)
	case "escape":
		return runEscape(rest, stdin, stdout, stderr)
	case "unescape":
		return runUnescape(rest, stdin, stdout, stderr)
	default:
		fmt.Fprintf(stderr, "string %s: invalid subcommand\n", cmd)
		return 1
	}
}

func writeHelp(w io.Writer) {
	fmt.Fprint(w, `Usage: string <command> [options] [STRING ...]

Commands:
  collect    Collect strings into one output
  escape     Escape strings for safe use in various contexts
  join       Join strings with a separator
  join0      Join strings with NUL separator
  length     Print the length of each string
  lower      Convert strings to lowercase
  match      Match strings against a pattern
  pad        Pad strings to a fixed width
  repeat     Repeat a string
  replace    Replace a pattern in strings
  shorten    Shorten strings to a maximum width
  split      Split strings by a separator
  split0     Split strings by NUL
  sub        Extract substrings
  trim       Trim whitespace or characters
  unescape   Unescape strings from encoded formats
  upper      Convert strings to uppercase

Use 'string <command> --help' for more information about a specific command.
`)
}

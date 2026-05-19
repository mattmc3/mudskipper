package stringcmd

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"rsc.io/getopt"
)

const usageCollect = `Usage: string collect [-h] [-a] [-N] [STRING ...]

  Collect all strings into a single output.

Options:
  -a, --allow-empty        Exit 0 even if result is empty
  -N, --no-trim-newlines   Preserve trailing newlines
  -h, --help               Show this help message
`

func runCollect(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := getopt.NewFlagSet("collect", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	allowEmpty := fs.Bool("allow-empty", false, "")
	noTrim := fs.Bool("no-trim-newlines", false, "")
	fs.Aliases("h", "help", "a", "allow-empty", "N", "no-trim-newlines")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: collect: %v\n", err)
		return 1
	}
	if *help {
		fmt.Fprint(stdout, usageCollect)
		return 0
	}

	var collected string
	if inputs := fs.Args(); len(inputs) > 0 {
		collected = strings.Join(inputs, "\n")
	} else {
		b, _ := io.ReadAll(stdin)
		collected = string(b)
	}

	if !*noTrim {
		collected = strings.TrimRight(collected, "\r\n")
	}

	if len(collected) == 0 {
		if *allowEmpty {
			return 0
		}
		return 1
	}

	if *noTrim {
		fmt.Fprint(stdout, collected)
	} else {
		fmt.Fprintln(stdout, collected)
	}
	return 0
}

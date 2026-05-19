package containscmd

import (
	"fmt"
	"io"
)

const usage = `Usage: contains [-h] [-i | --index] KEY [VALUE ...]

Test whether KEY appears in the list of VALUEs.
Exits 0 if found, 1 if not.

Options:
  -h, --help           Show this help message
  -i, --index          Print the 1-based index of the first match instead of just exiting
`

func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "contains: Key not specified")
		return 1
	}

	if args[0] == "-h" || args[0] == "--help" {
		fmt.Fprint(stdout, usage)
		return 0
	}

	printIndex := false
	if args[0] == "-i" || args[0] == "--index" {
		printIndex = true
		args = args[1:]
		if len(args) == 0 {
			fmt.Fprintln(stderr, "contains: Key not specified")
			return 1
		}
	}

	// -- allows searching for values that look like flags (e.g. contains -- -- a b --)
	if args[0] == "--" {
		args = args[1:]
		if len(args) == 0 {
			fmt.Fprintln(stderr, "contains: Key not specified")
			return 1
		}
	}

	value := args[0]
	list := args[1:]

	for i, item := range list {
		if item == value {
			if printIndex {
				fmt.Fprintln(stdout, i+1)
			}
			return 0
		}
	}
	return 1
}

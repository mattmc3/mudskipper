package containscmd

import (
	"fmt"
	"io"
)

func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "contains: Key not specified")
		return 1
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

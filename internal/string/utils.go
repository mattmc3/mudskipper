package stringcmd

import (
	"bufio"
	"io"
)

func readLines(r io.Reader) []string {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func inputStrings(args []string, stdin io.Reader) []string {
	if len(args) > 0 {
		return args
	}
	return readLines(stdin)
}

// newStringIter returns a function that yields one string at a time, allowing
// early exit without buffering all of stdin. Needed for --quiet early-exit behavior.
func newStringIter(args []string, stdin io.Reader) func() (string, bool) {
	if len(args) > 0 {
		i := 0
		return func() (string, bool) {
			if i >= len(args) {
				return "", false
			}
			s := args[i]
			i++
			return s, true
		}
	}
	scanner := bufio.NewScanner(stdin)
	return func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true
	}
}

package countcmd

import (
	"bufio"
	"fmt"
	"io"
)

func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	n := len(args)
	if n == 0 {
		scanner := bufio.NewScanner(stdin)
		for scanner.Scan() {
			n++
		}
	}
	fmt.Fprintln(stdout, n)
	if n > 0 {
		return 0
	}
	return 1
}

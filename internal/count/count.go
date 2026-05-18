package countcmd

import (
	"bufio"
	"fmt"
	"io"
)

func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	n := len(args)
	// Always count newlines from stdin (like wc -l), added to arg count.
	r := bufio.NewReader(stdin)
	for {
		b, err := r.ReadByte()
		if err != nil {
			break
		}
		if b == '\n' {
			n++
		}
	}
	fmt.Fprintln(stdout, n)
	if n > 0 {
		return 0
	}
	return 1
}

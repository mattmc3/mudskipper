package main

import (
	"io"
	"os"
	"strings"

	mathcmd "github.com/mattmc3/mudskipper/internal/math"
)

func main() {
	// Don't read stdin if it's a TTY — math with no args should error, not block.
	stdin := io.Reader(os.Stdin)
	fi, err := os.Stdin.Stat()
	if err != nil || (fi.Mode()&os.ModeCharDevice != 0) {
		stdin = strings.NewReader("")
	}
	os.Exit(mathcmd.Run(os.Args[1:], stdin, os.Stdout, os.Stderr))
}

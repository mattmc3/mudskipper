package main

import (
	"io"
	"os"
	"strings"

	countcmd "github.com/mattmc3/mudskipper/internal/count"
)

func main() {
	// count reads both args and stdin; skip stdin if it's a TTY to avoid blocking.
	stdin := io.Reader(os.Stdin)
	fi, err := os.Stdin.Stat()
	if err != nil || (fi.Mode()&os.ModeCharDevice != 0) {
		stdin = strings.NewReader("")
	}
	os.Exit(countcmd.Run(os.Args[1:], stdin, os.Stdout, os.Stderr))
}

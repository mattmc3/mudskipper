package main

import (
	"bufio"
	"os"

	stringcmd "github.com/mattmc3/string/internal/string"
)

func main() {
	bw := bufio.NewWriter(os.Stdout)
	code := stringcmd.Run(os.Args[1:], os.Stdin, bw, os.Stderr)
	bw.Flush()
	os.Exit(code)
}

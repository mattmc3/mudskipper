package main

import (
	"os"

	countcmd "github.com/mattmc3/string/internal/count"
)

func main() {
	os.Exit(countcmd.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

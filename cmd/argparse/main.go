package main

import (
	"os"

	argparsecmd "github.com/mattmc3/mudskipper/internal/argparse"
)

func main() {
	os.Exit(argparsecmd.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

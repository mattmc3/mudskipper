package main

import (
	"os"

	pathcmd "github.com/mattmc3/mudskipper/internal/path"
)

func main() {
	os.Exit(pathcmd.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

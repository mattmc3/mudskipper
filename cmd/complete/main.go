package main

import (
	"os"

	completecmd "github.com/mattmc3/mudskipper/internal/complete"
)

func main() {
	os.Exit(completecmd.Run(os.Args[1:], os.Stdout, os.Stderr))
}

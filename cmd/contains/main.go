package main

import (
	"os"

	containscmd "github.com/mattmc3/string/internal/contains"
)

func main() {
	os.Exit(containscmd.Run(os.Args[1:], os.Stdout, os.Stderr))
}

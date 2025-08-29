// filepath: /Users/matt/Projects/mattmc3/mudskipper/src/contains.go
package main

import (
	"fmt"
	"os"
)

const helpText = "Usage: contains [-i|--index|-0|--index0] NEEDLE HAYSTACK...\n" +
	"  -i, --index      Print the 1-based index of NEEDLE in HAYSTACK\n" +
	"  -0, --index0     Print the 0-based index of NEEDLE in HAYSTACK\n" +
	"  -h, --help       Show this help message\n" +
	"Returns 0 if NEEDLE is found, 1 otherwise.\n\n"

func main() {
	args := os.Args[1:]
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(helpText)
		return
	}

	indexFlag := false
	indexStart := 1
	start := 0

	if len(args) > 0 {
		if args[0] == "-i" || args[0] == "--index" {
			indexFlag = true
			indexStart = 1
			start = 1
		} else if args[0] == "-0" || args[0] == "--index0" {
			indexFlag = true
			indexStart = 0
			start = 1
		}
	}

	if len(args) <= start { // no needle
		fmt.Fprintln(os.Stderr, "contains: Key not specified")
		os.Exit(2)
	}

	needle := args[start]
	// Iterate remaining arguments after needle.
	// Relative index rel corresponds to haystack index starting at 0.
	for rel, item := range args[start+1:] {
		if item == needle {
			if indexFlag {
				fmt.Println(rel + indexStart)
			}
			os.Exit(0)
		}
	}
	os.Exit(1)
}

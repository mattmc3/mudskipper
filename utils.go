package main

import (
	"bufio"
	"io"
)

func readLines(r io.Reader) []string {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func inputStrings(args []string, stdin io.Reader) []string {
	if len(args) > 0 {
		return args
	}
	return readLines(stdin)
}

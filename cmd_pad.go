package main

import (
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"rsc.io/getopt"
)

func runPad(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := getopt.NewFlagSet("pad", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	right := fs.Bool("right", false, "")
	center := fs.Bool("center", false, "")
	charStr := fs.String("char", "", "")
	widthStr := fs.String("width", "", "")
	fs.Aliases("h", "help", "r", "right", "C", "center", "c", "char", "w", "width")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: pad: %v\n", err)
		return 1
	}
	if *help {
		fmt.Fprintln(stdout, "Usage: string pad [-h] [-r] [-C] [(-c | --char) CHAR] [(-w | --width) INTEGER] [STRING ...]")
		fmt.Fprintln(stdout, "")
		fmt.Fprintln(stdout, "  Pad strings to a fixed width.")
		fmt.Fprintln(stdout, "")
		fmt.Fprintln(stdout, "Options:")
		fmt.Fprintln(stdout, "  -r, --right         Pad on the right (left-align)")
		fmt.Fprintln(stdout, "  -C, --center        Center the string")
		fmt.Fprintln(stdout, "  -c, --char CHAR     Character to use for padding (default: space)")
		fmt.Fprintln(stdout, "  -w, --width INT     Target width (default: longest string)")
		fmt.Fprintln(stdout, "  -h, --help          Show this help message")
		return 0
	}

	padChar := ' '
	if *charStr != "" {
		runes := []rune(*charStr)
		if len(runes) != 1 {
			fmt.Fprintf(stderr, "error: pad: Padding should be a character '%s'\n", *charStr)
			return 1
		}
		if unicode.IsControl(runes[0]) {
			fmt.Fprintf(stderr, "error: pad: Invalid padding character of width zero '%s'\n", *charStr)
			return 1
		}
		padChar = runes[0]
	}

	width := -1
	if *widthStr != "" {
		n, err := strconv.Atoi(*widthStr)
		if err != nil || n < 0 {
			fmt.Fprintf(stderr, "error: pad: Invalid width value '%s'\n", *widthStr)
			return 1
		}
		width = n
	}

	strs := inputStrings(fs.Args(), stdin)
	if len(strs) == 0 {
		return 1
	}

	targetWidth := width
	for _, s := range strs {
		if len([]rune(s)) > targetWidth {
			targetWidth = len([]rune(s))
		}
	}
	if targetWidth < 0 {
		targetWidth = 0
	}

	for _, s := range strs {
		sRunes := []rune(s)
		sLen := len(sRunes)
		if sLen >= targetWidth {
			fmt.Fprintln(stdout, s)
			continue
		}
		total := targetWidth - sLen
		var result string
		if *center {
			var leftPad, rightPad int
			if *right {
				leftPad = total / 2
			} else {
				leftPad = (total + 1) / 2
			}
			rightPad = total - leftPad
			result = strings.Repeat(string(padChar), leftPad) + s + strings.Repeat(string(padChar), rightPad)
		} else if *right {
			result = s + strings.Repeat(string(padChar), total)
		} else {
			result = strings.Repeat(string(padChar), total) + s
		}
		fmt.Fprintln(stdout, result)
	}

	return 0
}

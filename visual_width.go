package main

import (
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
)

func skipAnsi(s string, i int) int {
	i++
	if i < len(s) && s[i] == '[' {
		i++
		for i < len(s) && !isASCIILetter(s[i]) {
			i++
		}
		if i < len(s) {
			i++
		}
	} else if i < len(s) {
		i++
	}
	return i
}

func isASCIILetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func visualWidthOf(s string) int {
	w := 0
	for i := 0; i < len(s); {
		if s[i] == 0x1b {
			i = skipAnsi(s, i)
			continue
		}
		if s[i] == '\b' {
			if w > 0 {
				w--
			}
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		w += runewidth.RuneWidth(r)
		i += size
	}
	return w
}

func visualWidthOfLines(s string) []int {
	var result []int
	pos := 0
	for i := 0; i < len(s); {
		if s[i] == '\n' {
			result = append(result, pos)
			pos = 0
			i++
			continue
		}
		if s[i] == '\r' {
			pos = 0
			i++
			continue
		}
		if s[i] == 0x1b {
			i = skipAnsi(s, i)
			continue
		}
		if s[i] == '\b' {
			if pos > 0 {
				pos--
			}
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		pos += runewidth.RuneWidth(r)
		i += size
	}
	return append(result, pos)
}

func visualTakeLeft(s string, targetWidth int) string {
	w := 0
	for i := 0; i < len(s); {
		if s[i] == '\b' {
			if w > 0 {
				w--
			}
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		rw := runewidth.RuneWidth(r)
		if rw == 0 {
			i += size
			continue
		}
		if w >= targetWidth {
			return s[:i]
		}
		w += rw
		i += size
	}
	return s
}

func visualTakeRight(s string, targetWidth int) string {
	w := 0
	i := len(s)
	for i > 0 {
		r, size := utf8.DecodeLastRuneInString(s[:i])
		rw := runewidth.RuneWidth(r)
		if rw > 0 {
			if w >= targetWidth {
				return s[i:]
			}
			w += rw
		}
		i -= size
	}
	return s[i:]
}

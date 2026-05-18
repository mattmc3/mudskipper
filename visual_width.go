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

type vseg struct {
	byteStart, byteEnd int
	width              int
}

func buildVSegs(s string) []vseg {
	var segs []vseg
	for i := 0; i < len(s); {
		if s[i] == 0x1b {
			end := skipAnsi(s, i)
			segs = append(segs, vseg{i, end, 0})
			i = end
			continue
		}
		if s[i] == '\b' {
			segs = append(segs, vseg{i, i + 1, -1})
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		segs = append(segs, vseg{i, i + size, runewidth.RuneWidth(r)})
		i += size
	}
	return segs
}

func visualTakeLeft(s string, targetWidth int) string {
	segs := buildVSegs(s)
	w := 0
	lastIdx := -1
	for i, seg := range segs {
		if seg.width == -1 {
			if w > 0 {
				w--
			}
			if w < targetWidth {
				lastIdx = i
			}
			continue
		}
		if seg.width == 0 {
			if w < targetWidth {
				lastIdx = i
			}
			continue
		}
		if w >= targetWidth {
			break
		}
		w += seg.width
		lastIdx = i
	}
	if lastIdx < 0 {
		return ""
	}
	return s[:segs[lastIdx].byteEnd]
}

func visualTakeRight(s string, targetWidth int) string {
	segs := buildVSegs(s)
	w := 0
	for j := len(segs) - 1; j >= 0; j-- {
		sw := segs[j].width
		if sw <= 0 {
			continue
		}
		if w >= targetWidth {
			return s[segs[j].byteEnd:]
		}
		w += sw
	}
	return s
}

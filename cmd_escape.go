package main

import (
	"flag"
	"fmt"
	"html"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"rsc.io/getopt"
)

var validEscapeStyles = map[string]bool{
	"script": true, "url": true, "html": true, "regex": true, "var": true,
}

func runEscape(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := getopt.NewFlagSet("escape", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	noQuoted := fs.Bool("no-quoted", false, "")
	style := fs.String("style", "script", "")
	fs.Aliases("h", "help", "n", "no-quoted")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: escape: %v\n", err)
		return 1
	}
	if *help {
		fmt.Fprint(stdout, `Usage: string escape [-h] [-n] [--style=STYLE] [STRING ...]

  Escape strings for safe use in various contexts.

Options:
  -n, --no-quoted    Skip quoting strings that don't need it (script style only)
  --style=STYLE      script (default), url, html, regex, var
  -h, --help         Show this help message
`)
		return 0
	}
	if !validEscapeStyles[*style] {
		fmt.Fprintf(stderr, "error: escape: Invalid escape style '%s'\n", *style)
		return 1
	}

	any := false
	for _, s := range inputStrings(fs.Args(), stdin) {
		fmt.Fprintln(stdout, escapeString(s, *style, *noQuoted))
		any = true
	}
	if any {
		return 0
	}
	return 1
}

func runUnescape(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := getopt.NewFlagSet("unescape", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	style := fs.String("style", "script", "")
	fs.Aliases("h", "help")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: unescape: %v\n", err)
		return 1
	}
	if *help {
		fmt.Fprint(stdout, `Usage: string unescape [-h] [--style=STYLE] [STRING ...]

  Unescape strings from various encoded formats.

Options:
  --style=STYLE    script (default), url, html, regex, var
  -h, --help       Show this help message
`)
		return 0
	}
	if !validEscapeStyles[*style] {
		fmt.Fprintf(stderr, "error: unescape: Invalid style value '%s'\n", *style)
		return 1
	}

	any := false
	for _, s := range inputStrings(fs.Args(), stdin) {
		result, err := unescapeString(s, *style)
		if err != nil {
			fmt.Fprintf(stderr, "error: unescape failed: %v\n", err)
			return 1
		}
		fmt.Fprintln(stdout, result)
		any = true
	}
	if any {
		return 0
	}
	return 1
}

func escapeString(s, style string, noQuoted bool) string {
	switch style {
	case "url":
		return strings.ReplaceAll(url.QueryEscape(s), "+", "%20")
	case "html":
		return html.EscapeString(s)
	case "regex":
		escaped := regexp.QuoteMeta(s)
		escaped = strings.ReplaceAll(escaped, "\n", `\n`)
		escaped = strings.ReplaceAll(escaped, "\r", `\r`)
		escaped = strings.ReplaceAll(escaped, "\t", `\t`)
		return escaped
	case "var":
		return escapeVar(s)
	default:
		return escapeScript(s, noQuoted)
	}
}

func unescapeString(s, style string) (string, error) {
	switch style {
	case "url":
		return url.PathUnescape(s)
	case "html":
		return html.UnescapeString(s), nil
	case "regex":
		return unescapeRegex(s), nil
	case "var":
		return unescapeVar(s), nil
	default:
		return unescapeScript(s), nil
	}
}

func escapeScript(s string, noQuoted bool) string {
	safe := len(s) > 0
	for _, c := range s {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c != '-' && c != '_' && c != '.' && c != '/' {
			safe = false
			break
		}
	}
	if safe {
		if noQuoted {
			return s
		}
		return "'" + s + "'"
	}
	if noQuoted {
		return escapeScriptNoQuoted(s)
	}
	return escapeScriptQuoted(s)
}

func escapeScriptQuoted(s string) string {
	var sb strings.Builder
	inQuote := false
	for _, c := range s {
		if isScriptControl(c) {
			if inQuote {
				sb.WriteByte('\'')
				inQuote = false
			}
			sb.WriteString(escapeControl(c))
		} else {
			if !inQuote {
				sb.WriteByte('\'')
				inQuote = true
			}
			if c == '\'' {
				sb.WriteByte('\\')
			}
			sb.WriteRune(c)
		}
	}
	if inQuote {
		sb.WriteByte('\'')
	}
	if sb.Len() == 0 {
		sb.WriteString("''")
	}
	return sb.String()
}

func escapeScriptNoQuoted(s string) string {
	var sb strings.Builder
	var prev rune
	first := true
	for _, c := range s {
		if isScriptControl(c) {
			sb.WriteString(escapeControl(c))
		} else if c == '#' && (first || prev == ' ' || prev == '\t') {
			sb.WriteString("\\#")
		} else if needsBackslash(c) {
			sb.WriteByte('\\')
			sb.WriteRune(c)
		} else {
			sb.WriteRune(c)
		}
		prev = c
		first = false
	}
	return sb.String()
}

func isScriptControl(c rune) bool { return c < 0x20 || c == 0x7f }

func escapeControl(c rune) string {
	n := int(c)
	if n >= 1 && n <= 26 {
		return fmt.Sprintf("\\c%c", 'a'+n-1)
	}
	if n == 27 {
		return "\\e"
	}
	if n == 127 {
		return "\\x7f"
	}
	return fmt.Sprintf("\\x%02x", n)
}

func needsBackslash(c rune) bool {
	switch c {
	case ' ', '\t', '"', '\'', '\\', '$', '~', '|', '&', ';', '(', ')', '[', ']', '{', '}', '<', '>':
		return true
	}
	return false
}

func unescapeScript(s string) string {
	var sb strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '\'' {
			i++
			for i < len(s) && s[i] != '\'' {
				if s[i] == '\\' && i+1 < len(s) && (s[i+1] == '\'' || s[i+1] == '\\') {
					sb.WriteByte(s[i+1])
					i += 2
				} else {
					sb.WriteByte(s[i])
					i++
				}
			}
			if i < len(s) {
				i++
			}
		} else if s[i] == '\\' && i+1 < len(s) {
			sb.WriteByte(s[i+1])
			i += 2
		} else {
			sb.WriteByte(s[i])
			i++
		}
	}
	return sb.String()
}

func escapeVar(s string) string {
	var sb strings.Builder
	for _, c := range s {
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			sb.WriteRune(c)
		} else if c == '_' {
			sb.WriteString("__")
		} else {
			fmt.Fprintf(&sb, "_%02X_", int(c))
		}
	}
	return sb.String()
}

func unescapeVar(s string) string {
	var sb strings.Builder
	i := 0
	for i < len(s) {
		if s[i] != '_' {
			sb.WriteByte(s[i])
			i++
			continue
		}
		if i+1 < len(s) && s[i+1] == '_' {
			sb.WriteByte('_')
			i += 2
			continue
		}
		j := i + 1
		for j < len(s) && s[j] != '_' {
			b := s[j]
			if !((b >= '0' && b <= '9') || (b >= 'A' && b <= 'F') || (b >= 'a' && b <= 'f')) {
				break
			}
			j++
		}
		if j > i+1 && j < len(s) && s[j] == '_' {
			if n, err := strconv.ParseInt(s[i+1:j], 16, 32); err == nil {
				sb.WriteRune(rune(n))
				i = j + 1
				continue
			}
		}
		sb.WriteByte(s[i])
		i++
	}
	return sb.String()
}

func unescapeRegex(s string) string {
	var sb strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '\\' && i+1 < len(s) {
			sb.WriteByte(s[i+1])
			i += 2
		} else {
			sb.WriteByte(s[i])
			i++
		}
	}
	return sb.String()
}

package main

import (
	"testing"
)

func TestEscape_script_wraps_in_single_quotes(t *testing.T) {
	exit, stdout, _ := runCmd("escape", "hello world")
	assertExit(t, 0, exit)
	assertLines(t, []string{"'hello world'"}, lines(stdout))
}

func TestEscape_script_safe_string_still_quoted(t *testing.T) {
	exit, stdout, _ := runCmd("escape", "hello")
	assertExit(t, 0, exit)
	assertLines(t, []string{"'hello'"}, lines(stdout))
}

func TestEscape_script_no_quoted_skips_safe_string(t *testing.T) {
	exit, stdout, _ := runCmd("escape", "-n", "hello")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hello"}, lines(stdout))
}

func TestEscape_script_no_quoted_uses_backslash_for_unsafe_string(t *testing.T) {
	exit, stdout, _ := runCmd("escape", "-n", "hello world")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hello\\ world"}, lines(stdout))
}

func TestEscape_script_embeds_single_quote(t *testing.T) {
	exit, stdout, _ := runCmd("escape", "it's")
	assertExit(t, 0, exit)
	assertLines(t, []string{"'it\\'s'"}, lines(stdout))
}

func TestEscape_url_encodes_special_chars(t *testing.T) {
	exit, stdout, _ := runCmd("escape", "--style=url", "hello world")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hello%20world"}, lines(stdout))
}

func TestEscape_html_encodes_angle_brackets(t *testing.T) {
	exit, stdout, _ := runCmd("escape", "--style=html", "<b>bold</b>")
	assertExit(t, 0, exit)
	assertLines(t, []string{"&lt;b&gt;bold&lt;/b&gt;"}, lines(stdout))
}

func TestEscape_regex_escapes_metacharacters(t *testing.T) {
	exit, stdout, _ := runCmd("escape", "--style=regex", `a.b*c`)
	assertExit(t, 0, exit)
	assertLines(t, []string{`a\.b\*c`}, lines(stdout))
}

func TestEscape_var_encodes_non_alphanumeric(t *testing.T) {
	exit, stdout, _ := runCmd("escape", "--style=var", "hello world")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hello_20_world"}, lines(stdout))
}

func TestEscape_empty_returns_1(t *testing.T) {
	exit, _, _ := runWithStdin("", "escape")
	assertExit(t, 1, exit)
}

func TestEscape_help_shows_usage(t *testing.T) {
	exit, stdout, _ := runCmd("escape", "--help")
	assertExit(t, 0, exit)
	assertContains(t, "Usage:", stdout)
}

func TestEscape_stdin_processes_lines(t *testing.T) {
	exit, stdout, _ := runWithStdin("hello\nworld\n", "escape", "-n")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hello", "world"}, lines(stdout))
}

func TestUnescape_script_removes_single_quotes(t *testing.T) {
	exit, stdout, _ := runCmd("unescape", "'hello world'")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hello world"}, lines(stdout))
}

func TestUnescape_script_handles_embedded_quote(t *testing.T) {
	exit, stdout, _ := runCmd("unescape", "'it'\\''s'")
	assertExit(t, 0, exit)
	assertLines(t, []string{"it's"}, lines(stdout))
}

func TestUnescape_url_decodes(t *testing.T) {
	exit, stdout, _ := runCmd("unescape", "--style=url", "hello%20world")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hello world"}, lines(stdout))
}

func TestUnescape_html_decodes(t *testing.T) {
	exit, stdout, _ := runCmd("unescape", "--style=html", "&lt;b&gt;bold&lt;/b&gt;")
	assertExit(t, 0, exit)
	assertLines(t, []string{"<b>bold</b>"}, lines(stdout))
}

func TestUnescape_regex_unescapes(t *testing.T) {
	exit, stdout, _ := runCmd("unescape", "--style=regex", `a\.b\*c`)
	assertExit(t, 0, exit)
	assertLines(t, []string{"a.b*c"}, lines(stdout))
}

func TestUnescape_var_decodes(t *testing.T) {
	exit, stdout, _ := runCmd("unescape", "--style=var", "hello_20_world")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hello world"}, lines(stdout))
}

func TestUnescape_empty_returns_1(t *testing.T) {
	exit, _, _ := runWithStdin("", "unescape")
	assertExit(t, 1, exit)
}

func TestUnescape_help_shows_usage(t *testing.T) {
	exit, stdout, _ := runCmd("unescape", "--help")
	assertExit(t, 0, exit)
	assertContains(t, "Usage:", stdout)
}

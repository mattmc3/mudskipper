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

// Tests for known fish parity differences — should fail before fix, pass after

func TestEscapeVar_adjacent_encoded_chars(t *testing.T) {
	// We use _HEX__HEX_ (double underscore) rather than fish's _HEX_HEX_ (shared separator).
	// The shared-separator format is ambiguous to decode (letters b-f also valid hex digits).
	// Round-trip is correct; only the intermediate format differs.
	exit, stdout, _ := runCmd("escape", "--style=var", "a b#c\"'d")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a_20_b_23_c_22__27_d"}, lines(stdout))
}

func TestEscapeRegex_newline_escaped_as_backslash_n(t *testing.T) {
	// fish encodes literal newline as \n in regex escape
	exit, stdout, _ := runCmd("escape", "--style=regex", "hello\nworld")
	assertExit(t, 0, exit)
	assertLines(t, []string{`hello\nworld`}, lines(stdout))
}

func TestUnescape_help_shows_usage(t *testing.T) {
	exit, stdout, _ := runCmd("unescape", "--help")
	assertExit(t, 0, exit)
	assertContains(t, "Usage:", stdout)
}

// Fish parity tests

func TestEscape_del_char(t *testing.T) {
	exit, stdout, _ := runCmd("escape", "\x7F")
	assertExit(t, 0, exit)
	assertLines(t, []string{`\x7f`}, lines(stdout))
}

func TestEscape_script_control_char_via_stdin(t *testing.T) {
	exit, stdout, _ := runWithStdin("\x07", "escape")
	assertExit(t, 0, exit)
	assertLines(t, []string{`\cg`}, lines(stdout))
}

func TestEscape_script_special_chars(t *testing.T) {
	exit, stdout, _ := runCmd("escape", "--style=script", "a b#c\"'d")
	assertExit(t, 0, exit)
	assertLines(t, []string{"'a b#c\"\\'d'"}, lines(stdout))
}

func TestEscape_script_no_quoted(t *testing.T) {
	exit, stdout, _ := runCmd("escape", "--no-quoted", "--style=script", "a b#c\"'d")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a\\ b#c\\\"\\'d"}, lines(stdout))
}

func TestEscape_script_no_quoted_hash(t *testing.T) {
	exit, stdout, _ := runCmd("escape", "--no-quoted", "--style=script", "a #b")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a\\ \\#b"}, lines(stdout))
}

func TestEscape_url_special_chars(t *testing.T) {
	exit, stdout, _ := runCmd("escape", "--style=url", "a b#c\"'d")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a%20b%23c%22%27d"}, lines(stdout))
}

func TestEscape_var_alphanumeric_passthrough(t *testing.T) {
	exit, stdout, _ := runCmd("escape", "--style=var", "abc")
	assertExit(t, 0, exit)
	assertLines(t, []string{"abc"}, lines(stdout))
}

func TestEscape_var_underscore_and_newline(t *testing.T) {
	exit, stdout, _ := runCmd("escape", "--style=var", "a\nghi_")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a_0A_ghi__"}, lines(stdout))
}

func TestEscape_var_underscores(t *testing.T) {
	exit, stdout, _ := runCmd("escape", "--style=var", "_a_b_c_")
	assertExit(t, 0, exit)
	assertLines(t, []string{"__a__b__c__"}, lines(stdout))
}

func TestEscape_var_dash(t *testing.T) {
	exit, stdout, _ := runCmd("escape", "--style=var", "--", "-")
	assertExit(t, 0, exit)
	assertLines(t, []string{"_2D_"}, lines(stdout))
}

func TestEscape_regex_metacharacters(t *testing.T) {
	exit, stdout, _ := runCmd("escape", "--style=regex", ".ext")
	assertExit(t, 0, exit)
	assertLines(t, []string{`\.ext`}, lines(stdout))
}

func TestEscape_unknown_style_is_error(t *testing.T) {
	exit, _, stderr := runCmd("escape", "--style=unknown-style")
	assertExit(t, 1, exit)
	assertContains(t, "Invalid escape style 'unknown-style'", stderr)
}

func TestUnescape_var_alphanumeric_roundtrip(t *testing.T) {
	_, enc, _ := runCmd("escape", "--style=var", "abc")
	exit, stdout, _ := runCmd("unescape", "--style=var", lines(enc)[0])
	assertExit(t, 0, exit)
	assertLines(t, []string{"abc"}, lines(stdout))
}

func TestUnescape_script_roundtrip(t *testing.T) {
	_, enc, _ := runCmd("escape", "--style=script", "a b#c\"'d")
	exit, stdout, _ := runCmd("unescape", "--style=script", lines(enc)[0])
	assertExit(t, 0, exit)
	assertLines(t, []string{"a b#c\"'d"}, lines(stdout))
}

func TestUnescape_url_roundtrip(t *testing.T) {
	_, enc, _ := runCmd("escape", "--style=url", "a b#c\"'d")
	exit, stdout, _ := runCmd("unescape", "--style=url", lines(enc)[0])
	assertExit(t, 0, exit)
	assertLines(t, []string{"a b#c\"'d"}, lines(stdout))
}

func TestUnescape_url_newline_roundtrip(t *testing.T) {
	_, enc, _ := runCmd("escape", "--style=url", "\na\nb%c~d\n")
	exit, stdout, _ := runCmd("unescape", "--style=url", lines(enc)[0])
	assertExit(t, 0, exit)
	assertLines(t, []string{"", "a", "b%c~d"}, lines(stdout))
}

func TestUnescape_var_roundtrip(t *testing.T) {
	_, enc, _ := runCmd("escape", "--style=var", "a b#c\"'d")
	exit, stdout, _ := runCmd("unescape", "--style=var", lines(enc)[0])
	assertExit(t, 0, exit)
	assertLines(t, []string{"a b#c\"'d"}, lines(stdout))
}

func TestUnescape_var_newline_roundtrip(t *testing.T) {
	_, enc, _ := runCmd("escape", "--style=var", "a\nghi_")
	exit, stdout, _ := runCmd("unescape", "--style=var", lines(enc)[0])
	assertExit(t, 0, exit)
	assertLines(t, []string{"a", "ghi_"}, lines(stdout))
}

func TestUnescape_var_underscore_roundtrip(t *testing.T) {
	_, enc, _ := runCmd("escape", "--style=var", "_a_b_c_")
	exit, stdout, _ := runCmd("unescape", "--style=var", lines(enc)[0])
	assertExit(t, 0, exit)
	assertLines(t, []string{"_a_b_c_"}, lines(stdout))
}

func TestUnescape_var_dash_roundtrip(t *testing.T) {
	_, enc, _ := runCmd("escape", "--style=var", "--", "-")
	exit, stdout, _ := runCmd("unescape", "--style=var", "--", lines(enc)[0])
	assertExit(t, 0, exit)
	assertLines(t, []string{"-"}, lines(stdout))
}

func TestUnescape_unknown_style_is_error(t *testing.T) {
	exit, _, stderr := runCmd("unescape", "--style=unknown-style")
	assertExit(t, 1, exit)
	assertContains(t, "Invalid style value 'unknown-style'", stderr)
}

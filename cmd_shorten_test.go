package main

import (
	"testing"
)

func TestShorten_truncates_right_by_default(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "-m", "5", "hello world")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hell…"}, lines(stdout))
}

func TestShorten_truncates_left_with_flag(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "-m", "5", "-l", "hello world")
	assertExit(t, 0, exit)
	assertLines(t, []string{"…orld"}, lines(stdout))
}

func TestShorten_custom_ellipsis(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "-m", "7", "-c", "...", "hello world")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hell..."}, lines(stdout))
}

func TestShorten_no_change_returns_0(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "-m", "20", "hello")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hello"}, lines(stdout))
}

func TestShorten_exact_length_unchanged(t *testing.T) {
	exit, _, _ := runCmd("shorten", "-m", "5", "hello")
	assertExit(t, 0, exit)
}

func TestShorten_multiple_strings(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "-m", "4", "hello", "hi", "world")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hel…", "hi", "wor…"}, lines(stdout))
}

func TestShorten_quiet_suppresses_output(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "-q", "-m", "3", "hello")
	assertExit(t, 0, exit)
	assertEmpty(t, stdout)
}

func TestShorten_no_newline_on_last(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "-N", "-m", "4", "hello")
	assertExit(t, 0, exit)
	assertStr(t, "hel…", stdout)
}

func TestShorten_max_shorter_than_ellipsis_truncates_content(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "-m", "1", "-c", "...", "hello")
	assertExit(t, 0, exit)
	assertLines(t, []string{"h"}, lines(stdout))
}

func TestShorten_help_shows_usage(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "--help")
	assertExit(t, 0, exit)
	assertContains(t, "Usage:", stdout)
}

func TestStdin_shorten(t *testing.T) {
	exit, stdout, _ := runWithStdin("hello world\nhi\n", "shorten", "-m", "6")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hello…", "hi"}, lines(stdout))
}

// Fish parity tests

func TestShorten_no_truncation_needed(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "-m", "3", "foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foo"}, lines(stdout))
}

func TestShorten_auto_width_from_shortest_string(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "foo", "foobar")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foo", "fo…"}, lines(stdout))
}

func TestShorten_max_zero_returns_all_as_is(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "-m0", "foo", "bar", "asodjsaoidj")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foo", "bar", "asodjsaoidj"}, lines(stdout))
}

func TestShorten_custom_single_char_ellipsis(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "-c", "w", "foo", "foobar")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foo", "fow"}, lines(stdout))
}

func TestShorten_empty_ellipsis(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "--max", "2", "--char", "", "foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"fo"}, lines(stdout))
}

func TestShorten_quiet_no_change_returns_1(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "-m", "2", "-q", "12")
	assertExit(t, 1, exit)
	assertEmpty(t, stdout)
}

func TestShorten_seq_auto_width_single_char_ellipsis(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "-c", "x", "2", "4", "8", "16", "32", "64", "128", "256", "512", "1024")
	assertExit(t, 0, exit)
	assertLines(t, []string{"2", "4", "8", "x", "x", "x", "x", "x", "x", "x"}, lines(stdout))
}

func TestShorten_basic_truncation(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "-m", "2", "foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"f…"}, lines(stdout))
}

func TestShorten_truncates_to_max(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "-m", "5", "foobar")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foob…"}, lines(stdout))
}

func TestShorten_ellipsis_longer_than_width_truncates_to_fit(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "-m", "5", "--char", "........", "foobar")
	assertExit(t, 0, exit)
	assertLines(t, []string{"fooba"}, lines(stdout))
}

func TestShorten_custom_ellipsis_3char(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "--max", "4", "-c", "///", "foobar")
	assertExit(t, 0, exit)
	assertLines(t, []string{"f///"}, lines(stdout))
}

func TestShorten_negative_max_is_error(t *testing.T) {
	exit, _, stderr := runCmd("shorten", "--max=-1", "--char", "", "foo")
	assertExit(t, 1, exit)
	assertContains(t, "Invalid max value '-1'", stderr)
}

func TestShorten_multiple_strings_same_max(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "-m4", "foobar", "bananarama")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foo…", "ban…"}, lines(stdout))
}

func TestShorten_stdin_auto_width_from_first_nonempty_line(t *testing.T) {
	exit, stdout, _ := runWithStdin("\n1. line\n2. another line\n3. third line", "shorten")
	assertExit(t, 0, exit)
	assertLines(t, []string{"", "1. line", "2. ano…", "3. thi…"}, lines(stdout))
}

func TestShorten_stdin_auto_width_left(t *testing.T) {
	exit, stdout, _ := runWithStdin("\n1. line\n2. another line\n3. third line", "shorten", "--left")
	assertExit(t, 0, exit)
	assertLines(t, []string{"", "1. line", "…r line", "…d line"}, lines(stdout))
}

func TestShorten_left_quiet_no_change_returns_1(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "-l", "-m", "2", "-q", "12")
	assertExit(t, 1, exit)
	assertEmpty(t, stdout)
}

func TestShorten_left_truncation(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "-l", "-m", "4", "foobar")
	assertExit(t, 0, exit)
	assertLines(t, []string{"…bar"}, lines(stdout))
}

func TestShorten_left_custom_ellipsis(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "-l", "-m", "4", "-c", "//", "foobar")
	assertExit(t, 0, exit)
	assertLines(t, []string{"//ar"}, lines(stdout))
}

func TestShorten_leading_backspace(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "abc", "\bab", "ab", "abcdef")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a…", "\bab", "ab", "a…"}, lines(stdout))
}

func TestShorten_leading_bell(t *testing.T) {
	exit, stdout, _ := runCmd("shorten", "abc", "\aab", "ab", "abcdef")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a…", "\aab", "ab", "a…"}, lines(stdout))
}

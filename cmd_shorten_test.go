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

package main

import (
	"testing"
)

func TestRepeat_basic(t *testing.T) {
	exit, stdout, _ := runCmd("repeat", "-n", "3", "ab")
	assertExit(t, 0, exit)
	assertLines(t, []string{"ababab"}, lines(stdout))
}

func TestRepeat_multiple_strings(t *testing.T) {
	exit, stdout, _ := runCmd("repeat", "-n", "2", "ab", "cd")
	assertExit(t, 0, exit)
	assertLines(t, []string{"abab", "cdcd"}, lines(stdout))
}

func TestRepeat_max_truncates(t *testing.T) {
	exit, stdout, _ := runCmd("repeat", "-n", "3", "-m", "4", "ab")
	assertExit(t, 0, exit)
	assertLines(t, []string{"abab"}, lines(stdout))
}

func TestRepeat_no_newline_suppresses_trailing_newline(t *testing.T) {
	exit, stdout, _ := runCmd("repeat", "-n", "2", "-N", "ab")
	assertExit(t, 0, exit)
	assertStr(t, "abab", stdout)
}

func TestRepeat_no_newline_only_suppresses_last(t *testing.T) {
	exit, stdout, _ := runCmd("repeat", "-n", "2", "-N", "ab", "cd")
	assertExit(t, 0, exit)
	assertStr(t, "abab\ncdcd", stdout)
}

func TestRepeat_quiet_suppresses_output(t *testing.T) {
	exit, stdout, _ := runCmd("repeat", "-n", "3", "-q", "ab")
	assertExit(t, 0, exit)
	assertEmpty(t, stdout)
}

func TestRepeat_returns_1_when_count_zero(t *testing.T) {
	exit, _, _ := runCmd("repeat", "-n", "0", "ab")
	assertExit(t, 1, exit)
}

func TestRepeat_returns_1_with_no_strings(t *testing.T) {
	exit, _, _ := runCmd("repeat", "-n", "3")
	assertExit(t, 1, exit)
}

func TestRepeat_help_shows_usage(t *testing.T) {
	exit, stdout, _ := runCmd("repeat", "--help")
	assertExit(t, 0, exit)
	assertContains(t, "Usage:", stdout)
}

func TestStdin_repeat(t *testing.T) {
	exit, stdout, _ := runWithStdin("ab\ncd\n", "repeat", "-n", "2")
	assertExit(t, 0, exit)
	assertLines(t, []string{"abab", "cdcd"}, lines(stdout))
}

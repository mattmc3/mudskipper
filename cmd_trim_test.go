package main

import (
	"testing"
)

func TestTrim_both_sides_by_default(t *testing.T) {
	exit, stdout, _ := runCmd("trim", "  hello  ")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hello"}, lines(stdout))
}

func TestTrim_left_only(t *testing.T) {
	exit, stdout, _ := runCmd("trim", "-l", "  hello  ")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hello  "}, lines(stdout))
}

func TestTrim_right_only(t *testing.T) {
	exit, stdout, _ := runCmd("trim", "-r", "  hello  ")
	assertExit(t, 0, exit)
	assertLines(t, []string{"  hello"}, lines(stdout))
}

func TestTrim_custom_chars(t *testing.T) {
	exit, stdout, _ := runCmd("trim", "-c", "xy", "xxhelloyx")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hello"}, lines(stdout))
}

func TestTrim_quiet_suppresses_output(t *testing.T) {
	exit, stdout, _ := runCmd("trim", "-q", "  hello  ")
	assertExit(t, 0, exit)
	assertEmpty(t, stdout)
}

func TestTrim_returns_1_when_nothing_trimmed(t *testing.T) {
	exit, _, _ := runCmd("trim", "hello")
	assertExit(t, 1, exit)
}

func TestTrim_multiple_strings(t *testing.T) {
	exit, stdout, _ := runCmd("trim", "  foo  ", "  bar  ")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foo", "bar"}, lines(stdout))
}

func TestTrim_help_shows_usage(t *testing.T) {
	exit, stdout, _ := runCmd("trim", "--help")
	assertExit(t, 0, exit)
	assertContains(t, "Usage:", stdout)
}

func TestStdin_trim(t *testing.T) {
	exit, stdout, _ := runWithStdin("  foo  \n  bar  \n", "trim")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foo", "bar"}, lines(stdout))
}

func TestStdin_trim_with_flag(t *testing.T) {
	exit, stdout, _ := runWithStdin("  foo  \n  bar  \n", "trim", "-l")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foo  ", "bar  "}, lines(stdout))
}

// Fish parity tests

func TestTrim_both_sides(t *testing.T) {
	exit, stdout, _ := runCmd("trim", " abc  ")
	assertExit(t, 0, exit)
	assertLines(t, []string{"abc"}, lines(stdout))
}

func TestTrim_right_custom_chars(t *testing.T) {
	exit, stdout, _ := runCmd("trim", "--right", "--chars=yz", "xyzzy", "zany")
	assertExit(t, 0, exit)
	assertLines(t, []string{"x", "zan"}, lines(stdout))
}

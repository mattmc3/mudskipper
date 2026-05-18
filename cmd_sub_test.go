package main

import (
	"testing"
)

func TestSub_no_args_returns_full_string(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "hello")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hello"}, lines(stdout))
}

func TestSub_start_from_position(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "-s", "2", "hello")
	assertExit(t, 0, exit)
	assertLines(t, []string{"ello"}, lines(stdout))
}

func TestSub_start_from_end(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "-s", "-1", "hello")
	assertExit(t, 0, exit)
	assertLines(t, []string{"o"}, lines(stdout))
}

func TestSub_end_position(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "-e", "3", "hello")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hel"}, lines(stdout))
}

func TestSub_end_from_end(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "-e", "-2", "hello")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hel"}, lines(stdout))
}

func TestSub_length_limits_output(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "-l", "3", "hello")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hel"}, lines(stdout))
}

func TestSub_start_and_length(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "-s", "2", "-l", "3", "hello")
	assertExit(t, 0, exit)
	assertLines(t, []string{"ell"}, lines(stdout))
}

func TestSub_end_and_length_mutually_exclusive(t *testing.T) {
	exit, _, stderr := runCmd("sub", "-e", "5", "-l", "2", "hello")
	assertExit(t, 1, exit)
	assertContains(t, "mutually exclusive", stderr)
}

func TestSub_start_beyond_length_returns_empty(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "-s", "99", "hello")
	assertExit(t, 1, exit)
	assertLines(t, []string{""}, lines(stdout))
}

func TestSub_quiet_suppresses_output(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "-q", "-l", "3", "hello")
	assertExit(t, 0, exit)
	assertEmpty(t, stdout)
}

func TestSub_multiple_strings(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "-l", "2", "hello", "world")
	assertExit(t, 0, exit)
	assertLines(t, []string{"he", "wo"}, lines(stdout))
}

func TestSub_help_shows_usage(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "--help")
	assertExit(t, 0, exit)
	assertContains(t, "Usage:", stdout)
}

func TestStdin_sub(t *testing.T) {
	exit, stdout, _ := runWithStdin("hello\nworld\n", "sub", "-l", "3")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hel", "wor"}, lines(stdout))
}

// Fish parity tests

func TestSub_start_i64_min_clamps_to_beginning(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "--start", "-9223372036854775808", "abc")
	assertExit(t, 0, exit)
	assertLines(t, []string{"abc"}, lines(stdout))
}

func TestSub_end_positive(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "--end=3", "abcde")
	assertExit(t, 0, exit)
	assertLines(t, []string{"abc"}, lines(stdout))
}

func TestSub_start_clamps_to_beginning(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "-s", "-100", "-e", "-2", "abcde")
	assertExit(t, 0, exit)
	assertLines(t, []string{"abc"}, lines(stdout))
}

func TestSub_start_zero_is_error(t *testing.T) {
	exit, _, stderr := runCmd("sub", "--start=0", "abc")
	assertExit(t, 1, exit)
	assertContains(t, "Invalid start value '0'", stderr)
}

func TestSub_length_from_start(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "--length", "2", "abcde")
	assertExit(t, 0, exit)
	assertLines(t, []string{"ab"}, lines(stdout))
}

func TestSub_negative_start_positive_end(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "-s", "-5", "-e", "2", "abcde")
	assertExit(t, 0, exit)
	assertLines(t, []string{"ab"}, lines(stdout))
}

func TestSub_start_and_length_fish(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "-s", "2", "-l", "2", "abcde")
	assertExit(t, 0, exit)
	assertLines(t, []string{"bc"}, lines(stdout))
}

func TestSub_start_and_negative_end(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "--start=2", "--end=-2", "abcde")
	assertExit(t, 0, exit)
	assertLines(t, []string{"bc"}, lines(stdout))
}

func TestSub_negative_length_is_error(t *testing.T) {
	exit, _, stderr := runCmd("sub", "--length=-1", "abcde")
	assertExit(t, 1, exit)
	assertContains(t, "Invalid length value '-1'", stderr)
}

func TestSub_negative_start_from_end(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "--start=-2", "abcde")
	assertExit(t, 0, exit)
	assertLines(t, []string{"de"}, lines(stdout))
}

func TestSub_end_negative(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "--end=-4", "abcde")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a"}, lines(stdout))
}

func TestSub_end_zero_is_error(t *testing.T) {
	exit, _, stderr := runCmd("sub", "--end=0", "abcde")
	assertExit(t, 1, exit)
	assertContains(t, "Invalid end value '0'", stderr)
}

func TestSub_negative_start_and_negative_end(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "-s", "-5", "-e", "-2", "abcdefgh")
	assertExit(t, 0, exit)
	assertLines(t, []string{"def"}, lines(stdout))
}

func TestSub_end_before_start_returns_empty(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "-s", "-50", "-e", "-100", "abcde")
	assertExit(t, 1, exit)
	assertLines(t, []string{""}, lines(stdout))
}

func TestSub_start_after_end_returns_empty(t *testing.T) {
	exit, stdout, _ := runCmd("sub", "-s", "2", "-e", "-5", "abcde")
	assertExit(t, 1, exit)
	assertLines(t, []string{""}, lines(stdout))
}

func TestSub_end_and_length_together_is_error(t *testing.T) {
	exit, _, stderr := runCmd("sub", "-s", "2", "-e", "-5", "-l", "3", "abcde")
	assertExit(t, 1, exit)
	assertContains(t, "--end and --length are mutually exclusive", stderr)
}

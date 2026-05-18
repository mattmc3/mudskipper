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

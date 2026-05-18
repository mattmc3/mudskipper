package main

import (
	"testing"
)

func TestReplace_literal_first_only(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "o", "0", "foobar")
	assertExit(t, 0, exit)
	assertLines(t, []string{"f0obar"}, lines(stdout))
}

func TestReplace_literal_all(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "-a", "o", "0", "foobar")
	assertExit(t, 0, exit)
	assertLines(t, []string{"f00bar"}, lines(stdout))
}

func TestReplace_no_match_returns_1(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "x", "y", "foobar")
	assertExit(t, 1, exit)
	assertLines(t, []string{"foobar"}, lines(stdout))
}

func TestReplace_ignore_case(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "-i", "FOO", "baz", "foobar")
	assertExit(t, 0, exit)
	assertLines(t, []string{"bazbar"}, lines(stdout))
}

func TestReplace_regex_basic(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "-r", "o+", "0", "foobar")
	assertExit(t, 0, exit)
	assertLines(t, []string{"f0bar"}, lines(stdout))
}

func TestReplace_regex_all(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "-r", "-a", "o", "0", "foobar")
	assertExit(t, 0, exit)
	assertLines(t, []string{"f00bar"}, lines(stdout))
}

func TestReplace_regex_backreference(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "-r", "(foo)(bar)", "$2$1", "foobar")
	assertExit(t, 0, exit)
	assertLines(t, []string{"barfoo"}, lines(stdout))
}

func TestReplace_regex_ignore_case(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "-r", "-i", "FOO", "baz", "foobar")
	assertExit(t, 0, exit)
	assertLines(t, []string{"bazbar"}, lines(stdout))
}

func TestReplace_filter_only_prints_changed(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "-f", "o", "0", "foobar", "hello", "baz")
	assertExit(t, 0, exit)
	assertLines(t, []string{"f0obar", "hell0"}, lines(stdout))
}

func TestReplace_max_matches_limits_replacements(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "-a", "-m", "2", "o", "0", "foooobar")
	assertExit(t, 0, exit)
	assertLines(t, []string{"f00oobar"}, lines(stdout))
}

func TestReplace_quiet_suppresses_output(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "-q", "o", "0", "foobar")
	assertExit(t, 0, exit)
	assertEmpty(t, stdout)
}

func TestReplace_multiple_strings(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "a", "x", "bar", "baz", "nope")
	assertExit(t, 0, exit)
	assertLines(t, []string{"bxr", "bxz", "nope"}, lines(stdout))
}

func TestReplace_missing_args_returns_error(t *testing.T) {
	exit, _, stderr := runCmd("replace", "only-pattern")
	assertExit(t, 1, exit)
	assertContains(t, "replacement", stderr)
}

func TestReplace_help_shows_usage(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "--help")
	assertExit(t, 0, exit)
	assertContains(t, "Usage:", stdout)
}

func TestStdin_replace(t *testing.T) {
	exit, stdout, _ := runWithStdin("foobar\nhello\n", "replace", "-a", "o", "0")
	assertExit(t, 0, exit)
	assertLines(t, []string{"f00bar", "hell0"}, lines(stdout))
}

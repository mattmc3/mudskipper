package main

import (
	"testing"
)

func TestLower_converts_to_lowercase(t *testing.T) {
	exit, stdout, _ := runCmd("lower", "Foo", "BAR", "baz")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foo", "bar", "baz"}, lines(stdout))
}

func TestLower_returns_1_when_nothing_changed(t *testing.T) {
	exit, _, _ := runCmd("lower", "foo", "bar")
	assertExit(t, 1, exit)
}

func TestLower_quiet_suppresses_output(t *testing.T) {
	exit, stdout, _ := runCmd("lower", "-q", "FOO")
	assertExit(t, 0, exit)
	assertEmpty(t, stdout)
}

func TestLower_help_shows_usage(t *testing.T) {
	exit, stdout, _ := runCmd("lower", "--help")
	assertExit(t, 0, exit)
	assertContains(t, "Usage:", stdout)
}

func TestUpper_converts_to_uppercase(t *testing.T) {
	exit, stdout, _ := runCmd("upper", "Foo", "bar", "BAZ")
	assertExit(t, 0, exit)
	assertLines(t, []string{"FOO", "BAR", "BAZ"}, lines(stdout))
}

func TestUpper_returns_1_when_nothing_changed(t *testing.T) {
	exit, _, _ := runCmd("upper", "FOO", "BAR")
	assertExit(t, 1, exit)
}

func TestUpper_quiet_suppresses_output(t *testing.T) {
	exit, stdout, _ := runCmd("upper", "-q", "foo")
	assertExit(t, 0, exit)
	assertEmpty(t, stdout)
}

func TestUpper_help_shows_usage(t *testing.T) {
	exit, stdout, _ := runCmd("upper", "--help")
	assertExit(t, 0, exit)
	assertContains(t, "Usage:", stdout)
}

func TestStdin_lower(t *testing.T) {
	exit, stdout, _ := runWithStdin("Foo\nBAR\nbaz\n", "lower")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foo", "bar", "baz"}, lines(stdout))
}

func TestStdin_upper(t *testing.T) {
	exit, stdout, _ := runWithStdin("Foo\nbar\nBAZ\n", "upper")
	assertExit(t, 0, exit)
	assertLines(t, []string{"FOO", "BAR", "BAZ"}, lines(stdout))
}

func TestStdin_args_take_priority_over_stdin(t *testing.T) {
	exit, stdout, _ := runWithStdin("IGNORED\n", "lower", "Foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foo"}, lines(stdout))
}

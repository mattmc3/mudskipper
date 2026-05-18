package main

import (
	"testing"
)

func TestJoin_basic(t *testing.T) {
	exit, stdout, _ := runCmd("join", ",", "a", "b", "c")
	assertExit(t, 0, exit)
	assertStr(t, "a,b,c\n", stdout)
}

func TestJoin_single_string_returns_1(t *testing.T) {
	exit, stdout, _ := runCmd("join", ",", "a")
	assertExit(t, 1, exit)
	assertStr(t, "a\n", stdout)
}

func TestJoin_no_strings_returns_1(t *testing.T) {
	exit, _, _ := runWithStdin("", "join", ",")
	assertExit(t, 1, exit)
}

func TestJoin_no_empty_filters_empty_strings(t *testing.T) {
	exit, stdout, _ := runCmd("join", "-n", "+", "a", "b", "", "c")
	assertExit(t, 0, exit)
	assertStr(t, "a+b+c\n", stdout)
}

func TestJoin_quiet_suppresses_output(t *testing.T) {
	exit, stdout, _ := runCmd("join", "-q", ",", "a", "b")
	assertExit(t, 0, exit)
	assertEmpty(t, stdout)
}

func TestJoin_empty_sep(t *testing.T) {
	exit, stdout, _ := runCmd("join", "", "a", "b", "c")
	assertExit(t, 0, exit)
	assertStr(t, "abc\n", stdout)
}

func TestJoin_missing_sep_returns_error(t *testing.T) {
	exit, _, stderr := runCmd("join")
	assertExit(t, 1, exit)
	assertContains(t, "separator", stderr)
}

func TestJoin_help_shows_usage(t *testing.T) {
	exit, stdout, _ := runCmd("join", "--help")
	assertExit(t, 0, exit)
	assertContains(t, "Usage:", stdout)
}

func TestStdin_join(t *testing.T) {
	exit, stdout, _ := runWithStdin("a\nb\nc\n", "join", "...")
	assertExit(t, 0, exit)
	assertStr(t, "a...b...c\n", stdout)
}

func TestJoin0_nul_separated_with_trailing_nul(t *testing.T) {
	exit, stdout, _ := runCmd("join0", "a", "b", "c")
	assertExit(t, 0, exit)
	assertStr(t, "a\x00b\x00c\x00", stdout)
}

func TestJoin0_roundtrips_with_split0(t *testing.T) {
	_, joined, _ := runCmd("join0", "a", "b", "c")
	exit, stdout, _ := runWithStdin(joined, "split0")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a", "b", "c"}, lines(stdout))
}

func TestJoin0_help_shows_usage(t *testing.T) {
	exit, stdout, _ := runCmd("join0", "--help")
	assertExit(t, 0, exit)
	assertContains(t, "Usage:", stdout)
}

// Fish parity tests

func TestJoin_seq_stdin_with_ellipsis(t *testing.T) {
	exit, stdout, _ := runWithStdin("1\n2\n3\n", "join", "...")
	assertExit(t, 0, exit)
	assertLines(t, []string{"1...2...3"}, lines(stdout))
}

func TestJoin_no_args_is_error(t *testing.T) {
	exit, _, stderr := runCmd("join")
	assertExit(t, 1, exit)
	assertContains(t, "join requires a separator", stderr)
}

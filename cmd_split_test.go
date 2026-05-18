package main

import (
	"testing"
)

func TestSplit_basic(t *testing.T) {
	exit, stdout, _ := runCmd("split", ".", "example.com")
	assertExit(t, 0, exit)
	assertLines(t, []string{"example", "com"}, lines(stdout))
}

func TestSplit_no_match_returns_1(t *testing.T) {
	exit, stdout, _ := runCmd("split", ".", "example")
	assertExit(t, 1, exit)
	assertLines(t, []string{"example"}, lines(stdout))
}

func TestSplit_multiple_strings(t *testing.T) {
	exit, stdout, _ := runCmd("split", ",", "a,b", "c,d")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a", "b", "c", "d"}, lines(stdout))
}

func TestSplit_max_limits_splits(t *testing.T) {
	exit, stdout, _ := runCmd("split", "-m", "1", "/", "/usr/local/bin")
	assertExit(t, 0, exit)
	assertLines(t, []string{"", "usr/local/bin"}, lines(stdout))
}

func TestSplit_right_splits_from_right(t *testing.T) {
	exit, stdout, _ := runCmd("split", "-r", "-m", "1", "/", "/usr/local/bin/fish")
	assertExit(t, 0, exit)
	assertLines(t, []string{"/usr/local/bin", "fish"}, lines(stdout))
}

func TestSplit_no_empty_filters_empty_parts(t *testing.T) {
	exit, stdout, _ := runCmd("split", "-n", ",", "a,,b")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a", "b"}, lines(stdout))
}

func TestSplit_empty_sep_chars(t *testing.T) {
	exit, stdout, _ := runCmd("split", "", "abc")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a", "b", "c"}, lines(stdout))
}

func TestSplit_fields_selects_fields(t *testing.T) {
	exit, stdout, _ := runCmd("split", "-f", "1,3", ",", "a,b,c")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a", "c"}, lines(stdout))
}

func TestSplit_fields_range(t *testing.T) {
	exit, stdout, _ := runCmd("split", "-f", "2-4", ",", "a,b,c,d,e")
	assertExit(t, 0, exit)
	assertLines(t, []string{"b", "c", "d"}, lines(stdout))
}

func TestSplit_fields_missing_returns_1(t *testing.T) {
	exit, _, _ := runCmd("split", "-f", "5", ",", "a,b,c")
	assertExit(t, 1, exit)
}

func TestSplit_fields_allow_empty_skips_missing(t *testing.T) {
	exit, stdout, _ := runCmd("split", "-f", "1,5", "-a", ",", "a,b,c")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a"}, lines(stdout))
}

func TestSplit_quiet_suppresses_output(t *testing.T) {
	exit, stdout, _ := runCmd("split", "-q", ".", "a.b")
	assertExit(t, 0, exit)
	assertEmpty(t, stdout)
}

func TestSplit_missing_sep_returns_error(t *testing.T) {
	exit, _, stderr := runCmd("split")
	assertExit(t, 1, exit)
	assertContains(t, "separator", stderr)
}

func TestSplit_help_shows_usage(t *testing.T) {
	exit, stdout, _ := runCmd("split", "--help")
	assertExit(t, 0, exit)
	assertContains(t, "Usage:", stdout)
}

func TestStdin_split(t *testing.T) {
	exit, stdout, _ := runWithStdin("a,b\nc,d\n", "split", ",")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a", "b", "c", "d"}, lines(stdout))
}

func TestSplit0_nul_separated(t *testing.T) {
	exit, stdout, _ := runWithStdin("a\x00b\x00c\x00", "split0")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a", "b", "c"}, lines(stdout))
}

func TestSplit0_no_trailing_empty_from_nul(t *testing.T) {
	exit, stdout, _ := runWithStdin("a\x00b\x00", "split0")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a", "b"}, lines(stdout))
}

func TestSplit0_help_shows_usage(t *testing.T) {
	exit, stdout, _ := runCmd("split0", "--help")
	assertExit(t, 0, exit)
	assertContains(t, "Usage:", stdout)
}

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

// Fish parity tests

func TestSplit_fields_out_of_range_returns_1(t *testing.T) {
	exit, _, _ := runCmd("split", "--fields=2,9", "", "abc")
	assertExit(t, 1, exit)
}

func TestSplit_fields_range_and_reverse_range(t *testing.T) {
	exit, stdout, _ := runCmd("split", "--fields=1-3,5,9-7", "", "123456789")
	assertExit(t, 0, exit)
	assertLines(t, []string{"1", "2", "3", "5", "9", "8", "7"}, lines(stdout))
}

func TestSplit_empty_delimiter_splits_chars(t *testing.T) {
	exit, stdout, _ := runCmd("split", "", "abc")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a", "b", "c"}, lines(stdout))
}

func TestSplit_f_shorthand_for_fields(t *testing.T) {
	exit, stdout, _ := runCmd("split", "-f1", " ", "a b", "c d")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a", "c"}, lines(stdout))
}

func TestSplit_fields_single(t *testing.T) {
	exit, stdout, _ := runCmd("split", "--fields=2", "", "abc")
	assertExit(t, 0, exit)
	assertLines(t, []string{"b"}, lines(stdout))
}

func TestSplit_fields_multiple_out_of_order(t *testing.T) {
	exit, stdout, _ := runCmd("split", "--fields=3,2", "", "abc")
	assertExit(t, 0, exit)
	assertLines(t, []string{"c", "b"}, lines(stdout))
}

func TestSplit_allow_empty_with_fields(t *testing.T) {
	exit, stdout, _ := runCmd("split", "--allow-empty", "--fields=2,9", "", "abc")
	assertExit(t, 0, exit)
	assertLines(t, []string{"b"}, lines(stdout))
}

func TestSplit_no_args_is_error(t *testing.T) {
	exit, _, stderr := runCmd("split")
	assertExit(t, 1, exit)
	assertContains(t, "split requires a separator", stderr)
}

func TestSplit_negative_max_is_error(t *testing.T) {
	exit, _, stderr := runCmd("split", "--max=-1", "12", "AB12CD")
	assertExit(t, 1, exit)
	assertContains(t, "Invalid max value '-1'", stderr)
}

func TestSplit_fields_invalid_spec_is_error(t *testing.T) {
	exit, _, stderr := runCmd("split", "--fields=2-3-,9", "", "a")
	assertExit(t, 1, exit)
	assertContains(t, "invalid field spec", stderr)
}

func TestSplit_fields_zero_is_error(t *testing.T) {
	exit, _, stderr := runCmd("split", "--fields=0", "", "c")
	assertExit(t, 1, exit)
	assertContains(t, "Invalid fields value '0'", stderr)
}

func TestSplit_fields_inverted_range_is_error(t *testing.T) {
	exit, _, stderr := runCmd("split", "--fields=1-0", "", "d")
	assertExit(t, 1, exit)
	assertContains(t, "Invalid range value for field '1-0'", stderr)
}

func TestSplit_fields_zero_start_range_is_error(t *testing.T) {
	exit, _, stderr := runCmd("split", "--fields=0-1", "", "e")
	assertExit(t, 1, exit)
	assertContains(t, "Invalid", stderr)
}

func TestSplit_fields_alpha_is_error(t *testing.T) {
	exit, _, stderr := runCmd("split", "--fields=a", "", "h")
	assertExit(t, 1, exit)
	assertContains(t, "invalid field spec", stderr)
}

func TestSplit_allow_empty_without_fields_is_error(t *testing.T) {
	exit, _, stderr := runCmd("split", "--allow-empty", "", "abc")
	assertExit(t, 1, exit)
	assertContains(t, "--allow-empty is only valid with --fields", stderr)
}

func TestSplit0_basic(t *testing.T) {
	exit, stdout, _ := runWithStdin("a\x00b", "split0")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a", "b"}, lines(stdout))
}

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

// Fish parity tests

func TestRepeat_quiet_returns_0(t *testing.T) {
	exit, stdout, _ := runCmd("repeat", "-n2", "-q", "foo")
	assertExit(t, 0, exit)
	assertEmpty(t, stdout)
}

func TestRepeat_zero_count_returns_1(t *testing.T) {
	exit, _, _ := runCmd("repeat", "-n0", "foo")
	assertExit(t, 1, exit)
}

func TestRepeat_zero_max_returns_1(t *testing.T) {
	exit, _, _ := runCmd("repeat", "-m0")
	assertExit(t, 1, exit)
}

func TestRepeat_n2(t *testing.T) {
	exit, stdout, _ := runCmd("repeat", "-n", "2", "foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foofoo"}, lines(stdout))
}

func TestRepeat_count_long_flag(t *testing.T) {
	exit, stdout, _ := runCmd("repeat", "--count", "2", "foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foofoo"}, lines(stdout))
}

func TestRepeat_positional_count(t *testing.T) {
	exit, stdout, _ := runCmd("repeat", "2", "foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foofoo"}, lines(stdout))
}

func TestRepeat_stdin_with_count(t *testing.T) {
	exit, stdout, _ := runWithStdin("foo\n", "repeat", "-n", "2")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foofoo"}, lines(stdout))
}

func TestRepeat_stdin_with_positional_count(t *testing.T) {
	exit, stdout, _ := runWithStdin("foo\n", "repeat", "2")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foofoo"}, lines(stdout))
}

func TestRepeat_no_args_is_error(t *testing.T) {
	exit, _, _ := runCmd("repeat")
	assertExit(t, 1, exit)
}

func TestRepeat_invalid_positional_count_is_error(t *testing.T) {
	exit, _, stderr := runCmd("repeat", "foo")
	assertExit(t, 1, exit)
	assertContains(t, "Invalid count value 'foo'", stderr)
}

func TestRepeat_no_newline_flag(t *testing.T) {
	exit, stdout, _ := runCmd("repeat", "-n1", "-N", "there is ")
	assertExit(t, 0, exit)
	assertStr(t, "there is ", stdout)
}

func TestRepeat_max_limits_output(t *testing.T) {
	exit, stdout, _ := runCmd("repeat", "-n10", "-m4", "foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foof"}, lines(stdout))
}

func TestRepeat_max_only_no_count(t *testing.T) {
	exit, stdout, _ := runCmd("repeat", "-m4", "foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foof"}, lines(stdout))
}

func TestRepeat_max_5_from_count_10(t *testing.T) {
	exit, stdout, _ := runCmd("repeat", "-n10", "--max", "5", "foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foofo"}, lines(stdout))
}

func TestRepeat_max_larger_than_result(t *testing.T) {
	exit, stdout, _ := runCmd("repeat", "-n3", "-m20", "foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foofoofoo"}, lines(stdout))
}

func TestRepeat_multiple_strings_fish(t *testing.T) {
	exit, stdout, _ := runCmd("repeat", "-n", "5", "a", "b", "c")
	assertExit(t, 0, exit)
	assertLines(t, []string{"aaaaa", "bbbbb", "ccccc"}, lines(stdout))
}

func TestRepeat_multiple_strings_with_max(t *testing.T) {
	exit, stdout, _ := runCmd("repeat", "-n", "5", "--max", "4", "123", "456", "789")
	assertExit(t, 0, exit)
	assertLines(t, []string{"1231", "4564", "7897"}, lines(stdout))
}

func TestRepeat_multiple_strings_with_empty_and_max(t *testing.T) {
	exit, stdout, _ := runCmd("repeat", "-n", "5", "--max", "4", "123", "", "789")
	assertExit(t, 0, exit)
	assertLines(t, []string{"1231", "", "7897"}, lines(stdout))
}

func TestRepeat_empty_string_returns_1(t *testing.T) {
	exit, _, _ := runCmd("repeat", "-n3", "")
	assertExit(t, 1, exit)
}

func TestRepeat_negative_count_is_error(t *testing.T) {
	exit, _, stderr := runCmd("repeat", "-n-1", "foo")
	assertExit(t, 1, exit)
	assertContains(t, "Invalid count value '-1'", stderr)
}

func TestRepeat_max_negative_is_error(t *testing.T) {
	exit, _, stderr := runCmd("repeat", "-m-1", "foo")
	assertExit(t, 1, exit)
	assertContains(t, "Invalid max value '-1'", stderr)
}

func TestRepeat_non_integer_count_is_error(t *testing.T) {
	exit, _, stderr := runCmd("repeat", "-n", "notanumber", "foo")
	assertExit(t, 1, exit)
	assertContains(t, "Invalid count value 'notanumber'", stderr)
}

func TestRepeat_max_invalid_is_error(t *testing.T) {
	exit, _, stderr := runCmd("repeat", "-m", "notanumber", "foo")
	assertExit(t, 1, exit)
	assertContains(t, "Invalid max value 'notanumber'", stderr)
}

func TestRepeat_n_missing_arg_is_error(t *testing.T) {
	exit, _, stderr := runCmd("repeat", "-n")
	assertExit(t, 1, exit)
	assertContains(t, "-n", stderr)
}

func TestRepeat_max_produces_exact_length(t *testing.T) {
	_, r, _ := runCmd("repeat", "-m", "5000", "aab")
	exit, stdout, _ := runCmd("length", r[:len(r)-1])
	assertExit(t, 0, exit)
	assertLines(t, []string{"5000"}, lines(stdout))
}

func TestRepeat_m17_aab_length_is_17(t *testing.T) {
	_, r, _ := runCmd("repeat", "-m", "17", "aab")
	exit, stdout, _ := runCmd("length", r[:len(r)-1])
	assertExit(t, 0, exit)
	assertLines(t, []string{"17"}, lines(stdout))
}

func TestRepeat_count_produces_correct_length(t *testing.T) {
	_, r, _ := runCmd("repeat", "-n", "17", "aab")
	exit, stdout, _ := runCmd("length", r[:len(r)-1])
	assertExit(t, 0, exit)
	assertLines(t, []string{"51"}, lines(stdout))
}

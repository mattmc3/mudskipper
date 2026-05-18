package main

import (
	"testing"
)

func TestMatch_glob_exact(t *testing.T) {
	exit, stdout, _ := runCmd("match", "foo", "foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foo"}, lines(stdout))
}

func TestMatch_glob_star(t *testing.T) {
	exit, stdout, _ := runCmd("match", "foo*", "foobar")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foobar"}, lines(stdout))
}

func TestMatch_glob_question(t *testing.T) {
	exit, stdout, _ := runCmd("match", "f?o", "foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foo"}, lines(stdout))
}

func TestMatch_glob_char_class(t *testing.T) {
	exit, stdout, _ := runCmd("match", "[abc]oo", "boo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"boo"}, lines(stdout))
}

func TestMatch_glob_negated_class(t *testing.T) {
	exit, stdout, _ := runCmd("match", "[!abc]oo", "zoo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"zoo"}, lines(stdout))
}

func TestMatch_glob_no_match_returns_1(t *testing.T) {
	exit, _, _ := runCmd("match", "foo", "bar")
	assertExit(t, 1, exit)
}

func TestMatch_glob_ignore_case(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-i", "FOO*", "foobar")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foobar"}, lines(stdout))
}

func TestMatch_regex_basic(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-r", "fo+", "foobar")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foo"}, lines(stdout))
}

func TestMatch_regex_groups_only(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-r", "-g", "f(o+)(bar)", "foobar")
	assertExit(t, 0, exit)
	assertLines(t, []string{"oo", "bar"}, lines(stdout))
}

func TestMatch_regex_groups_only_no_groups_returns_1(t *testing.T) {
	exit, _, _ := runCmd("match", "-r", "-g", "foobar", "foobar")
	assertExit(t, 1, exit)
}

func TestMatch_all_finds_multiple(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-r", "-a", "o+", "foobaroo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"oo", "oo"}, lines(stdout))
}

func TestMatch_entire_prints_full_string(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-r", "-e", "fo+", "foobar")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foobar"}, lines(stdout))
}

func TestMatch_index_prints_position(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-r", "-n", "fo+", "xfoobar")
	assertExit(t, 0, exit)
	assertLines(t, []string{"2 3"}, lines(stdout))
}

func TestMatch_invert_prints_non_matching(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-v", "foo", "foo", "bar", "baz")
	assertExit(t, 0, exit)
	assertLines(t, []string{"bar", "baz"}, lines(stdout))
}

func TestMatch_quiet_suppresses_output(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-q", "foo*", "foobar")
	assertExit(t, 0, exit)
	assertEmpty(t, stdout)
}

func TestMatch_max_matches_limits_results(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-r", "-a", "-m", "2", "o", "foooo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"o", "o"}, lines(stdout))
}

func TestMatch_missing_pattern_returns_error(t *testing.T) {
	exit, _, stderr := runCmd("match")
	assertExit(t, 1, exit)
	assertContains(t, "pattern", stderr)
}

func TestMatch_help_shows_usage(t *testing.T) {
	exit, stdout, _ := runCmd("match", "--help")
	assertExit(t, 0, exit)
	assertContains(t, "Usage:", stdout)
}

func TestStdin_match_glob(t *testing.T) {
	exit, stdout, _ := runWithStdin("foobar\nbaz\nfoo\n", "match", "foo*")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foobar", "foo"}, lines(stdout))
}

func TestStdin_match_regex(t *testing.T) {
	exit, stdout, _ := runWithStdin("foobar\nbaz\n", "match", "-r", "fo+")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foo"}, lines(stdout))
}

// Fish parity tests

func TestMatch_regex_invert_filters_matching_strings(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-r", "-v", "c.*", "dog", "can", "cat", "diz")
	assertExit(t, 0, exit)
	assertLines(t, []string{"dog", "diz"}, lines(stdout))
}

func TestMatch_glob_invert_filters_matching_strings(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-v", "c*", "dog", "can", "cat", "diz")
	assertExit(t, 0, exit)
	assertLines(t, []string{"dog", "diz"}, lines(stdout))
}

func TestMatch_regex_alternation(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-r", "cat|dog|fish", "nice dog")
	assertExit(t, 0, exit)
	assertLines(t, []string{"dog"}, lines(stdout))
}

func TestMatch_regex_invert_max1_via_stdin(t *testing.T) {
	exit, stdout, _ := runWithStdin("dog\ncat\nbat\nhog\n", "match", "-r", "-v", "-m1", "at$")
	assertExit(t, 0, exit)
	assertLines(t, []string{"dog"}, lines(stdout))
}

func TestMatch_quiet_regex_invert_returns_0_when_any_match(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-q", "-r", "-v", "c.*", "dog", "can", "cat", "diz")
	assertExit(t, 0, exit)
	assertEmpty(t, stdout)
}

func TestMatch_quiet_glob_invert_returns_0_when_any_match(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-q", "-v", "c*", "dog", "can", "cat", "diz")
	assertExit(t, 0, exit)
	assertEmpty(t, stdout)
}

func TestMatch_regex_invert_single_nonmatch_returns_0(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-r", "-v", "x", "y")
	assertExit(t, 0, exit)
	assertLines(t, []string{"y"}, lines(stdout))
}

func TestMatch_quiet_regex_invert_nonmatch_returns_0(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-q", "-r", "-v", "x", "y")
	assertExit(t, 0, exit)
	assertEmpty(t, stdout)
}

func TestMatch_glob_invert_returns_1_when_all_match(t *testing.T) {
	exit, _, _ := runCmd("match", "-v", "d*", "dog", "dan", "dat", "diz")
	assertExit(t, 1, exit)
}

func TestMatch_quiet_glob_invert_returns_1_when_all_match(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-q", "-v", "d*", "dog", "dan", "dat", "diz")
	assertExit(t, 1, exit)
	assertEmpty(t, stdout)
}

func TestMatch_regex_invert_single_match_returns_1(t *testing.T) {
	exit, _, _ := runCmd("match", "-r", "-v", "x", "x")
	assertExit(t, 1, exit)
}

func TestMatch_quiet_regex_invert_match_returns_1(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-q", "-r", "-v", "x", "x")
	assertExit(t, 1, exit)
	assertEmpty(t, stdout)
}

func TestMatch_invert_and_groups_only_is_error(t *testing.T) {
	exit, _, stderr := runCmd("match", "-v", "-g", "foo", "foo")
	assertExit(t, 1, exit)
	assertContains(t, "--invert and --groups-only", stderr)
}

func TestMatch_no_args_is_error(t *testing.T) {
	exit, _, stderr := runCmd("match")
	assertExit(t, 1, exit)
	assertContains(t, "match requires a pattern", stderr)
}

func TestMatch_groups_only_all_matches_stdin(t *testing.T) {
	exit, stdout, _ := runWithStdin("foo1x foo2x foo3x\n", "match", "-arg", `foo(\d)x`)
	assertExit(t, 0, exit)
	assertLines(t, []string{"1", "2", "3"}, lines(stdout))
}

func TestMatch_groups_only_no_match(t *testing.T) {
	exit, _, _ := runCmd("match", "-r", "--groups-only", "(.+)fish", "fish")
	assertExit(t, 1, exit)
}

func TestMatch_regex_capture_groups(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-r", `(\d\d?):(\d\d):(\d\d)`, "2:34:56")
	assertExit(t, 0, exit)
	assertLines(t, []string{"2:34:56", "2", "34", "56"}, lines(stdout))
}

func TestMatch_groups_only_multiple_groups(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-r", "--groups-only", "(.+)fish(.*)", "catfishcolor")
	assertExit(t, 0, exit)
	assertLines(t, []string{"cat", "color"}, lines(stdout))
}

func TestMatch_max_matches_limits_glob_results(t *testing.T) {
	exit, stdout, _ := runWithStdin("dog\ncat\nbat\ngnat\n", "match", "-m2", "*at")
	assertExit(t, 0, exit)
	assertLines(t, []string{"cat", "bat"}, lines(stdout))
}

func TestMatch_groups_only_shellfish(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-rg", "(.*)fish", "shellfish")
	assertExit(t, 0, exit)
	assertLines(t, []string{"shell"}, lines(stdout))
}

func TestMatch_groups_only_fish_empty_capture(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-rg", "(.*)fish", "fish")
	assertExit(t, 0, exit)
	assertLines(t, []string{""}, lines(stdout))
}

func TestMatch_groups_only_banana_no_match(t *testing.T) {
	exit, _, _ := runCmd("match", "-rg", "(.*)fish", "banana")
	assertExit(t, 1, exit)
}

func TestMatch_groups_only_stdin(t *testing.T) {
	exit, stdout, _ := runWithStdin("foo bar baz\n", "match", "-rg", "foo (bar) baz")
	assertExit(t, 0, exit)
	assertLines(t, []string{"bar"}, lines(stdout))
}

func TestMatch_max_zero_is_error(t *testing.T) {
	exit, _, stderr := runCmd("match", "-m0", "foo")
	assertExit(t, 1, exit)
	assertContains(t, "Invalid max matches value '0'", stderr)
}

func TestMatch_max_matches_overflow_is_error(t *testing.T) {
	exit, _, stderr := runCmd("match", "-m999999999999999999999999999999999999999", "foo")
	assertExit(t, 1, exit)
	assertContains(t, "Invalid max matches value", stderr)
}

func TestMatch_entire_regex_with_group(t *testing.T) {
	exit, stdout, _ := runCmd("match", "--entire", "-r", "a*b([xy]+)", "abxc")
	assertExit(t, 0, exit)
	assertLines(t, []string{"abxc", "x"}, lines(stdout))
}

func TestMatch_entire_and_index_is_error(t *testing.T) {
	exit, _, stderr := runCmd("match", "--entire", "--index", "foo", "foo")
	assertExit(t, 1, exit)
	assertContains(t, "--entire and --index", stderr)
}

func TestMatch_regex_with_group_outputs_match_and_group(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-r", "a*b([xy]+)", "abxc")
	assertExit(t, 0, exit)
	assertLines(t, []string{"abx", "x"}, lines(stdout))
}

func TestMatch_quiet_entire_regex_returns_0(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-q", "-e", "-r", "asd", "asd")
	assertExit(t, 0, exit)
	assertEmpty(t, stdout)
}

func TestMatch_entire_quiet_returns_0(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-e", "-q", "asd", "asd")
	assertExit(t, 0, exit)
	assertEmpty(t, stdout)
}

func TestMatch_glob_star_matches_any(t *testing.T) {
	exit, stdout, _ := runCmd("match", "*", "a")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a"}, lines(stdout))
}

func TestMatch_glob_with_wildcard(t *testing.T) {
	exit, stdout, _ := runCmd("match", "a*b", "axxb")
	assertExit(t, 0, exit)
	assertLines(t, []string{"axxb"}, lines(stdout))
}

func TestMatch_glob_case_insensitive(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-i", "a**B", "Axxb")
	assertExit(t, 0, exit)
	assertLines(t, []string{"Axxb"}, lines(stdout))
}

func TestMatch_glob_question_mark(t *testing.T) {
	exit, stdout, _ := runWithStdin("ok?\n", "match", "*?")
	assertExit(t, 0, exit)
	assertLines(t, []string{"ok?"}, lines(stdout))
}

func TestMatch_entire_contains_substring(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-e", "x", "abc", "dxf", "xyz", "jkx", "x", "z")
	assertExit(t, 0, exit)
	assertLines(t, []string{"dxf", "xyz", "jkx", "x"}, lines(stdout))
}

func TestMatch_glob_exact_only(t *testing.T) {
	exit, stdout, _ := runCmd("match", "x", "abc", "dxf", "xyz", "jkx", "x", "z")
	assertExit(t, 0, exit)
	assertLines(t, []string{"x"}, lines(stdout))
}

func TestMatch_entire_regex_with_capture_groups(t *testing.T) {
	exit, stdout, _ := runCmd("match", "--entire", "-r", "a*b([xy]+)", "abc", "abxc", "bye", "aaabyz", "kaabxz", "abbxy", "abcx", "caabxyxz")
	assertExit(t, 0, exit)
	assertLines(t, []string{"abxc", "x", "bye", "y", "aaabyz", "y", "kaabxz", "x", "abbxy", "xy", "caabxyxz", "xyx"}, lines(stdout))
}

func TestMatch_regex_with_capture_groups_output(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-r", "a*b([xy]+)", "abc", "abxc", "bye", "aaabyz", "kaabxz", "abbxy", "abcx", "caabxyxz")
	assertExit(t, 0, exit)
	assertLines(t, []string{"abx", "x", "by", "y", "aaaby", "y", "aabx", "x", "bxy", "xy", "aabxyx", "xyx"}, lines(stdout))
}

func TestMatch_entire_regex_returns_full_string(t *testing.T) {
	exit, stdout, _ := runCmd("match", "--entire", "-r", "a*b[xy]+", "abc", "abxc", "bye", "aaabyz", "kaabxz", "abbxy", "abcx", "caabxyxz")
	assertExit(t, 0, exit)
	assertLines(t, []string{"abxc", "bye", "aaabyz", "kaabxz", "abbxy", "caabxyxz"}, lines(stdout))
}

func TestMatch_regex_returns_matched_portion(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-r", "a*b[xy]+", "abc", "abxc", "bye", "aaabyz", "kaabxz", "abbxy", "abcx", "caabxyxz")
	assertExit(t, 0, exit)
	assertLines(t, []string{"abx", "by", "aaaby", "aabx", "bxy", "aabxyx"}, lines(stdout))
}

func TestMatch_entire_empty_pattern_matches_all(t *testing.T) {
	exit, stdout, _ := runCmd("match", "--entire", "", "banana")
	assertExit(t, 0, exit)
	assertLines(t, []string{"banana"}, lines(stdout))
}

func TestMatch_regex_all_with_index(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-r", "-a", "-n", "at", "ratatat")
	assertExit(t, 0, exit)
	assertLines(t, []string{"2 2", "4 2", "6 2"}, lines(stdout))
}

func TestMatch_regex_case_insensitive_hex(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-r", "-i", "0x[0-9a-f]{1,8}", "int magic = 0xBadC0de;")
	assertExit(t, 0, exit)
	assertLines(t, []string{"0xBadC0de"}, lines(stdout))
}

func TestMatch_regex_invert_all_match_returns_1(t *testing.T) {
	exit, _, _ := runCmd("match", "-r", "-v", "[dcantg].*", "dog", "can", "cat", "diz")
	assertExit(t, 1, exit)
}

func TestMatch_glob_invert_star_matches_all_returns_1(t *testing.T) {
	exit, _, _ := runCmd("match", "-v", "*", "dog", "can", "cat", "diz")
	assertExit(t, 1, exit)
}

func TestMatch_regex_invert_with_index(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-r", "-v", "-n", "a", "bbb")
	assertExit(t, 0, exit)
	assertLines(t, []string{"1 3"}, lines(stdout))
}

func TestMatch_regex_empty_capture_groups(t *testing.T) {
	exit, stdout, _ := runCmd("match", "-r", `^([ugoa]*)([=+-]?)([rwx]*)$`, "=r")
	assertExit(t, 0, exit)
	assertLines(t, []string{"=r", "", "=", "r"}, lines(stdout))
}

func TestMatch_regex_compile_error(t *testing.T) {
	exit, _, stderr := runCmd("match", "-r", "[", "a[sd")
	assertExit(t, 1, exit)
	assertContains(t, "error:", stderr)
}

func TestMatch_utf_mode_regex_is_error(t *testing.T) {
	exit, _, stderr := runCmd("match", "-r", "(*UTF).*", "aaa")
	assertExit(t, 1, exit)
	assertContains(t, "error:", stderr)
}

func TestMatch_unknown_option_is_error(t *testing.T) {
	exit, _, stderr := runCmd("match", "--unknown-opt")
	assertExit(t, 1, exit)
	assertContains(t, "not defined", stderr)
}

func TestMatch_regex_with_equals_is_error(t *testing.T) {
	exit, _, stderr := runCmd("match", "--regex=abc")
	assertExit(t, 1, exit)
	assertContains(t, "parse error", stderr)
}

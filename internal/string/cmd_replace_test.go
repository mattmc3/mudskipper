package stringcmd

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

// Fish parity tests

func TestReplace_regex_max_matches_1_no_filter(t *testing.T) {
	exit, stdout, _ := runWithStdin("dog\ncat\nbat\n", "replace", "-r", "--max-matches", "1", "^c", "h")
	assertExit(t, 0, exit)
	assertLines(t, []string{"dog", "hat", "bat"}, lines(stdout))
}

func TestReplace_literal_basic(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "is", "was", "blue is my favorite")
	assertExit(t, 0, exit)
	assertLines(t, []string{"blue was my favorite"}, lines(stdout))
}

func TestReplace_literal_multiple_strings(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "3rd", "last", "1st", "2nd", "3rd")
	assertExit(t, 0, exit)
	assertLines(t, []string{"1st", "2nd", "last"}, lines(stdout))
}

func TestReplace_all_spaces_with_underscore(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "-a", " ", "_", "spaces to underscores")
	assertExit(t, 0, exit)
	assertLines(t, []string{"spaces_to_underscores"}, lines(stdout))
}

func TestReplace_regex_all_non_numeric(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "-r", "-a", `[^\d.]+`, " ", "0 one two 3.14 four 5x")
	assertExit(t, 0, exit)
	assertLines(t, []string{"0 3.14 5 "}, lines(stdout))
}

func TestReplace_regex_backreference_and_literal_dollar(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "-r", `(\w+)\s+(\w+)`, "$2 $1 $$", "left right")
	assertExit(t, 0, exit)
	assertLines(t, []string{"right left $"}, lines(stdout))
}

func TestReplace_regex_insert_newline(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "-r", `\s*newline\s*`, "\n", "put a newline here")
	assertExit(t, 0, exit)
	assertLines(t, []string{"put a", "here"}, lines(stdout))
}

func TestReplace_regex_all_double_chars(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "-r", "-a", `(\w)`, "$1$1", "ab")
	assertExit(t, 0, exit)
	assertLines(t, []string{"aabb"}, lines(stdout))
}

func TestReplace_quiet_no_match_returns_1(t *testing.T) {
	exit, stdout, _ := runWithStdin("a\n", "replace", "-q", "b", "c")
	assertExit(t, 1, exit)
	assertEmpty(t, stdout)
}

func TestReplace_regex_quiet_no_match_returns_1(t *testing.T) {
	exit, stdout, _ := runWithStdin("a\n", "replace", "-rq", "b", "c")
	assertExit(t, 1, exit)
	assertEmpty(t, stdout)
}

func TestReplace_filter_prints_only_changed(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "--filter", "x", "X", "abc", "axc", "x", "def", "jkx")
	assertExit(t, 0, exit)
	assertLines(t, []string{"aXc", "X", "jkX"}, lines(stdout))
}

func TestReplace_regex_filter_prints_only_changed(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "--regex", "-f", `\d`, "X", "1bc", "axc", "2", "d3f", "jk4", "xyz")
	assertExit(t, 0, exit)
	assertLines(t, []string{"Xbc", "X", "dXf", "jkX"}, lines(stdout))
}

func TestReplace_filter_no_match_returns_1(t *testing.T) {
	exit, _, _ := runCmd("replace", "--filter", "y", "Y", "abc", "axc", "x", "def", "jkx")
	assertExit(t, 1, exit)
}

func TestReplace_regex_filter_no_match_returns_1(t *testing.T) {
	exit, _, _ := runCmd("replace", "--regex", "-f", "Z", "X", "1bc", "axc", "2", "d3f", "jk4", "xyz")
	assertExit(t, 1, exit)
}

func TestReplace_utf_mode_regex_is_error(t *testing.T) {
	exit, _, stderr := runCmd("replace", "-r", "(*UTF).*", "replacement", "aaa")
	assertExit(t, 1, exit)
	assertContains(t, "error:", stderr)
}

func TestReplace_regex_unmatched_group_is_empty(t *testing.T) {
	exit, stdout, _ := runCmd("replace", "-r", `a(b.+)?z`, `a:$1z`, "az")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a:z"}, lines(stdout))
}

func TestReplace_regex_all_chars_with_asterisk(t *testing.T) {
	exit, stdout, _ := runWithStdin("my-password", "replace", "-ra", ".", "*")
	assertExit(t, 0, exit)
	assertLines(t, []string{"***********"}, lines(stdout))
}

func TestReplace_max_matches_non_integer_is_error(t *testing.T) {
	exit, _, stderr := runCmd("replace", "--max-matches", "abc")
	assertExit(t, 1, exit)
	assertContains(t, "Invalid max matches value 'abc'", stderr)
}

func TestReplace_max_matches_negative_is_error(t *testing.T) {
	exit, _, stderr := runCmd("replace", "--max-matches", "-1")
	assertExit(t, 1, exit)
	assertContains(t, "Invalid max matches value '-1'", stderr)
}

func TestReplace_max_matches_overflow_is_error(t *testing.T) {
	exit, _, stderr := runCmd("replace", "--max-matches", "99999999999999999999")
	assertExit(t, 1, exit)
	assertContains(t, "Invalid max matches value", stderr)
}

// Test for $1z template ambiguity — should fail before fix, pass after

func TestReplace_regex_group_ref_followed_by_alphanum(t *testing.T) {
	// Go ExpandString reads $1z as group "1z"; should treat as group 1 + literal "z"
	exit, stdout, _ := runCmd("replace", "-r", `a(b.+)?z`, `a:$1z`, "az")
	assertExit(t, 0, exit)
	assertLines(t, []string{"a:z"}, lines(stdout))
}

func TestReplace_no_args_is_error(t *testing.T) {
	exit, _, stderr := runCmd("replace")
	assertExit(t, 1, exit)
	assertContains(t, "replace requires", stderr)
}

func TestReplace_one_arg_is_error(t *testing.T) {
	exit, _, stderr := runCmd("replace", "one")
	assertExit(t, 1, exit)
	assertContains(t, "replace requires", stderr)
}

package stringcmd

import (
	"strings"
	"testing"
)

func TestString_help(t *testing.T) {
	for _, flag := range []string{"-h", "--help", "help"} {
		code, out, _ := runCmd(flag)
		if code != 0 {
			t.Errorf("%s: exit want 0, got %d", flag, code)
		}
		if !strings.Contains(out, "Usage:") {
			t.Errorf("%s: stdout want Usage, got %q", flag, out)
		}
	}
}

func TestCollect_joins_args_with_newlines(t *testing.T) {
	exit, stdout, _ := runCmd("collect", "a", "b", "c")
	assertExit(t, 0, exit)
	assertStr(t, "a\nb\nc\n", stdout)
}

func TestCollect_stdin_trims_trailing_newline(t *testing.T) {
	exit, stdout, _ := runWithStdin("a\nb\nc\n", "collect")
	assertExit(t, 0, exit)
	assertStr(t, "a\nb\nc\n", stdout)
}

func TestCollect_no_trim_preserves_trailing_newlines(t *testing.T) {
	exit, stdout, _ := runWithStdin("a\nb\nc\n", "collect", "-N")
	assertExit(t, 0, exit)
	assertStr(t, "a\nb\nc\n", stdout)
}

func TestCollect_empty_returns_1(t *testing.T) {
	exit, stdout, _ := runWithStdin("", "collect")
	assertExit(t, 1, exit)
	assertEmpty(t, stdout)
}

func TestCollect_empty_allow_empty_returns_0(t *testing.T) {
	exit, _, _ := runWithStdin("", "collect", "-a")
	assertExit(t, 0, exit)
}

func TestCollect_help_shows_usage(t *testing.T) {
	exit, stdout, _ := runCmd("collect", "--help")
	assertExit(t, 0, exit)
	assertContains(t, "Usage:", stdout)
}

// Fish parity tests

func TestCollect_no_args_returns_1(t *testing.T) {
	exit, _, _ := runCmd("collect")
	assertExit(t, 1, exit)
}

func TestCollect_single_arg_returns_0(t *testing.T) {
	exit, _, _ := runCmd("collect", "a")
	assertExit(t, 0, exit)
}

func TestCollect_no_trim_empty_string_returns_1(t *testing.T) {
	exit, _, _ := runCmd("collect", "-N", "")
	assertExit(t, 1, exit)
}

func TestCollect_only_newlines_returns_1(t *testing.T) {
	exit, _, _ := runCmd("collect", "\n\n")
	assertExit(t, 1, exit)
}

func TestCollect_N_with_stdin_newline_returns_0(t *testing.T) {
	exit, _, _ := runWithStdin("\n", "collect", "-N")
	assertExit(t, 0, exit)
}

func TestCollect_trim_stdin_only_newline_returns_1(t *testing.T) {
	exit, _, _ := runWithStdin("\n", "collect")
	assertExit(t, 1, exit)
}

func TestCollect_allow_empty_no_args_returns_0(t *testing.T) {
	exit, _, _ := runCmd("collect", "--allow-empty")
	assertExit(t, 0, exit)
}

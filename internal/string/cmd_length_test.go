package stringcmd

import (
	"testing"
)

func TestLength_single_string(t *testing.T) {
	exit, stdout, _ := runCmd("length", "hello")
	assertExit(t, 0, exit)
	assertLines(t, []string{"5"}, lines(stdout))
}

func TestLength_multiple_strings(t *testing.T) {
	exit, stdout, _ := runCmd("length", "foo", "hello", "ab")
	assertExit(t, 0, exit)
	assertLines(t, []string{"3", "5", "2"}, lines(stdout))
}

func TestLength_empty_string_returns_1(t *testing.T) {
	exit, stdout, _ := runCmd("length", "")
	assertExit(t, 1, exit)
	assertLines(t, []string{"0"}, lines(stdout))
}

func TestLength_quiet_suppresses_output(t *testing.T) {
	exit, stdout, _ := runCmd("length", "-q", "hello")
	assertExit(t, 0, exit)
	assertEmpty(t, stdout)
}

func TestLength_quiet_empty_returns_1(t *testing.T) {
	exit, _, _ := runCmd("length", "-q", "")
	assertExit(t, 1, exit)
}

func TestLength_visible_strips_ansi(t *testing.T) {
	exit, stdout, _ := runCmd("length", "-V", "\x1b[31mhello\x1b[0m")
	assertExit(t, 0, exit)
	assertLines(t, []string{"5"}, lines(stdout))
}

func TestLength_visible_no_ansi_same_as_normal(t *testing.T) {
	exit, stdout, _ := runCmd("length", "-V", "hello")
	assertExit(t, 0, exit)
	assertLines(t, []string{"5"}, lines(stdout))
}

func TestLength_help_shows_usage(t *testing.T) {
	exit, stdout, _ := runCmd("length", "--help")
	assertExit(t, 0, exit)
	assertContains(t, "Usage:", stdout)
}

func TestStdin_length(t *testing.T) {
	exit, stdout, _ := runWithStdin("foo\nhello\n", "length")
	assertExit(t, 0, exit)
	assertLines(t, []string{"3", "5"}, lines(stdout))
}

// Fish parity tests

func TestLength_hello_world(t *testing.T) {
	exit, stdout, _ := runCmd("length", "hello, world")
	assertExit(t, 0, exit)
	assertLines(t, []string{"12"}, lines(stdout))
}

func TestLength_quiet_empty_string_returns_1(t *testing.T) {
	exit, stdout, _ := runCmd("length", "-q", "")
	assertExit(t, 1, exit)
	assertEmpty(t, stdout)
}

func TestLength_no_args_returns_1(t *testing.T) {
	exit, _, _ := runCmd("length")
	assertExit(t, 1, exit)
}

func TestLength_visible_ignores_color(t *testing.T) {
	exit, stdout, _ := runCmd("length", "--visible", "\x1b[31mabc")
	assertExit(t, 0, exit)
	assertLines(t, []string{"3"}, lines(stdout))
}

func TestLength_visible_multiline(t *testing.T) {
	exit, stdout, _ := runCmd("length", "--visible", "a\x1b[34mb\ncde")
	assertExit(t, 0, exit)
	assertLines(t, []string{"2", "3"}, lines(stdout))
}

func TestLength_visible_carriage_return(t *testing.T) {
	exit, stdout, _ := runCmd("length", "--visible", "\x1b[0mabcdef\rfooba\x1b[31mraaa")
	assertExit(t, 0, exit)
	assertLines(t, []string{"9"}, lines(stdout))
}

func TestLength_visible_backspace_then_char(t *testing.T) {
	exit, stdout, _ := runCmd("length", "--visible", "f")
	assertExit(t, 0, exit)
	assertLines(t, []string{"1"}, lines(stdout))
}

func TestLength_visible_empty_returns_1(t *testing.T) {
	exit, stdout, _ := runCmd("length", "--visible", "")
	assertExit(t, 1, exit)
	assertLines(t, []string{"0"}, lines(stdout))
}

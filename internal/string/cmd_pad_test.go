package stringcmd

import (
	"testing"
)

func TestPad_default_pads_left(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "-w", "5", "hi")
	assertExit(t, 0, exit)
	assertLines(t, []string{"   hi"}, lines(stdout))
}

func TestPad_right_pads_right(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "-r", "-w", "5", "hi")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hi   "}, lines(stdout))
}

func TestPad_center_centers_string(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "-C", "-w", "7", "hi")
	assertExit(t, 0, exit)
	assertLines(t, []string{"   hi  "}, lines(stdout))
}

func TestPad_custom_char(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "-w", "5", "-c", ".", "hi")
	assertExit(t, 0, exit)
	assertLines(t, []string{"...hi"}, lines(stdout))
}

func TestPad_auto_width_uses_longest(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "hi", "hello")
	assertExit(t, 0, exit)
	assertLines(t, []string{"   hi", "hello"}, lines(stdout))
}

func TestPad_string_at_width_unchanged(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "-w", "5", "hello")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hello"}, lines(stdout))
}

func TestPad_string_longer_than_width_unchanged(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "-w", "3", "hello")
	assertExit(t, 0, exit)
	assertLines(t, []string{"hello"}, lines(stdout))
}

func TestPad_invalid_char_returns_error(t *testing.T) {
	exit, _, stderr := runCmd("pad", "-c", "ab", "hi")
	assertExit(t, 1, exit)
	assertContains(t, "error", stderr)
}

func TestPad_help_shows_usage(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "--help")
	assertExit(t, 0, exit)
	assertContains(t, "Usage:", stdout)
}

func TestStdin_pad(t *testing.T) {
	exit, stdout, _ := runWithStdin("hi\nhello\n", "pad", "-w", "7")
	assertExit(t, 0, exit)
	assertLines(t, []string{"     hi", "  hello"}, lines(stdout))
}

// Fish parity tests

func TestPad_no_width_returns_string_as_is(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foo"}, lines(stdout))
}

func TestPad_center_no_width_returns_string_as_is(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "-C", "foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foo"}, lines(stdout))
}

func TestPad_right_with_width_and_char(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "-r", "-w", "7", "--char", "-", "foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foo----"}, lines(stdout))
}

func TestPad_left_with_width_and_char(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "--width", "7", "-c", "=", "foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"====foo"}, lines(stdout))
}

func TestPad_center_with_width_and_char(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "--width", "7", "-c", "=", "-C", "foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"==foo=="}, lines(stdout))
}

func TestPad_center_even_width_bias_left(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "--width", "8", "-c", "=", "-C", "foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"===foo=="}, lines(stdout))
}

func TestPad_center_right_even_width_bias_right(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "--width", "8", "-c", "=", "-C", "-r", "foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"==foo==="}, lines(stdout))
}

func TestPad_right_width_10(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "--width", "10", "--right", "foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"foo       "}, lines(stdout))
}

func TestPad_right_center_width_10(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "--width", "10", "--right", "--center", "foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"   foo    "}, lines(stdout))
}

func TestPad_center_width_10(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "--width", "10", "--center", "foo")
	assertExit(t, 0, exit)
	assertLines(t, []string{"    foo   "}, lines(stdout))
}

func TestPad_auto_width_multiple_strings(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "-c", ".", "long", "longer", "longest")
	assertExit(t, 0, exit)
	assertLines(t, []string{"...long", ".longer", "longest"}, lines(stdout))
}

func TestPad_center_auto_width_multiple_strings(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "-c", ".", "-C", "long", "longer", "longest")
	assertExit(t, 0, exit)
	assertLines(t, []string{"..long.", ".longer", "longest"}, lines(stdout))
}

func TestPad_center_right_auto_width_multiple_strings(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "-c", ".", "-C", "-r", "long", "longer", "longest")
	assertExit(t, 0, exit)
	assertLines(t, []string{".long..", "longer.", "longest"}, lines(stdout))
}

func TestPad_width_overruled_by_longest_string(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "-c_", "--width", "5", "longer-than-width-param", "x")
	assertExit(t, 0, exit)
	assertLines(t, []string{"longer-than-width-param", "______________________x"}, lines(stdout))
}

func TestPad_center_width_overruled_by_longest_string(t *testing.T) {
	exit, stdout, _ := runCmd("pad", "-c_", "--width", "5", "--center", "longer-than-width-param", "x")
	assertExit(t, 0, exit)
	assertLines(t, []string{"longer-than-width-param", "___________x___________"}, lines(stdout))
}

func TestPad_multi_char_padding_is_error(t *testing.T) {
	exit, _, stderr := runCmd("pad", "-c", "ab", "-w4", ".")
	assertExit(t, 1, exit)
	assertContains(t, "Padding should be a character", stderr)
}

func TestPad_zero_width_char_is_error(t *testing.T) {
	exit, _, stderr := runCmd("pad", "-c", "", ".")
	assertExit(t, 1, exit)
	assertContains(t, "Invalid padding character", stderr)
}

func TestPad_negative_width_is_error(t *testing.T) {
	exit, _, stderr := runCmd("pad", "--width=-1", "foo")
	assertExit(t, 1, exit)
	assertContains(t, "Invalid width value '-1'", stderr)
}

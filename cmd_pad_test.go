package main

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

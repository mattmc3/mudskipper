package stringcmd

import (
	"testing"
)

func TestVisualWidthOf_plain_string(t *testing.T) {
	assertEqual(t, 5, visualWidthOf("hello"))
}

func TestVisualWidthOf_empty_string(t *testing.T) {
	assertEqual(t, 0, visualWidthOf(""))
}

func TestVisualWidthOf_strips_ansi_escape(t *testing.T) {
	assertEqual(t, 3, visualWidthOf("\x1b[31mabc\x1b[0m"))
}

func TestVisualWidthOf_backspace_decrements(t *testing.T) {
	assertEqual(t, 2, visualWidthOf("abc\b"))
}

func TestVisualWidthOf_backspace_floors_at_zero(t *testing.T) {
	assertEqual(t, 0, visualWidthOf("\b\b\b"))
}

func TestVisualWidthOf_bell_is_zero_width(t *testing.T) {
	assertEqual(t, 3, visualWidthOf("\afoo"))
}

func TestVisualWidthOf_control_chars_are_zero_width(t *testing.T) {
	assertEqual(t, 3, visualWidthOf("\x07\x0c\x0efoo"))
}

func TestVisualWidthOfLines_single_line(t *testing.T) {
	assertInts(t, []int{3}, visualWidthOfLines("foo"))
}

func TestVisualWidthOfLines_splits_on_newline(t *testing.T) {
	assertInts(t, []int{3, 2}, visualWidthOfLines("foo\nab"))
}

func TestVisualWidthOfLines_carriage_return_resets_position(t *testing.T) {
	assertInts(t, []int{9}, visualWidthOfLines("abcdef\rfooba\x1b[31mraaa"))
}

func TestVisualWidthOfLines_ansi_ignored(t *testing.T) {
	assertInts(t, []int{2}, visualWidthOfLines("a\x1b[34mb"))
}

func TestVisualWidthOfLines_backspace(t *testing.T) {
	assertInts(t, []int{0}, visualWidthOfLines("\b"))
}

func TestVisualTakeLeft_normal_string(t *testing.T) {
	assertStr(t, "fo", visualTakeLeft("foobar", 2))
}

func TestVisualTakeLeft_zero_width(t *testing.T) {
	assertStr(t, "", visualTakeLeft("foobar", 0))
}

func TestVisualTakeLeft_backspace_in_string(t *testing.T) {
	assertStr(t, "\ba", visualTakeLeft("\babc", 1))
}

func TestVisualTakeLeft_control_chars_included_free(t *testing.T) {
	assertStr(t, "\afoo", visualTakeLeft("\afoobar", 3))
}

func TestVisualTakeRight_normal_string(t *testing.T) {
	assertStr(t, "ar", visualTakeRight("foobar", 2))
}

func TestVisualTakeRight_zero_width(t *testing.T) {
	assertStr(t, "", visualTakeRight("foobar", 0))
}

func TestVisualTakeRight_full_string(t *testing.T) {
	assertStr(t, "foo", visualTakeRight("foo", 3))
}

func TestVisualTakeLeft_ansi_sequence_not_counted(t *testing.T) {
	assertStr(t, "\x1b[31mhel", visualTakeLeft("\x1b[31mhello\x1b[0m", 3))
}

func TestVisualTakeLeft_ansi_sequence_zero_width_included(t *testing.T) {
	assertStr(t, "\x1b[31mhello", visualTakeLeft("\x1b[31mhello\x1b[0m", 5))
}

func TestVisualTakeRight_ansi_sequence_not_counted(t *testing.T) {
	assertStr(t, "llo\x1b[0m", visualTakeRight("\x1b[31mhello\x1b[0m", 3))
}

func TestVisualTakeRight_ansi_only_prefix(t *testing.T) {
	assertStr(t, "\x1b[31mhello\x1b[0m", visualTakeRight("\x1b[31mhello\x1b[0m", 5))
}

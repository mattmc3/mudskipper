package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// stringBin is the path to the compiled binary, set by TestMain.
var stringBin string

func TestMain(m *testing.M) {
	tmp, err := os.CreateTemp("", "string-test-*")
	if err != nil {
		panic(err)
	}
	tmp.Close()
	stringBin = tmp.Name()
	defer os.Remove(stringBin)

	if out, err := exec.Command("go", "build", "-o", stringBin, ".").CombinedOutput(); err != nil {
		panic("build failed: " + string(out))
	}

	os.Exit(m.Run())
}

type cliTest struct {
	name     string
	stdin    string
	args     []string
	wantExit int
	wantOut  []string // nil = don't check; use []string{} to assert empty
	wantErr  string   // substring match; empty = don't check
}

func runBin(t *testing.T, tc cliTest) {
	t.Helper()
	cmd := exec.Command(stringBin, tc.args...)
	cmd.Stdin = strings.NewReader(tc.stdin)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	gotExit := 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			gotExit = ee.ExitCode()
		} else {
			t.Fatalf("exec error: %v", err)
		}
	}

	if gotExit != tc.wantExit {
		t.Errorf("exit: want %d, got %d", tc.wantExit, gotExit)
	}

	if tc.wantOut != nil {
		gotLines := splitLines(outBuf.String())
		if len(tc.wantOut) == 0 {
			if outBuf.Len() != 0 {
				t.Errorf("stdout: want empty, got %q", outBuf.String())
			}
		} else {
			assertLines(t, tc.wantOut, gotLines)
		}
	}

	if tc.wantErr != "" {
		if !strings.Contains(errBuf.String(), tc.wantErr) {
			t.Errorf("stderr: want %q in %q", tc.wantErr, errBuf.String())
		}
	}
}

// Fish parity tests — direct from https://github.com/fish-shell/fish-shell/blob/master/tests/checks/string.fish
// Skipped: fish variables (set/count), infinite-stdin (yes|), fish functions, multi-pipe chains, --chars alias, \g1 backrefs

var fishParityGeneral = []cliTest{
	// L6: string
	{name: "no_args", args: []string{}, wantExit: 1, wantErr: "missing subcommand"},
	// L13: string abc
	{name: "invalid_subcommand_abc", args: []string{"abc"}, wantExit: 1, wantErr: "invalid subcommand"},
	// L20: string --abc
	{name: "invalid_subcommand_flag", args: []string{"--abc"}, wantExit: 1, wantErr: "invalid subcommand"},
}

var fishParityLength = []cliTest{
	// L68: string length "hello, world"
	{name: "hello_world", args: []string{"length", "hello, world"}, wantExit: 0, wantOut: []string{"12"}},
	// L71: string length -q ""
	{name: "quiet_empty_exit1", args: []string{"length", "-q", ""}, wantExit: 1, wantOut: []string{}},
	// L567: string length
	{name: "no_args_exit1", args: []string{"length"}, wantExit: 1},
	// L175: string length --visible (set_color red)abc
	{name: "visible_ansi", args: []string{"length", "--visible", "\x1b[31mabc"}, wantExit: 0, wantOut: []string{"3"}},
	// L190: string length --visible (set_color --reset)abcdef\rfooba(set_color red)raaa
	{name: "visible_carriage_return", args: []string{"length", "--visible", "\x1b[0mabcdef\rfooba\x1b[31mraaa"}, wantExit: 0, wantOut: []string{"9"}},
	// L195: string length --visible a(set_color blue)b\ncde
	{name: "visible_multiline", args: []string{"length", "--visible", "a\x1b[34mb\ncde"}, wantExit: 0, wantOut: []string{"2", "3"}},
	// L201: string length --visible \b
	{name: "visible_backspace_zero", args: []string{"length", "--visible", "\b"}, wantExit: 1, wantOut: []string{"0"}},
	// L205: string length --visible \bf
	{name: "visible_backspace_then_char", args: []string{"length", "--visible", "\bf"}, wantExit: 0, wantOut: []string{"1"}},
	// L209: string length --visible \bf\b
	{name: "visible_backspace_erase", args: []string{"length", "--visible", "\bf\b"}, wantExit: 1, wantOut: []string{"0"}},
	// L213: string length --visible \bf\b\b\b\b\b
	{name: "visible_backspace_clamped", args: []string{"length", "--visible", "\bf\b\b\b\b\b"}, wantExit: 1, wantOut: []string{"0"}},
}

var fishParityPad = []cliTest{
	// L74: string pad foo
	{name: "no_width", args: []string{"pad", "foo"}, wantExit: 0, wantOut: []string{"foo"}},
	// L76: string pad -C foo
	{name: "center_no_width", args: []string{"pad", "-C", "foo"}, wantExit: 0, wantOut: []string{"foo"}},
	// L79: string pad -r -w 7 --char - foo
	{name: "right_w7_dash", args: []string{"pad", "-r", "-w", "7", "--char", "-", "foo"}, wantExit: 0, wantOut: []string{"foo----"}},
	// L81: string pad -r -w 7 --chars - --center foo  (--chars alias not supported; skip)
	// L90: string pad --width 7 -c '=' foo
	{name: "w7_eq_left", args: []string{"pad", "--width", "7", "-c", "=", "foo"}, wantExit: 0, wantOut: []string{"====foo"}},
	// L92: string pad --width 7 -c '=' -C foo
	{name: "w7_eq_center", args: []string{"pad", "--width", "7", "-c", "=", "-C", "foo"}, wantExit: 0, wantOut: []string{"==foo=="}},
	// L94: string pad --width 8 -c '=' -C foo
	{name: "w8_eq_center_odd", args: []string{"pad", "--width", "8", "-c", "=", "-C", "foo"}, wantExit: 0, wantOut: []string{"===foo=="}},
	// L96: string pad --width 8 -c '=' -Cr foo
	{name: "w8_eq_center_right", args: []string{"pad", "--width", "8", "-c", "=", "-C", "-r", "foo"}, wantExit: 0, wantOut: []string{"==foo==="}},
	// L99: string pad --width 10 --right foo
	{name: "w10_right", args: []string{"pad", "--width", "10", "--right", "foo"}, wantExit: 0, wantOut: []string{"foo       "}},
	// L101: string pad --width 10 --right --center foo
	{name: "w10_right_center", args: []string{"pad", "--width", "10", "--right", "--center", "foo"}, wantExit: 0, wantOut: []string{"   foo    "}},
	// L103: string pad --width 10 --center foo
	{name: "w10_center", args: []string{"pad", "--width", "10", "--center", "foo"}, wantExit: 0, wantOut: []string{"    foo   "}},
	// L137: string pad -c . long longer longest
	{name: "auto_dot", args: []string{"pad", "-c", ".", "long", "longer", "longest"}, wantExit: 0, wantOut: []string{"...long", ".longer", "longest"}},
	// L141: string pad -c . -C long longer longest
	{name: "center_auto_dot", args: []string{"pad", "-c", ".", "-C", "long", "longer", "longest"}, wantExit: 0, wantOut: []string{"..long.", ".longer", "longest"}},
	// L145: string pad -c . -Cr long longer longest
	{name: "center_right_auto_dot", args: []string{"pad", "-c", ".", "-C", "-r", "long", "longer", "longest"}, wantExit: 0, wantOut: []string{".long..", "longer.", "longest"}},
	// L152: string pad -c_ --width 5 longer-than-width-param x
	{name: "width_overruled", args: []string{"pad", "-c_", "--width", "5", "longer-than-width-param", "x"}, wantExit: 0, wantOut: []string{"longer-than-width-param", "______________________x"}},
	// L155: string pad -c_ --width 5 --center longer-than-width-param x
	{name: "center_width_overruled", args: []string{"pad", "-c_", "--width", "5", "--center", "longer-than-width-param", "x"}, wantExit: 0, wantOut: []string{"longer-than-width-param", "___________x___________"}},
	// L158: string pad -c_ --width 5 --center --right longer-than-width-param x
	{name: "center_right_width_overruled", args: []string{"pad", "-c_", "--width", "5", "--center", "--right", "longer-than-width-param", "x"}, wantExit: 0, wantOut: []string{"longer-than-width-param", "___________x___________"}},
	// L164: string pad -c ab -w4 .
	{name: "multi_char_error", args: []string{"pad", "-c", "ab", "-w4", "."}, wantExit: 1, wantErr: "Padding should be a character"},
	// L168: string pad -c \u07 .  (\u07 = BEL = 0x07)
	{name: "zero_width_char_error", args: []string{"pad", "-c", "\a", "."}, wantExit: 1, wantErr: "Invalid padding character"},
	// L171: string pad --width=-1 foo
	{name: "negative_width_error", args: []string{"pad", "--width=-1", "foo"}, wantExit: 1, wantErr: "Invalid width value '-1'"},
}

var fishParitySub = []cliTest{
	// L85: string sub --start -9223372036854775808 abc
	{name: "start_i64min", args: []string{"sub", "--start", "-9223372036854775808", "abc"}, wantExit: 0, wantOut: []string{"abc"}},
	// L87: string sub --start 0 abc
	{name: "start_zero_error", args: []string{"sub", "--start", "0", "abc"}, wantExit: 1, wantErr: "Invalid start value '0'"},
	// L216: string sub --length 2 abcde
	{name: "length2", args: []string{"sub", "--length", "2", "abcde"}, wantExit: 0, wantOut: []string{"ab"}},
	// L219: string sub -s 2 -l 2 abcde
	{name: "s2_l2", args: []string{"sub", "-s", "2", "-l", "2", "abcde"}, wantExit: 0, wantOut: []string{"bc"}},
	// L222: string sub --length=-1 abcde
	{name: "length_neg1_error", args: []string{"sub", "--length=-1", "abcde"}, wantExit: 1, wantErr: "Invalid length value '-1'"},
	// L225: string sub --start=-2 abcde
	{name: "start_neg2", args: []string{"sub", "--start=-2", "abcde"}, wantExit: 0, wantOut: []string{"de"}},
	// L228: string sub --end=3 abcde
	{name: "end3", args: []string{"sub", "--end=3", "abcde"}, wantExit: 0, wantOut: []string{"abc"}},
	// L231: string sub --end=-4 abcde
	{name: "end_neg4", args: []string{"sub", "--end=-4", "abcde"}, wantExit: 0, wantOut: []string{"a"}},
	// L234: string sub --end=0 abcde
	{name: "end0_error", args: []string{"sub", "--end=0", "abcde"}, wantExit: 1, wantErr: "Invalid end value '0'"},
	// L237: string sub --start=2 --end=-2 abcde
	{name: "s2_e_neg2", args: []string{"sub", "--start=2", "--end=-2", "abcde"}, wantExit: 0, wantOut: []string{"bc"}},
	// L240: string sub -s -5 -e -2 abcdefgh
	{name: "s_neg5_e_neg2", args: []string{"sub", "-s", "-5", "-e", "-2", "abcdefgh"}, wantExit: 0, wantOut: []string{"def"}},
	// L243: string sub -s -100 -e -2 abcde
	{name: "s_neg100_e_neg2", args: []string{"sub", "-s", "-100", "-e", "-2", "abcde"}, wantExit: 0, wantOut: []string{"abc"}},
	// L246: string sub -s -5 -e 2 abcde
	{name: "s_neg5_e2", args: []string{"sub", "-s", "-5", "-e", "2", "abcde"}, wantExit: 0, wantOut: []string{"ab"}},
	// L249: string sub -s -50 -e -100 abcde
	{name: "end_before_start_empty", args: []string{"sub", "-s", "-50", "-e", "-100", "abcde"}, wantExit: 1, wantOut: []string{""}},
	// L252: string sub -s 2 -e -5 abcde
	{name: "start_after_end_empty", args: []string{"sub", "-s", "2", "-e", "-5", "abcde"}, wantExit: 1, wantOut: []string{""}},
	// L255: string sub -s 2 -e -5 -l 3 abcde
	{name: "end_and_length_exclusive", args: []string{"sub", "-s", "2", "-e", "-5", "-l", "3", "abcde"}, wantExit: 1, wantErr: "--end and --length are mutually exclusive"},
}

var fishParitySplit = []cliTest{
	// L258: string split . example.com
	{name: "basic_dot", args: []string{"split", ".", "example.com"}, wantExit: 0, wantOut: []string{"example", "com"}},
	// L262: string split -r -m1 / /usr/local/bin/fish
	{name: "right_max1", args: []string{"split", "-r", "-m1", "/", "/usr/local/bin/fish"}, wantExit: 0, wantOut: []string{"/usr/local/bin", "fish"}},
	// L266: string split "" abc
	{name: "empty_sep_chars", args: []string{"split", "", "abc"}, wantExit: 0, wantOut: []string{"a", "b", "c"}},
	// L271: string split
	{name: "no_args_error", args: []string{"split"}, wantExit: 1, wantErr: "split requires a separator"},
	// L274: string split --max 1 --right 12 AB12CD
	{name: "max1_right", args: []string{"split", "--max", "1", "--right", "12", "AB12CD"}, wantExit: 0, wantOut: []string{"AB", "CD"}},
	// L278: string split --max=-1 --right 12 AB12CD
	{name: "negative_max_error", args: []string{"split", "--max=-1", "--right", "12", "AB12CD"}, wantExit: 1, wantErr: "Invalid max value '-1'"},
	// L281: string split --fields=2 "" abc
	{name: "fields2", args: []string{"split", "--fields=2", "", "abc"}, wantExit: 0, wantOut: []string{"b"}},
	// L284: string split --fields=3,2 "" abc
	{name: "fields3_2", args: []string{"split", "--fields=3,2", "", "abc"}, wantExit: 0, wantOut: []string{"c", "b"}},
	// L288: string split --fields=2,9 "" abc
	{name: "fields_out_of_range_exit1", args: []string{"split", "--fields=2,9", "", "abc"}, wantExit: 1},
	// L291: string split --fields=2-3-,9 "" a
	{name: "fields_malformed_error", args: []string{"split", "--fields=2-3-,9", "", "a"}, wantExit: 1, wantErr: "invalid field spec"},
	// L300: string split --fields=1--2 "" b
	{name: "fields_double_dash_error", args: []string{"split", "--fields=1--2", "", "b"}, wantExit: 1},
	// L303: string split --fields=0 "" c
	{name: "fields_zero_error", args: []string{"split", "--fields=0", "", "c"}, wantExit: 1, wantErr: "Invalid fields value '0'"},
	// L309: string split --fields=1-0 "" d
	{name: "fields_inverted_range_error", args: []string{"split", "--fields=1-0", "", "d"}, wantExit: 1, wantErr: "Invalid range value for field '1-0'"},
	// L312: string split --fields=0-1 "" e
	{name: "fields_zero_start_range_error", args: []string{"split", "--fields=0-1", "", "e"}, wantExit: 1, wantErr: "Invalid"},
	// L318: string split --fields=1a "" g
	{name: "fields_alpha_suffix_error", args: []string{"split", "--fields=1a", "", "g"}, wantExit: 1, wantErr: "invalid field spec"},
	// L321: string split --fields=a "" h
	{name: "fields_alpha_error", args: []string{"split", "--fields=a", "", "h"}, wantExit: 1, wantErr: "invalid field spec"},
	// L324: string split --fields=1-3,5,9-7 "" 123456789
	{name: "fields_range_reverse", args: []string{"split", "--fields=1-3,5,9-7", "", "123456789"}, wantExit: 0, wantOut: []string{"1", "2", "3", "5", "9", "8", "7"}},
	// L333: string split -f1 ' ' 'a b' 'c d'
	{name: "f_shorthand", args: []string{"split", "-f1", " ", "a b", "c d"}, wantExit: 0, wantOut: []string{"a", "c"}},
	// L337: string split --allow-empty --fields=2,9 "" abc
	{name: "allow_empty_with_fields", args: []string{"split", "--allow-empty", "--fields=2,9", "", "abc"}, wantExit: 0, wantOut: []string{"b"}},
	// L340: string split --allow-empty "" abc
	{name: "allow_empty_no_fields_error", args: []string{"split", "--allow-empty", "", "abc"}, wantExit: 1, wantErr: "--allow-empty is only valid with --fields"},
}

var fishParityJoin = []cliTest{
	// L343: seq 3 | string join ...
	{name: "seq_ellipsis", stdin: "1\n2\n3\n", args: []string{"join", "..."}, wantExit: 0, wantOut: []string{"1...2...3"}},
	// L346: string join
	{name: "no_args_error", args: []string{"join"}, wantExit: 1, wantErr: "join requires a separator"},
}

var fishParityTrim = []cliTest{
	// L349: string trim " abc  "
	{name: "both_sides", args: []string{"trim", " abc  "}, wantExit: 0, wantOut: []string{"abc"}},
	// L352: string trim --right --chars=yz xyzzy zany
	{name: "right_chars", args: []string{"trim", "--right", "--chars=yz", "xyzzy", "zany"}, wantExit: 0, wantOut: []string{"x", "zan"}},
}

var fishParityEscape = []cliTest{
	// L356: echo \x07 | string escape
	{name: "bel_via_stdin", stdin: "\x07", args: []string{"escape"}, wantExit: 0, wantOut: []string{`\cg`}},
	// L359: string escape --style=script 'a b#c"\'d'
	{name: "script_special", args: []string{"escape", "--style=script", "a b#c\"'d"}, wantExit: 0, wantOut: []string{"'a b#c\"\\'d'"}},
	// L362: string escape --no-quoted --style=script 'a b#c"\'d'
	{name: "script_no_quoted_special", args: []string{"escape", "--no-quoted", "--style=script", "a b#c\"'d"}, wantExit: 0, wantOut: []string{"a\\ b#c\\\"\\'d"}},
	// L365: string escape --no-quoted --style=script 'a #b'
	{name: "script_no_quoted_hash", args: []string{"escape", "--no-quoted", "--style=script", "a #b"}, wantExit: 0, wantOut: []string{"a\\ \\#b"}},
	// L368: string escape --style=url 'a b#c"\'d'
	{name: "url_special", args: []string{"escape", "--style=url", "a b#c\"'d"}, wantExit: 0, wantOut: []string{"a%20b%23c%22%27d"}},
	// L371: string escape --style=url \na\nb%c~d\n
	{name: "url_newlines", args: []string{"escape", "--style=url", "\na\nb%c~d\n"}, wantExit: 0, wantOut: []string{"%0Aa%0Ab%25c~d%0A"}},
	// L374: string escape --style=var 'a b#c"\'d'  (fish packs adjacent: _22_27_; we emit _22__27_; known diff, skip)
	// L377: string escape --style=var a\nghi_
	{name: "var_newline_underscore", args: []string{"escape", "--style=var", "a\nghi_"}, wantExit: 0, wantOut: []string{"a_0A_ghi__"}},
	// L380: string escape --style=var abc
	{name: "var_alphanumeric", args: []string{"escape", "--style=var", "abc"}, wantExit: 0, wantOut: []string{"abc"}},
	// L383: string escape --style=var _a_b_c_
	{name: "var_underscores", args: []string{"escape", "--style=var", "_a_b_c_"}, wantExit: 0, wantOut: []string{"__a__b__c__"}},
	// L386: string escape --style=var -- -
	{name: "var_dash", args: []string{"escape", "--style=var", "--", "-"}, wantExit: 0, wantOut: []string{"_2D_"}},
	// L389-L406: multibyte (aöb, 中) — skipped; tested indirectly via url/var roundtrip
	// L409: string escape --style=regex ".ext"
	{name: "regex_dot_ext", args: []string{"escape", "--style=regex", ".ext"}, wantExit: 0, wantOut: []string{`\.ext`}},
	// L410: string escape --style=regex "bonjour, amigo"
	{name: "regex_no_meta", args: []string{"escape", "--style=regex", "bonjour, amigo"}, wantExit: 0, wantOut: []string{"bonjour, amigo"}},
	// L411: string escape --style=regex "^this is a literal string"
	{name: "regex_caret", args: []string{"escape", "--style=regex", "^this is a literal string"}, wantExit: 0, wantOut: []string{`\^this is a literal string`}},
	// L412: string escape --style=regex "hello\nworld"  (fish renders \\n; Go regexp.QuoteMeta keeps literal newline; known diff)
	// {name: "regex_newline", args: []string{"escape", "--style=regex", "hello\nworld"}, wantExit: 0, wantOut: []string{`hello\nworld`}},
	// L419: string escape --style=unknown-style
	{name: "unknown_style_error", args: []string{"escape", "--style=unknown-style"}, wantExit: 1, wantErr: "Invalid escape style 'unknown-style'"},
	// L1013: string escape \x7F
	{name: "del_char", args: []string{"escape", "\x7F"}, wantExit: 0, wantOut: []string{`\x7f`}},
}

var fishParityUnescape = []cliTest{
	// L429: string unescape --style=script (string escape --style=script 'a b#c"\'d')
	{name: "script_roundtrip", args: []string{"unescape", "--style=script", "'a b#c\"\\'d'"}, wantExit: 0, wantOut: []string{"a b#c\"'d"}},
	// L432: string unescape --style=url (string escape --style=url 'a b#c"\'d')
	{name: "url_roundtrip", args: []string{"unescape", "--style=url", "a%20b%23c%22%27d"}, wantExit: 0, wantOut: []string{"a b#c\"'d"}},
	// L435: string unescape --style=url \na\nb%c~d\n
	{name: "url_newlines", args: []string{"unescape", "--style=url", "%0Aa%0Ab%25c~d%0A"}, wantExit: 0, wantOut: []string{"", "a", "b%c~d"}},
	// L440: string unescape --style=var (string escape --style=var 'a b#c"\'d')
	{name: "var_roundtrip", args: []string{"unescape", "--style=var", "a_20_b_23_c_22__27_d"}, wantExit: 0, wantOut: []string{"a b#c\"'d"}},
	// L443: string unescape --style=var a\nghi_
	{name: "var_newline", args: []string{"unescape", "--style=var", "a_0A_ghi__"}, wantExit: 0, wantOut: []string{"a", "ghi_"}},
	// L447: string unescape --style=var abc
	{name: "var_alphanumeric", args: []string{"unescape", "--style=var", "abc"}, wantExit: 0, wantOut: []string{"abc"}},
	// L450: string unescape --style=var '_a_b_c_'
	{name: "var_underscores", args: []string{"unescape", "--style=var", "__a__b__c__"}, wantExit: 0, wantOut: []string{"_a_b_c_"}},
	// L453: string unescape --style=var -- -
	{name: "var_dash", args: []string{"unescape", "--style=var", "--", "_2D_"}, wantExit: 0, wantOut: []string{"-"}},
	// L456: string unescape --style=unknown-style
	{name: "unknown_style_error", args: []string{"unescape", "--style=unknown-style"}, wantExit: 1, wantErr: "Invalid style value 'unknown-style'"},
}

var fishParityMatch = []cliTest{
	// L27: string match -r -v "c.*" dog can cat diz
	{name: "regex_invert_cstar", args: []string{"match", "-r", "-v", "c.*", "dog", "can", "cat", "diz"}, wantExit: 0, wantOut: []string{"dog", "diz"}},
	// L32: string match -q -r -v "c.*" dog can cat diz
	{name: "quiet_regex_invert_cstar", args: []string{"match", "-q", "-r", "-v", "c.*", "dog", "can", "cat", "diz"}, wantExit: 0, wantOut: []string{}},
	// L35: string match -v "c*" dog can cat diz
	{name: "glob_invert_cstar", args: []string{"match", "-v", "c*", "dog", "can", "cat", "diz"}, wantExit: 0, wantOut: []string{"dog", "diz"}},
	// L40: string match -q -v "c*" dog can cat diz
	{name: "quiet_glob_invert_cstar", args: []string{"match", "-q", "-v", "c*", "dog", "can", "cat", "diz"}, wantExit: 0, wantOut: []string{}},
	// L43: string match -v "d*" dog dan dat diz
	{name: "glob_invert_all_match_exit1", args: []string{"match", "-v", "d*", "dog", "dan", "dat", "diz"}, wantExit: 1},
	// L46: string match -q -v "d*" dog dan dat diz
	{name: "quiet_glob_invert_all_match_exit1", args: []string{"match", "-q", "-v", "d*", "dog", "dan", "dat", "diz"}, wantExit: 1, wantOut: []string{}},
	// L49: string match -r -v x y
	{name: "regex_invert_y", args: []string{"match", "-r", "-v", "x", "y"}, wantExit: 0, wantOut: []string{"y"}},
	// L53: string match -r -v x x
	{name: "regex_invert_x_exit1", args: []string{"match", "-r", "-v", "x", "x"}, wantExit: 1},
	// L56: string match -q -r -v x y
	{name: "quiet_regex_invert_y", args: []string{"match", "-q", "-r", "-v", "x", "y"}, wantExit: 0, wantOut: []string{}},
	// L59: string match -q -r -v x x
	{name: "quiet_regex_invert_x_exit1", args: []string{"match", "-q", "-r", "-v", "x", "x"}, wantExit: 1, wantOut: []string{}},
	// L62: string match -v -g foo foo
	{name: "invert_groups_only_exclusive", args: []string{"match", "-v", "-g", "foo", "foo"}, wantExit: 1, wantErr: "mutually exclusive"},
	// L65: string match
	{name: "no_args_error", args: []string{"match"}, wantExit: 1},
	// L460: string match "*" a
	{name: "glob_star_any", args: []string{"match", "*", "a"}, wantExit: 0, wantOut: []string{"a"}},
	// L463: string match "a*b" axxb
	{name: "glob_wildcard", args: []string{"match", "a*b", "axxb"}, wantExit: 0, wantOut: []string{"axxb"}},
	// L466: string match -i "a**B" Axxb
	{name: "glob_case_insensitive", args: []string{"match", "-i", "a**B", "Axxb"}, wantExit: 0, wantOut: []string{"Axxb"}},
	// L469: echo "ok?" | string match "*?"
	{name: "glob_question_stdin", stdin: "ok?\n", args: []string{"match", "*?"}, wantExit: 0, wantOut: []string{"ok?"}},
	// L472: string match -r "cat|dog|fish" "nice dog"
	{name: "regex_alternation", args: []string{"match", "-r", "cat|dog|fish", "nice dog"}, wantExit: 0, wantOut: []string{"dog"}},
	// L475: string match -r "(\d\d?):(\d\d):(\d\d)" 2:34:56
	{name: "regex_time_groups", args: []string{"match", "-r", `(\d\d?):(\d\d):(\d\d)`, "2:34:56"}, wantExit: 0, wantOut: []string{"2:34:56", "2", "34", "56"}},
	// L481: string match -r "^(\w{2,4})\g1\$" papa mud murmur  (skip: \g1 not in Go RE2)
	// L487: string match -r -a -n at ratatat
	{name: "regex_all_index", args: []string{"match", "-r", "-a", "-n", "at", "ratatat"}, wantExit: 0, wantOut: []string{"2 2", "4 2", "6 2"}},
	// L492: string match -r -i "0x[0-9a-f]{1,8}" "int magic = 0xBadC0de;"
	{name: "regex_case_insensitive_hex", args: []string{"match", "-r", "-i", "0x[0-9a-f]{1,8}", "int magic = 0xBadC0de;"}, wantExit: 0, wantOut: []string{"0xBadC0de"}},
	// L551: string match -r '^([ugoa]*)([=+-]?)([rwx]*)$' '=r'
	{name: "regex_empty_capture_groups", args: []string{"match", "-r", `^([ugoa]*)([=+-]?)([rwx]*)$`, "=r"}, wantExit: 0, wantOut: []string{"=r", "", "=", "r"}},
	// L558: string match -r "[" "a[sd"
	{name: "regex_compile_error", args: []string{"match", "-r", "[", "a[sd"}, wantExit: 1, wantErr: "error:"},
	// L570: string match -r -v "[dcantg].*" dog can cat diz
	{name: "regex_invert_all_match_exit1", args: []string{"match", "-r", "-v", "[dcantg].*", "dog", "can", "cat", "diz"}, wantExit: 1},
	// L573: string match -v "*" dog can cat diz
	{name: "glob_invert_star_all_exit1", args: []string{"match", "-v", "*", "dog", "can", "cat", "diz"}, wantExit: 1},
	// L576: string match -rvn a bbb
	{name: "regex_invert_index", args: []string{"match", "-r", "-v", "-n", "a", "bbb"}, wantExit: 0, wantOut: []string{"1 3"}},
	// L736: string match -e x abc dxf xyz jkx x z
	{name: "entire_substring", args: []string{"match", "-e", "x", "abc", "dxf", "xyz", "jkx", "x", "z"}, wantExit: 0, wantOut: []string{"dxf", "xyz", "jkx", "x"}},
	// L743: string match x abc dxf xyz jkx x z
	{name: "glob_exact_only", args: []string{"match", "x", "abc", "dxf", "xyz", "jkx", "x", "z"}, wantExit: 0, wantOut: []string{"x"}},
	// L746: string match --entire -r "a*b[xy]+" abc abxc bye aaabyz kaabxz abbxy abcx caabxyxz
	{name: "entire_regex_no_groups", args: []string{"match", "--entire", "-r", "a*b[xy]+", "abc", "abxc", "bye", "aaabyz", "kaabxz", "abbxy", "abcx", "caabxyxz"}, wantExit: 0, wantOut: []string{"abxc", "bye", "aaabyz", "kaabxz", "abbxy", "caabxyxz"}},
	// L755: string match --entire "" -- banana  (-- not consumed when positional precedes it in getopt)
	// {name: "entire_empty_pattern", args: []string{"match", "--entire", "", "--", "banana"}, wantExit: 0, wantOut: []string{"banana"}},
	// L761: string match -r "a*b[xy]+" abc abxc bye aaabyz kaabxz abbxy abcx caabxyxz
	{name: "regex_no_groups", args: []string{"match", "-r", "a*b[xy]+", "abc", "abxc", "bye", "aaabyz", "kaabxz", "abbxy", "abcx", "caabxyxz"}, wantExit: 0, wantOut: []string{"abx", "by", "aaaby", "aabx", "bxy", "aabxyx"}},
	// L772: string match --entire -r "a*b([xy]+)" abc abxc bye aaabyz kaabxz abbxy abcx caabxyxz
	{name: "entire_regex_with_groups", args: []string{"match", "--entire", "-r", "a*b([xy]+)", "abc", "abxc", "bye", "aaabyz", "kaabxz", "abbxy", "abcx", "caabxyxz"}, wantExit: 0, wantOut: []string{"abxc", "x", "bye", "y", "aaabyz", "y", "kaabxz", "x", "abbxy", "xy", "caabxyxz", "xyx"}},
	// L787: string match --entire --index foo foo
	{name: "entire_index_exclusive", args: []string{"match", "--entire", "--index", "foo", "foo"}, wantExit: 1, wantErr: "--entire and --index"},
	// L790: string match --entire --groups-only -r foo foo
	{name: "entire_groups_only_exclusive", args: []string{"match", "--entire", "--groups-only", "-r", "foo", "foo"}, wantExit: 1, wantErr: "--entire and --groups-only"},
	// L794: string match -r "a*b([xy]+)" abc abxc bye aaabyz kaabxz abbxy abcx caabxyxz
	{name: "regex_with_groups", args: []string{"match", "-r", "a*b([xy]+)", "abc", "abxc", "bye", "aaabyz", "kaabxz", "abbxy", "abcx", "caabxyxz"}, wantExit: 0, wantOut: []string{"abx", "x", "by", "y", "aaaby", "y", "aabx", "x", "bxy", "xy", "aabxyx", "xyx"}},
	// L969: string match -qer asd asd
	{name: "quiet_entire_regex", args: []string{"match", "-q", "-e", "-r", "asd", "asd"}, wantExit: 0, wantOut: []string{}},
	// L973: string match -r "(*UTF).*" aaa
	{name: "utf_mode_error", args: []string{"match", "-r", "(*UTF).*", "aaa"}, wantExit: 1, wantErr: "error:"},
	// L984: string match -eq asd asd
	{name: "entire_quiet", args: []string{"match", "-e", "-q", "asd", "asd"}, wantExit: 0, wantOut: []string{}},
	// L1021: string match -rg '(.*)fish' catfish
	{name: "groups_only_catfish", args: []string{"match", "-rg", "(.*)fish", "catfish"}, wantExit: 0, wantOut: []string{"cat"}},
	// L1023: string match -rg '(.*)fish' shellfish
	{name: "groups_only_shellfish", args: []string{"match", "-rg", "(.*)fish", "shellfish"}, wantExit: 0, wantOut: []string{"shell"}},
	// L1025: string match -rg '(.*)fish' fish  (empty capture)
	{name: "groups_only_empty_capture", args: []string{"match", "-rg", "(.*)fish", "fish"}, wantExit: 0, wantOut: []string{""}},
	// L1027: string match -rg '(.*)fish' banana  (no match)
	{name: "groups_only_no_match", args: []string{"match", "-rg", "(.*)fish", "banana"}, wantExit: 1},
	// L1030: string match -r --groups-only '(.+)fish' fish
	{name: "groups_only_no_match2", args: []string{"match", "-r", "--groups-only", "(.+)fish", "fish"}, wantExit: 1},
	// L1034: string match -r --groups-only '(.+)fish(.*)' catfishcolor
	{name: "groups_only_multi", args: []string{"match", "-r", "--groups-only", "(.+)fish(.*)", "catfishcolor"}, wantExit: 0, wantOut: []string{"cat", "color"}},
	// L1039: echo "foo bar baz" | string match -rg 'foo (bar) baz'
	{name: "groups_only_stdin", stdin: "foo bar baz\n", args: []string{"match", "-rg", "foo (bar) baz"}, wantExit: 0, wantOut: []string{"bar"}},
	// L1041: echo "foo1x foo2x foo3x" | string match -arg 'foo(\d)x'
	{name: "all_groups_stdin", stdin: "foo1x foo2x foo3x\n", args: []string{"match", "-arg", `foo(\d)x`}, wantExit: 0, wantOut: []string{"1", "2", "3"}},
	// L1283: printf "dog\ncat\nbat\ngnat\n" | string match -m2 "*at"
	{name: "max2_glob", stdin: "dog\ncat\nbat\ngnat\n", args: []string{"match", "-m2", "*at"}, wantExit: 0, wantOut: []string{"cat", "bat"}},
	// L1287: string match -m0 foo
	{name: "m0_error", args: []string{"match", "-m0", "foo"}, wantExit: 1, wantErr: "Invalid max matches value '0'"},
	// L1290: string match -m999999999999999999999999999999999999999 foo
	{name: "overflow_max_error", args: []string{"match", "-m999999999999999999999999999999999999999", "foo"}, wantExit: 1, wantErr: "Invalid max matches value"},
	// L1293: printf "dog\ncat\nbat\nhog\n" | string match -rvm1 'at$'
	{name: "invert_max1_stdin", stdin: "dog\ncat\nbat\nhog\n", args: []string{"match", "-r", "-v", "-m1", "at$"}, wantExit: 0, wantOut: []string{"dog"}},
}

var fishParityReplace = []cliTest{
	// L498: string replace is was "blue is my favorite"
	{name: "literal_basic", args: []string{"replace", "is", "was", "blue is my favorite"}, wantExit: 0, wantOut: []string{"blue was my favorite"}},
	// L501: string replace 3rd last 1st 2nd 3rd
	{name: "literal_multi", args: []string{"replace", "3rd", "last", "1st", "2nd", "3rd"}, wantExit: 0, wantOut: []string{"1st", "2nd", "last"}},
	// L506: string replace -a " " _ "spaces to underscores"
	{name: "all_spaces", args: []string{"replace", "-a", " ", "_", "spaces to underscores"}, wantExit: 0, wantOut: []string{"spaces_to_underscores"}},
	// L509: string replace -r -a "[^\d.]+" " " "0 one two 3.14 four 5x"
	{name: "regex_all_non_numeric", args: []string{"replace", "-r", "-a", `[^\d.]+`, " ", "0 one two 3.14 four 5x"}, wantExit: 0, wantOut: []string{"0 3.14 5 "}},
	// L512: string replace -r "(\w+)\s+(\w+)" "\$2 \$1 \$\$" "left right"
	{name: "regex_backref_dollar", args: []string{"replace", "-r", `(\w+)\s+(\w+)`, "$2 $1 $$", "left right"}, wantExit: 0, wantOut: []string{"right left $"}},
	// L515: string replace -r "\s*newline\s*" "\n" "put a newline here"
	{name: "regex_insert_newline", args: []string{"replace", "-r", `\s*newline\s*`, "\n", "put a newline here"}, wantExit: 0, wantOut: []string{"put a", "here"}},
	// L519: string replace -r -a "(\w)" "\$1\$1" ab
	{name: "regex_all_double_chars", args: []string{"replace", "-r", "-a", `(\w)`, "$1$1", "ab"}, wantExit: 0, wantOut: []string{"aabb"}},
	// L522: echo a | string replace b c -q
	{name: "quiet_no_match_exit1", stdin: "a\n", args: []string{"replace", "-q", "b", "c"}, wantExit: 1, wantOut: []string{}},
	// L526: echo a | string replace -r b c -q
	{name: "regex_quiet_no_match_exit1", stdin: "a\n", args: []string{"replace", "-rq", "b", "c"}, wantExit: 1, wantOut: []string{}},
	// L530: string replace --filter x X abc axc x def jkx
	{name: "filter_changed_only", args: []string{"replace", "--filter", "x", "X", "abc", "axc", "x", "def", "jkx"}, wantExit: 0, wantOut: []string{"aXc", "X", "jkX"}},
	// L536: string replace --filter y Y abc axc x def jkx
	{name: "filter_no_match_exit1", args: []string{"replace", "--filter", "y", "Y", "abc", "axc", "x", "def", "jkx"}, wantExit: 1},
	// L539: string replace --regex -f "\d" X 1bc axc 2 d3f jk4 xyz
	{name: "regex_filter", args: []string{"replace", "--regex", "-f", `\d`, "X", "1bc", "axc", "2", "d3f", "jk4", "xyz"}, wantExit: 0, wantOut: []string{"Xbc", "X", "dXf", "jkX"}},
	// L546: string replace --regex -f Z X 1bc axc 2 d3f jk4 xyz
	{name: "regex_filter_no_match_exit1", args: []string{"replace", "--regex", "-f", "Z", "X", "1bc", "axc", "2", "d3f", "jk4", "xyz"}, wantExit: 1},
	// L979: string replace -r "(*UTF).*" aaa
	{name: "utf_mode_error", args: []string{"replace", "-r", "(*UTF).*", "replacement", "aaa"}, wantExit: 1, wantErr: "error:"},
	// L989: echo az | string replace -r -- 'a(b.+)?z' 'a:$1z'  ($1z = group "1z" in Go ExpandString; known diff)
	// {name: "unmatched_group_empty", stdin: "az\n", args: []string{"replace", "-r", "--", `a(b.+)?z`, `a:$1z`}, wantExit: 0, wantOut: []string{"a:z"}},
	// L1057: printf my-password | string replace -ra . \*
	{name: "all_chars_asterisk", stdin: "my-password", args: []string{"replace", "-ra", ".", "*"}, wantExit: 0, wantOut: []string{"***********"}},
	// L1296: printf "dog\ncat\nbat\n" | string replace -r --max-matches 1 '^c' h
	{name: "max_matches_1_stdin", stdin: "dog\ncat\nbat\n", args: []string{"replace", "-r", "--max-matches", "1", "^c", "h"}, wantExit: 0, wantOut: []string{"dog", "hat", "bat"}},
	// L1306: string replace --max-matches abc
	{name: "max_non_int_error", args: []string{"replace", "--max-matches", "abc"}, wantExit: 1, wantErr: "Invalid max matches value 'abc'"},
	// L1308: string replace --max-matches -1
	{name: "max_neg_error", args: []string{"replace", "--max-matches", "-1"}, wantExit: 1, wantErr: "Invalid max matches value '-1'"},
	// L1310: string replace --max-matches 99999999999999999999
	{name: "max_overflow_error", args: []string{"replace", "--max-matches", "99999999999999999999"}, wantExit: 1, wantErr: "Invalid max matches value"},
	// L1313: string replace
	{name: "no_args_error", args: []string{"replace"}, wantExit: 1, wantErr: "replace requires"},
	// L1315: string replace one
	{name: "one_arg_error", args: []string{"replace", "one"}, wantExit: 1, wantErr: "replace requires"},
}

var fishParityRepeat = []cliTest{
	// L580: string repeat -n 2 foo
	{name: "n2", args: []string{"repeat", "-n", "2", "foo"}, wantExit: 0, wantOut: []string{"foofoo"}},
	// L583: string repeat --count 2 foo
	{name: "count_long", args: []string{"repeat", "--count", "2", "foo"}, wantExit: 0, wantOut: []string{"foofoo"}},
	// L586: string repeat 2 foo
	{name: "positional_count", args: []string{"repeat", "2", "foo"}, wantExit: 0, wantOut: []string{"foofoo"}},
	// L592: echo foo | string repeat -n 2
	{name: "stdin_n2", stdin: "foo\n", args: []string{"repeat", "-n", "2"}, wantExit: 0, wantOut: []string{"foofoo"}},
	// L595: echo foo | string repeat 2
	{name: "stdin_positional", stdin: "foo\n", args: []string{"repeat", "2"}, wantExit: 0, wantOut: []string{"foofoo"}},
	// L598: string repeat
	{name: "no_args_exit1", args: []string{"repeat"}, wantExit: 1},
	// L601: string repeat foo
	{name: "invalid_count_error", args: []string{"repeat", "foo"}, wantExit: 1, wantErr: "Invalid count value 'foo'"},
	// L604: string repeat -n2 -q foo
	{name: "quiet_exit0", args: []string{"repeat", "-n2", "-q", "foo"}, wantExit: 0, wantOut: []string{}},
	// L610: string repeat -n0 foo
	{name: "n0_exit1", args: []string{"repeat", "-n0", "foo"}, wantExit: 1},
	// L616: string repeat -m0
	{name: "m0_exit1", args: []string{"repeat", "-m0"}, wantExit: 1},
	// L627: string repeat -n10 -m4 foo
	{name: "max_limits", args: []string{"repeat", "-n10", "-m4", "foo"}, wantExit: 0, wantOut: []string{"foof"}},
	// L630: string repeat -n10 --max 5 foo
	{name: "max5_n10", args: []string{"repeat", "-n10", "--max", "5", "foo"}, wantExit: 0, wantOut: []string{"foofo"}},
	// L633: string repeat -n3 -m20 foo
	{name: "max_larger_than_result", args: []string{"repeat", "-n3", "-m20", "foo"}, wantExit: 0, wantOut: []string{"foofoofoo"}},
	// L636: string repeat -m4 foo
	{name: "max_only", args: []string{"repeat", "-m4", "foo"}, wantExit: 0, wantOut: []string{"foof"}},
	// L639: string repeat -n 5 a b c
	{name: "multi_strings", args: []string{"repeat", "-n", "5", "a", "b", "c"}, wantExit: 0, wantOut: []string{"aaaaa", "bbbbb", "ccccc"}},
	// L644: string repeat -n 5 --max 4 123 456 789
	{name: "multi_with_max", args: []string{"repeat", "-n", "5", "--max", "4", "123", "456", "789"}, wantExit: 0, wantOut: []string{"1231", "4564", "7897"}},
	// L649: string repeat -n 5 --max 4 123 '' 789
	{name: "multi_empty_with_max", args: []string{"repeat", "-n", "5", "--max", "4", "123", "", "789"}, wantExit: 0, wantOut: []string{"1231", "", "7897"}},
	// L676: string repeat -n-1 foo
	{name: "neg_count_error", args: []string{"repeat", "-n-1", "foo"}, wantExit: 1, wantErr: "Invalid count value '-1'"},
	// L679: string repeat -m-1 foo
	{name: "neg_max_error", args: []string{"repeat", "-m-1", "foo"}, wantExit: 1, wantErr: "Invalid max value '-1'"},
	// L701: string repeat -n3 ""
	{name: "empty_string_exit1", args: []string{"repeat", "-n3", ""}, wantExit: 1},
}

var fishParityShorten = []cliTest{
	// L1061: string shorten -m 3 foo
	{name: "no_truncation", args: []string{"shorten", "-m", "3", "foo"}, wantExit: 0, wantOut: []string{"foo"}},
	// L1063: string shorten -m 2 foo
	{name: "basic", args: []string{"shorten", "-m", "2", "foo"}, wantExit: 0, wantOut: []string{"f…"}},
	// L1066: string shorten -m 5 foobar
	{name: "truncate_to_5", args: []string{"shorten", "-m", "5", "foobar"}, wantExit: 0, wantOut: []string{"foob…"}},
	// L1070: string shorten -m 2 -q 12
	{name: "quiet_no_change_exit1", args: []string{"shorten", "-m", "2", "-q", "12"}, wantExit: 1, wantOut: []string{}},
	// L1073: string shorten -lm 2 -q 12
	{name: "left_quiet_no_change_exit1", args: []string{"shorten", "-l", "-m", "2", "-q", "12"}, wantExit: 1, wantOut: []string{}},
	// L1078: string shorten -m 5 --char ........ foobar
	{name: "ellipsis_longer_than_width", args: []string{"shorten", "-m", "5", "--char", "........", "foobar"}, wantExit: 0, wantOut: []string{"fooba"}},
	// L1081: string shorten --max 4 -c /// foobar
	{name: "custom_ellipsis_3char", args: []string{"shorten", "--max", "4", "-c", "///", "foobar"}, wantExit: 0, wantOut: []string{"f///"}},
	// L1087: string shorten --max 2 --char "" foo
	{name: "empty_ellipsis", args: []string{"shorten", "--max", "2", "--char", "", "foo"}, wantExit: 0, wantOut: []string{"fo"}},
	// L1090: string shorten --max=-1 --char "" foo
	{name: "neg_max_error", args: []string{"shorten", "--max=-1", "--char", "", "foo"}, wantExit: 1, wantErr: "Invalid max value '-1'"},
	// L1093: string shorten foo foobar
	{name: "auto_width", args: []string{"shorten", "foo", "foobar"}, wantExit: 0, wantOut: []string{"foo", "fo…"}},
	// L1190: string shorten -m4 foobar bananarama
	{name: "multi_same_max", args: []string{"shorten", "-m4", "foobar", "bananarama"}, wantExit: 0, wantOut: []string{"foo…", "ban…"}},
	// L1226: string shorten -m0 foo bar asodjsaoidj
	{name: "max0_all_as_is", args: []string{"shorten", "-m0", "foo", "bar", "asodjsaoidj"}, wantExit: 0, wantOut: []string{"foo", "bar", "asodjsaoidj"}},
}

// fishParityQuietEarlyExit verifies --quiet exits 0 on first match without consuming all stdin.
// Uses real infinite pipes so the test actually hangs if early-exit is broken.
// L993-L1003: yes | string match -q y / string length -q / string replace -q y n
var fishParityQuietEarlyExit = []struct{ name string; args []string }{
	{"match_q", []string{"match", "-q", "y"}},
	{"length_q", []string{"length", "-q"}},
	{"replace_q", []string{"replace", "-q", "y", "n"}},
}

func TestFishParity(t *testing.T) {
	run := func(group string, tests []cliTest) {
		t.Run(group, func(t *testing.T) {
			for _, tc := range tests {
				t.Run(tc.name, func(t *testing.T) {
					runBin(t, tc)
				})
			}
		})
	}
	run("general", fishParityGeneral)
	run("length", fishParityLength)
	run("pad", fishParityPad)
	run("sub", fishParitySub)
	run("split", fishParitySplit)
	run("join", fishParityJoin)
	run("trim", fishParityTrim)
	run("escape", fishParityEscape)
	run("unescape", fishParityUnescape)
	run("match", fishParityMatch)
	run("replace", fishParityReplace)
	run("repeat", fishParityRepeat)
	run("shorten", fishParityShorten)
}

// TestQuietEarlyExit uses real infinite pipes to prove --quiet exits without consuming all stdin.
// If early-exit is broken, each sub-test hangs until the 2s timeout kills it and fails.
func TestQuietEarlyExit(t *testing.T) {
	for _, tc := range fishParityQuietEarlyExit {
		t.Run(tc.name, func(t *testing.T) {
			pr, pw := io.Pipe()
			go func() {
				for {
					if _, err := pw.Write([]byte("y\n")); err != nil {
						return
					}
				}
			}()

			cmd := exec.Command(stringBin, tc.args...)
			cmd.Stdin = pr

			done := make(chan error, 1)
			go func() { done <- cmd.Run() }()

			select {
			case err := <-done:
				pr.Close()
				pw.Close()
				if ee, ok := err.(*exec.ExitError); ok {
					if ee.ExitCode() != 0 {
						t.Errorf("exit: want 0, got %d", ee.ExitCode())
					}
				} else if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			case <-time.After(2 * time.Second):
				cmd.Process.Kill()
				pr.Close()
				pw.Close()
				t.Error("command hung reading infinite stdin; --quiet early-exit broken")
			}
		})
	}
}

package argparsecmd

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// argparseBin is the path to the compiled binary, set by TestMain.
var argparseBin string

func TestMain(m *testing.M) {
	tmp, err := os.CreateTemp("", "argparse-test-*")
	if err != nil {
		panic(err)
	}
	tmp.Close()
	argparseBin = tmp.Name()
	defer os.Remove(argparseBin)

	if out, err := exec.Command("go", "build", "-o", argparseBin, "github.com/mattmc3/mudskipper/cmd/argparse").CombinedOutput(); err != nil {
		panic("build failed: " + string(out))
	}

	os.Exit(m.Run())
}

type cliTest struct {
	name     string
	stdin    string
	args     []string
	wantExit int
	wantOut  []string // nil = don't check; []string{} = assert empty
	wantErr  string   // substring match; empty = don't check
}

func runBin(t *testing.T, tc cliTest) {
	t.Helper()
	cmd := exec.Command(argparseBin, tc.args...)
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
		if len(tc.wantOut) == 0 {
			if outBuf.Len() != 0 {
				t.Errorf("stdout: want empty, got %q", outBuf.String())
			}
		} else {
			for _, want := range tc.wantOut {
				if !strings.Contains(outBuf.String(), want) {
					t.Errorf("stdout: want %q in output, got:\n%s", want, outBuf.String())
				}
			}
		}
	}

	if tc.wantErr != "" {
		if !strings.Contains(errBuf.String(), tc.wantErr) {
			t.Errorf("stderr: want %q in %q", tc.wantErr, errBuf.String())
		}
	}
}

// argparse CLI tests — real binary execution, verifying eval-able output

var argparseDispatch = []cliTest{
	{name: "help_short", args: []string{"-h"}, wantExit: 0, wantOut: []string{"Usage:"}},
	{name: "help_long", args: []string{"--help"}, wantExit: 0, wantOut: []string{"Usage:"}},
	{name: "missing_separator", args: []string{"h/help"}, wantExit: 1, wantErr: "Missing -- separator"},
	{name: "no_specs_no_args", args: []string{"--shell=bash", "--"}, wantExit: 0, wantOut: []string{"set --"}},
	{name: "empty_args", args: []string{"--shell=bash", "h/help", "--"}, wantExit: 0, wantOut: []string{"set --"}},
}

var argparseBoolFlags = []cliTest{
	// Both _flag_h and _flag_help set — key fish compat
	{name: "both_vars_set_long", args: []string{"--shell=bash", "h/help", "--", "--help"}, wantExit: 0,
		wantOut: []string{"_flag_help=--help", "_flag_h=--help"}},
	{name: "both_vars_set_short", args: []string{"--shell=bash", "h/help", "--", "-h"}, wantExit: 0,
		wantOut: []string{"_flag_help=-h", "_flag_h=-h"}},
	// Value is the flag string seen, not "1"
	{name: "bool_value_is_flag_string", args: []string{"--shell=bash", "v/verbose", "--", "--verbose"}, wantExit: 0,
		wantOut: []string{"_flag_verbose=--verbose"}},
	// Short-only spec
	{name: "short_only", args: []string{"--shell=bash", "v", "--", "-v"}, wantExit: 0,
		wantOut: []string{"_flag_v=-v"}},
	// Long-only spec
	{name: "long_only", args: []string{"--shell=bash", "/verbose", "--", "--verbose"}, wantExit: 0,
		wantOut: []string{"_flag_verbose=--verbose"}},
}

var argparseStringFlags = []cliTest{
	{name: "string_long", args: []string{"--shell=bash", "n/name=", "--", "--name", "alice"}, wantExit: 0,
		wantOut: []string{"_flag_name=alice", "_flag_n=alice"}},
	{name: "string_short", args: []string{"--shell=bash", "n/name=", "--", "-n", "alice"}, wantExit: 0,
		wantOut: []string{"_flag_name=alice"}},
	{name: "string_equals_syntax", args: []string{"--shell=bash", "n/name=", "--", "--name=alice"}, wantExit: 0,
		wantOut: []string{"_flag_name=alice"}},
	{name: "string_quoted", args: []string{"--shell=bash", "n/name=", "--", "--name", "hello world"}, wantExit: 0,
		wantOut: []string{"'hello world'"}},
	{name: "string_missing_value_error", args: []string{"--shell=bash", "n/name=", "--", "--name"}, wantExit: 1,
		wantErr: "option requires an argument"},
}

var argparseMultiFlags = []cliTest{
	{name: "multi_repeat", args: []string{"--shell=bash", "t/tag=+", "--", "--tag", "a", "--tag", "b"}, wantExit: 0,
		wantOut: []string{"_flag_tag"}},
	{name: "optional_with_value", args: []string{"--shell=bash", "c/color=?", "--", "--color=red"}, wantExit: 0,
		wantOut: []string{"_flag_color=red"}},
	{name: "optional_without_value", args: []string{"--shell=bash", "c/color=?", "--", "--color"}, wantExit: 0,
		wantOut: []string{"_flag_color="}},
}

var argparseRemaining = []cliTest{
	{name: "remaining_after_flags", args: []string{"--shell=bash", "v/verbose", "--", "-v", "foo", "bar"}, wantExit: 0,
		wantOut: []string{"set -- foo bar"}},
	{name: "double_dash_stops", args: []string{"--shell=bash", "v/verbose", "--", "-v", "--", "--not-a-flag"}, wantExit: 0,
		wantOut: []string{"--not-a-flag"}},
	{name: "stop_nonopt", args: []string{"--shell=bash", "--stop-nonopt", "v/verbose", "--", "foo", "-v"}, wantExit: 0,
		wantOut: []string{"set -- foo -v"}},
}

var argparseLimits = []cliTest{
	{name: "min_args_pass", args: []string{"--shell=bash", "--min-args=1", "--", "foo"}, wantExit: 0},
	{name: "min_args_fail", args: []string{"--shell=bash", "--min-args=2", "--", "one"}, wantExit: 1,
		wantErr: "at least 2"},
	{name: "max_args_pass", args: []string{"--shell=bash", "--max-args=2", "--", "a", "b"}, wantExit: 0},
	{name: "max_args_fail", args: []string{"--shell=bash", "--max-args=1", "--", "a", "b", "c"}, wantExit: 1,
		wantErr: "at most 1"},
}

var argparseShells = []cliTest{
	{name: "bash", args: []string{"--shell=bash", "v/verbose", "n/name=", "--", "--verbose", "--name", "bob", "rest"}, wantExit: 0,
		wantOut: []string{"_flag_verbose=--verbose", "_flag_name=bob", "set -- rest"}},
	{name: "zsh", args: []string{"--shell=zsh", "v/verbose", "n/name=", "--", "-v", "--name", "bob"}, wantExit: 0,
		wantOut: []string{"_flag_verbose=-v", "_flag_name=bob", "set --"}},
	{name: "fish", args: []string{"--shell=fish", "v/verbose", "n/name=", "--", "-v", "--name", "bob", "rest"}, wantExit: 0,
		wantOut: []string{"set -l _flag_verbose", "set -l _flag_v", "set -l _flag_name 'bob'", "set -l argv 'rest'"}},
	{name: "elvish", args: []string{"--shell=elvish", "v/verbose", "--", "--verbose"}, wantExit: 0,
		wantOut: []string{"var _flag_verbose = $true"}},
}

var argparseErrors = []cliTest{
	{name: "unknown_long_flag", args: []string{"--shell=bash", "v/verbose", "--", "--unknown"}, wantExit: 1,
		wantErr: "unknown option"},
	{name: "unknown_short_flag", args: []string{"--shell=bash", "v/verbose", "--", "-x"}, wantExit: 1,
		wantErr: "unknown option"},
	{name: "named_error_message", args: []string{"--shell=bash", "--name=myscript", "v/verbose", "--", "-x"}, wantExit: 1,
		wantErr: "myscript"},
}

// Reference tests — equivalent to https://github.com/fish-shell/fish-shell/blob/master/tests/checks/argparse.fish
// Not matching fish output exactly (different shell, no argv_opts, local scoping differs).
// Tests use --shell=bash --no-local for readable output without the `local` keyword.
// Skipped: implicit int flags (#-val), !validation scripts, --ignore-unknown, --move-unknown.

var argparseRef = []cliTest{
	// argparse.fish L12: argparse (no args, no --) → Missing -- separator
	{name: "L12_no_args_no_sep", args: []string{}, wantExit: 1, wantErr: "Missing -- separator"},
	// argparse.fish L20: argparse h/help (no --) → Missing -- separator
	{name: "L20_spec_no_sep", args: []string{"h/help"}, wantExit: 1, wantErr: "Missing -- separator"},
	// argparse.fish L74: argparse h/ → empty long name → invalid spec
	{name: "L74_invalid_spec_empty_long", args: []string{"h/", "--"}, wantExit: 1, wantErr: "invalid option spec"},

	// argparse.fish L82-98: --min-args and --max-args
	{name: "L84_min_args_fail", args: []string{"--shell=bash", "--no-local", "--name", "min-max", "--min-args", "1", "h/help", "--"}, wantExit: 1, wantErr: "at least 1"},
	{name: "L87_min_max_pass", args: []string{"--shell=bash", "--no-local", "--min-args", "1", "--max-args", "3", "h/help", "--", "arg1"}, wantExit: 0},
	{name: "L91_max_args_fail", args: []string{"--shell=bash", "--no-local", "--name", "min-max", "--min-args", "1", "--max-args", "3", "h/help", "--", "arg1", "arg2", "-h", "arg3", "arg4"}, wantExit: 1, wantErr: "at most 3"},
	{name: "L95_max_args_fail2", args: []string{"--shell=bash", "--no-local", "--name", "min-max", "--max-args", "1", "h/help", "--", "arg1", "arg2"}, wantExit: 1, wantErr: "at most 1"},

	// argparse.fish L180-183: no args passes
	{name: "L180_no_args_passes", args: []string{"--shell=bash", "--no-local", "h/help", "--"}, wantExit: 0},
	// argparse.fish L185-191: one non-flag arg → remaining
	{name: "L185_one_nonopt_arg", args: []string{"--shell=bash", "--no-local", "h/help", "--", "help"}, wantExit: 0, wantOut: []string{"set -- help"}},

	// argparse.fish L193-201: five args, two matching --help and -h
	// fish: _flag_h = '--help' '-h'; _flag_help = '--help' '-h'; argv = help me 'a lot more'
	// ours: both vars set, value = last seen flag; remaining has 3 args
	{name: "L193_bool_repeated", args: []string{"--shell=bash", "--no-local", "h/help", "--", "help", "--help", "me", "-h", "a lot more"}, wantExit: 0,
		wantOut: []string{"_flag_help", "_flag_h", "set -- help me 'a lot more'"}},

	// argparse.fish L203-217: required, optional, and multiple flags
	{name: "L203_required_optional_multi", args: []string{"--shell=bash", "--no-local", "h/help", "a/abc=", "d/def=?", "g/ghk=+", "--",
		"help", "--help", "me", "--ghk=g1", "--abc=ABC", "--ghk", "g2", "-d", "-g", "g3"}, wantExit: 0,
		wantOut: []string{"_flag_abc=ABC", "_flag_a=ABC", "_flag_g", "_flag_help", "set -- help me"}},

	// argparse.fish L219-229: --stop-nonopt
	{name: "L219_stop_nonopt", args: []string{"--shell=bash", "--no-local", "--stop-nonopt", "h/help", "a/abc=", "--",
		"-a", "A1", "-h", "--abc", "A2", "non-opt", "second non-opt", "--help"}, wantExit: 0,
		wantOut: []string{"_flag_abc=A2", "_flag_h=-h", "set -- non-opt 'second non-opt' --help"}},

	// argparse.fish L270-279: short-only bool flags
	{name: "L270_short_only_bool", args: []string{"--shell=bash", "--no-local", "C", "v", "--", "-C", "-v", "arg1", "-v", "arg2"}, wantExit: 0,
		wantOut: []string{"_flag_C=-C", "_flag_v=-v", "set -- arg1 arg2"}},

	// argparse.fish L281-290: short-only string flag + verbose
	{name: "L281_short_only_string", args: []string{"--shell=bash", "--no-local", "x=", "v/verbose", "--", "--verbose", "arg1", "-v", "-x", "arg2"}, wantExit: 0,
		wantOut: []string{"_flag_x=arg2", "_flag_verbose", "set -- arg1"}},

	// argparse.fish L359-364: --name used in error messages
	{name: "L359_name_in_errors", args: []string{"--shell=bash", "--no-local", "--name=myscript", "a/alpha", "--", "--banana"}, wantExit: 1, wantErr: "myscript"},

	// argparse.fish L368-379: --ignore-unknown (-i) passes unknown options to remaining
	{name: "L368_ignore_unknown", args: []string{"--shell=bash", "--no-local", "-i", "a=+", "b=+", "--",
		"-a", "alpha", "-b", "bravo", "-t", "tango", "-a", "aaaa", "--wurst"}, wantExit: 0,
		wantOut: []string{"_flag_a", "alpha", "aaaa", "set -- -t tango --wurst"}},
	{name: "L368_ignore_unknown_exit0", args: []string{"--shell=bash", "--no-local", "-i", "a=+", "b=+", "--",
		"-a", "alpha", "-t", "tango"}, wantExit: 0},
}

func TestArgparseCLI(t *testing.T) {
	run := func(group string, tests []cliTest) {
		t.Run(group, func(t *testing.T) {
			for _, tc := range tests {
				t.Run(tc.name, func(t *testing.T) {
					runBin(t, tc)
				})
			}
		})
	}
	run("dispatch", argparseDispatch)
	run("bool", argparseBoolFlags)
	run("string", argparseStringFlags)
	run("multi", argparseMultiFlags)
	run("remaining", argparseRemaining)
	run("limits", argparseLimits)
	run("shells", argparseShells)
	run("errors", argparseErrors)
	run("reference", argparseRef)
}

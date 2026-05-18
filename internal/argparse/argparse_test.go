package argparsecmd

import (
	"strings"
	"testing"
)

func run(args ...string) (int, string, string) {
	var out, err strings.Builder
	code := Run(args, strings.NewReader(""), &out, &err)
	return code, out.String(), err.String()
}

func TestArgparse_missing_separator(t *testing.T) {
	code, _, stderr := run("h/help")
	if code != 1 {
		t.Errorf("exit: want 1, got %d", code)
	}
	if !strings.Contains(stderr, "Missing -- separator") {
		t.Errorf("stderr: want 'Missing -- separator', got %q", stderr)
	}
}

func TestArgparse_no_specs_no_args(t *testing.T) {
	code, out, _ := run("--shell=bash", "--")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if out != "set --\n" {
		t.Errorf("stdout: want 'set --\\n', got %q", out)
	}
}

func TestArgparse_bool_flag_long(t *testing.T) {
	// Default: local scoping. Both _flag_verbose and _flag_v set.
	code, out, _ := run("--shell=bash", "v/verbose", "--", "--verbose")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if !strings.Contains(out, "local _flag_verbose=--verbose") {
		t.Errorf("stdout: want 'local _flag_verbose=--verbose', got %q", out)
	}
	if !strings.Contains(out, "local _flag_v=--verbose") {
		t.Errorf("stdout: want 'local _flag_v=--verbose', got %q", out)
	}
}

func TestArgparse_bool_flag_short(t *testing.T) {
	code, out, _ := run("--shell=bash", "v/verbose", "--", "-v")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if !strings.Contains(out, "local _flag_verbose=-v") {
		t.Errorf("stdout: want 'local _flag_verbose=-v', got %q", out)
	}
}

func TestArgparse_bash_no_local(t *testing.T) {
	code, out, _ := run("--shell=bash", "--no-local", "v/verbose", "--", "--verbose")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if strings.Contains(out, "local ") {
		t.Errorf("stdout: --no-local should not use local, got %q", out)
	}
	if !strings.Contains(out, "_flag_verbose=--verbose") {
		t.Errorf("stdout: want _flag_verbose=--verbose, got %q", out)
	}
}

func TestArgparse_string_flag(t *testing.T) {
	code, out, _ := run("--shell=bash", "n/name=", "--", "--name", "alice")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if !strings.Contains(out, "_flag_name=alice") {
		t.Errorf("stdout: want _flag_name=alice, got %q", out)
	}
}

func TestArgparse_string_flag_equals(t *testing.T) {
	code, out, _ := run("--shell=bash", "n/name=", "--", "--name=alice")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if !strings.Contains(out, "_flag_name=alice") {
		t.Errorf("stdout: want _flag_name=alice, got %q", out)
	}
}

func TestArgparse_remaining_args(t *testing.T) {
	code, out, _ := run("--shell=bash", "v/verbose", "--", "-v", "foo", "bar")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if !strings.Contains(out, "set -- foo bar") {
		t.Errorf("stdout: want 'set -- foo bar', got %q", out)
	}
}

func TestArgparse_unknown_flag_error(t *testing.T) {
	code, _, stderr := run("--shell=bash", "v/verbose", "--", "--unknown")
	if code != 1 {
		t.Errorf("exit: want 1, got %d", code)
	}
	if !strings.Contains(stderr, "unknown option") {
		t.Errorf("stderr: want 'unknown option', got %q", stderr)
	}
}

func TestArgparse_fish_output(t *testing.T) {
	// Default: set -l (local)
	code, out, _ := run("--shell=fish", "v/verbose", "n/name=", "--", "-v", "--name", "bob", "rest")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if !strings.Contains(out, "set -l _flag_verbose") {
		t.Errorf("stdout: missing 'set -l _flag_verbose', got %q", out)
	}
	if !strings.Contains(out, "set -l _flag_v") {
		t.Errorf("stdout: missing 'set -l _flag_v', got %q", out)
	}
	if !strings.Contains(out, "set -l _flag_name 'bob'") {
		t.Errorf("stdout: missing set -l _flag_name, got %q", out)
	}
	if !strings.Contains(out, "set -l argv 'rest'") {
		t.Errorf("stdout: missing 'set -l argv', got %q", out)
	}
}

func TestArgparse_fish_no_local(t *testing.T) {
	code, out, _ := run("--shell=fish", "--no-local", "v/verbose", "--", "-v")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if strings.Contains(out, "set -l") {
		t.Errorf("stdout: --no-local should not use set -l, got %q", out)
	}
	if !strings.Contains(out, "set _flag_verbose") {
		t.Errorf("stdout: want 'set _flag_verbose', got %q", out)
	}
}

func TestArgparse_elvish_output(t *testing.T) {
	// Default: var (local)
	code, out, _ := run("--shell=elvish", "v/verbose", "--", "--verbose")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if !strings.Contains(out, "var _flag_verbose = $true") {
		t.Errorf("stdout: want 'var _flag_verbose = $true', got %q", out)
	}
}

func TestArgparse_elvish_no_local(t *testing.T) {
	code, out, _ := run("--shell=elvish", "--no-local", "v/verbose", "--", "--verbose")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if !strings.Contains(out, "set _flag_verbose = $true") {
		t.Errorf("stdout: want 'set _flag_verbose = $true', got %q", out)
	}
}

func TestArgparse_multi_flag(t *testing.T) {
	code, out, _ := run("--shell=bash", "t/tag=+", "--", "--tag", "a", "--tag", "b")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if !strings.Contains(out, "_flag_tag=") {
		t.Errorf("stdout: want _flag_tag, got %q", out)
	}
}

func TestArgparse_quoted_value(t *testing.T) {
	code, out, _ := run("--shell=bash", "n/name=", "--", "--name", "hello world")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if !strings.Contains(out, "'hello world'") {
		t.Errorf("stdout: want quoted value, got %q", out)
	}
}

func TestArgparse_double_dash_stops_parsing(t *testing.T) {
	code, out, _ := run("--shell=bash", "v/verbose", "--", "--verbose", "--", "--not-a-flag")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if !strings.Contains(out, "--not-a-flag") {
		t.Errorf("stdout: want --not-a-flag in remaining, got %q", out)
	}
}

func TestArgparse_short_only_spec(t *testing.T) {
	code, out, _ := run("--shell=bash", "v", "--", "-v")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if !strings.Contains(out, "_flag_v=-v") {
		t.Errorf("stdout: want _flag_v=-v, got %q", out)
	}
}

func TestArgparse_long_only_spec(t *testing.T) {
	code, out, _ := run("--shell=bash", "/verbose", "--", "--verbose")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if !strings.Contains(out, "_flag_verbose=--verbose") {
		t.Errorf("stdout: want _flag_verbose=--verbose, got %q", out)
	}
}

func TestArgparse_stop_nonopt(t *testing.T) {
	code, out, _ := run("--shell=bash", "--stop-nonopt", "v/verbose", "--", "foo", "-v")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	// -v after non-option should be in remaining, not parsed
	if strings.Contains(out, "_flag_verbose") {
		t.Errorf("stdout: -v after nonopt should not be parsed, got %q", out)
	}
	if !strings.Contains(out, "foo") {
		t.Errorf("stdout: foo should be in remaining, got %q", out)
	}
}

func TestArgparse_min_args_error(t *testing.T) {
	code, _, stderr := run("--shell=bash", "--min-args=2", "--", "only-one")
	if code != 1 {
		t.Errorf("exit: want 1, got %d", code)
	}
	if !strings.Contains(stderr, "at least 2") {
		t.Errorf("stderr: want 'at least 2', got %q", stderr)
	}
}

func TestArgparse_max_args_error(t *testing.T) {
	code, _, stderr := run("--shell=bash", "--max-args=1", "--", "a", "b", "c")
	if code != 1 {
		t.Errorf("exit: want 1, got %d", code)
	}
	if !strings.Contains(stderr, "at most 1") {
		t.Errorf("stderr: want 'at most 1', got %q", stderr)
	}
}

func TestArgparse_both_short_and_long_set(t *testing.T) {
	// Key fish compat: h/help sets both _flag_h and _flag_help
	code, out, _ := run("--shell=bash", "h/help", "--", "--help")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if !strings.Contains(out, "_flag_help=--help") {
		t.Errorf("stdout: want _flag_help, got %q", out)
	}
	if !strings.Contains(out, "_flag_h=--help") {
		t.Errorf("stdout: want _flag_h (both set), got %q", out)
	}
}

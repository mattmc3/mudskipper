package pathcmd

import (
	"io"
	"strings"
	"testing"
)

func runCmd(args ...string) (int, string, string) {
	var out, err strings.Builder
	code := Run(args, strings.NewReader(""), &out, &err)
	return code, out.String(), err.String()
}

func TestPath_help(t *testing.T) {
	for _, flag := range []string{"-h", "--help"} {
		code, out, _ := runCmd(flag)
		if code != 0 {
			t.Errorf("%s: exit want 0, got %d", flag, code)
		}
		if !strings.Contains(out, "Usage:") {
			t.Errorf("%s: stdout want Usage, got %q", flag, out)
		}
	}
}

func TestPath_no_subcommand(t *testing.T) {
	code, _, stderr := runCmd()
	if code != 1 {
		t.Errorf("exit: want 1, got %d", code)
	}
	if !strings.Contains(stderr, "missing subcommand") {
		t.Errorf("stderr: want 'missing subcommand', got %q", stderr)
	}
}

func TestPath_invalid_subcommand(t *testing.T) {
	code, _, stderr := runCmd("invalid-subcmd")
	if code != 1 {
		t.Errorf("exit: want 1, got %d", code)
	}
	if !strings.Contains(stderr, "invalid subcommand") {
		t.Errorf("stderr: want 'invalid subcommand', got %q", stderr)
	}
}

// extension tests

func TestExtension(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"/", ""},
		{"/.", ""},
		{"/.foo", ""},
		{"/foo", ""},
		{"/foo.txt", ".txt"},
		{"/foo.txt/bar", ""},
		{".", ""},
		{"..", ""},
		{"./foo.mp4", ".mp4"},
		{"../banana", ""},
		{"~/.config", ""},
		{"~/.config.d", ".d"},
		{"~/.config.", "."},
	}
	for _, tt := range tests {
		got := extension(tt.input)
		if got != tt.want {
			t.Errorf("extension(%q): want %q, got %q", tt.input, tt.want, got)
		}
	}
}

// basename tests

func TestBasename_basic(t *testing.T) {
	code, out, _ := runCmd("basename", "./foo.mp4")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if out != "foo.mp4\n" {
		t.Errorf("stdout: want foo.mp4, got %q", out)
	}
}

func TestBasename_strip_extension(t *testing.T) {
	code, out, _ := runCmd("basename", "-E", "./foo.mp4")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if out != "foo\n" {
		t.Errorf("stdout: want foo, got %q", out)
	}
}

func TestBasename_trailing_slash(t *testing.T) {
	code, out, _ := runCmd("basename", "/usr/bin/")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if out != "bin\n" {
		t.Errorf("stdout: want bin, got %q", out)
	}
}

// dirname tests

func TestDirname_basic(t *testing.T) {
	code, out, _ := runCmd("dirname", "./foo.mp4")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if out != ".\n" {
		t.Errorf("stdout: want ., got %q", out)
	}
}

// normalize tests

func TestNormalize_double_slash(t *testing.T) {
	code, out, _ := runCmd("normalize", "/usr/bin//../../etc/fish")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if out != "/etc/fish\n" {
		t.Errorf("stdout: want /etc/fish, got %q", out)
	}
}

func TestNormalize_dash_prefix(t *testing.T) {
	code, out, _ := runCmd("normalize", "--", "-/foo")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if out != "./-/foo\n" {
		t.Errorf("stdout: want ./-/foo, got %q", out)
	}
}

func TestNormalize_relative_dash(t *testing.T) {
	code, out, _ := runCmd("normalize", "--", "../-foo")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if out != "../-foo\n" {
		t.Errorf("stdout: want ../-foo, got %q", out)
	}
}

// sort tests

func TestSort_basic(t *testing.T) {
	code, out, _ := runCmd("sort", "c", "a", "b")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if out != "a\nb\nc\n" {
		t.Errorf("stdout: want a/b/c, got %q", out)
	}
}

func TestSort_reverse(t *testing.T) {
	code, out, _ := runCmd("sort", "-r", "a", "b", "c")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if out != "c\nb\na\n" {
		t.Errorf("stdout: want c/b/a, got %q", out)
	}
}

func TestSort_unique(t *testing.T) {
	code, out, _ := runCmd("sort", "-u", "b", "a", "b", "c", "a")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if out != "a\nb\nc\n" {
		t.Errorf("stdout: want a/b/c, got %q", out)
	}
}

func TestSort_key_basename(t *testing.T) {
	code, out, _ := runCmd("sort", "--key=basename", "z/b", "a/c", "m/a")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if out != "m/a\nz/b\na/c\n" {
		t.Errorf("stdout: got %q", out)
	}
}

func TestSort_invalid_key(t *testing.T) {
	code, _, stderr := runCmd("sort", "--key=invalid-key")
	if code != 1 {
		t.Errorf("exit: want 1, got %d", code)
	}
	if !strings.Contains(stderr, "Invalid sort key") {
		t.Errorf("stderr: want 'Invalid sort key', got %q", stderr)
	}
}

// change-extension tests

func TestChangeExtension_replace(t *testing.T) {
	code, out, _ := runCmd("change-extension", "wmv", "./foo.mp4")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if out != "./foo.wmv\n" {
		t.Errorf("stdout: want ./foo.wmv, got %q", out)
	}
}

func TestChangeExtension_remove(t *testing.T) {
	code, out, _ := runCmd("change-extension", "", "./foo.mp4")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if out != "./foo\n" {
		t.Errorf("stdout: want ./foo, got %q", out)
	}
}

func TestChangeExtension_no_ext_input(t *testing.T) {
	code, out, _ := runCmd("change-extension", "", "../banana")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if out != "../banana\n" {
		t.Errorf("stdout: want ../banana, got %q", out)
	}
}

// Pipe-chain tests — verify subcommands compose via stdin/stdout

func pipe(args1 []string, args2 []string) (int, string) {
	var mid, out strings.Builder
	Run(args1, strings.NewReader(""), &mid, io.Discard)
	code := Run(args2, strings.NewReader(mid.String()), &out, io.Discard)
	return code, out.String()
}

func TestPipeChain_dirname_twice(t *testing.T) {
	// echo /a/b/c/d | path dirname | path dirname → /a/b
	code, out := pipe([]string{"dirname", "/a/b/c/d"}, []string{"dirname"})
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if out != "/a/b\n" {
		t.Errorf("stdout: want /a/b, got %q", out)
	}
}

func TestPipeChain_basename_then_strip_extension(t *testing.T) {
	// echo ./releases/foo-1.2.3.tar.gz | path basename | path basename -E → foo-1.2.3.tar
	code, out := pipe([]string{"basename", "./releases/foo-1.2.3.tar.gz"}, []string{"basename", "-E"})
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if out != "foo-1.2.3.tar\n" {
		t.Errorf("stdout: want foo-1.2.3.tar, got %q", out)
	}
}

func TestPipeChain_resolve_dotdot(t *testing.T) {
	// echo /a/b/c/../.. | path resolve → /a  (resolve handles .. without filesystem for abs paths)
	// Uses normalize instead since we don't have /a/b on disk
	code, out := pipe([]string{"normalize", "/a/b/c/../.."}, []string{"normalize"})
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if out != "/a\n" {
		t.Errorf("stdout: want /a, got %q", out)
	}
}

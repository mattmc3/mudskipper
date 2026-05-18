package pathcmd

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// pathBin is the path to the compiled binary, set by TestMain.
var pathBin string

func TestMain(m *testing.M) {
	tmp, err := os.CreateTemp("", "path-test-*")
	if err != nil {
		panic(err)
	}
	tmp.Close()
	pathBin = tmp.Name()
	defer os.Remove(pathBin)

	if out, err := exec.Command("go", "build", "-o", pathBin, "github.com/mattmc3/mudskipper/cmd/path").CombinedOutput(); err != nil {
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
	cmd := exec.Command(pathBin, tc.args...)
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
			got := strings.TrimRight(outBuf.String(), "\n")
			want := strings.Join(tc.wantOut, "\n")
			if got != want {
				t.Errorf("stdout: want %q, got %q", want, got)
			}
		}
	}

	if tc.wantErr != "" {
		if !strings.Contains(errBuf.String(), tc.wantErr) {
			t.Errorf("stderr: want %q in %q", tc.wantErr, errBuf.String())
		}
	}
}

// Fish parity tests — direct from https://github.com/fish-shell/fish-shell/blob/master/tests/checks/path.fish

var fishParityExtension = []cliTest{
	// path.fish L7: path extension /
	{name: "root_no_ext", args: []string{"extension", "/"}, wantExit: 1, wantOut: []string{""}},
	// path.fish L13: path extension /.
	{name: "dot_no_ext", args: []string{"extension", "/."}, wantExit: 1, wantOut: []string{""}},
	// path.fish L19: path extension /.foo
	{name: "hidden_no_ext", args: []string{"extension", "/.foo"}, wantExit: 1, wantOut: []string{""}},
	// path.fish L24: path extension /foo
	{name: "no_ext", args: []string{"extension", "/foo"}, wantExit: 1, wantOut: []string{""}},
	// path.fish L28: path extension /foo.txt
	{name: "txt_ext", args: []string{"extension", "/foo.txt"}, wantExit: 0, wantOut: []string{".txt"}},
	// path.fish L32: path extension /foo.txt/bar
	{name: "dir_component_no_ext", args: []string{"extension", "/foo.txt/bar"}, wantExit: 1, wantOut: []string{""}},
	// path.fish L36: path extension . ..  (both print empty lines, exit 1)
	{name: "dot_dotdot_no_ext", args: []string{"extension", ".", ".."}, wantExit: 1},
	// path.fish L40: path extension ./foo.mp4
	{name: "relative_mp4", args: []string{"extension", "./foo.mp4"}, wantExit: 0, wantOut: []string{".mp4"}},
	// path.fish L42: path extension ../banana
	{name: "relative_no_ext", args: []string{"extension", "../banana"}, wantExit: 1, wantOut: []string{""}},
	// path.fish L54: path extension ~/.config.
	{name: "trailing_dot", args: []string{"extension", "~/.config."}, wantExit: 0, wantOut: []string{"."}},
}

var fishParityChangeExtension = []cliTest{
	// path.fish L60: path change-extension '' ./foo.mp4
	{name: "remove_ext", args: []string{"change-extension", "", "./foo.mp4"}, wantExit: 0, wantOut: []string{"./foo"}},
	// path.fish L62: path change-extension wmv ./foo.mp4
	{name: "replace_ext_no_dot", args: []string{"change-extension", "wmv", "./foo.mp4"}, wantExit: 0, wantOut: []string{"./foo.wmv"}},
	// path.fish L64: path change-extension .wmv ./foo.mp4
	{name: "replace_ext_with_dot", args: []string{"change-extension", ".wmv", "./foo.mp4"}, wantExit: 0, wantOut: []string{"./foo.wmv"}},
	// path.fish L66: path change-extension '' ../banana
	{name: "remove_ext_no_ext", args: []string{"change-extension", "", "../banana"}, wantExit: 0, wantOut: []string{"../banana"}},
}

var fishParityBasename = []cliTest{
	// path.fish L76: path basename ./foo.mp4
	{name: "relative", args: []string{"basename", "./foo.mp4"}, wantExit: 0, wantOut: []string{"foo.mp4"}},
	// path.fish L78: path basename ../banana
	{name: "parent_rel", args: []string{"basename", "../banana"}, wantExit: 0, wantOut: []string{"banana"}},
	// path.fish L80: path basename /usr/bin/
	{name: "trailing_slash", args: []string{"basename", "/usr/bin/"}, wantExit: 0, wantOut: []string{"bin"}},
	// path.fish L437: path basename -E foo.txt /usr/local/foo.bar /foo.tar.gz
	{name: "strip_ext_multi", args: []string{"basename", "-E", "foo.txt", "/usr/local/foo.bar", "/foo.tar.gz"}, wantExit: 0, wantOut: []string{"foo", "foo", "foo.tar"}},
}

var fishParityDirname = []cliTest{
	// path.fish L82: path dirname ./foo.mp4
	{name: "relative", args: []string{"dirname", "./foo.mp4"}, wantExit: 0, wantOut: []string{"."}},
	// path.fish L84: path dirname ../banana → ..
	{name: "parent_rel", args: []string{"dirname", "../banana"}, wantExit: 0, wantOut: []string{".."}},
	// path.fish L86: path dirname /usr/bin/ → /usr
	{name: "trailing_slash", args: []string{"dirname", "/usr/bin/"}, wantExit: 0, wantOut: []string{"/usr"}},
}

var fishParityNormalize = []cliTest{
	// path.fish L319: path normalize /usr/bin//../../etc/fish
	{name: "double_slash_dotdot", args: []string{"normalize", "/usr/bin//../../etc/fish"}, wantExit: 0, wantOut: []string{"/etc/fish"}},
	// path.fish L322: path normalize /bin//bash
	{name: "double_slash", args: []string{"normalize", "/bin//bash"}, wantExit: 0, wantOut: []string{"/bin/bash"}},
	// path.fish L327: path normalize -- -/foo -foo/foo
	{name: "dash_prefix", args: []string{"normalize", "--", "-/foo", "-foo/foo"}, wantExit: 0, wantOut: []string{"./-/foo", "./-foo/foo"}},
	// path.fish L330: path normalize -- ../-foo
	{name: "relative_dash_no_prefix", args: []string{"normalize", "--", "../-foo"}, wantExit: 0, wantOut: []string{"../-foo"}},
}

var fishParitySort = []cliTest{
	// path.fish L474: path sort --key=invalid-key
	{name: "invalid_key_error", args: []string{"sort", "--key=invalid-key"}, wantExit: 1, wantErr: "Invalid sort key"},
}

var fishParityDispatch = []cliTest{
	// path.fish L477: path
	{name: "missing_subcommand", args: []string{}, wantExit: 1, wantErr: "missing subcommand"},
	// path.fish L484: path invalid-subcmd
	{name: "invalid_subcommand", args: []string{"invalid-subcmd"}, wantExit: 1, wantErr: "invalid subcommand"},
	// path.fish L455: path filter -t invalid_type
	{name: "filter_invalid_type", args: []string{"filter", "-t", "invalid_type"}, wantExit: 1, wantErr: "Invalid type"},
	// path.fish L468: path change-extension (no args)
	{name: "change_extension_no_args", args: []string{"change-extension"}, wantExit: 1, wantErr: "missing argument"},
}

func TestFishParityPath(t *testing.T) {
	run := func(group string, tests []cliTest) {
		t.Run(group, func(t *testing.T) {
			for _, tc := range tests {
				t.Run(tc.name, func(t *testing.T) {
					runBin(t, tc)
				})
			}
		})
	}
	run("dispatch", fishParityDispatch)
	run("extension", fishParityExtension)
	run("change_extension", fishParityChangeExtension)
	run("basename", fishParityBasename)
	run("dirname", fishParityDirname)
	run("normalize", fishParityNormalize)
	run("sort", fishParitySort)
}

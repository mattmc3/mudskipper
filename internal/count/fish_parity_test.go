package countcmd

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// countBin is the path to the compiled binary, set by TestMain.
var countBin string

func TestMain(m *testing.M) {
	tmp, err := os.CreateTemp("", "count-test-*")
	if err != nil {
		panic(err)
	}
	tmp.Close()
	countBin = tmp.Name()
	defer os.Remove(countBin)

	if out, err := exec.Command("go", "build", "-o", countBin, "github.com/mattmc3/mudskipper/cmd/count").CombinedOutput(); err != nil {
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
	cmd := exec.Command(countBin, tc.args...)
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

// Fish parity tests from https://github.com/fish-shell/fish-shell/blob/master/tests/checks/count.fish
var fishParityCount = []cliTest{
	// count.fish L7: count
	{name: "no_args", args: []string{}, wantExit: 1, wantOut: []string{"0"}},
	// count.fish L11: count x
	{name: "one_arg", args: []string{"x"}, wantExit: 0, wantOut: []string{"1"}},
	// count.fish L15: count x y
	{name: "two_args", args: []string{"x", "y"}, wantExit: 0, wantOut: []string{"2"}},
	// count.fish L19: count -h  (no flag processing)
	{name: "flag_h_is_an_arg", args: []string{"-h"}, wantExit: 0, wantOut: []string{"1"}},
	// count.fish L20: count --help  (no flag processing)
	{name: "flag_help_is_an_arg", args: []string{"--help"}, wantExit: 0, wantOut: []string{"1"}},
	// count.fish L21: count --
	{name: "double_dash_is_an_arg", args: []string{"--"}, wantExit: 0, wantOut: []string{"1"}},
	// count.fish L22: count -- abc
	{name: "double_dash_plus_arg", args: []string{"--", "abc"}, wantExit: 0, wantOut: []string{"2"}},
	// count.fish L23: count def -- abc
	{name: "three_args_with_double_dash", args: []string{"def", "--", "abc"}, wantExit: 0, wantOut: []string{"3"}},
	// count.fish L44: printf '%s\n' 1 2 3 4 5 | count 6 7 8 9 10
	{name: "stdin_plus_args", stdin: "1\n2\n3\n4\n5\n", args: []string{"6", "7", "8", "9", "10"}, wantExit: 0, wantOut: []string{"10"}},
	// count.fish L48: echo -n 0 | count  (counts newlines, not bytes)
	{name: "stdin_no_newline", stdin: "0", args: []string{}, wantExit: 1, wantOut: []string{"0"}},
	// count.fish L51: echo 1 | count
	{name: "stdin_one_line", stdin: "1\n", args: []string{}, wantExit: 0, wantOut: []string{"1"}},
}

func TestFishParityCount(t *testing.T) {
	for _, tc := range fishParityCount {
		t.Run(tc.name, func(t *testing.T) {
			runBin(t, tc)
		})
	}
}

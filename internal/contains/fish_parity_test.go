package containscmd

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// containsBin is the path to the compiled binary, set by TestMain.
var containsBin string

func TestMain(m *testing.M) {
	tmp, err := os.CreateTemp("", "contains-test-*")
	if err != nil {
		panic(err)
	}
	tmp.Close()
	containsBin = tmp.Name()
	defer os.Remove(containsBin)

	if out, err := exec.Command("go", "build", "-o", containsBin, "github.com/mattmc3/mudskipper/cmd/contains").CombinedOutput(); err != nil {
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
	cmd := exec.Command(containsBin, tc.args...)
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

// Fish parity tests from https://github.com/fish-shell/fish-shell/blob/master/tests/checks/basic.fish
var fishParityContains = []cliTest{
	// basic.fish L228: contains -i string a b c string d
	{name: "index_found", args: []string{"-i", "string", "a", "b", "c", "string", "d"}, wantExit: 0, wantOut: []string{"4"}},
	// basic.fish L230: contains -i string a b c d
	{name: "index_not_found", args: []string{"-i", "string", "a", "b", "c", "d"}, wantExit: 1, wantOut: []string{}},
	// basic.fish L232: contains -i -- string a b c string d
	{name: "index_with_double_dash", args: []string{"-i", "--", "string", "a", "b", "c", "string", "d"}, wantExit: 0, wantOut: []string{"4"}},
	// basic.fish L234: contains -i -- -- a b c
	{name: "double_dash_value_not_found", args: []string{"-i", "--", "--", "a", "b", "c"}, wantExit: 1, wantOut: []string{}},
	// basic.fish L236: contains -i -- -- a b c -- v
	{name: "double_dash_value_found", args: []string{"-i", "--", "--", "a", "b", "c", "--", "v"}, wantExit: 0, wantOut: []string{"4"}},
	// basic.fish L238: contains
	{name: "no_args_error", args: []string{}, wantExit: 1, wantErr: "Key not specified"},
}

func TestFishParityContains(t *testing.T) {
	for _, tc := range fishParityContains {
		t.Run(tc.name, func(t *testing.T) {
			runBin(t, tc)
		})
	}
}

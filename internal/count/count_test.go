package countcmd

import (
	"strings"
	"testing"
)

func run(args []string, stdin string) (int, string, string) {
	var out, err strings.Builder
	code := Run(args, strings.NewReader(stdin), &out, &err)
	return code, out.String(), err.String()
}

func TestCount_args(t *testing.T) {
	code, out, _ := run([]string{"a", "b", "c"}, "")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if out != "3\n" {
		t.Errorf("stdout: want \"3\", got %q", out)
	}
}

func TestCount_stdin_lines(t *testing.T) {
	code, out, _ := run(nil, "a\nb\nc\n")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if out != "3\n" {
		t.Errorf("stdout: want \"3\", got %q", out)
	}
}

func TestCount_zero_args_empty_stdin_returns_1(t *testing.T) {
	code, out, _ := run(nil, "")
	if code != 1 {
		t.Errorf("exit: want 1, got %d", code)
	}
	if out != "0\n" {
		t.Errorf("stdout: want \"0\", got %q", out)
	}
}

func TestCount_single_arg(t *testing.T) {
	code, out, _ := run([]string{"x"}, "")
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if out != "1\n" {
		t.Errorf("stdout: want \"1\", got %q", out)
	}
}

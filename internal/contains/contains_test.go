package containscmd

import (
	"strings"
	"testing"
)

func run(args []string) (int, string, string) {
	var out, err strings.Builder
	code := Run(args, &out, &err)
	return code, out.String(), err.String()
}

func TestContains_help(t *testing.T) {
	for _, flag := range []string{"-h", "--help"} {
		code, out, _ := run([]string{flag})
		if code != 0 {
			t.Errorf("%s: exit want 0, got %d", flag, code)
		}
		if !strings.Contains(out, "Usage:") {
			t.Errorf("%s: stdout want Usage, got %q", flag, out)
		}
	}
}

func TestContains_found(t *testing.T) {
	code, _, _ := run([]string{"foo", "bar", "baz", "foo"})
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
}

func TestContains_not_found(t *testing.T) {
	code, _, _ := run([]string{"foo", "bar", "baz"})
	if code != 1 {
		t.Errorf("exit: want 1, got %d", code)
	}
}

func TestContains_index(t *testing.T) {
	code, out, _ := run([]string{"-i", "foo", "bar", "foo", "baz"})
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
	if out != "2\n" {
		t.Errorf("stdout: want \"2\", got %q", out)
	}
}

func TestContains_index_not_found(t *testing.T) {
	code, out, _ := run([]string{"-i", "x", "a", "b", "c"})
	if code != 1 {
		t.Errorf("exit: want 1, got %d", code)
	}
	if out != "" {
		t.Errorf("stdout: want empty, got %q", out)
	}
}

func TestContains_no_args_returns_1(t *testing.T) {
	code, _, _ := run(nil)
	if code != 1 {
		t.Errorf("exit: want 1, got %d", code)
	}
}

func TestContains_double_dash_separator(t *testing.T) {
	code, _, _ := run([]string{"--", "-foo", "a", "-foo", "b"})
	if code != 0 {
		t.Errorf("exit: want 0, got %d", code)
	}
}

package main

import (
	"strings"
	"testing"
)

func runCmd(args ...string) (exitCode int, stdout, stderr string) {
	return runWithStdin("", args...)
}

func runWithStdin(stdin string, args ...string) (exitCode int, stdout, stderr string) {
	var out, err strings.Builder
	exitCode = run(args, strings.NewReader(stdin), &out, &err)
	return exitCode, out.String(), err.String()
}

func lines(stdout string) []string {
	return splitLines(stdout)
}

func assertExit(t *testing.T, want, got int) {
	t.Helper()
	if got != want {
		t.Errorf("exit code: want %d, got %d", want, got)
	}
}

func assertLines(t *testing.T, want, got []string) {
	t.Helper()
	if len(want) != len(got) {
		t.Errorf("lines: want %v, got %v", want, got)
		return
	}
	for i := range want {
		if want[i] != got[i] {
			t.Errorf("line %d: want %q, got %q", i, want[i], got[i])
		}
	}
}

func assertEmpty(t *testing.T, s string) {
	t.Helper()
	if s != "" {
		t.Errorf("expected empty, got %q", s)
	}
}

func assertContains(t *testing.T, substr, s string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Errorf("expected %q to contain %q", s, substr)
	}
}

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

func splitLines(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.TrimRight(s, "\n")
	return strings.Split(s, "\n")
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

func assertEqual(t *testing.T, want, got int) {
	t.Helper()
	if got != want {
		t.Errorf("want %d, got %d", want, got)
	}
}

func assertInts(t *testing.T, want, got []int) {
	t.Helper()
	if len(want) != len(got) {
		t.Errorf("want %v, got %v", want, got)
		return
	}
	for i := range want {
		if want[i] != got[i] {
			t.Errorf("[%d]: want %d, got %d", i, want[i], got[i])
		}
	}
}

func assertStr(t *testing.T, want, got string) {
	t.Helper()
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

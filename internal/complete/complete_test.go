package completecmd

import (
	"strings"
	"testing"
)

func run(args ...string) (int, string, string) {
	var out, err strings.Builder
	code := Run(args, &out, &err)
	return code, out.String(), err.String()
}

func TestComplete_help(t *testing.T) {
	for _, flag := range []string{"-h", "--help"} {
		code, out, _ := run(flag)
		if code != 0 {
			t.Errorf("%s: exit want 0, got %d", flag, code)
		}
		if !strings.Contains(out, "Usage:") {
			t.Errorf("%s: want Usage in stdout, got %q", flag, out)
		}
	}
}

func TestComplete_init_zsh(t *testing.T) {
	code, out, _ := run("init", "zsh")
	if code != 0 {
		t.Fatalf("exit want 0, got %d", code)
	}
	if !strings.Contains(out, "function complete()") {
		t.Errorf("want complete function, got %q", out)
	}
	if !strings.Contains(out, "--shell=zsh") {
		t.Errorf("want --shell=zsh in wrapper, got %q", out)
	}
}

func TestComplete_init_bash(t *testing.T) {
	code, out, _ := run("init", "bash")
	if code != 0 {
		t.Fatalf("exit want 0, got %d", code)
	}
	if !strings.Contains(out, "function compx()") {
		t.Errorf("want compx function, got %q", out)
	}
	if !strings.Contains(out, "--shell=bash") {
		t.Errorf("want --shell=bash in wrapper, got %q", out)
	}
}

func TestComplete_init_missing_shell(t *testing.T) {
	code, _, stderr := run("init")
	if code == 0 {
		t.Errorf("want non-zero exit")
	}
	if !strings.Contains(stderr, "shell argument required") {
		t.Errorf("want error message, got %q", stderr)
	}
}

func TestComplete_zsh_long_flag(t *testing.T) {
	code, out, _ := run("--shell=zsh", "-c", "git", "-l", "pull", "-d", "Fetch and merge")
	if code != 0 {
		t.Fatalf("exit want 0, got %d", code)
	}
	if !strings.Contains(out, "--pull") {
		t.Errorf("want --pull in output, got %q", out)
	}
	if !strings.Contains(out, "Fetch and merge") {
		t.Errorf("want description in output, got %q", out)
	}
	if !strings.Contains(out, "compdef") {
		t.Errorf("want compdef in output, got %q", out)
	}
	if !strings.Contains(out, "_complete_git") {
		t.Errorf("want _complete_git in output, got %q", out)
	}
}

func TestComplete_zsh_short_and_long(t *testing.T) {
	code, out, _ := run("--shell=zsh", "-c", "git", "-s", "p", "-l", "pull", "-d", "Fetch")
	if code != 0 {
		t.Fatalf("exit want 0, got %d", code)
	}
	if !strings.Contains(out, "{-p,--pull}") {
		t.Errorf("want combined short/long, got %q", out)
	}
}

func TestComplete_zsh_arguments(t *testing.T) {
	code, out, _ := run("--shell=zsh", "-c", "git", "-l", "format", "-a", "json yaml toml", "-d", "Output format")
	if code != 0 {
		t.Fatalf("exit want 0, got %d", code)
	}
	if !strings.Contains(out, "json yaml toml") {
		t.Errorf("want word list in output, got %q", out)
	}
}

func TestComplete_zsh_positional(t *testing.T) {
	code, out, _ := run("--shell=zsh", "-c", "git", "-a", "pull push fetch")
	if code != 0 {
		t.Fatalf("exit want 0, got %d", code)
	}
	if !strings.Contains(out, "pull push fetch") {
		t.Errorf("want positional args in output, got %q", out)
	}
}

func TestComplete_bash_long_flag(t *testing.T) {
	code, out, _ := run("--shell=bash", "-c", "git", "-l", "pull", "-d", "Fetch and merge")
	if code != 0 {
		t.Fatalf("exit want 0, got %d", code)
	}
	if !strings.Contains(out, "--pull") {
		t.Errorf("want --pull in output, got %q", out)
	}
	if !strings.Contains(out, "builtin complete") {
		t.Errorf("want builtin complete in output, got %q", out)
	}
	if !strings.Contains(out, "_complete_git") {
		t.Errorf("want _complete_git in output, got %q", out)
	}
}

func TestComplete_missing_shell(t *testing.T) {
	code, _, stderr := run("-c", "git", "-l", "pull")
	if code == 0 {
		t.Errorf("want non-zero exit")
	}
	if !strings.Contains(stderr, "--shell required") {
		t.Errorf("want --shell error, got %q", stderr)
	}
}

func TestComplete_missing_command(t *testing.T) {
	code, _, stderr := run("--shell=zsh", "-l", "pull")
	if code == 0 {
		t.Errorf("want non-zero exit")
	}
	if !strings.Contains(stderr, "required") {
		t.Errorf("want required error, got %q", stderr)
	}
}

func TestComplete_zsh_old_option(t *testing.T) {
	code, out, _ := run("--shell=zsh", "-c", "java", "-o", "verbose", "-d", "Be verbose")
	if code != 0 {
		t.Fatalf("exit want 0, got %d", code)
	}
	if !strings.Contains(out, "-verbose") {
		t.Errorf("want -verbose in output, got %q", out)
	}
}

func TestComplete_zsh_condition(t *testing.T) {
	code, out, _ := run("--shell=zsh", "-c", "git", "-l", "push", "-d", "Push", "-n", "[[ $CURRENT == 2 ]]")
	if code != 0 {
		t.Fatalf("exit want 0, got %d", code)
	}
	if !strings.Contains(out, "_complete_git_cond_opts") {
		t.Errorf("want cond_opts in output, got %q", out)
	}
	if !strings.Contains(out, "CURRENT == 2") {
		t.Errorf("want condition in output, got %q", out)
	}
}

func TestComplete_zsh_erase(t *testing.T) {
	code, out, _ := run("--shell=zsh", "-c", "git", "-e")
	if code != 0 {
		t.Fatalf("exit want 0, got %d", code)
	}
	if !strings.Contains(out, "unset") {
		t.Errorf("want unset in output, got %q", out)
	}
	if !strings.Contains(out, "compdef -d git") {
		t.Errorf("want compdef -d in output, got %q", out)
	}
}

func TestComplete_zsh_wraps(t *testing.T) {
	code, out, _ := run("--shell=zsh", "-c", "mygit", "-w", "git")
	if code != 0 {
		t.Fatalf("exit want 0, got %d", code)
	}
	if !strings.Contains(out, "_complete_git") {
		t.Errorf("want _complete_git in output, got %q", out)
	}
	if !strings.Contains(out, "mygit") {
		t.Errorf("want mygit in output, got %q", out)
	}
}

func TestComplete_zsh_exclusive(t *testing.T) {
	code, out, _ := run("--shell=zsh", "-c", "cmd", "-l", "output", "-x")
	if code != 0 {
		t.Fatalf("exit want 0, got %d", code)
	}
	// -x means requireParam + noFiles → action should be "( )" not _files
	if !strings.Contains(out, "( )") {
		t.Errorf("want empty action for exclusive, got %q", out)
	}
}

func TestComplete_zsh_path(t *testing.T) {
	code, out, _ := run("--shell=zsh", "-p", "/usr/bin/python3", "-l", "version", "-d", "Show version")
	if code != 0 {
		t.Fatalf("exit want 0, got %d", code)
	}
	if !strings.Contains(out, "_complete_python3") {
		t.Errorf("want _complete_python3 in output, got %q", out)
	}
	if !strings.Contains(out, "/usr/bin/python3") {
		t.Errorf("want path in compdef, got %q", out)
	}
}

func TestComplete_bash_erase(t *testing.T) {
	code, out, _ := run("--shell=bash", "-c", "git", "-e")
	if code != 0 {
		t.Fatalf("exit want 0, got %d", code)
	}
	if !strings.Contains(out, "builtin complete -r git") {
		t.Errorf("want complete -r in output, got %q", out)
	}
}

func TestComplete_bash_wraps(t *testing.T) {
	code, out, _ := run("--shell=bash", "-c", "mygit", "-w", "git")
	if code != 0 {
		t.Fatalf("exit want 0, got %d", code)
	}
	if !strings.Contains(out, "_complete_git") {
		t.Errorf("want _complete_git in output, got %q", out)
	}
}

func TestComplete_do_complete_unsupported(t *testing.T) {
	code, _, stderr := run("--shell=zsh", "-C", "git ")
	if code == 0 {
		t.Errorf("want non-zero exit for -C")
	}
	if !strings.Contains(stderr, "Fish-specific") {
		t.Errorf("want Fish-specific error, got %q", stderr)
	}
}

func TestComplete_unsupported_shell(t *testing.T) {
	code, _, stderr := run("init", "fish")
	if code == 0 {
		t.Errorf("want non-zero exit")
	}
	if !strings.Contains(stderr, "unsupported shell") {
		t.Errorf("want unsupported shell error, got %q", stderr)
	}
}

package completecmd

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

const usage = `Usage: complete [--shell SHELL] [-c CMD] [-s SHORT] [-l LONG] [-o OLD] [-d DESC]
                [-a ARGS] [-n COND] [-w WRAPS] [-p PATH] [-e] [-f] [-F] [-r] [-x] [-k]
       complete init SHELL

Emit shell completion definitions in Fish-compatible syntax.

Subcommands:
  init SHELL    Emit shell wrapper function for SHELL (zsh, bash)

Options:
  --shell SHELL              Target shell: zsh bash
  -c, --command CMD          Command to complete for
  -p, --path PATH            Path to command (uses basename as command name)
  -s, --short-option S       POSIX short option (single char, no dash)
  -l, --long-option L        GNU long option (no dashes)
  -o, --old-option O         Old-style single-dash long option
  -d, --description DESC     Description shown in completion menu
  -a, --arguments ARGS       Space-separated completion candidates or flag values
  -n, --condition COND       Shell expression; completion added only if it exits 0
  -w, --wraps CMD            Inherit completions from CMD
  -e, --erase                Remove completions for the command
  -f, --no-files             Suppress file path completion
  -F, --force-files          Always offer file path completion
  -r, --require-parameter    Flag requires a value
  -x, --exclusive            Require parameter and suppress file completion (-r -f)
  -k, --keep-order           Keep candidate order instead of sorting alphabetically
  -h, --help                 Show this help message
`

type spec struct {
	command      string
	path         string
	shortOption  string
	longOption   string
	oldOption    string
	description  string
	arguments    string
	condition    string
	wraps        string
	noFiles      bool
	forceFiles   bool
	requireParam bool
	keepOrder    bool
	erase        bool
}

// effectiveCommand returns the command name derived from -c or -p.
func (sp spec) effectiveCommand() string {
	if sp.command != "" {
		return sp.command
	}
	if sp.path != "" {
		return filepath.Base(sp.path)
	}
	return ""
}

// compdefTarget returns the value to pass to compdef/complete — path if set, else command.
func (sp spec) compdefTarget() string {
	if sp.path != "" {
		return sp.path
	}
	return sp.command
}

func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Fprint(stdout, usage)
		return 0
	}

	if args[0] == "init" {
		if len(args) < 2 {
			fmt.Fprintln(stderr, "complete init: shell argument required")
			return 1
		}
		return emitInit(args[1], stdout, stderr)
	}

	shell, sp, ok := parseArgs(args, stderr)
	if !ok {
		return 1
	}
	if shell == "" {
		fmt.Fprintln(stderr, "complete: --shell required")
		return 1
	}

	switch strings.ToLower(shell) {
	case "zsh":
		return emitZsh(sp, stdout, stderr)
	case "bash":
		return emitBash(sp, stdout, stderr)
	default:
		fmt.Fprintf(stderr, "complete: unsupported shell %q (supported: zsh bash)\n", shell)
		return 1
	}
}

func parseArgs(args []string, stderr io.Writer) (shell string, sp spec, ok bool) {
	i := 0
	for i < len(args) {
		arg := args[i]
		get := func() (string, bool) {
			i++
			if i >= len(args) {
				fmt.Fprintf(stderr, "complete: %s requires a value\n", arg)
				return "", false
			}
			return args[i], true
		}
		switch {
		case strings.HasPrefix(arg, "--shell="):
			shell = arg[8:]
		case arg == "--shell":
			v, ok2 := get()
			if !ok2 {
				return shell, sp, false
			}
			shell = v
		case arg == "-c" || arg == "--command":
			v, ok2 := get()
			if !ok2 {
				return shell, sp, false
			}
			sp.command = v
		case strings.HasPrefix(arg, "--command="):
			sp.command = arg[10:]
		case arg == "-p" || arg == "--path":
			v, ok2 := get()
			if !ok2 {
				return shell, sp, false
			}
			sp.path = v
		case strings.HasPrefix(arg, "--path="):
			sp.path = arg[7:]
		case arg == "-s" || arg == "--short-option":
			v, ok2 := get()
			if !ok2 {
				return shell, sp, false
			}
			sp.shortOption = v
		case strings.HasPrefix(arg, "--short-option="):
			sp.shortOption = arg[15:]
		case arg == "-l" || arg == "--long-option":
			v, ok2 := get()
			if !ok2 {
				return shell, sp, false
			}
			sp.longOption = v
		case strings.HasPrefix(arg, "--long-option="):
			sp.longOption = arg[14:]
		case arg == "-o" || arg == "--old-option":
			v, ok2 := get()
			if !ok2 {
				return shell, sp, false
			}
			sp.oldOption = v
		case strings.HasPrefix(arg, "--old-option="):
			sp.oldOption = arg[13:]
		case arg == "-d" || arg == "--description":
			v, ok2 := get()
			if !ok2 {
				return shell, sp, false
			}
			sp.description = v
		case strings.HasPrefix(arg, "--description="):
			sp.description = arg[14:]
		case arg == "-a" || arg == "--arguments":
			v, ok2 := get()
			if !ok2 {
				return shell, sp, false
			}
			sp.arguments = v
		case strings.HasPrefix(arg, "--arguments="):
			sp.arguments = arg[12:]
		case arg == "-n" || arg == "--condition":
			v, ok2 := get()
			if !ok2 {
				return shell, sp, false
			}
			sp.condition = v
		case strings.HasPrefix(arg, "--condition="):
			sp.condition = arg[12:]
		case arg == "-w" || arg == "--wraps":
			v, ok2 := get()
			if !ok2 {
				return shell, sp, false
			}
			sp.wraps = v
		case strings.HasPrefix(arg, "--wraps="):
			sp.wraps = arg[8:]
		case arg == "-e" || arg == "--erase":
			sp.erase = true
		case arg == "-f" || arg == "--no-files":
			sp.noFiles = true
		case arg == "-F" || arg == "--force-files":
			sp.forceFiles = true
		case arg == "-r" || arg == "--require-parameter":
			sp.requireParam = true
		case arg == "-x" || arg == "--exclusive":
			sp.requireParam = true
			sp.noFiles = true
		case arg == "-k" || arg == "--keep-order":
			sp.keepOrder = true
		case arg == "-C" || arg == "--do-complete":
			fmt.Fprintln(stderr, "complete: -C/--do-complete is Fish-specific and not supported for other shells")
			return shell, sp, false
		case arg == "--color", arg == "--escape":
			// cosmetic/fish-only, no-op
			if !strings.Contains(arg, "=") {
				i++ // skip value if next arg isn't a flag
				if i < len(args) && !strings.HasPrefix(args[i], "-") {
					// consumed
				} else {
					i-- // wasn't a value
				}
			}
		default:
			fmt.Fprintf(stderr, "complete: unknown option %s\n", arg)
			return shell, sp, false
		}
		i++
	}
	return shell, sp, true
}

func cmdToIdent(cmd string) string {
	r := strings.NewReplacer("-", "_", " ", "_", ".", "_", "/", "_")
	return r.Replace(cmd)
}

package argparsecmd

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

const usage = `Usage: argparse [-h] [-n NAME] [--shell SHELL] [-s] [-i] [-S] [-x FLAGS] [-N MIN] [-X MAX] SPEC... -- [ARG...]

Parse command-line arguments and emit shell variable assignments.

Options:
  -h, --help               Show this help message
  -n, --name NAME          Command name used in error messages (default: argparse)
      --shell SHELL        Output shell syntax: fish bash zsh sh elvish nushell osh ysh (default: fish)
  -s, --stop-nonopt        Stop parsing flags at first non-option argument
      --no-local           Emit global assignments instead of local/set -l
  -i, --ignore-unknown     Pass unknown flags through to remaining args
  -S, --strict-longopts    Require exact long flag matches (no prefix abbreviation)
  -x, --exclusive FLAGS    Comma-separated group of mutually exclusive flags; may repeat
  -N, --min-args N         Require at least N remaining arguments
  -X, --max-args N         Allow at most N remaining arguments

Spec format:
  short/long               Boolean flag (-s / --long)
  short/long=              Required value flag
  short/long=?             Optional value flag
  short/long=+             Required value flag, may repeat
  short/long=*             Optional value flag, may repeat
`

type flagKind int

const (
	flagBool     flagKind = iota // boolean, no value
	flagStr                      // --name VALUE (required, last wins)
	flagOptional                 // --name [VALUE] (optional, last wins)
	flagMulti                    // --name VALUE (required, all saved)
	flagOptMulti                 // --name [VALUE] (optional, all saved) — fish =*
)

type flagSpec struct {
	short   string
	long    string
	varName string // derived from long name (or short if no long)
	kind    flagKind
}

// allVarNames returns all variable names this flag should set.
// For h/help: ["help", "h"]. For /help: ["help"]. For h: ["h"].
func (fs *flagSpec) allVarNames() []string {
	if fs.short != "" && fs.long != "" && fs.varName != fs.short {
		return []string{fs.varName, fs.short}
	}
	return []string{fs.varName}
}

type parseResult struct {
	flags     map[string][]string // varName → values (bool = ["-v"] or ["--verbose"])
	remaining []string
}

type parseOptions struct {
	name          string
	shell         string
	stopNonopt    bool
	noLocal       bool
	ignoreUnknown bool
	strictLong    bool
	exclusive     [][]string
	minArgs       int
	maxArgs       int // -1 = unlimited
}

func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	for _, a := range args {
		if a == "--" {
			break
		}
		if a == "-h" || a == "--help" {
			fmt.Fprint(stdout, usage)
			return 0
		}
	}

	opts := parseOptions{
		name:    "argparse",
		shell:   "fish",
		maxArgs: -1,
	}

	// Find -- separator first so we can scan the full pre-separator region.
	sepIdx := -1
	for j, a := range args {
		if a == "--" {
			sepIdx = j
			break
		}
	}
	if sepIdx < 0 {
		fmt.Fprintf(stderr, "%s: Missing -- separator\n", opts.name)
		return 1
	}

	// Walk args[0:sepIdx], consuming recognized options and collecting specs.
	var specArgs []string
	i := 0
	for i < sepIdx {
		arg := args[i]
		switch {
		case strings.HasPrefix(arg, "--shell="):
			opts.shell = arg[8:]
			i++
		case arg == "--shell":
			i++
			if i >= sepIdx {
				fmt.Fprintln(stderr, "argparse: --shell requires a value")
				return 1
			}
			opts.shell = args[i]
			i++
		case strings.HasPrefix(arg, "--name="):
			opts.name = arg[7:]
			i++
		case arg == "--name" || arg == "-n":
			i++
			if i >= sepIdx {
				fmt.Fprintln(stderr, "argparse: --name requires a value")
				return 1
			}
			opts.name = args[i]
			i++
		case arg == "--stop-nonopt" || arg == "-s":
			opts.stopNonopt = true
			i++
		case arg == "--no-local":
			opts.noLocal = true
			i++
		case arg == "--ignore-unknown" || arg == "-i":
			opts.ignoreUnknown = true
			i++
		case arg == "--strict-longopts" || arg == "-S":
			opts.strictLong = true
			i++
		case arg == "--exclusive" || arg == "-x":
			i++
			if i >= sepIdx {
				fmt.Fprintln(stderr, "argparse: --exclusive requires a value")
				return 1
			}
			opts.exclusive = append(opts.exclusive, strings.Split(args[i], ","))
			i++
		case strings.HasPrefix(arg, "--exclusive="):
			opts.exclusive = append(opts.exclusive, strings.Split(arg[12:], ","))
			i++
		case strings.HasPrefix(arg, "--min-args="):
			n, err := strconv.Atoi(arg[11:])
			if err != nil {
				fmt.Fprintf(stderr, "argparse: invalid --min-args value: %s\n", arg[11:])
				return 1
			}
			opts.minArgs = n
			i++
		case arg == "--min-args" || arg == "-N":
			i++
			if i >= sepIdx {
				fmt.Fprintln(stderr, "argparse: --min-args requires a value")
				return 1
			}
			n, err := strconv.Atoi(args[i])
			if err != nil {
				fmt.Fprintf(stderr, "argparse: invalid --min-args value: %s\n", args[i])
				return 1
			}
			opts.minArgs = n
			i++
		case strings.HasPrefix(arg, "--max-args="):
			n, err := strconv.Atoi(arg[11:])
			if err != nil {
				fmt.Fprintf(stderr, "argparse: invalid --max-args value: %s\n", arg[11:])
				return 1
			}
			opts.maxArgs = n
			i++
		case arg == "--max-args" || arg == "-X":
			i++
			if i >= sepIdx {
				fmt.Fprintln(stderr, "argparse: --max-args requires a value")
				return 1
			}
			n, err := strconv.Atoi(args[i])
			if err != nil {
				fmt.Fprintf(stderr, "argparse: invalid --max-args value: %s\n", args[i])
				return 1
			}
			opts.maxArgs = n
			i++
		case strings.HasPrefix(arg, "-"):
			fmt.Fprintf(stderr, "argparse: unknown option %s\n", arg)
			return 1
		default:
			specArgs = append(specArgs, arg)
			i++
		}
	}

	specs, err := parseSpecs(specArgs)
	if err != nil {
		fmt.Fprintf(stderr, "%s: %v\n", opts.name, err)
		return 1
	}

	for _, group := range opts.exclusive {
		for _, name := range group {
			if findSpec(specs, name, name) == nil {
				fmt.Fprintf(stderr, "%s: -x: unknown flag %q\n", opts.name, name)
				return 1
			}
		}
	}

	result, ok := parseArgs(args[sepIdx+1:], specs, opts, stderr)
	if !ok {
		return 1
	}

	// Validate arg counts
	if len(result.remaining) < opts.minArgs {
		fmt.Fprintf(stderr, "%s: expected at least %d arguments, got %d\n", opts.name, opts.minArgs, len(result.remaining))
		return 1
	}
	if opts.maxArgs >= 0 && len(result.remaining) > opts.maxArgs {
		fmt.Fprintf(stderr, "%s: expected at most %d arguments, got %d\n", opts.name, opts.maxArgs, len(result.remaining))
		return 1
	}

	for _, group := range opts.exclusive {
		var seen []string
		for _, name := range group {
			fs := findSpec(specs, name, name)
			if _, ok := result.flags[fs.varName]; ok {
				seen = append(seen, flagLabel(fs))
			}
		}
		if len(seen) > 1 {
			fmt.Fprintf(stderr, "%s: %s cannot be used together\n", opts.name, strings.Join(seen, " and "))
			return 1
		}
	}

	if err := emitShellCode(result, specs, opts.shell, opts.noLocal, stdout); err != nil {
		fmt.Fprintf(stderr, "argparse: %v\n", err)
		return 1
	}
	return 0
}

func parseSpecs(raw []string) ([]*flagSpec, error) {
	var specs []*flagSpec
	for _, s := range raw {
		fs, err := parseOneSpec(s)
		if err != nil {
			return nil, err
		}
		specs = append(specs, fs)
	}
	return specs, nil
}

func parseOneSpec(s string) (*flagSpec, error) {
	orig := s
	kind := flagBool

	switch {
	case strings.HasSuffix(s, "=+"):
		kind = flagMulti
		s = s[:len(s)-2]
	case strings.HasSuffix(s, "=*"):
		kind = flagOptMulti
		s = s[:len(s)-2]
	case strings.HasSuffix(s, "=?"):
		kind = flagOptional
		s = s[:len(s)-2]
	case strings.HasSuffix(s, "="):
		kind = flagStr
		s = s[:len(s)-1]
	}

	for _, c := range s {
		if c != '/' && !isAlphaNum(byte(c)) && c != '-' && c != '_' {
			return nil, fmt.Errorf("invalid option spec %q at char %q", orig, string(c))
		}
	}

	fs := &flagSpec{kind: kind}
	if idx := strings.Index(s, "/"); idx >= 0 {
		fs.short = s[:idx]
		fs.long = s[idx+1:]
		if fs.long == "" {
			return nil, fmt.Errorf("invalid option spec %q", orig)
		}
	} else if len(s) == 0 {
		return nil, fmt.Errorf("An option spec must have at least a short or a long flag")
	} else if len(s) == 1 {
		fs.short = s
	} else {
		fs.long = s
	}

	if fs.long != "" {
		fs.varName = strings.ReplaceAll(fs.long, "-", "_")
	} else {
		fs.varName = fs.short
	}

	return fs, nil
}

func findSpec(specs []*flagSpec, short, long string) *flagSpec {
	for _, s := range specs {
		if short != "" && s.short == short {
			return s
		}
		if long != "" && s.long == long {
			return s
		}
	}
	return nil
}

// findSpecLong finds by exact or unambiguous prefix (GNU-style, used unless --strict-longopts).
func findSpecLong(specs []*flagSpec, name string) *flagSpec {
	if fs := findSpec(specs, "", name); fs != nil {
		return fs
	}
	var match *flagSpec
	for _, s := range specs {
		if s.long != "" && strings.HasPrefix(s.long, name) {
			if match != nil {
				return nil // ambiguous
			}
			match = s
		}
	}
	return match
}

func parseArgs(args []string, specs []*flagSpec, opts parseOptions, stderr io.Writer) (*parseResult, bool) {
	res := &parseResult{flags: make(map[string][]string)}

	setFlag := func(fs *flagSpec, val string) {
		for _, name := range fs.allVarNames() {
			res.flags[name] = append(res.flags[name], val)
		}
	}

	i := 0
	for i < len(args) {
		arg := args[i]

		if arg == "--" {
			res.remaining = append(res.remaining, args[i+1:]...)
			break
		}

		if !strings.HasPrefix(arg, "-") {
			if opts.stopNonopt {
				res.remaining = append(res.remaining, args[i:]...)
				break
			}
			res.remaining = append(res.remaining, arg)
			i++
			continue
		}

		if strings.HasPrefix(arg, "--") {
			// Long option
			name := arg[2:]
			var val string
			hasVal := false
			if eq := strings.Index(name, "="); eq >= 0 {
				val = name[eq+1:]
				name = name[:eq]
				hasVal = true
			}
			seenFlag := "--" + name

			var fs *flagSpec
			if opts.strictLong {
				fs = findSpec(specs, "", name)
			} else {
				fs = findSpecLong(specs, name)
			}
			if fs == nil {
				if opts.ignoreUnknown {
					res.remaining = append(res.remaining, arg)
					i++
					continue
				}
				if hasVal {
					fmt.Fprintf(stderr, "%s: --%s=%s: unknown option\n", opts.name, name, val)
				} else {
					fmt.Fprintf(stderr, "%s: unknown option --%s\n", opts.name, name)
				}
				return nil, false
			}
			if hasVal && fs.kind == flagBool {
				fmt.Fprintf(stderr, "%s: %s: option does not take an argument\n", opts.name, seenFlag)
				return nil, false
			}
			switch fs.kind {
			case flagBool:
				setFlag(fs, seenFlag)
			case flagStr, flagMulti:
				if !hasVal {
					i++
					if i >= len(args) {
						fmt.Fprintf(stderr, "%s: %s: option requires an argument\n", opts.name, seenFlag)
						return nil, false
					}
					val = args[i]
				}
				setFlag(fs, val)
			case flagOptional, flagOptMulti:
				if hasVal {
					setFlag(fs, val)
				} else {
					setFlag(fs, "")
				}
			}
			if fs.kind == flagStr {
				for _, n := range fs.allVarNames() {
					v := res.flags[n]
					if len(v) > 1 {
						res.flags[n] = v[len(v)-1:]
					}
				}
			}
			i++
		} else {
			// Short option — may be grouped: -abc = -a -b -c
			j := 1
			for j < len(arg) {
				ch := string(arg[j])
				fs := findSpec(specs, ch, "")
				if fs == nil {
					if opts.ignoreUnknown {
						// Keep remaining chars as a new short-flag arg
						res.remaining = append(res.remaining, "-"+arg[j:])
						break
					}
					fmt.Fprintf(stderr, "%s: unknown option -%s\n", opts.name, ch)
					return nil, false
				}
				switch fs.kind {
				case flagBool:
					setFlag(fs, "-"+ch)
					j++
				case flagStr, flagMulti:
					if j+1 < len(arg) {
						// Attached value: -nVALUE
						setFlag(fs, arg[j+1:])
					} else {
						i++
						if i >= len(args) {
							fmt.Fprintf(stderr, "%s: -%s: option requires an argument\n", opts.name, ch)
							return nil, false
						}
						setFlag(fs, args[i])
					}
					if fs.kind == flagStr {
						for _, n := range fs.allVarNames() {
							v := res.flags[n]
							if len(v) > 1 {
								res.flags[n] = v[len(v)-1:]
							}
						}
					}
					j = len(arg) // consumed rest
				case flagOptional, flagOptMulti:
					if j+1 < len(arg) {
						setFlag(fs, arg[j+1:])
					} else {
						setFlag(fs, "")
					}
					j = len(arg)
				}
			}
			i++
		} // end else (short options)
	} // end for i

	return res, true
}

func emitShellCode(res *parseResult, specs []*flagSpec, shell string, noLocal bool, w io.Writer) error {
	switch strings.ToLower(shell) {
	case "fish":
		emitFish(res, noLocal, w)
	case "elvish", "elv":
		emitElvish(res, noLocal, w)
	case "sh":
		emitSh(res, specs, noLocal, w)
	case "nushell", "nu":
		emitNushell(res, specs, w)
	case "bash", "zsh", "osh": // osh is oils.pub OSH — bash-compatible
		emitBash(res, specs, noLocal, w)
	case "ysh": // oils.pub YSH — var declarations, no local keyword needed
		emitYsh(res, specs, noLocal, w)
	default:
		return fmt.Errorf("unsupported shell %q (supported: bash zsh sh fish elvish nushell osh ysh)", shell)
	}
	return nil
}

// emitSh emits POSIX sh output.
// - No arrays: multi-values stored RS-delimited (0x1E). Iterate: IFS=$(printf '\036'); for item in $_flag_x
// - Uses `local` by default (works in functions). Pass --no-local for global assignment at top level.
func emitSh(res *parseResult, specs []*flagSpec, noLocal bool, w io.Writer) {
	local := "local "
	if noLocal {
		local = ""
	}

	for varName, vals := range res.flags {
		kind := kindFor(varName, specs)
		if (kind == flagMulti || kind == flagOptMulti) && len(vals) > 1 {
			// RS-delimited: IFS=$(printf '\036'); for item in $_flag_x; do ...
			fmt.Fprintf(w, "%s_flag_%s=%s\n", local, varName, shellQuote(strings.Join(vals, "\x1e")))
		} else {
			fmt.Fprintf(w, "%s_flag_%s=%s\n", local, varName, shellQuote(vals[len(vals)-1]))
		}
	}
	args := make([]string, len(res.remaining))
	for i, r := range res.remaining {
		args[i] = shellQuote(r)
	}
	if len(args) > 0 {
		fmt.Fprintf(w, "set -- %s\n", strings.Join(args, " "))
	} else {
		fmt.Fprintln(w, "set --")
	}
}

// emitNushell emits Nushell output using let bindings and lists.
// Note: Nushell variables are always block-scoped; --no-local has no effect.
// Remaining args are stored in _args; Nushell has no `set --` equivalent.
func emitNushell(res *parseResult, specs []*flagSpec, w io.Writer) {
	for varName, vals := range res.flags {
		kind := kindFor(varName, specs)
		switch {
		case (kind == flagMulti || kind == flagOptMulti) && len(vals) > 1:
			quoted := make([]string, len(vals))
			for i, v := range vals {
				quoted[i] = nuQuote(v)
			}
			fmt.Fprintf(w, "let _flag_%s = [%s]\n", varName, strings.Join(quoted, " "))
		case kind == flagBool:
			fmt.Fprintf(w, "let _flag_%s = true\n", varName)
		default:
			fmt.Fprintf(w, "let _flag_%s = %s\n", varName, nuQuote(vals[len(vals)-1]))
		}
	}
	args := make([]string, len(res.remaining))
	for i, r := range res.remaining {
		args[i] = nuQuote(r)
	}
	fmt.Fprintf(w, "let _args = [%s]\n", strings.Join(args, " "))
}

// emitYsh emits YSH (oils.pub) output.
// var = local to proc (default); setglobal = module-level (--no-local).
// No set -- equivalent; remaining args stored in _args list.
func emitYsh(res *parseResult, specs []*flagSpec, noLocal bool, w io.Writer) {
	decl := "var"
	if noLocal {
		decl = "setglobal"
	}
	for varName, vals := range res.flags {
		kind := kindFor(varName, specs)
		switch {
		case (kind == flagMulti || kind == flagOptMulti) && len(vals) > 1:
			quoted := make([]string, len(vals))
			for i, v := range vals {
				quoted[i] = nuQuote(v)
			}
			fmt.Fprintf(w, "%s _flag_%s = [%s]\n", decl, varName, strings.Join(quoted, ", "))
		case kind == flagBool:
			fmt.Fprintf(w, "%s _flag_%s = true\n", decl, varName)
		default:
			fmt.Fprintf(w, "%s _flag_%s = %s\n", decl, varName, nuQuote(vals[len(vals)-1]))
		}
	}
	args := make([]string, len(res.remaining))
	for i, r := range res.remaining {
		args[i] = nuQuote(r)
	}
	fmt.Fprintf(w, "%s _args = [%s]\n", decl, strings.Join(args, ", "))
}

func nuQuote(s string) string {
	if s == "" {
		return `""`
	}
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return `"` + s + `"`
}

// emitBash emits bash/zsh output using local and arrays.
func emitBash(res *parseResult, specs []*flagSpec, noLocal bool, w io.Writer) {
	local := "local "
	if noLocal {
		local = ""
	}
	for varName, vals := range res.flags {
		kind := kindFor(varName, specs)
		if (kind == flagMulti || kind == flagOptMulti) && len(vals) > 1 {
			quoted := make([]string, len(vals))
			for i, v := range vals {
				quoted[i] = shellQuote(v)
			}
			fmt.Fprintf(w, "%s_flag_%s=( %s )\n", local, varName, strings.Join(quoted, " "))
		} else {
			fmt.Fprintf(w, "%s_flag_%s=%s\n", local, varName, shellQuote(vals[len(vals)-1]))
		}
	}
	args := make([]string, len(res.remaining))
	for i, r := range res.remaining {
		args[i] = shellQuote(r)
	}
	if len(args) > 0 {
		fmt.Fprintf(w, "set -- %s\n", strings.Join(args, " "))
	} else {
		fmt.Fprintln(w, "set --")
	}
}

func emitFish(res *parseResult, noLocal bool, w io.Writer) {
	setCmd := "set -l"
	if noLocal {
		setCmd = "set"
	}
	for varName, vals := range res.flags {
		parts := make([]string, len(vals))
		for i, v := range vals {
			parts[i] = fishQuote(v)
		}
		fmt.Fprintf(w, "%s _flag_%s %s\n", setCmd, varName, strings.Join(parts, " "))
	}
	args := make([]string, len(res.remaining))
	for i, r := range res.remaining {
		args[i] = fishQuote(r)
	}
	if len(args) > 0 {
		fmt.Fprintf(w, "%s argv %s\n", setCmd, strings.Join(args, " "))
	} else {
		fmt.Fprintf(w, "%s argv\n", setCmd)
	}
}

func emitElvish(res *parseResult, noLocal bool, w io.Writer) {
	decl := "var"
	if noLocal {
		decl = "set"
	}
	for varName, vals := range res.flags {
		if len(vals) > 1 {
			quoted := make([]string, len(vals))
			for i, v := range vals {
				quoted[i] = elvishQuote(v)
			}
			fmt.Fprintf(w, "%s _flag_%s = [%s]\n", decl, varName, strings.Join(quoted, " "))
		} else if vals[0] == "-"+varName || vals[0] == "--"+varName {
			fmt.Fprintf(w, "%s _flag_%s = $true\n", decl, varName)
		} else {
			fmt.Fprintf(w, "%s _flag_%s = %s\n", decl, varName, elvishQuote(vals[0]))
		}
	}
	args := make([]string, len(res.remaining))
	for i, r := range res.remaining {
		args[i] = elvishQuote(r)
	}
	fmt.Fprintf(w, "%s args = [%s]\n", decl, strings.Join(args, " "))
}

func kindFor(varName string, specs []*flagSpec) flagKind {
	for _, s := range specs {
		if s.varName == varName || s.short == varName {
			return s.kind
		}
	}
	return flagBool
}

func shellQuote(s string) string {
	if s == "" {
		return "''"
	}
	safe := true
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !isAlphaNum(c) && c != '-' && c != '_' && c != '.' && c != '/' && c != ':' && c != '@' {
			safe = false
			break
		}
	}
	if safe {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func fishQuote(s string) string {
	if s == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(s, "'", "\\'") + "'"
}

func elvishQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

func isAlphaNum(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func flagLabel(fs *flagSpec) string {
	if fs.short != "" && fs.long != "" {
		return "-" + fs.short + "/--" + fs.long
	}
	if fs.short != "" {
		return "-" + fs.short
	}
	return "--" + fs.long
}




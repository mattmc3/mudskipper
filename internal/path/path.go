package pathcmd

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	goSort "sort"
	"strconv"
	"strings"
	"time"

	"rsc.io/getopt"
)

const usage = `Usage: path <subcommand> [options] [PATH ...]

Subcommands:
  basename        Strip directory and optionally suffix from a path
  change-extension  Replace the file extension
  dirname         Return the directory component of a path
  extension       Return the file extension
  filter          Filter paths by type
  is              Test path properties
  mtime           Print last-modified time of a path
  normalize       Resolve symlinks and canonicalize a path
  resolve         Resolve a relative path to absolute
  sort            Sort a list of paths

Use 'path <subcommand> --help' for more information about a specific subcommand.
`

const usageBasename = `Usage: path basename [-h] [-E] [-q] [-Z] [PATH ...]

  Print the non-directory part of each PATH.

Options:
  -E, --strip-extension   Strip file extension from result
  -q, --quiet             Suppress output; exit 0 if any path, 1 if none
  -Z, --null-out          Separate output with NUL instead of newline
  -h, --help              Show this help message
`

const usageDirname = `Usage: path dirname [-h] [-q] [-Z] [PATH ...]

  Print the directory part of each PATH.

Options:
  -q, --quiet      Suppress output; exit 0 if any path, 1 if none
  -Z, --null-out   Separate output with NUL instead of newline
  -h, --help       Show this help message
`

const usageExtension = `Usage: path extension [-h] [-q] [-Z] [PATH ...]

  Print the extension of each PATH (including leading dot).
  Exits 0 if at least one path has an extension, 1 otherwise.

Options:
  -q, --quiet      Suppress output
  -Z, --null-out   Separate output with NUL instead of newline
  -h, --help       Show this help message
`

const usageChangeExtension = `Usage: path change-extension [-h] [-Z] EXT [PATH ...]

  Replace or remove the extension of each PATH.

Options:
  -Z, --null-out   Separate output with NUL instead of newline
  -h, --help       Show this help message
`

const usageNormalize = `Usage: path normalize [-h] [-q] [-Z] [PATH ...]

  Normalize each PATH (resolve . and .. components, collapse slashes).
  Paths starting with '-' are prefixed with './' to prevent flag confusion.

Options:
  -q, --quiet      Suppress output
  -Z, --null-out   Separate output with NUL instead of newline
  -h, --help       Show this help message
`

const usageResolve = `Usage: path resolve [-h] [-q] [-Z] [PATH ...]

  Resolve each PATH to an absolute path, resolving symlinks.
  Non-existent paths are made absolute relative to CWD.

Options:
  -q, --quiet      Suppress output
  -Z, --null-out   Separate output with NUL instead of newline
  -h, --help       Show this help message
`

const usageSort = `Usage: path sort [-h] [-r] [-u] [--key=KEY] [-Z] [PATH ...]

  Sort paths.

Options:
  -r, --reverse    Reverse sort order
  -u, --unique     Remove duplicates (by sort key)
  --key=KEY        Sort key: path (default), basename, dirname
  -Z, --null-out   Separate output with NUL instead of newline
  -h, --help       Show this help message
`

const usageFilter = `Usage: path filter [-h] [-v] [-q] [-f|-d|-l] [-t TYPE] [-r|-w|-x] [-p PERM] [-Z] [PATH ...]

  Filter paths by existence and optional criteria.

Options:
  -v, --invert      Print paths that do NOT match
  -f                Match regular files
  -d                Match directories
  -l                Match symlinks
  -t, --type TYPE   Match type: file,dir,link,block,char,fifo,socket
  -r                Match readable paths
  -w                Match writable paths
  -x                Match executable paths
  -p, --perm PERM   Match permission: read,write,exec,suid,sgid
  -q, --quiet       Suppress output; exit 0 if any match
  -Z, --null-out    Separate output with NUL instead of newline
  -h, --help        Show this help message
`

const usageIs = `Usage: path is [-h] [-v] [-f|-d|-l] [-t TYPE] [-r|-w|-x] [-p PERM] [PATH ...]

  Test if paths match criteria. Sets exit status only; no output.

Options:
  -v, --invert    Test that paths do NOT match
  -f              Test for regular files
  -d              Test for directories
  -l              Test for symlinks
  -t, --type      Match type: file,dir,link,...
  -r/-w/-x        Test readable/writable/executable
  -p, --perm      Match permission: read,write,exec,suid,sgid
  -h, --help      Show this help message
`

const usageMtime = `Usage: path mtime [-h] [--relative] [-R] [PATH ...]

  Print modification time of each PATH as Unix timestamp.

Options:
  --relative   Print seconds since modification (age)
  -R, --reverse  Sort oldest first
  -h, --help     Show this help message
`

func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "path: missing subcommand")
		return 1
	}
	cmd, rest := args[0], args[1:]
	switch cmd {
	case "-h", "--help":
		fmt.Fprint(stdout, usage)
		return 0
	case "basename":
		return runBasename(rest, stdin, stdout, stderr)
	case "dirname":
		return runDirname(rest, stdin, stdout, stderr)
	case "extension":
		return runExtension(rest, stdin, stdout, stderr)
	case "change-extension":
		return runChangeExtension(rest, stdin, stdout, stderr)
	case "filter":
		return runFilter(rest, stdin, stdout, stderr)
	case "is":
		return runIs(rest, stdin, stdout, stderr)
	case "normalize":
		return runNormalize(rest, stdin, stdout, stderr)
	case "resolve":
		return runResolve(rest, stdin, stdout, stderr)
	case "sort":
		return runSortCmd(rest, stdin, stdout, stderr)
	case "mtime":
		return runMtime(rest, stdin, stdout, stderr)
	default:
		fmt.Fprintf(stderr, "path %s: invalid subcommand\n", cmd)
		return 1
	}
}

func inputPaths(args []string, stdin io.Reader) []string {
	if len(args) > 0 {
		return args
	}
	var paths []string
	scanner := bufio.NewScanner(stdin)
	for scanner.Scan() {
		paths = append(paths, scanner.Text())
	}
	return paths
}

func writePaths(w io.Writer, paths []string, nullOut bool) {
	for _, p := range paths {
		if nullOut {
			fmt.Fprint(w, p+"\x00")
		} else {
			fmt.Fprintln(w, p)
		}
	}
}

// extension returns the file extension of the basename (including leading dot),
// or "" if none. A leading dot in the basename (hidden file) is not an extension.
func extension(p string) string {
	base := filepath.Base(p)
	if base == "." || base == ".." {
		return ""
	}
	dot := strings.LastIndex(base, ".")
	if dot <= 0 {
		return ""
	}
	return base[dot:]
}

func runBasename(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := getopt.NewFlagSet("basename", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	stripExt := fs.Bool("strip-extension", false, "")
	quiet := fs.Bool("quiet", false, "")
	nullOut := fs.Bool("null-out", false, "")
	fs.Aliases("h", "help", "E", "strip-extension", "q", "quiet", "Z", "null-out")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: path basename: %v\n", err)
		return 1
	}
	if *help {
		fmt.Fprint(stdout, usageBasename)
		return 0
	}
	paths := inputPaths(fs.Args(), stdin)
	if len(paths) == 0 {
		return 1
	}
	for _, p := range paths {
		b := filepath.Base(p)
		if *stripExt {
			ext := extension(b)
			if ext != "" {
				b = b[:len(b)-len(ext)]
			}
		}
		if !*quiet {
			if *nullOut {
				fmt.Fprint(stdout, b+"\x00")
			} else {
				fmt.Fprintln(stdout, b)
			}
		}
	}
	return 0
}

func runDirname(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := getopt.NewFlagSet("dirname", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	quiet := fs.Bool("quiet", false, "")
	nullOut := fs.Bool("null-out", false, "")
	fs.Aliases("h", "help", "q", "quiet", "Z", "null-out")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: path dirname: %v\n", err)
		return 1
	}
	if *help {
		fmt.Fprint(stdout, usageDirname)
		return 0
	}
	paths := inputPaths(fs.Args(), stdin)
	if len(paths) == 0 {
		return 1
	}
	for _, p := range paths {
		// Clean before Dir so trailing slashes don't confuse: dirname /usr/bin/ → /usr
		dir := filepath.Dir(filepath.Clean(p))
		if !*quiet {
			if *nullOut {
				fmt.Fprint(stdout, dir+"\x00")
			} else {
				fmt.Fprintln(stdout, dir)
			}
		}
	}
	return 0
}

func runExtension(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := getopt.NewFlagSet("extension", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	quiet := fs.Bool("quiet", false, "")
	nullOut := fs.Bool("null-out", false, "")
	fs.Aliases("h", "help", "q", "quiet", "Z", "null-out")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: path extension: %v\n", err)
		return 1
	}
	if *help {
		fmt.Fprint(stdout, usageExtension)
		return 0
	}
	paths := inputPaths(fs.Args(), stdin)
	any := false
	for _, p := range paths {
		ext := extension(p)
		if !*quiet {
			if *nullOut {
				fmt.Fprint(stdout, ext+"\x00")
			} else {
				fmt.Fprintln(stdout, ext)
			}
		}
		if ext != "" {
			any = true
		}
	}
	if any {
		return 0
	}
	return 1
}

func runChangeExtension(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := getopt.NewFlagSet("change-extension", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	nullOut := fs.Bool("null-out", false, "")
	fs.Aliases("h", "help", "Z", "null-out")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: path change-extension: %v\n", err)
		return 1
	}
	if *help {
		fmt.Fprint(stdout, usageChangeExtension)
		return 0
	}
	rest := fs.Args()
	if len(rest) == 0 {
		fmt.Fprintln(stderr, "path change-extension: missing argument")
		return 1
	}
	newExt := rest[0]
	// Ensure extension starts with a dot (unless empty)
	if newExt != "" && !strings.HasPrefix(newExt, ".") {
		newExt = "." + newExt
	}
	paths := inputPaths(rest[1:], stdin)
	if len(paths) == 0 {
		return 1
	}
	for _, p := range paths {
		// Preserve the original dir prefix (filepath.Dir normalizes "./foo" → ".")
		slash := strings.LastIndex(p, "/")
		var prefix, base string
		if slash < 0 {
			prefix, base = "", p
		} else {
			prefix, base = p[:slash+1], p[slash+1:]
		}
		ext := extension(base)
		stem := base
		if ext != "" {
			stem = base[:len(base)-len(ext)]
		}
		result := prefix + stem + newExt
		if *nullOut {
			fmt.Fprint(stdout, result+"\x00")
		} else {
			fmt.Fprintln(stdout, result)
		}
	}
	return 0
}

func runNormalize(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := getopt.NewFlagSet("normalize", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	quiet := fs.Bool("quiet", false, "")
	nullOut := fs.Bool("null-out", false, "")
	fs.Aliases("h", "help", "q", "quiet", "Z", "null-out")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: path normalize: %v\n", err)
		return 1
	}
	if *help {
		fmt.Fprint(stdout, usageNormalize)
		return 0
	}
	paths := inputPaths(fs.Args(), stdin)
	if len(paths) == 0 {
		return 1
	}
	for _, p := range paths {
		clean := filepath.Clean(p)
		if strings.HasPrefix(clean, "-") {
			clean = "./" + clean
		}
		if !*quiet {
			if *nullOut {
				fmt.Fprint(stdout, clean+"\x00")
			} else {
				fmt.Fprintln(stdout, clean)
			}
		}
	}
	return 0
}

func runResolve(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := getopt.NewFlagSet("resolve", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	quiet := fs.Bool("quiet", false, "")
	nullOut := fs.Bool("null-out", false, "")
	fs.Aliases("h", "help", "q", "quiet", "Z", "null-out")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: path resolve: %v\n", err)
		return 1
	}
	if *help {
		fmt.Fprint(stdout, usageResolve)
		return 0
	}
	paths := inputPaths(fs.Args(), stdin)
	if len(paths) == 0 {
		return 1
	}
	for _, p := range paths {
		resolved, err := filepath.EvalSymlinks(p)
		if err != nil {
			// Non-existent: make absolute relative to CWD
			abs, err2 := filepath.Abs(p)
			if err2 != nil {
				abs = p
			}
			resolved = filepath.Clean(abs)
		}
		if !*quiet {
			if *nullOut {
				fmt.Fprint(stdout, resolved+"\x00")
			} else {
				fmt.Fprintln(stdout, resolved)
			}
		}
	}
	return 0
}

var validSortKeys = map[string]bool{"basename": true, "dirname": true, "path": true, "version": true}

func runSortCmd(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := getopt.NewFlagSet("sort", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	reverse := fs.Bool("reverse", false, "")
	unique := fs.Bool("unique", false, "")
	key := fs.String("key", "path", "")
	nullOut := fs.Bool("null-out", false, "")
	fs.Aliases("h", "help", "r", "reverse", "u", "unique", "Z", "null-out")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: path sort: %v\n", err)
		return 1
	}
	if *help {
		fmt.Fprint(stdout, usageSort)
		return 0
	}
	if !validSortKeys[*key] {
		fmt.Fprintf(stderr, "path sort: Invalid sort key '%s'\n", *key)
		return 1
	}
	paths := inputPaths(fs.Args(), stdin)
	if len(paths) == 0 {
		return 0
	}

	keyFunc := func(p string) string {
		switch *key {
		case "basename":
			return filepath.Base(p)
		case "dirname":
			return filepath.Dir(p)
		default:
			return p
		}
	}

	goSort.SliceStable(paths, func(i, j int) bool {
		ki, kj := keyFunc(paths[i]), keyFunc(paths[j])
		if *reverse {
			return ki > kj
		}
		return ki < kj
	})

	if *unique {
		seen := make(map[string]bool)
		var deduped []string
		for _, p := range paths {
			k := keyFunc(p)
			if !seen[k] {
				seen[k] = true
				deduped = append(deduped, p)
			}
		}
		paths = deduped
	}

	writePaths(stdout, paths, *nullOut)
	return 0
}

var validTypes = map[string]bool{
	"file": true, "f": true,
	"dir": true, "d": true,
	"link": true, "l": true,
	"block": true, "char": true, "fifo": true, "socket": true,
}

var validPerms = map[string]bool{
	"read": true, "r": true,
	"write": true, "w": true,
	"exec": true, "x": true,
	"suid": true, "sgid": true,
}

func matchesType(fi os.FileInfo, lfi os.FileInfo, types []string) bool {
	if len(types) == 0 {
		return true
	}
	for _, t := range types {
		switch t {
		case "file", "f":
			if fi.Mode().IsRegular() {
				return true
			}
		case "dir", "d":
			if fi.Mode().IsDir() {
				return true
			}
		case "link", "l":
			if lfi != nil && lfi.Mode()&os.ModeSymlink != 0 {
				return true
			}
		case "block":
			if fi.Mode()&os.ModeDevice != 0 && fi.Mode()&os.ModeCharDevice == 0 {
				return true
			}
		case "char":
			if fi.Mode()&os.ModeCharDevice != 0 {
				return true
			}
		case "fifo":
			if fi.Mode()&os.ModeNamedPipe != 0 {
				return true
			}
		case "socket":
			if fi.Mode()&os.ModeSocket != 0 {
				return true
			}
		}
	}
	return false
}

func matchesPerm(fi os.FileInfo, perms []string) bool {
	if len(perms) == 0 {
		return true
	}
	mode := fi.Mode()
	for _, p := range perms {
		switch p {
		case "read", "r":
			if mode&0444 == 0 {
				return false
			}
		case "write", "w":
			if mode&0222 == 0 {
				return false
			}
		case "exec", "x":
			if mode&0111 == 0 {
				return false
			}
		case "suid":
			if mode&os.ModeSetuid == 0 {
				return false
			}
		case "sgid":
			if mode&os.ModeSetgid == 0 {
				return false
			}
		}
	}
	return true
}

func parseCommaList(s string, valid map[string]bool, errPrefix string, stderr io.Writer) ([]string, bool) {
	if s == "" {
		return nil, true
	}
	parts := strings.Split(s, ",")
	for _, p := range parts {
		if !valid[strings.TrimSpace(p)] {
			fmt.Fprintf(stderr, "%s '%s'\n", errPrefix, p)
			return nil, false
		}
	}
	return parts, true
}

func runFilter(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := getopt.NewFlagSet("filter", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	invert := fs.Bool("invert", false, "")
	quiet := fs.Bool("quiet", false, "")
	nullOut := fs.Bool("null-out", false, "")
	typeStr := fs.String("type", "", "")
	permStr := fs.String("perm", "", "")
	isFile := fs.Bool("f", false, "")
	isDir := fs.Bool("d", false, "")
	isLink := fs.Bool("l", false, "")
	isRead := fs.Bool("r", false, "")
	isWrite := fs.Bool("w", false, "")
	isExec := fs.Bool("x", false, "")
	fs.Aliases("h", "help", "v", "invert", "q", "quiet", "Z", "null-out",
		"t", "type", "p", "perm")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: path filter: %v\n", err)
		return 1
	}
	if *help {
		fmt.Fprint(stdout, usageFilter)
		return 0
	}

	// Build type filter from shorthand flags
	typeSpec := *typeStr
	if *isFile {
		typeSpec = appendCSV(typeSpec, "file")
	}
	if *isDir {
		typeSpec = appendCSV(typeSpec, "dir")
	}
	if *isLink {
		typeSpec = appendCSV(typeSpec, "link")
	}
	types, ok := parseCommaList(typeSpec, validTypes, "path filter: Invalid type", stderr)
	if !ok {
		return 1
	}

	// Build perm filter from shorthand flags
	permSpec := *permStr
	if *isRead {
		permSpec = appendCSV(permSpec, "read")
	}
	if *isWrite {
		permSpec = appendCSV(permSpec, "write")
	}
	if *isExec {
		permSpec = appendCSV(permSpec, "exec")
	}
	perms, ok := parseCommaList(permSpec, validPerms, "path filter: Invalid permission", stderr)
	if !ok {
		return 1
	}

	paths := inputPaths(fs.Args(), stdin)
	matched := false

	for _, p := range paths {
		// Normalize paths starting with - for output
		displayP := p
		clean := filepath.Clean(p)
		if strings.HasPrefix(clean, "-") {
			displayP = "./" + clean
		}

		fi, err := os.Stat(p)
		lfi, _ := os.Lstat(p)
		exists := err == nil

		passes := exists && matchesType(fi, lfi, types) && matchesPerm(fi, perms)
		if *invert {
			passes = !passes
		}

		if passes {
			matched = true
			if !*quiet {
				if *nullOut {
					fmt.Fprint(stdout, displayP+"\x00")
				} else {
					fmt.Fprintln(stdout, displayP)
				}
			}
		}
	}

	if matched {
		return 0
	}
	return 1
}

func appendCSV(s, item string) string {
	if s == "" {
		return item
	}
	return s + "," + item
}

func runIs(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	// is = filter with output suppressed; exit 0 if any match
	fs := getopt.NewFlagSet("is", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	invert := fs.Bool("invert", false, "")
	typeStr := fs.String("type", "", "")
	permStr := fs.String("perm", "", "")
	isFile := fs.Bool("f", false, "")
	isDir := fs.Bool("d", false, "")
	isLink := fs.Bool("l", false, "")
	isRead := fs.Bool("r", false, "")
	isWrite := fs.Bool("w", false, "")
	isExec := fs.Bool("x", false, "")
	fs.Aliases("h", "help", "v", "invert", "t", "type", "p", "perm")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: path is: %v\n", err)
		return 1
	}
	if *help {
		fmt.Fprint(stdout, usageIs)
		return 0
	}

	typeSpec := *typeStr
	if *isFile {
		typeSpec = appendCSV(typeSpec, "file")
	}
	if *isDir {
		typeSpec = appendCSV(typeSpec, "dir")
	}
	if *isLink {
		typeSpec = appendCSV(typeSpec, "link")
	}
	types, ok := parseCommaList(typeSpec, validTypes, "path is: Invalid type", stderr)
	if !ok {
		return 1
	}

	permSpec := *permStr
	if *isRead {
		permSpec = appendCSV(permSpec, "read")
	}
	if *isWrite {
		permSpec = appendCSV(permSpec, "write")
	}
	if *isExec {
		permSpec = appendCSV(permSpec, "exec")
	}
	perms, ok := parseCommaList(permSpec, validPerms, "path is: Invalid permission", stderr)
	if !ok {
		return 1
	}

	paths := inputPaths(fs.Args(), stdin)
	for _, p := range paths {
		fi, err := os.Stat(p)
		lfi, _ := os.Lstat(p)
		exists := err == nil
		passes := exists && matchesType(fi, lfi, types) && matchesPerm(fi, perms)
		if *invert {
			passes = !passes
		}
		if passes {
			return 0
		}
	}
	return 1
}

func runMtime(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := getopt.NewFlagSet("mtime", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	help := fs.Bool("help", false, "")
	relative := fs.Bool("relative", false, "")
	reverse := fs.Bool("reverse", false, "")
	fs.Aliases("h", "help", "R", "reverse")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: path mtime: %v\n", err)
		return 1
	}
	if *help {
		fmt.Fprint(stdout, usageMtime)
		return 0
	}
	paths := inputPaths(fs.Args(), stdin)
	if len(paths) == 0 {
		return 1
	}

	type entry struct {
		path  string
		mtime time.Time
	}
	var entries []entry
	for _, p := range paths {
		fi, err := os.Stat(p)
		if err != nil {
			fmt.Fprintf(stderr, "path mtime: %v\n", err)
			return 1
		}
		entries = append(entries, entry{p, fi.ModTime()})
	}

	if len(entries) > 1 {
		goSort.SliceStable(entries, func(i, j int) bool {
			if *reverse {
				return entries[i].mtime.Before(entries[j].mtime)
			}
			return entries[j].mtime.Before(entries[i].mtime)
		})
	}

	now := time.Now()
	for _, e := range entries {
		if *relative {
			fmt.Fprintln(stdout, strconv.FormatInt(int64(now.Sub(e.mtime).Seconds()), 10))
		} else {
			fmt.Fprintln(stdout, strconv.FormatInt(e.mtime.Unix(), 10))
		}
	}
	return 0
}

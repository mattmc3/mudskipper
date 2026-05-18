# fish-utils

Fish shell utilities ported to Go for use in bash, zsh, sh, and other shells.

## Commands

### string

Manipulate strings. A port of Fish's [`string` builtin](https://fishshell.com/docs/current/cmds/string.html).

```
string collect [-a | --allow-empty] [-N | --no-trim-newlines] [STRING ...]
string escape [-n | --no-quoted] [--style=STYLE] [STRING ...]
string join [-q | --quiet] [-n | --no-empty] [--] SEP [STRING ...]
string join0 [-q | --quiet] [-n | --no-empty] [--] [STRING ...]
string length [-q | --quiet] [-V | --visible] [STRING ...]
string lower [-q | --quiet] [STRING ...]
string match [-a | --all] [-e | --entire] [-i | --ignore-case]
             [-g | --groups-only] [-r | --regex] [-n | --index]
             [-q | --quiet] [-v | --invert] [(-m | --max-matches) MAX]
             PATTERN [STRING ...]
string pad [-r | --right] [-C | --center] [(-c | --char) CHAR]
           [(-w | --width) INTEGER] [STRING ...]
string repeat [(-n | --count) COUNT] [(-m | --max) MAX]
              [-N | --no-newline] [-q | --quiet] [STRING ...]
string replace [-a | --all] [-f | --filter] [-i | --ignore-case]
               [-r | --regex] [(-m | --max-matches) MAX] [-q | --quiet]
               PATTERN REPLACEMENT [STRING ...]
string shorten [(-c | --char) CHARS] [(-m | --max) INTEGER]
               [-N | --no-newline] [-l | --left] [-q | --quiet] [STRING ...]
string split [(-f | --fields) FIELDS [-a | --allow-empty]] [(-m | --max) MAX]
             [-n | --no-empty] [-q | --quiet] [-r | --right] SEP [STRING ...]
string split0 [(-f | --fields) FIELDS [-a | --allow-empty]] [(-m | --max) MAX]
              [-n | --no-empty] [-q | --quiet] [-r | --right] [STRING ...]
string sub [(-s | --start) START] [(-e | --end) END] [(-l | --length) LENGTH]
           [-q | --quiet] [STRING ...]
string trim [-l | --left] [-r | --right] [(-c | --chars) CHARS]
            [-q | --quiet] [STRING ...]
string unescape [--style=STYLE] [STRING ...]
string upper [-q | --quiet] [STRING ...]
```

### count

Count arguments or stdin lines. A port of Fish's [`count` builtin](https://fishshell.com/docs/current/cmds/count.html).

```
count [ARGUMENT ...]
```

Counts arguments. With no arguments, counts newlines from stdin (like `wc -l`). When both arguments and stdin are available, counts both. Exits 0 if count > 0, 1 if count == 0. No flag processing — `-h`, `--`, etc. are counted as arguments.

### contains

Check if a value is in a list. A port of Fish's [`contains` builtin](https://fishshell.com/docs/current/cmds/contains.html).

```
contains [-i | --index] [--] VALUE [LIST ...]
```

Exits 0 if VALUE is found in LIST, 1 otherwise. With `-i`, prints the 1-based index of the first match. Use `--` to search for values that look like flags.

### path

Manipulate file paths. A port of Fish's [`path` builtin](https://fishshell.com/docs/current/cmds/path.html).

```
path basename [-E | --strip-extension] [-q] [-Z] [PATH ...]
path change-extension EXT [PATH ...]
path dirname [-q] [-Z] [PATH ...]
path extension [-q] [-Z] [PATH ...]
path filter [-v] [-f|-d|-l] [-t TYPE] [-r|-w|-x] [-p PERM] [-q] [-Z] [PATH ...]
path is [-v] [-f|-d|-l] [-t TYPE] [-r|-w|-x] [-p PERM] [PATH ...]
path mtime [--relative] [-R] [PATH ...]
path normalize [-q] [-Z] [PATH ...]
path resolve [-q] [-Z] [PATH ...]
path sort [-r] [-u] [--key=KEY] [-Z] [PATH ...]
```

Paths starting with `-` are prefixed with `./` in normalize/filter output to prevent flag confusion.

### math

Evaluate arithmetic expressions. Extends Fish's [`math` builtin](https://fishshell.com/docs/current/cmds/math.html) with additional functions.

```
math [-s SCALE | --scale=N] [-b BASE | --base=N] EXPRESSION ...
```

Reads from stdin if no expression given. Supports:
- Arithmetic: `+` `-` `*` `/` `%` `^`
- Functions: `sin` `cos` `tan` `asin` `acos` `atan` `atan2` `sinh` `cosh` `tanh` `sqrt` `exp` `ln` `log` `log2` `log10` `abs` `floor` `ceil` `round` `roundf` `floorf` `ceilf` `pow` `max` `min` `fac` `ncr` `npr` `gcd` `bitand` `bitor` `bitxor`
- Constants: `pi` `e` `tau` `inf`
- Base output: `--base=16` (hex with `0x` prefix), `--base=octal`, `--base=2`
- Note: `log` is base-10; use `ln` for natural log. Comparison operators (`<` `>` `<=` `>=`) are not supported — use `test` instead.

### argparse

Parse command-line arguments and emit eval-able shell code. Deviates from Fish's [`argparse`](https://fishshell.com/docs/current/cmds/argparse.html) to support multiple shells.

```
argparse [--shell=SHELL] [--name=NAME] [--no-local]
         [--stop-nonopt] [--ignore-unknown] [--strict-longopts]
         [--min-args=N] [--max-args=N]
         SPEC ... -- ARG ...
```

Emits `eval`-able code that sets `_flag_NAME` variables and resets `$@` (or `$argv` in fish). Use inside functions for automatic scoping; at top level variables are global.

**Shell detection**: auto-detected from `$ZSH_VERSION`, `$BASH_VERSION`, or `$SHELL`. Override with `--shell=bash|zsh|sh|fish|elvish`.

**Spec format** (same as fish):
- `h/help` — short `-h` and long `--help`, boolean
- `/help` — long only
- `h` — short only
- `n/name=` — required value
- `n/name=?` — optional value
- `n/name=+` — required, repeatable
- `n/name=*` — optional, repeatable

**Usage**:
```bash
# bash/zsh — inside a function
myfunc() {
    eval "$(argparse 'h/help' 'v/verbose' 'n/name=' -- "$@")"
    [[ -n "$_flag_help" ]] && { echo "Usage: ..."; return; }
    echo "name: $_flag_name, verbose: $_flag_verbose"
}

# sh — works in functions; use --no-local at top level
eval "$(argparse --shell=sh --no-local 'v/verbose' -- "$@")"
```

**Multi-value flags in sh**: stored RS-delimited (ASCII 0x1E). Iterate with:
```sh
IFS=$(printf '\036'); for item in $_flag_names; do echo "$item"; done; unset IFS
```

## Description

Input is taken from command-line arguments unless stdin is connected to a pipe or file, in which case it is read line by line.

Most commands accept `-q` / `--quiet` to suppress output while preserving exit status.

These utilities aim for [fish shell](https://fishshell.com) parity where applicable. Some behaviors intentionally deviate to better support bash/zsh/sh.

## See also

- [fish `string` docs](https://fishshell.com/docs/current/cmds/string.html)
- [fish `count` docs](https://fishshell.com/docs/current/cmds/count.html)
- [fish `contains` docs](https://fishshell.com/docs/current/cmds/contains.html)
- [fish `path` docs](https://fishshell.com/docs/current/cmds/path.html)
- [fish `math` docs](https://fishshell.com/docs/current/cmds/math.html)
- [fish `argparse` docs](https://fishshell.com/docs/current/cmds/argparse.html)

# fish-utils

Fish shell utilities ported to Go for use in bash, zsh, and other shells.

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

Prints the number of arguments, or lines from stdin if no arguments are given.
Exits 0 if count > 0, 1 if count == 0.

### contains

Check if a value is in a list. A port of Fish's [`contains` builtin](https://fishshell.com/docs/current/cmds/contains.html).

```
contains [-i | --index] VALUE [LIST ...]
```

Exits 0 if VALUE is found in LIST, 1 otherwise. With `-i`, prints the 1-based index of the first match.

## Description

STRING/argument input is taken from the command line unless stdin is connected to a pipe or file, in which case it is read line by line.

Most commands accept `-q` / `--quiet` to suppress output while preserving exit status.

These utilities aim for [fish shell](https://fishshell.com) parity. See individual command help (`command --help`) for details.

## See also

- [fish `string` docs](https://fishshell.com/docs/current/cmds/string.html)
- [fish `count` docs](https://fishshell.com/docs/current/cmds/count.html)
- [fish `contains` docs](https://fishshell.com/docs/current/cmds/contains.html)

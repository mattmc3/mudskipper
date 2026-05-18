# string - manipulate strings

A shameless clone of Fish's `string` utility for use in other shells.

## Synopsis

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

## Description

`string` performs operations on strings.

STRING arguments are taken from the command line unless standard input is connected to a pipe or a file, in which case they are read from standard input, one STRING per line. It is an error to supply STRING arguments on the command line and on standard input.

Most subcommands accept a `-q` or `--quiet` switch, which suppresses the usual output but exits with the documented status.

This tool is modeled after the [fish shell `string` builtin](https://fishshell.com/docs/current/cmds/string.html) and aims for parity with its behavior.

## See also

- [fish `string` documentation](https://fishshell.com/docs/current/cmds/string.html)
- [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) — build, test, and project structure

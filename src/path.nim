import os, ./mudskipper_command

# PathCommand inherits from MudskipperCommand
type PathCommand* = ref object of MudskipperCommand

method usage*(self: PathCommand) =
  self.output("""path basename GENERAL_OPTIONS [(-E | --no-extension)] [PATH ...]
path dirname GENERAL_OPTIONS  [PATH ...]
path extension GENERAL_OPTIONS [PATH ...]
path filter GENERAL_OPTIONS [-v | --invert]
    [-d] [-f] [-l] [-r] [-w] [-x]
    [(-t | --type) TYPE] [(-p | --perm) PERMISSION] [PATH ...]
path is GENERAL_OPTIONS [(-v | --invert)] [(-t | --type) TYPE]
    [-d] [-f] [-l] [-r] [-w] [-x]
    [(-p | --perm) PERMISSION] [PATH ...]
path mtime GENERAL_OPTIONS [(-R | --relative)] [PATH ...]
path normalize GENERAL_OPTIONS [PATH ...]
path resolve GENERAL_OPTIONS [PATH ...]
path change-extension GENERAL_OPTIONS EXTENSION [PATH ...]
path sort GENERAL_OPTIONS [-r | --reverse]
    [-u | --unique] [--key=basename|dirname|path] [PATH ...]

GENERAL_OPTIONS
    [-z | --null-in] [-Z | --null-out] [-q | --quiet]
""")

proc basename*(self: PathCommand, args: seq[string]): int =
  # path basename [(-E|--no-extension)] [PATH ...]; args now excludes subcommand
  var noExt = false
  var paths: seq[string] = @[]
  for a in args:
    if a == "-E" or a == "--no-extension": noExt = true
    elif a == "-h" or a == "--help":
      self.output("Usage: path basename [(-E|--no-extension)] [PATH ...]\n")
      return 0
    else: paths.add a
  for p in paths:
    let (_, name, ext) = splitFile(p)
    if noExt: self.output(name & "\n")
    else: self.output(name & ext & "\n")
  return 0

proc dirname*(self: PathCommand, args: seq[string]): int =
  # path dirname [PATH ...]; args now excludes subcommand
  var paths: seq[string] = @[]
  for a in args:
    if a == "-h" or a == "--help":
      self.output("Usage: path dirname [PATH ...]\n")
      return 0
    else: paths.add a
  for p in paths: self.output(parentDir(p) & "\n")
  return 0

proc extension*(self: PathCommand, args: seq[string]): int =
  # path extension [PATH ...]; args now excludes subcommand
  var paths: seq[string] = @[]
  for a in args:
    if a == "-h" or a == "--help":
      self.output("Usage: path extension [PATH ...]\n")
      return 0
    else: paths.add a
  for p in paths:
    let (_, _, ext) = splitFile(p)
    self.output(ext & "\n")
  return 0

proc run*(self: PathCommand, args: seq[string]): int =
  if args.len == 0:
    self.outerr("path: missing subcommand\n")
    return 2
  let subcmd = args[0]
  let rest = if args.len > 1: args[1..^1] else: @[]
  case subcmd
  of "-h", "--help", "help":
    self.usage()
    return 0
  of "basename":
    return self.basename(rest)
  of "dirname":
    return self.dirname(rest)
  of "extension":
    return self.extension(rest)
  else:
    self.outerr("path: " & subcmd & ": invalid subcommand\n")
    return 2

when isMainModule:
  import os
  var cmd = PathCommand()
  quit cmd.run(commandLineParams())

import ./mudskipper_command
import std/parseopt

# ContainsCommand inherits base command
# Refactored command object (writer fields inherited)
type ContainsCommand* = ref object of MudskipperCommand

method usage*(self: ContainsCommand) =
  self.output("""
Usage: contains [-i|--index|-0|--index0] [--] NEEDLE HAYSTACK...
  -i, --index      Print the 1-based index of NEEDLE in HAYSTACK
  -0, --index0     Print the 0-based index of NEEDLE in HAYSTACK
  -h, --help       Show this help message
Returns 0 if NEEDLE is found, 1 otherwise.
""")

# Method performing the contains logic using parseopt
proc contains*(self: ContainsCommand, args: seq[string]): int =
  var indexFlag = false
  var indexStart = 1
  var positionals: seq[string]
  var optParser = initOptParser(args)

  for kind, key, val in optParser.getopt():
    case kind
    of cmdArgument:
      # Look for a non-flag and bail, prepending the arg we just consumed
      positionals = @[key] & optParser.remainingArgs
      break
    of cmdLongOption, cmdShortOption:
      case key
      of "h", "help":
        self.usage()
        return 0
      of "i", "index":
        indexFlag = true
        indexStart = 1
      of "0", "index0":
        indexFlag = true
        indexStart = 0
      of "":
        # Look for "--" and bail
        positionals = optParser.remainingArgs
        break
      else:
        let prefix = (if kind == cmdShortOption: "-" else: "--")
        self.outerr("contains: Unknown option " & prefix & key & "\n")
        return 2
    of cmdEnd:
      break

  if positionals.len == 0:
    self.outerr("contains: Key not specified\n")
    return 2
  let needle = positionals[0]
  for idx, val in positionals[1..^1]:
    if val == needle:
      if indexFlag:
        self.output($(idx + indexStart) & "\n")
      return 0
  return 1

# Convenience free function preserving previous API
proc contains*(args: seq[string]): int =
  var cmd = ContainsCommand()  # default writers nil -> stdout/stderr
  return cmd.contains(args)

when isMainModule:
  import os
  quit(contains(commandLineParams()))

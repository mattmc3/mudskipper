# Swappable writers
type Writer = proc (s: string)

# Refactored command object (writer fields exported so tests can assign)
type ContainsCommand* = object
  outWriter*: Writer
  errWriter*: Writer

# Fallback writer helpers
proc output(self: ContainsCommand, s: string) =
  if self.outWriter == nil: stdout.write(s) else: self.outWriter(s)
proc outerr(self: ContainsCommand, s: string) =
  if self.errWriter == nil: stderr.write(s) else: self.errWriter(s)

proc usage*(self: ContainsCommand) =
  self.output("""
Usage: contains [-i|--index|-0|--index0] [--] NEEDLE HAYSTACK...
  -i, --index      Print the 1-based index of NEEDLE in HAYSTACK
  -0, --index0     Print the 0-based index of NEEDLE in HAYSTACK
  -h, --help       Show this help message
Returns 0 if NEEDLE is found, 1 otherwise.
""")

# Method performing the contains logic
proc contains*(self: ContainsCommand, args: seq[string]): int =
  if args.len == 0:
    self.usage()
    return 0

  var indexFlag = false
  var indexStart = 1
  var start = 0

  # Option parsing via case for clarity
  case args[0]
  of "-h", "--help":
    self.usage()
    return 0
  of "-i", "--index":
    indexFlag = true
    indexStart = 1
    start = 1
  of "-0", "--index0":
    indexFlag = true
    indexStart = 0
    start = 1
  else:
    discard

  if args.len > start and args[start] == "--":
    inc start
  if args.len <= start:
    self.outerr("contains: Key not specified\n"); return 2
  let needle = args[start]
  let firstHay = start + 1

  # Scan haystack using for loop
  for rel in 0 ..< (args.len - firstHay):
    if args[firstHay + rel] == needle:
      if indexFlag:
        self.output($(rel + indexStart) & "\n")
      return 0
  return 1

# Convenience free function preserving previous API
proc contains*(args: seq[string]): int =
  var cmd: ContainsCommand
  return cmd.contains(args)


when isMainModule:
  import os
  quit(contains(commandLineParams()))

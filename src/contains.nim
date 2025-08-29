import os

proc printHelp() =
  echo """
Usage: contains [-i|--index|-0|--index0] NEEDLE HAYSTACK...
  -i, --index      Print the 1-based index of NEEDLE in HAYSTACK
  -0, --index0     Print the 0-based index of NEEDLE in HAYSTACK
  -h, --help       Show this help message
Returns 0 if NEEDLE is found, 1 otherwise.
"""

# Rewritten to avoid slicing / allocations and to early-exit
proc contains*(args: seq[string]): int =
  if args.len == 0 or args[0] == "-h" or args[0] == "--help":
    printHelp()
    return 0

  var indexFlag = false
  var indexStart = 1
  var start = 0

  if args.len > 0 and (args[0] == "-i" or args[0] == "--index"):
    indexFlag = true
    indexStart = 1
    start = 1
  elif args.len > 0 and (args[0] == "-0" or args[0] == "--index0"):
    indexFlag = true
    indexStart = 0
    start = 1

  if args.len <= start:
    stderr.write("contains: Key not specified\n")
    return 2

  let needle = args[start]
  var i = start + 1          # position in full args
  var rel = 0                # index within haystack list
  while i < args.len:
    if args[i] == needle:
      if indexFlag:
        stdout.write($(rel + indexStart), "\n")
      return 0
    inc i
    inc rel
  return 1

when isMainModule:
  quit(contains(commandLineParams()))
